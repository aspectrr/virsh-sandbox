package extract

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"virsh-sandbox/internal/workflow"
)

// Sanitizer handles filesystem sanitization for container usage.
type Sanitizer struct {
	// verbose enables detailed logging of sanitization steps.
	verbose bool
}

// SanitizerConfig configures the sanitizer.
type SanitizerConfig struct {
	// Verbose enables detailed logging.
	Verbose bool
}

// NewSanitizer creates a new Sanitizer with the given configuration.
func NewSanitizer(cfg SanitizerConfig) *Sanitizer {
	return &Sanitizer{
		verbose: cfg.Verbose,
	}
}

// SanitizeResult contains the result of filesystem sanitization.
type SanitizeResult struct {
	// SanitizedPath is the path to the sanitized filesystem copy.
	SanitizedPath string

	// RemovedPaths lists paths that were removed or neutralized.
	RemovedPaths []string

	// ModifiedPaths lists paths that were modified.
	ModifiedPaths []string

	// Cleanup is a function to remove the sanitized copy.
	Cleanup workflow.CleanupFunc
}

// SanitizeFilesystem creates a sanitized copy of the mounted filesystem
// suitable for container usage. It removes or neutralizes:
// - /boot directory
// - kernel modules (/lib/modules)
// - device nodes under /dev
// - fstab contents
// - swap references
// - systemd services that block container execution
func (s *Sanitizer) SanitizeFilesystem(ctx context.Context, sourcePath string, workDir string) (*SanitizeResult, error) {
	// Create a working directory for the sanitized copy
	sanitizedPath := filepath.Join(workDir, "sanitized")
	if err := os.MkdirAll(sanitizedPath, 0o755); err != nil {
		return nil, workflow.NewWorkflowError(
			workflow.StageSanitizeFS,
			workflow.ErrSanitizeFailed,
			fmt.Sprintf("failed to create sanitized directory: %v", err),
		)
	}

	result := &SanitizeResult{
		SanitizedPath: sanitizedPath,
		RemovedPaths:  make([]string, 0),
		ModifiedPaths: make([]string, 0),
	}

	// Copy the filesystem using rsync for efficiency
	// We exclude certain paths during copy rather than copying then deleting
	if err := s.copyFilesystem(ctx, sourcePath, sanitizedPath); err != nil {
		_ = os.RemoveAll(sanitizedPath)
		return nil, workflow.NewWorkflowError(
			workflow.StageSanitizeFS,
			workflow.ErrSanitizeFailed,
			fmt.Sprintf("failed to copy filesystem: %v", err),
		)
	}

	// Apply sanitization steps
	sanitizers := []struct {
		name string
		fn   func(ctx context.Context, rootPath string, result *SanitizeResult) error
	}{
		{"remove boot directory", s.removeBoot},
		{"remove kernel modules", s.removeKernelModules},
		{"clear device nodes", s.clearDeviceNodes},
		{"sanitize fstab", s.sanitizeFstab},
		{"remove swap references", s.removeSwapReferences},
		{"disable blocking systemd services", s.disableBlockingServices},
		{"set container environment marker", s.setContainerMarker},
	}

	for _, sanitizer := range sanitizers {
		if err := sanitizer.fn(ctx, sanitizedPath, result); err != nil {
			_ = os.RemoveAll(sanitizedPath)
			return nil, workflow.NewWorkflowError(
				workflow.StageSanitizeFS,
				workflow.ErrSanitizeFailed,
				fmt.Sprintf("%s failed: %v", sanitizer.name, err),
			)
		}
	}

	result.Cleanup = func() error {
		return os.RemoveAll(sanitizedPath)
	}

	return result, nil
}

// copyFilesystem copies the source filesystem to the destination,
// excluding paths that will be removed anyway.
func (s *Sanitizer) copyFilesystem(ctx context.Context, src, dst string) error {
	// Use rsync for efficient copying with exclusions
	// Exclude paths we're going to remove anyway to save time and space
	excludes := []string{
		"--exclude=/boot/*",
		"--exclude=/lib/modules/*",
		"--exclude=/dev/*",
		"--exclude=/proc/*",
		"--exclude=/sys/*",
		"--exclude=/run/*",
		"--exclude=/tmp/*",
		"--exclude=/var/tmp/*",
		"--exclude=/var/cache/*",
		"--exclude=/var/log/*",
		"--exclude=*.swap",
		"--exclude=/swapfile",
	}

	args := []string{
		"-a",           // archive mode (preserves permissions, ownership, etc.)
		"--hard-links", // preserve hard links
		"--acls",       // preserve ACLs
		"--xattrs",     // preserve extended attributes
		"--sparse",     // handle sparse files efficiently
	}
	args = append(args, excludes...)
	args = append(args, src+"/", dst+"/")

	cmd := exec.CommandContext(ctx, "rsync", args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		// Try with cp if rsync is not available
		return s.copyFilesystemFallback(ctx, src, dst)
	}

	return nil
}

// copyFilesystemFallback uses cp when rsync is not available.
func (s *Sanitizer) copyFilesystemFallback(ctx context.Context, src, dst string) error {
	args := []string{
		"-a",             // archive mode
		"--reflink=auto", // use copy-on-write if available
		src + "/.",
		dst + "/",
	}

	cmd := exec.CommandContext(ctx, "cp", args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("cp failed: %w: %s", err, stderr.String())
	}

	return nil
}

// removeBoot removes the /boot directory contents.
func (s *Sanitizer) removeBoot(ctx context.Context, rootPath string, result *SanitizeResult) error {
	bootPath := filepath.Join(rootPath, "boot")

	// Remove contents but keep the directory
	if err := s.clearDirectory(bootPath); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	result.RemovedPaths = append(result.RemovedPaths, "/boot/*")
	return nil
}

// removeKernelModules removes kernel modules from /lib/modules.
func (s *Sanitizer) removeKernelModules(ctx context.Context, rootPath string, result *SanitizeResult) error {
	modulesPath := filepath.Join(rootPath, "lib", "modules")

	if err := s.clearDirectory(modulesPath); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	result.RemovedPaths = append(result.RemovedPaths, "/lib/modules/*")
	return nil
}

// clearDeviceNodes removes all device nodes under /dev.
func (s *Sanitizer) clearDeviceNodes(ctx context.Context, rootPath string, result *SanitizeResult) error {
	devPath := filepath.Join(rootPath, "dev")

	// Remove contents but keep the directory
	if err := s.clearDirectory(devPath); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	// Create minimal /dev entries that containers expect
	// The container runtime will populate /dev properly
	if err := os.MkdirAll(devPath, 0o755); err != nil {
		return err
	}

	// Create /dev/null, /dev/zero, /dev/random placeholders
	// These are symlinks that the container runtime will handle
	devEntries := []string{"null", "zero", "random", "urandom", "tty", "console"}
	for _, entry := range devEntries {
		placeholder := filepath.Join(devPath, entry)
		// Create empty placeholder files
		f, err := os.Create(placeholder)
		if err != nil {
			continue // Non-fatal, container runtime will create these
		}
		f.Close()
	}

	result.RemovedPaths = append(result.RemovedPaths, "/dev/*")
	return nil
}

// sanitizeFstab clears or comments out /etc/fstab entries.
func (s *Sanitizer) sanitizeFstab(ctx context.Context, rootPath string, result *SanitizeResult) error {
	fstabPath := filepath.Join(rootPath, "etc", "fstab")

	content, err := os.ReadFile(fstabPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	// Comment out all mount entries, keeping only comments
	lines := strings.Split(string(content), "\n")
	var newLines []string
	newLines = append(newLines, "# fstab sanitized for container usage")
	newLines = append(newLines, "# Original entries commented out:")
	newLines = append(newLines, "")

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			newLines = append(newLines, line)
		} else {
			newLines = append(newLines, "# "+line)
		}
	}

	if err := os.WriteFile(fstabPath, []byte(strings.Join(newLines, "\n")), 0o644); err != nil {
		return err
	}

	result.ModifiedPaths = append(result.ModifiedPaths, "/etc/fstab")
	return nil
}

// removeSwapReferences removes or disables swap configuration.
func (s *Sanitizer) removeSwapReferences(ctx context.Context, rootPath string, result *SanitizeResult) error {
	// Remove swapfile if it exists
	swapfile := filepath.Join(rootPath, "swapfile")
	if _, err := os.Stat(swapfile); err == nil {
		if err := os.Remove(swapfile); err != nil {
			return err
		}
		result.RemovedPaths = append(result.RemovedPaths, "/swapfile")
	}

	// Remove any .swap files in root
	entries, err := os.ReadDir(rootPath)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".swap") {
			swapPath := filepath.Join(rootPath, entry.Name())
			if err := os.Remove(swapPath); err != nil {
				continue // Non-fatal
			}
			result.RemovedPaths = append(result.RemovedPaths, "/"+entry.Name())
		}
	}

	// Disable swap-related systemd units
	swapUnits := []string{
		"swap.target",
		"dev-*.swap",
	}

	// systemdPath := filepath.Join(rootPath, "etc", "systemd", "system")
	for _, unit := range swapUnits {
		// unitPath := filepath.Join(systemdPath, unit)
		// Create a masked symlink to /dev/null
		if err := s.maskSystemdUnit(rootPath, unit); err != nil {
			continue // Non-fatal
		}
	}

	return nil
}

// disableBlockingServices disables systemd services that block container execution.
func (s *Sanitizer) disableBlockingServices(ctx context.Context, rootPath string, result *SanitizeResult) error {
	// Services that commonly block container startup or are inappropriate
	blockingServices := []string{
		// Hardware/kernel related
		"systemd-modules-load.service",
		"systemd-sysctl.service",
		"systemd-udevd.service",
		"systemd-udev-trigger.service",
		"systemd-udev-settle.service",
		"kmod-static-nodes.service",
		"systemd-tmpfiles-setup-dev.service",

		// Filesystem related
		"systemd-remount-fs.service",
		"systemd-fsck@.service",
		"systemd-fsck-root.service",
		"local-fs.target",
		"local-fs-pre.target",

		// Network hardware related
		"NetworkManager-wait-online.service",
		"systemd-networkd-wait-online.service",

		// Other blocking services
		"plymouth-start.service",
		"plymouth-quit.service",
		"plymouth-quit-wait.service",
		"systemd-machine-id-commit.service",
		"systemd-firstboot.service",
		"systemd-random-seed.service",

		// Console/TTY related
		"getty@.service",
		"serial-getty@.service",
		"console-getty.service",
		"container-getty@.service",
		"systemd-ask-password-wall.service",
		"systemd-ask-password-console.service",
	}

	for _, service := range blockingServices {
		if err := s.maskSystemdUnit(rootPath, service); err != nil {
			continue // Non-fatal, service may not exist
		}
		result.ModifiedPaths = append(result.ModifiedPaths, "/etc/systemd/system/"+service)
	}

	return nil
}

// maskSystemdUnit masks a systemd unit by creating a symlink to /dev/null.
func (s *Sanitizer) maskSystemdUnit(rootPath, unitName string) error {
	systemdPath := filepath.Join(rootPath, "etc", "systemd", "system")
	if err := os.MkdirAll(systemdPath, 0o755); err != nil {
		return err
	}

	unitPath := filepath.Join(systemdPath, unitName)

	// Remove existing unit/symlink if present
	_ = os.Remove(unitPath)

	// Create symlink to /dev/null to mask the unit
	return os.Symlink("/dev/null", unitPath)
}

// setContainerMarker creates markers indicating container environment.
func (s *Sanitizer) setContainerMarker(ctx context.Context, rootPath string, result *SanitizeResult) error {
	// Create /run/container marker directory
	runPath := filepath.Join(rootPath, "run")
	if err := os.MkdirAll(runPath, 0o755); err != nil {
		return err
	}

	// Create /.dockerenv equivalent for container detection
	dockerenvPath := filepath.Join(rootPath, ".dockerenv")
	f, err := os.Create(dockerenvPath)
	if err != nil {
		return err
	}
	f.Close()
	result.ModifiedPaths = append(result.ModifiedPaths, "/.dockerenv")

	// Create /run/.containerenv for Podman detection
	containerenvPath := filepath.Join(runPath, ".containerenv")
	containerenvContent := `engine="podman"
name="vmclone"
`
	if err := os.WriteFile(containerenvPath, []byte(containerenvContent), 0o644); err != nil {
		return err
	}
	result.ModifiedPaths = append(result.ModifiedPaths, "/run/.containerenv")

	return nil
}

// clearDirectory removes all contents of a directory but keeps the directory itself.
func (s *Sanitizer) clearDirectory(dirPath string) error {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		entryPath := filepath.Join(dirPath, entry.Name())
		if err := os.RemoveAll(entryPath); err != nil {
			// Try to continue with other entries
			continue
		}
	}

	return nil
}
