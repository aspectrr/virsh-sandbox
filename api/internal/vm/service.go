package vm

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/google/uuid"

	"virsh-sandbox/internal/libvirt"
	"virsh-sandbox/internal/store"
)

// Service orchestrates libvirt operations and data persistence.
// It represents the main application layer for sandbox lifecycle, command exec,
// snapshotting, diffing, and artifact generation orchestration.
type Service struct {
	mgr       libvirt.Manager
	store     store.Store
	ssh       SSHRunner
	cfg       Config
	timeNowFn func() time.Time
}

// Config controls default VM parameters and timeouts used by the service.
type Config struct {
	// Default libvirt network name (e.g., "default") used when creating VMs.
	Network string

	// Default shape if not provided by callers.
	DefaultVCPUs    int
	DefaultMemoryMB int

	// CommandTimeout sets a default timeout for RunCommand when caller doesn't provide one.
	CommandTimeout time.Duration

	// IPDiscoveryTimeout controls how long StartSandbox waits for the VM IP (when requested).
	IPDiscoveryTimeout time.Duration
}

// Option configures the Service during construction.
type Option func(*Service)

// WithSSHRunner overrides the default SSH runner implementation.
func WithSSHRunner(r SSHRunner) Option {
	return func(s *Service) { s.ssh = r }
}

// WithTimeNow overrides the clock (useful for tests).
func WithTimeNow(fn func() time.Time) Option {
	return func(s *Service) { s.timeNowFn = fn }
}

// NewService constructs a VM service with the provided libvirt manager, store and config.
func NewService(mgr libvirt.Manager, st store.Store, cfg Config, opts ...Option) *Service {
	if cfg.DefaultVCPUs <= 0 {
		cfg.DefaultVCPUs = 2
	}
	if cfg.DefaultMemoryMB <= 0 {
		cfg.DefaultMemoryMB = 2048
	}
	if cfg.CommandTimeout <= 0 {
		cfg.CommandTimeout = 10 * time.Minute
	}
	if cfg.IPDiscoveryTimeout <= 0 {
		cfg.IPDiscoveryTimeout = 2 * time.Minute
	}
	s := &Service{
		mgr:       mgr,
		store:     st,
		cfg:       cfg,
		ssh:       &DefaultSSHRunner{},
		timeNowFn: time.Now,
	}
	for _, o := range opts {
		o(s)
	}
	return s
}

// CreateSandbox clones a VM from a base image and persists a Sandbox record.
//
// baseImage is the base qcow2 filename (e.g., ubuntu-22.04.qcow2) located in the host's base images dir.
// vmName is optional; if empty, a name will be generated.
// cpu and memoryMB are optional; if <=0 the service defaults are used.
func (s *Service) CreateSandbox(ctx context.Context, baseImage, agentID, vmName string, cpu, memoryMB int) (*store.Sandbox, error) {
	if strings.TrimSpace(baseImage) == "" {
		return nil, fmt.Errorf("baseImage is required")
	}
	if strings.TrimSpace(agentID) == "" {
		return nil, fmt.Errorf("agentID is required")
	}
	if cpu <= 0 {
		cpu = s.cfg.DefaultVCPUs
	}
	if memoryMB <= 0 {
		memoryMB = s.cfg.DefaultMemoryMB
	}
	if vmName == "" {
		vmName = fmt.Sprintf("sbx-%s", shortID())
	}

	jobID := fmt.Sprintf("JOB-%s", shortID())

	// Create the VM via libvirt manager
	_, err := s.mgr.CloneVM(ctx, baseImage, vmName, cpu, memoryMB, s.cfg.Network)
	if err != nil {
		return nil, fmt.Errorf("clone vm: %w", err)
	}

	sb := &store.Sandbox{
		ID:        fmt.Sprintf("SBX-%s", shortID()),
		JobID:     jobID,
		AgentID:   agentID,
		VMName:    vmName,
		BaseImage: baseImage,
		Network:   s.cfg.Network,
		State:     store.SandboxStateCreated,
		CreatedAt: s.timeNowFn().UTC(),
		UpdatedAt: s.timeNowFn().UTC(),
	}
	if err := s.store.CreateSandbox(ctx, sb); err != nil {
		return nil, fmt.Errorf("persist sandbox: %w", err)
	}
	return sb, nil
}

// InjectSSHKey injects a public key for a user into the VM disk prior to boot.
func (s *Service) InjectSSHKey(ctx context.Context, sandboxID, username, publicKey string) error {
	if strings.TrimSpace(sandboxID) == "" {
		return fmt.Errorf("sandboxID is required")
	}
	if strings.TrimSpace(username) == "" {
		return fmt.Errorf("username is required")
	}
	if strings.TrimSpace(publicKey) == "" {
		return fmt.Errorf("publicKey is required")
	}
	sb, err := s.store.GetSandbox(ctx, sandboxID)
	if err != nil {
		return err
	}
	if err := s.mgr.InjectSSHKey(ctx, sb.VMName, username, publicKey); err != nil {
		return fmt.Errorf("inject ssh key: %w", err)
	}
	sb.UpdatedAt = s.timeNowFn().UTC()
	return s.store.UpdateSandbox(ctx, sb)
}

// StartSandbox boots the VM and optionally waits for IP discovery.
// Returns the discovered IP if waitForIP is true and discovery succeeds (empty string otherwise).
func (s *Service) StartSandbox(ctx context.Context, sandboxID string, waitForIP bool) (string, error) {
	if strings.TrimSpace(sandboxID) == "" {
		return "", fmt.Errorf("sandboxID is required")
	}
	sb, err := s.store.GetSandbox(ctx, sandboxID)
	if err != nil {
		return "", err
	}

	if err := s.mgr.StartVM(ctx, sb.VMName); err != nil {
		_ = s.store.UpdateSandboxState(ctx, sb.ID, store.SandboxStateError, nil)
		return "", fmt.Errorf("start vm: %w", err)
	}

	// Update state -> STARTING
	if err := s.store.UpdateSandboxState(ctx, sb.ID, store.SandboxStateStarting, nil); err != nil {
		return "", err
	}

	var ip string
	if waitForIP {
		ip, err = s.mgr.GetIPAddress(ctx, sb.VMName, s.cfg.IPDiscoveryTimeout)
		if err != nil {
			// Still mark as running even if we couldn't discover the IP
			_ = s.store.UpdateSandboxState(ctx, sb.ID, store.SandboxStateRunning, nil)
			return "", fmt.Errorf("get ip: %w", err)
		}
		if err := s.store.UpdateSandboxState(ctx, sb.ID, store.SandboxStateRunning, &ip); err != nil {
			return "", err
		}
	} else {
		if err := s.store.UpdateSandboxState(ctx, sb.ID, store.SandboxStateRunning, nil); err != nil {
			return "", err
		}
	}

	return ip, nil
}

// StopSandbox gracefully shuts down the VM or forces if force is true.
func (s *Service) StopSandbox(ctx context.Context, sandboxID string, force bool) error {
	if strings.TrimSpace(sandboxID) == "" {
		return fmt.Errorf("sandboxID is required")
	}
	sb, err := s.store.GetSandbox(ctx, sandboxID)
	if err != nil {
		return err
	}
	if err := s.mgr.StopVM(ctx, sb.VMName, force); err != nil {
		return fmt.Errorf("stop vm: %w", err)
	}
	return s.store.UpdateSandboxState(ctx, sb.ID, store.SandboxStateStopped, sb.IPAddress)
}

// DestroySandbox forcibly destroys and undefines the VM and removes its workspace.
// The sandbox is then soft-deleted from the store.
func (s *Service) DestroySandbox(ctx context.Context, sandboxID string) error {
	if strings.TrimSpace(sandboxID) == "" {
		return fmt.Errorf("sandboxID is required")
	}
	sb, err := s.store.GetSandbox(ctx, sandboxID)
	if err != nil {
		return err
	}
	if err := s.mgr.DestroyVM(ctx, sb.VMName); err != nil {
		return fmt.Errorf("destroy vm: %w", err)
	}
	return s.store.DeleteSandbox(ctx, sandboxID)
}

// CreateSnapshot creates a snapshot and persists a Snapshot record.
func (s *Service) CreateSnapshot(ctx context.Context, sandboxID, name string, external bool) (*store.Snapshot, error) {
	if strings.TrimSpace(sandboxID) == "" || strings.TrimSpace(name) == "" {
		return nil, fmt.Errorf("sandboxID and name are required")
	}
	sb, err := s.store.GetSandbox(ctx, sandboxID)
	if err != nil {
		return nil, err
	}
	ref, err := s.mgr.CreateSnapshot(ctx, sb.VMName, name, external)
	if err != nil {
		return nil, fmt.Errorf("create snapshot: %w", err)
	}
	sn := &store.Snapshot{
		ID:        fmt.Sprintf("SNP-%s", shortID()),
		SandboxID: sb.ID,
		Name:      ref.Name,
		Kind:      snapshotKindFromString(ref.Kind),
		Ref:       ref.Ref,
		CreatedAt: s.timeNowFn().UTC(),
	}
	if err := s.store.CreateSnapshot(ctx, sn); err != nil {
		return nil, err
	}
	return sn, nil
}

// DiffSnapshots computes a normalized change set between two snapshots and persists a Diff.
// Note: This implementation currently aggregates command history into CommandsRun and
// leaves file/package/service diffs empty. A dedicated diff engine should populate these fields
// by mounting snapshots and computing differences.
func (s *Service) DiffSnapshots(ctx context.Context, sandboxID, from, to string) (*store.Diff, error) {
	if strings.TrimSpace(sandboxID) == "" || strings.TrimSpace(from) == "" || strings.TrimSpace(to) == "" {
		return nil, fmt.Errorf("sandboxID, from, to are required")
	}
	sb, err := s.store.GetSandbox(ctx, sandboxID)
	if err != nil {
		return nil, err
	}

	// Best-effort: get a plan (notes/instructions) from manager; ignore failure.
	_, _ = s.mgr.DiffSnapshot(ctx, sb.VMName, from, to)

	// For now, compose CommandsRun from command history as partial diff signal.
	cmds, err := s.store.ListCommands(ctx, sandboxID, &store.ListOptions{OrderBy: "started_at", Asc: true})
	if err != nil && !errors.Is(err, store.ErrNotFound) {
		return nil, fmt.Errorf("list commands: %w", err)
	}
	var cr []store.CommandSummary
	for _, c := range cmds {
		cr = append(cr, store.CommandSummary{
			Cmd:      c.Command,
			ExitCode: c.ExitCode,
			At:       c.EndedAt,
		})
	}

	diff := &store.Diff{
		ID:           fmt.Sprintf("DIF-%s", shortID()),
		SandboxID:    sandboxID,
		FromSnapshot: from,
		ToSnapshot:   to,
		DiffJSON: store.ChangeDiff{
			FilesModified:   []string{},
			FilesAdded:      []string{},
			FilesRemoved:    []string{},
			PackagesAdded:   []store.PackageInfo{},
			PackagesRemoved: []store.PackageInfo{},
			ServicesChanged: []store.ServiceChange{},
			CommandsRun:     cr,
		},
		CreatedAt: s.timeNowFn().UTC(),
	}
	if err := s.store.SaveDiff(ctx, diff); err != nil {
		return nil, err
	}
	return diff, nil
}

// RunCommand executes a command inside the sandbox via SSH.
// The username and privateKeyPath are required for SSH auth. The service obtains
// the VM IP from the sandbox record or discovers it via libvirt if missing.
func (s *Service) RunCommand(ctx context.Context, sandboxID, username, privateKeyPath, command string, timeout time.Duration, env map[string]string) (*store.Command, error) {
	if strings.TrimSpace(sandboxID) == "" {
		return nil, fmt.Errorf("sandboxID is required")
	}
	if strings.TrimSpace(username) == "" {
		return nil, fmt.Errorf("username is required")
	}
	if strings.TrimSpace(privateKeyPath) == "" {
		return nil, fmt.Errorf("privateKeyPath is required")
	}
	if strings.TrimSpace(command) == "" {
		return nil, fmt.Errorf("command is required")
	}
	if timeout <= 0 {
		timeout = s.cfg.CommandTimeout
	}

	sb, err := s.store.GetSandbox(ctx, sandboxID)
	if err != nil {
		return nil, err
	}
	ip := ""
	if sb.IPAddress != nil && *sb.IPAddress != "" {
		ip = *sb.IPAddress
	} else {
		ip, err = s.mgr.GetIPAddress(ctx, sb.VMName, s.cfg.IPDiscoveryTimeout)
		if err != nil {
			return nil, fmt.Errorf("discover ip: %w", err)
		}
		// Persist discovered IP for subsequent calls
		if err := s.store.UpdateSandboxState(ctx, sb.ID, sb.State, &ip); err != nil {
			return nil, fmt.Errorf("persist ip: %w", err)
		}
	}

	cmdID := fmt.Sprintf("CMD-%s", shortID())
	now := s.timeNowFn().UTC()

	// Encode environment for persistence.
	var envJSON *string
	if len(env) > 0 {
		b, _ := json.Marshal(env)
		tmp := string(b)
		envJSON = &tmp
	}

	stdout, stderr, code, runErr := s.ssh.Run(ctx, ip, username, privateKeyPath, commandWithEnv(command, env), timeout, env)

	cmd := &store.Command{
		ID:        cmdID,
		SandboxID: sandboxID,
		Command:   command,
		EnvJSON:   envJSON,
		Stdout:    stdout,
		Stderr:    stderr,
		ExitCode:  code,
		StartedAt: now,
		EndedAt:   s.timeNowFn().UTC(),
	}
	if err := s.store.SaveCommand(ctx, cmd); err != nil {
		return nil, fmt.Errorf("save command: %w", err)
	}

	if runErr != nil {
		return cmd, fmt.Errorf("ssh run: %w", runErr)
	}
	return cmd, nil
}

// SSHRunner executes commands on a remote host via SSH.
type SSHRunner interface {
	// Run executes command on user@addr using the provided private key file.
	// Returns stdout, stderr, and exit code. Implementations should use StrictHostKeyChecking=no
	// or a known_hosts strategy appropriate for ephemeral sandboxes.
	Run(ctx context.Context, addr, user, privateKeyPath, command string, timeout time.Duration, env map[string]string) (stdout, stderr string, exitCode int, err error)
}

// DefaultSSHRunner is a simple implementation backed by the system's ssh binary.
type DefaultSSHRunner struct{}

// Run implements SSHRunner.Run using the local ssh client.
// It disables strict host key checking and sets a connect timeout.
// It assumes the VM is reachable on the default SSH port (22).
func (r *DefaultSSHRunner) Run(ctx context.Context, addr, user, privateKeyPath, command string, timeout time.Duration, _ map[string]string) (string, string, int, error) {
	if _, ok := ctx.Deadline(); !ok && timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	args := []string{
		"-i", privateKeyPath,
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null",
		"-o", "ConnectTimeout=15",
		fmt.Sprintf("%s@%s", user, addr),
		"--",
		command,
	}

	var stdout, stderr bytes.Buffer
	cmd := exec.CommandContext(ctx, "ssh", args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	exitCode := 0
	if err != nil {
		// Best-effort extract exit code
		var ee *exec.ExitError
		if errors.As(err, &ee) && ee.ProcessState != nil {
			exitCode = ee.ProcessState.ExitCode()
		} else {
			exitCode = 255
		}
		return stdout.String(), stderr.String(), exitCode, err
	}
	if cmd.ProcessState != nil {
		exitCode = cmd.ProcessState.ExitCode()
	}
	return stdout.String(), stderr.String(), exitCode, nil
}

// Helpers

func snapshotKindFromString(k string) store.SnapshotKind {
	switch strings.ToUpper(k) {
	case "EXTERNAL":
		return store.SnapshotKindExternal
	default:
		return store.SnapshotKindInternal
	}
}

func shortID() string {
	id := uuid.NewString()
	if i := strings.IndexByte(id, '-'); i > 0 {
		return id[:i]
	}
	return id
}

func commandWithEnv(cmd string, env map[string]string) string {
	if len(env) == 0 {
		// Execute in login shell to emulate typical interactive environment
		return fmt.Sprintf("bash -lc %q", cmd)
	}
	var exports []string
	for k, v := range env {
		exports = append(exports, fmt.Sprintf(`export %s=%s`, safeShellIdent(k), shellQuote(v)))
	}
	preamble := strings.Join(exports, "; ") + "; "
	return fmt.Sprintf("bash -lc %q", preamble+cmd)
}

func shellQuote(s string) string {
	// Basic single-quote shell escaping
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}

func safeShellIdent(s string) string {
	// Allow alnum and underscore, replace others with underscore
	var b strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') ||
			(r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') ||
			r == '_' {
			b.WriteRune(r)
		} else {
			b.WriteRune('_')
		}
	}
	out := b.String()
	if out == "" {
		return "VAR"
	}
	return out
}
