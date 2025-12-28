package extract

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"virsh-sandbox/internal/workflow"
)

// Archiver handles creation of root filesystem archives.
type Archiver struct {
	// tarPath is the path to the tar binary.
	tarPath string
}

// ArchiverConfig configures the archiver.
type ArchiverConfig struct {
	// TarPath is the path to the tar binary.
	// If empty, "tar" is looked up in PATH.
	TarPath string
}

// NewArchiver creates a new Archiver with the given configuration.
func NewArchiver(cfg ArchiverConfig) *Archiver {
	tarPath := cfg.TarPath
	if tarPath == "" {
		tarPath = "tar"
	}
	return &Archiver{
		tarPath: tarPath,
	}
}

// ArchiveResult contains the result of creating an archive.
type ArchiveResult struct {
	// ArchivePath is the path to the created tar archive.
	ArchivePath string

	// Size is the size of the archive in bytes.
	Size int64

	// Cleanup is a function to remove the archive.
	Cleanup workflow.CleanupFunc
}

// CreateRootFSArchive creates a tar archive of the sanitized root filesystem.
// The archive preserves numeric ownership and extended attributes.
func (a *Archiver) CreateRootFSArchive(ctx context.Context, sourcePath string, workDir string) (*ArchiveResult, error) {
	// Generate archive filename with timestamp
	timestamp := time.Now().UTC().Format("20060102T150405Z")
	archiveName := fmt.Sprintf("rootfs-%s.tar", timestamp)
	archivePath := filepath.Join(workDir, archiveName)

	// Create the archive
	if err := a.createTarArchive(ctx, sourcePath, archivePath); err != nil {
		return nil, workflow.NewWorkflowError(
			workflow.StageCreateArchive,
			workflow.ErrArchiveFailed,
			fmt.Sprintf("failed to create archive: %v", err),
		)
	}

	// Get archive size
	info, err := os.Stat(archivePath)
	if err != nil {
		_ = os.Remove(archivePath)
		return nil, workflow.NewWorkflowError(
			workflow.StageCreateArchive,
			workflow.ErrArchiveFailed,
			fmt.Sprintf("failed to stat archive: %v", err),
		)
	}

	return &ArchiveResult{
		ArchivePath: archivePath,
		Size:        info.Size(),
		Cleanup: func() error {
			return os.Remove(archivePath)
		},
	}, nil
}

// createTarArchive creates a tar archive of the source directory.
func (a *Archiver) createTarArchive(ctx context.Context, sourcePath, archivePath string) error {
	// Build tar command with options:
	// -c: create archive
	// -f: output file
	// --numeric-owner: preserve numeric UID/GID (important for container images)
	// --xattrs: preserve extended attributes
	// --xattrs-include=*: include all xattrs
	// --acls: preserve ACLs (if supported)
	// --selinux: preserve SELinux contexts (if applicable)
	// -C: change to directory before archiving
	// .: archive current directory contents

	args := []string{
		"-cf", archivePath,
		"--numeric-owner",
	}

	// Check if tar supports xattrs
	if a.supportsXattrs(ctx) {
		args = append(args, "--xattrs", "--xattrs-include=*")
	}

	// Check if tar supports ACLs
	if a.supportsACLs(ctx) {
		args = append(args, "--acls")
	}

	// Check if tar supports SELinux
	if a.supportsSELinux(ctx) {
		args = append(args, "--selinux")
	}

	// Add source directory
	args = append(args, "-C", sourcePath, ".")

	cmd := exec.CommandContext(ctx, a.tarPath, args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("tar failed: %w: %s", err, stderr.String())
	}

	return nil
}

// supportsXattrs checks if tar supports --xattrs option.
func (a *Archiver) supportsXattrs(ctx context.Context) bool {
	cmd := exec.CommandContext(ctx, a.tarPath, "--xattrs", "--help")
	return cmd.Run() == nil
}

// supportsACLs checks if tar supports --acls option.
func (a *Archiver) supportsACLs(ctx context.Context) bool {
	cmd := exec.CommandContext(ctx, a.tarPath, "--acls", "--help")
	return cmd.Run() == nil
}

// supportsSELinux checks if tar supports --selinux option.
func (a *Archiver) supportsSELinux(ctx context.Context) bool {
	cmd := exec.CommandContext(ctx, a.tarPath, "--selinux", "--help")
	return cmd.Run() == nil
}

// ExtractArchive extracts a tar archive to the specified destination.
// This is useful for testing or container image import operations.
func (a *Archiver) ExtractArchive(ctx context.Context, archivePath, destPath string) error {
	// Ensure destination exists
	if err := os.MkdirAll(destPath, 0o755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	args := []string{
		"-xf", archivePath,
		"--numeric-owner",
		"-C", destPath,
	}

	// Add xattrs support if available
	if a.supportsXattrs(ctx) {
		args = append(args[:2], append([]string{"--xattrs", "--xattrs-include=*"}, args[2:]...)...)
	}

	cmd := exec.CommandContext(ctx, a.tarPath, args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("tar extraction failed: %w: %s", err, stderr.String())
	}

	return nil
}

// GetArchiveSize returns the size of an archive file.
func (a *Archiver) GetArchiveSize(archivePath string) (int64, error) {
	info, err := os.Stat(archivePath)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

// ListArchiveContents lists the contents of an archive.
// This is useful for verification and debugging.
func (a *Archiver) ListArchiveContents(ctx context.Context, archivePath string) ([]string, error) {
	args := []string{"-tf", archivePath}

	cmd := exec.CommandContext(ctx, a.tarPath, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("tar list failed: %w: %s", err, stderr.String())
	}

	// Parse output into list of files
	output := stdout.String()
	if output == "" {
		return []string{}, nil
	}

	lines := bytes.Split(stdout.Bytes(), []byte("\n"))
	files := make([]string, 0, len(lines))
	for _, line := range lines {
		if len(line) > 0 {
			files = append(files, string(line))
		}
	}

	return files, nil
}

// VerifyArchive performs basic verification of an archive.
func (a *Archiver) VerifyArchive(ctx context.Context, archivePath string) error {
	// Check file exists
	info, err := os.Stat(archivePath)
	if err != nil {
		return fmt.Errorf("archive not found: %w", err)
	}

	// Check file is not empty
	if info.Size() == 0 {
		return fmt.Errorf("archive is empty")
	}

	// Try to list contents to verify integrity
	args := []string{"-tf", archivePath}
	cmd := exec.CommandContext(ctx, a.tarPath, args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("archive verification failed: %w: %s", err, stderr.String())
	}

	return nil
}
