package libvirt

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

// Manager defines the VM orchestration operations we support against libvirt/KVM via virsh.
type Manager interface {
	// CloneVM creates a linked-clone VM from a golden base image and defines a libvirt domain for it.
	// cpu and memoryMB are the VM shape. network is the libvirt network name (e.g., "default").
	CloneVM(ctx context.Context, baseImage, newVMName string, cpu, memoryMB int, network string) (DomainRef, error)

	// InjectSSHKey injects an SSH public key for a user into the VM disk before boot.
	// The mechanism is determined by configuration (e.g., virt-customize or cloud-init seed).
	InjectSSHKey(ctx context.Context, vmName, username, publicKey string) error

	// StartVM boots a defined domain.
	StartVM(ctx context.Context, vmName string) error

	// StopVM gracefully shuts down a domain, or forces if force is true.
	StopVM(ctx context.Context, vmName string, force bool) error

	// DestroyVM undefines the domain and removes its workspace (overlay files, domain XML, seeds).
	// If the domain is running, it will be destroyed first.
	DestroyVM(ctx context.Context, vmName string) error

	// CreateSnapshot creates a snapshot with the given name.
	// If external is true, attempts a disk-only external snapshot.
	CreateSnapshot(ctx context.Context, vmName, snapshotName string, external bool) (SnapshotRef, error)

	// DiffSnapshot prepares a plan to compare two snapshots' filesystems.
	// The returned plan includes advice or prepared mounts where possible.
	DiffSnapshot(ctx context.Context, vmName, fromSnapshot, toSnapshot string) (*FSComparePlan, error)

	// GetIPAddress attempts to fetch the VM's primary IP via libvirt leases.
	GetIPAddress(ctx context.Context, vmName string, timeout time.Duration) (string, error)
}

// Config controls how the virsh-based manager interacts with the host.
type Config struct {
	LibvirtURI            string // e.g., qemu:///system
	BaseImageDir          string // e.g., /var/lib/libvirt/images/base
	WorkDir               string // e.g., /var/lib/libvirt/images/jobs
	DefaultNetwork        string // e.g., default
	SSHKeyInjectMethod    string // "virt-customize" or "cloud-init"
	CloudInitMetaTemplate string // optional meta-data template for cloud-init seed

	// Optional explicit paths to binaries; if empty these are looked up in PATH.
	VirshPath         string
	QemuImgPath       string
	VirtCustomizePath string
	QemuNbdPath       string

	// Domain defaults
	DefaultVCPUs    int
	DefaultMemoryMB int
}

// DomainRef is a minimal reference to a libvirt domain (VM).
type DomainRef struct {
	Name string
	UUID string
}

// SnapshotRef references a snapshot created for a domain.
type SnapshotRef struct {
	Name string
	// Kind: "INTERNAL" or "EXTERNAL"
	Kind string
	// Ref is driver-specific; could be an internal UUID or a file path for external snapshots.
	Ref string
}

// FSComparePlan describes a plan for diffing two snapshots' filesystems.
type FSComparePlan struct {
	VMName       string
	FromSnapshot string
	ToSnapshot   string

	// Best-effort mount points (if prepared); may be empty strings when not mounted automatically.
	FromMount string
	ToMount   string

	// Devices or files used; informative.
	FromRef string
	ToRef   string

	// Free-form notes with instructions if the manager couldn't mount automatically.
	Notes []string
}

// VirshManager implements Manager using virsh/qemu-img/qemu-nbd/virt-customize and simple domain XML.
type VirshManager struct {
	cfg Config
}

// NewVirshManager creates a new VirshManager with the provided config.
func NewVirshManager(cfg Config) *VirshManager {
	// Fill sensible defaults
	if cfg.DefaultVCPUs == 0 {
		cfg.DefaultVCPUs = 2
	}
	if cfg.DefaultMemoryMB == 0 {
		cfg.DefaultMemoryMB = 2048
	}
	return &VirshManager{cfg: cfg}
}

// NewFromEnv builds a Config from environment variables and returns a manager.
// LIBVIRT_URI, BASE_IMAGE_DIR, SANDBOX_WORKDIR, LIBVIRT_NETWORK, SSH_KEY_INJECT_METHOD
func NewFromEnv() *VirshManager {
	cfg := Config{
		LibvirtURI:         getenvDefault("LIBVIRT_URI", "qemu:///system"),
		BaseImageDir:       getenvDefault("BASE_IMAGE_DIR", "/var/lib/libvirt/images/base"),
		WorkDir:            getenvDefault("SANDBOX_WORKDIR", "/var/lib/libvirt/images/jobs"),
		DefaultNetwork:     getenvDefault("LIBVIRT_NETWORK", "default"),
		SSHKeyInjectMethod: getenvDefault("SSH_KEY_INJECT_METHOD", "virt-customize"),
		DefaultVCPUs:       intFromEnv("DEFAULT_VCPUS", 2),
		DefaultMemoryMB:    intFromEnv("DEFAULT_MEMORY_MB", 2048),
	}
	return NewVirshManager(cfg)
}

func (m *VirshManager) CloneVM(ctx context.Context, baseImage, newVMName string, cpu, memoryMB int, network string) (DomainRef, error) {
	if newVMName == "" {
		return DomainRef{}, fmt.Errorf("new VM name is required")
	}
	if baseImage == "" {
		return DomainRef{}, fmt.Errorf("base image is required")
	}
	if cpu <= 0 {
		cpu = m.cfg.DefaultVCPUs
	}
	if memoryMB <= 0 {
		memoryMB = m.cfg.DefaultMemoryMB
	}
	if network == "" {
		network = m.cfg.DefaultNetwork
	}

	basePath := filepath.Join(m.cfg.BaseImageDir, baseImage)
	if _, err := os.Stat(basePath); err != nil {
		return DomainRef{}, fmt.Errorf("base image not accessible: %s: %w", basePath, err)
	}

	jobDir := filepath.Join(m.cfg.WorkDir, newVMName)
	if err := os.MkdirAll(jobDir, 0o755); err != nil {
		return DomainRef{}, fmt.Errorf("create job dir: %w", err)
	}

	overlayPath := filepath.Join(jobDir, "disk-overlay.qcow2")
	qemuImg := m.binPath("qemu-img", m.cfg.QemuImgPath)
	if _, err := m.run(ctx, qemuImg, "create", "-f", "qcow2", "-F", "qcow2", "-b", basePath, overlayPath); err != nil {
		return DomainRef{}, fmt.Errorf("create overlay: %w", err)
	}

	// Create minimal domain XML referencing overlay disk and network.
	xmlPath := filepath.Join(jobDir, "domain.xml")
	xml, err := renderDomainXML(domainXMLParams{
		Name:      newVMName,
		MemoryMB:  memoryMB,
		VCPUs:     cpu,
		DiskPath:  overlayPath,
		Network:   network,
		BootOrder: []string{"hd", "cdrom", "network"},
	})
	if err != nil {
		return DomainRef{}, fmt.Errorf("render domain xml: %w", err)
	}
	if err := os.WriteFile(xmlPath, []byte(xml), 0o644); err != nil {
		return DomainRef{}, fmt.Errorf("write domain xml: %w", err)
	}

	// virsh define
	virsh := m.binPath("virsh", m.cfg.VirshPath)
	if _, err := m.run(ctx, virsh, "--connect", m.cfg.LibvirtURI, "define", xmlPath); err != nil {
		return DomainRef{}, fmt.Errorf("virsh define: %w", err)
	}

	// Fetch UUID
	out, err := m.run(ctx, virsh, "--connect", m.cfg.LibvirtURI, "domuuid", newVMName)
	if err != nil {
		// Best-effort: If domuuid fails, we still return Name.
		return DomainRef{Name: newVMName}, nil
	}
	return DomainRef{Name: newVMName, UUID: strings.TrimSpace(out)}, nil
}

func (m *VirshManager) InjectSSHKey(ctx context.Context, vmName, username, publicKey string) error {
	if vmName == "" {
		return fmt.Errorf("vmName is required")
	}
	if username == "" {
		username = defaultGuestUser(vmName)
	}
	if strings.TrimSpace(publicKey) == "" {
		return fmt.Errorf("publicKey is required")
	}

	jobDir := filepath.Join(m.cfg.WorkDir, vmName)
	overlay := filepath.Join(jobDir, "disk-overlay.qcow2")
	if _, err := os.Stat(overlay); err != nil {
		return fmt.Errorf("overlay not found for VM %s: %w", vmName, err)
	}

	switch strings.ToLower(m.cfg.SSHKeyInjectMethod) {
	case "virt-customize":
		// Requires libguestfs tools on host.
		virtCustomize := m.binPath("virt-customize", m.cfg.VirtCustomizePath)
		// Ensure account exists and inject key. This is offline before first boot.
		cmdArgs := []string{
			"-a", overlay,
			"--run-command", fmt.Sprintf("id -u %s >/dev/null 2>&1 || useradd -m -s /bin/bash %s", shEscape(username), shEscape(username)),
			"--ssh-inject", fmt.Sprintf("%s:string:%s", username, publicKey),
		}
		if _, err := m.run(ctx, virtCustomize, cmdArgs...); err != nil {
			return fmt.Errorf("virt-customize inject: %w", err)
		}
	case "cloud-init":
		// Build a NoCloud seed with the provided key and attach as CD-ROM.
		seedISO := filepath.Join(jobDir, "seed.iso")
		if err := m.buildCloudInitSeed(ctx, vmName, username, publicKey, seedISO); err != nil {
			return fmt.Errorf("build cloud-init seed: %w", err)
		}
		// Attach seed ISO to domain XML (adds a CDROM) and redefine the domain.
		xmlPath := filepath.Join(jobDir, "domain.xml")
		if err := m.attachISOToDomainXML(xmlPath, seedISO); err != nil {
			return fmt.Errorf("attach seed iso to domain xml: %w", err)
		}
		virsh := m.binPath("virsh", m.cfg.VirshPath)
		if _, err := m.run(ctx, virsh, "--connect", m.cfg.LibvirtURI, "define", xmlPath); err != nil {
			return fmt.Errorf("re-define domain with seed: %w", err)
		}
	default:
		return fmt.Errorf("unsupported SSHKeyInjectMethod: %s", m.cfg.SSHKeyInjectMethod)
	}
	return nil
}

func (m *VirshManager) StartVM(ctx context.Context, vmName string) error {
	if vmName == "" {
		return fmt.Errorf("vmName is required")
	}
	virsh := m.binPath("virsh", m.cfg.VirshPath)
	_, err := m.run(ctx, virsh, "--connect", m.cfg.LibvirtURI, "start", vmName)
	return err
}

func (m *VirshManager) StopVM(ctx context.Context, vmName string, force bool) error {
	if vmName == "" {
		return fmt.Errorf("vmName is required")
	}
	virsh := m.binPath("virsh", m.cfg.VirshPath)
	if force {
		_, err := m.run(ctx, virsh, "--connect", m.cfg.LibvirtURI, "destroy", vmName)
		return err
	}
	_, err := m.run(ctx, virsh, "--connect", m.cfg.LibvirtURI, "shutdown", vmName)
	return err
}

func (m *VirshManager) DestroyVM(ctx context.Context, vmName string) error {
	if vmName == "" {
		return fmt.Errorf("vmName is required")
	}
	virsh := m.binPath("virsh", m.cfg.VirshPath)
	// Best-effort destroy if running
	_, _ = m.run(ctx, virsh, "--connect", m.cfg.LibvirtURI, "destroy", vmName)
	// Undefine
	if _, err := m.run(ctx, virsh, "--connect", m.cfg.LibvirtURI, "undefine", vmName); err != nil {
		// continue to remove files even if undefine fails
		_ = err
	}
	// Remove workspace
	jobDir := filepath.Join(m.cfg.WorkDir, vmName)
	if err := os.RemoveAll(jobDir); err != nil {
		return fmt.Errorf("cleanup job dir: %w", err)
	}
	return nil
}

func (m *VirshManager) CreateSnapshot(ctx context.Context, vmName, snapshotName string, external bool) (SnapshotRef, error) {
	if vmName == "" || snapshotName == "" {
		return SnapshotRef{}, fmt.Errorf("vmName and snapshotName are required")
	}
	virsh := m.binPath("virsh", m.cfg.VirshPath)

	if external {
		// External disk-only snapshot.
		jobDir := filepath.Join(m.cfg.WorkDir, vmName)
		snapPath := filepath.Join(jobDir, fmt.Sprintf("snap-%s.qcow2", snapshotName))
		// NOTE: This is a simplified attempt; real-world disk-only snapshots may need
		// additional options and disk target identification.
		args := []string{
			"--connect", m.cfg.LibvirtURI, "snapshot-create-as", vmName, snapshotName,
			"--disk-only", "--atomic", "--no-metadata",
			"--diskspec", fmt.Sprintf("vda,file=%s", snapPath),
		}
		if _, err := m.run(ctx, virsh, args...); err != nil {
			return SnapshotRef{}, fmt.Errorf("external snapshot create: %w", err)
		}
		return SnapshotRef{Name: snapshotName, Kind: "EXTERNAL", Ref: snapPath}, nil
	}

	// Internal snapshot (managed by libvirt/qemu).
	if _, err := m.run(ctx, virsh, "--connect", m.cfg.LibvirtURI, "snapshot-create-as", vmName, snapshotName); err != nil {
		return SnapshotRef{}, fmt.Errorf("internal snapshot create: %w", err)
	}
	return SnapshotRef{Name: snapshotName, Kind: "INTERNAL", Ref: snapshotName}, nil
}

func (m *VirshManager) DiffSnapshot(ctx context.Context, vmName, fromSnapshot, toSnapshot string) (*FSComparePlan, error) {
	if vmName == "" || fromSnapshot == "" || toSnapshot == "" {
		return nil, fmt.Errorf("vmName, fromSnapshot and toSnapshot are required")
	}

	// Implementation shell:
	// Strategy options:
	// 1) For internal snapshots: use qemu-nbd with snapshot selection to mount and diff trees.
	// 2) For external snapshots: mount the two qcow2 snapshot files via qemu-nbd.
	//
	// Because snapshot storage varies, we return advisory plan data and notes.
	plan := &FSComparePlan{
		VMName:       vmName,
		FromSnapshot: fromSnapshot,
		ToSnapshot:   toSnapshot,
		Notes:        []string{},
	}

	// Attempt to detect external snapshot files in job dir.
	jobDir := filepath.Join(m.cfg.WorkDir, vmName)
	fromPath := filepath.Join(jobDir, fmt.Sprintf("snap-%s.qcow2", fromSnapshot))
	toPath := filepath.Join(jobDir, fmt.Sprintf("snap-%s.qcow2", toSnapshot))
	if fileExists(fromPath) && fileExists(toPath) {
		plan.FromRef = fromPath
		plan.ToRef = toPath
		plan.Notes = append(plan.Notes,
			"External snapshots detected. You can mount them with qemu-nbd and diff the trees.",
			fmt.Sprintf("sudo modprobe nbd max_part=16 && sudo qemu-nbd --connect=/dev/nbd0 %s", shEscape(fromPath)),
			fmt.Sprintf("sudo qemu-nbd --connect=/dev/nbd1 %s", shEscape(toPath)),
			"sudo mount /dev/nbd0p1 /mnt/from && sudo mount /dev/nbd1p1 /mnt/to",
			"Then run: sudo diff -ruN /mnt/from /mnt/to or use rsync --dry-run to list changes.",
			"Be sure to umount and disconnect nbd after.",
		)
		return plan, nil
	}

	// Fallback: internal snapshots guidance.
	plan.Notes = append(plan.Notes,
		"Internal snapshots assumed. Use qemu-nbd with -s to select snapshot, then mount and diff.",
		"For example: qemu-nbd may support --snapshot=<name> (varies by version) or use qemu-img to create temporary exports.",
		"Alternatively, boot the VM into each snapshot separately and export filesystem states.",
	)
	return plan, nil
}

func (m *VirshManager) GetIPAddress(ctx context.Context, vmName string, timeout time.Duration) (string, error) {
	if vmName == "" {
		return "", fmt.Errorf("vmName is required")
	}
	virsh := m.binPath("virsh", m.cfg.VirshPath)
	deadline := time.Now().Add(timeout)
	for {
		out, err := m.run(ctx, virsh, "--connect", m.cfg.LibvirtURI, "domifaddr", vmName, "--source", "lease")
		if err == nil {
			ip := parseDomIfAddrIPv4(out)
			if ip != "" {
				return ip, nil
			}
		}
		if time.Now().After(deadline) {
			break
		}
		time.Sleep(2 * time.Second)
	}
	return "", errors.New("ip address not found within timeout")
}

// --- Helpers ---

func (m *VirshManager) binPath(defaultName, override string) string {
	if override != "" {
		return override
	}
	return defaultName
}

func (m *VirshManager) run(ctx context.Context, bin string, args ...string) (string, error) {
	var stdout, stderr bytes.Buffer
	// Provide a default timeout if the context has none.
	if _, ok := ctx.Deadline(); !ok {
		ctx2, cancel := context.WithTimeout(ctx, 120*time.Second)
		defer cancel()
		ctx = ctx2
	}
	cmd := exec.CommandContext(ctx, bin, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	// Pass LIBVIRT_DEFAULT_URI for convenience when set.
	env := os.Environ()
	if m.cfg.LibvirtURI != "" {
		env = append(env, "LIBVIRT_DEFAULT_URI="+m.cfg.LibvirtURI)
	}
	cmd.Env = env

	err := cmd.Run()
	outStr := strings.TrimSpace(stdout.String())
	if err != nil {
		errStr := strings.TrimSpace(stderr.String())
		if errStr != "" {
			return outStr, fmt.Errorf("%s %s failed: %w: %s", bin, strings.Join(args, " "), err, errStr)
		}
		return outStr, fmt.Errorf("%s %s failed: %w", bin, strings.Join(args, " "), err)
	}
	return outStr, nil
}

func getenvDefault(k, def string) string {
	v := os.Getenv(k)
	if v == "" {
		return def
	}
	return v
}

func intFromEnv(k string, def int) int {
	v := os.Getenv(k)
	if v == "" {
		return def
	}
	var parsed int
	_, err := fmt.Sscanf(v, "%d", &parsed)
	if err != nil {
		return def
	}
	return parsed
}

func fileExists(p string) bool {
	st, err := os.Stat(p)
	return err == nil && !st.IsDir()
}

func shEscape(s string) string {
	// naive escape for use inside run-command; rely on controlled inputs.
	s = strings.ReplaceAll(s, `'`, `'\'\'`)
	return s
}

func defaultGuestUser(vmName string) string {
	// Heuristic default depending on distro naming conventions.
	// Adjust as needed by calling code.
	if strings.Contains(strings.ToLower(vmName), "ubuntu") {
		return "ubuntu"
	}
	if strings.Contains(strings.ToLower(vmName), "centos") || strings.Contains(strings.ToLower(vmName), "rhel") {
		return "centos"
	}
	return "cloud-user"
}

func parseDomIfAddrIPv4(s string) string {
	// virsh domifaddr output example:
	// Name       MAC address          Protocol     Address
	// ----------------------------------------------------------------------------
	// vnet0      52:54:00:6b:3c:86    ipv4         192.168.122.63/24
	lines := strings.Split(s, "\n")
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l == "" || strings.HasPrefix(l, "Name") || strings.HasPrefix(l, "-") {
			continue
		}
		parts := strings.Fields(l)
		if len(parts) >= 4 && parts[2] == "ipv4" {
			addr := parts[3]
			if i := strings.IndexByte(addr, '/'); i > 0 {
				return addr[:i]
			}
			return addr
		}
	}
	return ""
}

// --- Domain XML rendering ---

type domainXMLParams struct {
	Name      string
	MemoryMB  int
	VCPUs     int
	DiskPath  string
	Network   string
	BootOrder []string
}

func renderDomainXML(p domainXMLParams) (string, error) {
	// A minimal domain XML; adjust virtio model as needed by your environment.
	const tpl = `<?xml version="1.0" encoding="utf-8"?>
<domain type="kvm">
  <name>{{ .Name }}</name>
  <memory unit="MiB">{{ .MemoryMB }}</memory>
  <vcpu placement="static">{{ .VCPUs }}</vcpu>
  <os>
    <type arch="x86_64" machine="pc-q35-6.2">hvm</type>
    <boot dev="hd"/>
    <boot dev="cdrom"/>
  </os>
  <features>
    <acpi/>
    <apic/>
    <pae/>
  </features>
  <cpu mode="host-passthrough"/>
  <devices>
    <disk type="file" device="disk">
      <driver name="qemu" type="qcow2" cache="none"/>
      <source file="{{ .DiskPath }}"/>
      <target dev="vda" bus="virtio"/>
    </disk>
    <controller type="pci" model="pcie-root"/>
    <interface type="network">
      <source network="{{ .Network }}"/>
      <model type="virtio"/>
    </interface>
    <graphics type="vnc" autoport="yes" listen="0.0.0.0"/>
    <console type="pty"/>
    <input type="tablet" bus="usb"/>
    <rng model="virtio">
      <backend model="random">/dev/urandom</backend>
    </rng>
  </devices>
</domain>
`
	var b bytes.Buffer
	t := template.Must(template.New("domain").Parse(tpl))
	if err := t.Execute(&b, p); err != nil {
		return "", err
	}
	return b.String(), nil
}

// attachISOToDomainXML is a simple XML string replacement to add a CD-ROM pointing to seed ISO.
// For a production system, consider parsing XML and building a proper DOM.
func (m *VirshManager) attachISOToDomainXML(xmlPath, isoPath string) error {
	data, err := os.ReadFile(xmlPath)
	if err != nil {
		return err
	}
	xml := string(data)
	needle := "</devices>"
	cdrom := fmt.Sprintf(`
    <disk type="file" device="cdrom">
      <driver name="qemu" type="raw"/>
      <source file="%s"/>
      <target dev="sda" bus="sata"/>
      <readonly/>
    </disk>`, isoPath)
	if strings.Contains(xml, cdrom) {
		// already attached
		return nil
	}
	xml = strings.Replace(xml, needle, cdrom+"\n  "+needle, 1)
	return os.WriteFile(xmlPath, []byte(xml), 0o644)
}

// buildCloudInitSeed creates a NoCloud seed ISO with a single user and SSH key.
// Requires cloud-localds (cloud-image-utils) on the host if implemented via external tool.
// This implementation writes user-data/meta-data and attempts to use genisoimage or mkisofs.
func (m *VirshManager) buildCloudInitSeed(ctx context.Context, vmName, username, publicKey, outISO string) error {
	jobDir := filepath.Dir(outISO)
	userData := fmt.Sprintf(`#cloud-config
users:
  - name: %s
    sudo: ALL=(ALL) NOPASSWD:ALL
    groups: users, admin, sudo
    shell: /bin/bash
    ssh_authorized_keys:
      - %s
`, username, publicKey)

	metaData := fmt.Sprintf(`instance-id: %s
local-hostname: %s
`, vmName, vmName)

	userDataPath := filepath.Join(jobDir, "user-data")
	metaDataPath := filepath.Join(jobDir, "meta-data")
	if err := os.WriteFile(userDataPath, []byte(userData), 0o644); err != nil {
		return fmt.Errorf("write user-data: %w", err)
	}
	if err := os.WriteFile(metaDataPath, []byte(metaData), 0o644); err != nil {
		return fmt.Errorf("write meta-data: %w", err)
	}

	// Try cloud-localds if available
	if hasBin("cloud-localds") {
		if _, err := m.run(ctx, "cloud-localds", outISO, userDataPath, metaDataPath); err == nil {
			return nil
		}
	}

	// Fallback to genisoimage/mkisofs
	if hasBin("genisoimage") {
		// genisoimage -output seed.iso -volid cidata -joliet -rock user-data meta-data
		_, err := m.run(ctx, "genisoimage", "-output", outISO, "-volid", "cidata", "-joliet", "-rock", userDataPath, metaDataPath)
		return err
	}
	if hasBin("mkisofs") {
		_, err := m.run(ctx, "mkisofs", "-output", outISO, "-V", "cidata", "-J", "-R", userDataPath, metaDataPath)
		return err
	}

	return fmt.Errorf("cloud-init seed build tools not found: need cloud-localds or genisoimage/mkisofs")
}

func hasBin(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}
