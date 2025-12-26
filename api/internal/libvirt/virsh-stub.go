//go:build !libvirt
// +build !libvirt

package libvirt

import (
	"context"
	"errors"
	"time"
)

// ErrLibvirtNotAvailable is returned by all stub methods when libvirt support is not compiled in.
var ErrLibvirtNotAvailable = errors.New("libvirt support not available: rebuild with -tags libvirt")

// Manager defines the VM orchestration operations we support against libvirt/KVM via virsh.
type Manager interface {
	// CloneVM creates a linked-clone VM from a golden base image and defines a libvirt domain for it.
	// cpu and memoryMB are the VM shape. network is the libvirt network name (e.g., "default").
	CloneVM(ctx context.Context, baseImage, newVMName string, cpu, memoryMB int, network string) (DomainRef, error)

	// CloneFromVM creates a linked-clone VM from an existing VM's disk.
	// It looks up the source VM by name in libvirt, retrieves its disk path,
	// and creates an overlay pointing to that disk as the backing file.
	CloneFromVM(ctx context.Context, sourceVMName, newVMName string, cpu, memoryMB int, network string) (DomainRef, error)

	// InjectSSHKey injects an SSH public key for a user into the VM disk before boot.
	// The mechanism is determined by configuration (e.g., virt-customize or cloud-init seed).
	InjectSSHKey(ctx context.Context, sandboxName, username, publicKey string) error

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
// This is a stub implementation that returns errors when libvirt is not available.
type VirshManager struct {
	cfg Config
}

// NewVirshManager creates a new VirshManager with the provided config.
// Note: This stub implementation will return errors for all operations.
func NewVirshManager(cfg Config) *VirshManager {
	return &VirshManager{cfg: cfg}
}

// NewFromEnv builds a Config from environment variables and returns a manager.
// Note: This stub implementation will return errors for all operations.
func NewFromEnv() *VirshManager {
	cfg := Config{
		DefaultVCPUs:    2,
		DefaultMemoryMB: 2048,
	}
	return NewVirshManager(cfg)
}

// CloneVM is a stub that returns an error when libvirt is not available.
func (m *VirshManager) CloneVM(ctx context.Context, baseImage, newVMName string, cpu, memoryMB int, network string) (DomainRef, error) {
	return DomainRef{}, ErrLibvirtNotAvailable
}

// CloneFromVM is a stub that returns an error when libvirt is not available.
func (m *VirshManager) CloneFromVM(ctx context.Context, sourceVMName, newVMName string, cpu, memoryMB int, network string) (DomainRef, error) {
	return DomainRef{}, ErrLibvirtNotAvailable
}

// InjectSSHKey is a stub that returns an error when libvirt is not available.
func (m *VirshManager) InjectSSHKey(ctx context.Context, sandboxName, username, publicKey string) error {
	return ErrLibvirtNotAvailable
}

// StartVM is a stub that returns an error when libvirt is not available.
func (m *VirshManager) StartVM(ctx context.Context, vmName string) error {
	return ErrLibvirtNotAvailable
}

// StopVM is a stub that returns an error when libvirt is not available.
func (m *VirshManager) StopVM(ctx context.Context, vmName string, force bool) error {
	return ErrLibvirtNotAvailable
}

// DestroyVM is a stub that returns an error when libvirt is not available.
func (m *VirshManager) DestroyVM(ctx context.Context, vmName string) error {
	return ErrLibvirtNotAvailable
}

// CreateSnapshot is a stub that returns an error when libvirt is not available.
func (m *VirshManager) CreateSnapshot(ctx context.Context, vmName, snapshotName string, external bool) (SnapshotRef, error) {
	return SnapshotRef{}, ErrLibvirtNotAvailable
}

// DiffSnapshot is a stub that returns an error when libvirt is not available.
func (m *VirshManager) DiffSnapshot(ctx context.Context, vmName, fromSnapshot, toSnapshot string) (*FSComparePlan, error) {
	return nil, ErrLibvirtNotAvailable
}

// GetIPAddress is a stub that returns an error when libvirt is not available.
func (m *VirshManager) GetIPAddress(ctx context.Context, vmName string, timeout time.Duration) (string, error) {
	return "", ErrLibvirtNotAvailable
}
