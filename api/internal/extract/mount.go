package extract

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"virsh-sandbox/internal/workflow"
)

// MountManager handles qemu-nbd attachment and filesystem mounting.
type MountManager struct {
	// qemuNbdPath is the path to the qemu-nbd binary.
	qemuNbdPath string

	// nbdDeviceMu protects NBD device allocation.
	nbdDeviceMu sync.Mutex

	// usedNbdDevices tracks which NBD devices are currently in use.
	usedNbdDevices map[string]bool
}

// MountConfig configures the mount manager.
type MountConfig struct {
	// QemuNbdPath is the path to the qemu-nbd binary.
	// If empty, "qemu-nbd" is looked up in PATH.
	QemuNbdPath string
}

// MountResult contains the result of mounting a disk image.
type MountResult struct {
	// NBDDevice is the /dev/nbdX device the image is attached to.
	NBDDevice string

	// Partition is the partition device (e.g., /dev/nbd0p1).
	Partition string

	// MountPoint is the path where the filesystem is mounted.
	MountPoint string

	// Cleanup is a function that unmounts and disconnects everything.
	Cleanup workflow.CleanupFunc
}

// NewMountManager creates a new MountManager with the given configuration.
func NewMountManager(cfg MountConfig) *MountManager {
	qemuNbdPath := cfg.QemuNbdPath
	if qemuNbdPath == "" {
		qemuNbdPath = "qemu-nbd"
	}
	return &MountManager{
		qemuNbdPath:    qemuNbdPath,
		usedNbdDevices: make(map[string]bool),
	}
}

// MountDisk attaches a disk image via qemu-nbd and mounts the root filesystem.
// The returned MountResult contains a cleanup function that must be called
// to unmount and disconnect the NBD device.
func (m *MountManager) MountDisk(ctx context.Context, diskPath string, workDir string) (*MountResult, error) {
	// Verify disk exists
	if _, err := os.Stat(diskPath); err != nil {
		return nil, workflow.NewWorkflowError(
			workflow.StageMountDisk,
			workflow.ErrMountFailed,
			fmt.Sprintf("disk image not found: %s", diskPath),
		)
	}

	// Find an available NBD device
	nbdDevice, err := m.findAvailableNBDDevice()
	if err != nil {
		return nil, workflow.NewWorkflowError(
			workflow.StageMountDisk,
			workflow.ErrNBDAttachFailed,
			fmt.Sprintf("no available NBD device: %v", err),
		)
	}

	result := &MountResult{
		NBDDevice: nbdDevice,
	}

	// Track cleanup steps for rollback
	cleanups := workflow.NewCleanupStack()

	// Attach the disk to NBD device
	if err := m.attachNBD(ctx, diskPath, nbdDevice); err != nil {
		m.releaseNBDDevice(nbdDevice)
		return nil, workflow.NewWorkflowError(
			workflow.StageMountDisk,
			workflow.ErrNBDAttachFailed,
			fmt.Sprintf("failed to attach %s to %s: %v", diskPath, nbdDevice, err),
		)
	}
	cleanups.Push(func() error {
		return m.detachNBD(context.Background(), nbdDevice)
	})

	// Run partprobe to detect partitions
	if err := m.runPartprobe(ctx, nbdDevice); err != nil {
		_ = cleanups.ExecuteAll()
		m.releaseNBDDevice(nbdDevice)
		return nil, workflow.NewWorkflowError(
			workflow.StageMountDisk,
			workflow.ErrMountFailed,
			fmt.Sprintf("partprobe failed: %v", err),
		)
	}

	// Find the root partition
	partition, err := m.findRootPartition(ctx, nbdDevice)
	if err != nil {
		_ = cleanups.ExecuteAll()
		m.releaseNBDDevice(nbdDevice)
		return nil, workflow.NewWorkflowError(
			workflow.StageMountDisk,
			workflow.ErrMountFailed,
			fmt.Sprintf("failed to find root partition: %v", err),
		)
	}
	result.Partition = partition

	// Create mount point
	mountPoint := filepath.Join(workDir, "rootfs")
	if err := os.MkdirAll(mountPoint, 0755); err != nil {
		_ = cleanups.ExecuteAll()
		m.releaseNBDDevice(nbdDevice)
		return nil, workflow.NewWorkflowError(
			workflow.StageMountDisk,
			workflow.ErrMountFailed,
			fmt.Sprintf("failed to create mount point: %v", err),
		)
	}
	cleanups.Push(func() error {
		return os.RemoveAll(mountPoint)
	})

	// Mount the partition read-only
	if err := m.mountPartition(ctx, partition, mountPoint); err != nil {
		_ = cleanups.ExecuteAll()
		m.releaseNBDDevice(nbdDevice)
		return nil, workflow.NewWorkflowError(
			workflow.StageMountDisk,
			workflow.ErrMountFailed,
			fmt.Sprintf("failed to mount %s at %s: %v", partition, mountPoint, err),
		)
	}
	result.MountPoint = mountPoint

	// Build the final cleanup function that does everything in reverse order
	result.Cleanup = func() error {
		var errs []error

		// Unmount filesystem
		if err := m.unmount(context.Background(), mountPoint); err != nil {
			errs = append(errs, fmt.Errorf("unmount %s: %w", mountPoint, err))
		}

		// Remove mount point directory
		if err := os.RemoveAll(mountPoint); err != nil {
			errs = append(errs, fmt.Errorf("remove mount point: %w", err))
		}

		// Detach NBD device
		if err := m.detachNBD(context.Background(), nbdDevice); err != nil {
			errs = append(errs, fmt.Errorf("detach NBD: %w", err))
		}

		// Release the NBD device for reuse
		m.releaseNBDDevice(nbdDevice)

		if len(errs) > 0 {
			return fmt.Errorf("cleanup errors: %v", errs)
		}
		return nil
	}

	return result, nil
}

// findAvailableNBDDevice finds an available /dev/nbdX device.
func (m *MountManager) findAvailableNBDDevice() (string, error) {
	m.nbdDeviceMu.Lock()
	defer m.nbdDeviceMu.Unlock()

	// Check for nbd module
	if _, err := os.Stat("/sys/module/nbd"); os.IsNotExist(err) {
		return "", fmt.Errorf("nbd kernel module not loaded; run 'modprobe nbd max_part=16'")
	}

	// Try to find an available NBD device (typically nbd0 through nbd15)
	for i := 0; i < 16; i++ {
		device := fmt.Sprintf("/dev/nbd%d", i)

		// Skip if we're already using it
		if m.usedNbdDevices[device] {
			continue
		}

		// Check if device exists
		if _, err := os.Stat(device); os.IsNotExist(err) {
			continue
		}

		// Check if device is in use by examining its size
		sizePath := fmt.Sprintf("/sys/block/nbd%d/size", i)
		data, err := os.ReadFile(sizePath)
		if err != nil {
			continue
		}

		size, err := strconv.ParseInt(strings.TrimSpace(string(data)), 10, 64)
		if err != nil {
			continue
		}

		// Size of 0 means the device is not in use
		if size == 0 {
			m.usedNbdDevices[device] = true
			return device, nil
		}
	}

	return "", fmt.Errorf("all NBD devices are in use")
}

// releaseNBDDevice marks an NBD device as available for reuse.
func (m *MountManager) releaseNBDDevice(device string) {
	m.nbdDeviceMu.Lock()
	defer m.nbdDeviceMu.Unlock()
	delete(m.usedNbdDevices, device)
}

// attachNBD attaches a disk image to an NBD device using qemu-nbd.
func (m *MountManager) attachNBD(ctx context.Context, diskPath, nbdDevice string) error {
	// Connect the image to the NBD device
	// --read-only for safety, --connect to specify the device
	args := []string{
		"--read-only",
		"--connect", nbdDevice,
		"--format", "qcow2",
		diskPath,
	}

	cmd := exec.CommandContext(ctx, m.qemuNbdPath, args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("qemu-nbd failed: %w: %s", err, stderr.String())
	}

	// Wait a bit for the device to be ready
	time.Sleep(500 * time.Millisecond)

	return nil
}

// detachNBD disconnects an NBD device.
func (m *MountManager) detachNBD(ctx context.Context, nbdDevice string) error {
	args := []string{"--disconnect", nbdDevice}

	cmd := exec.CommandContext(ctx, m.qemuNbdPath, args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("qemu-nbd disconnect failed: %w: %s", err, stderr.String())
	}

	return nil
}

// runPartprobe runs partprobe to detect partitions on the NBD device.
func (m *MountManager) runPartprobe(ctx context.Context, nbdDevice string) error {
	cmd := exec.CommandContext(ctx, "partprobe", nbdDevice)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("partprobe failed: %w: %s", err, stderr.String())
	}

	// Wait for partition devices to appear
	time.Sleep(500 * time.Millisecond)

	return nil
}

// findRootPartition attempts to find the root partition on the NBD device.
// It looks for common partition layouts and returns the likely root partition.
func (m *MountManager) findRootPartition(ctx context.Context, nbdDevice string) (string, error) {
	// Get the device name without /dev/ prefix
	devName := filepath.Base(nbdDevice)

	// Check for partitions in /sys/block/<device>/
	sysPath := fmt.Sprintf("/sys/block/%s", devName)

	entries, err := os.ReadDir(sysPath)
	if err != nil {
		return "", fmt.Errorf("failed to read %s: %w", sysPath, err)
	}

	var partitions []string
	for _, entry := range entries {
		name := entry.Name()
		// Partition entries start with the device name
		if strings.HasPrefix(name, devName+"p") {
			partitions = append(partitions, "/dev/"+name)
		}
	}

	if len(partitions) == 0 {
		// No partitions found, might be a whole-disk filesystem
		// Try to mount the device directly
		return nbdDevice, nil
	}

	// Sort partitions and try to find the root partition
	// Typically:
	// - p1 is often /boot or EFI on modern systems
	// - p2 or p3 is often root
	// We'll try to identify by checking for common root filesystem indicators

	for _, partition := range partitions {
		// Use blkid to check filesystem type
		cmd := exec.CommandContext(ctx, "blkid", "-o", "value", "-s", "TYPE", partition)
		output, err := cmd.Output()
		if err != nil {
			continue
		}

		fsType := strings.TrimSpace(string(output))
		// Look for ext4, xfs, btrfs which are common root filesystems
		if fsType == "ext4" || fsType == "xfs" || fsType == "btrfs" || fsType == "ext3" {
			// This is likely the root partition
			// We could do more checks (mount and look for /etc, /bin, etc.)
			// but for now we'll use the first Linux filesystem we find
			// that isn't obviously a boot partition
			if !m.isBootPartition(ctx, partition) {
				return partition, nil
			}
		}
	}

	// If we couldn't find a definitive root, try the largest partition
	if len(partitions) > 0 {
		largest := partitions[0]
		var largestSize int64

		for _, partition := range partitions {
			partName := filepath.Base(partition)
			sizePath := fmt.Sprintf("/sys/block/%s/%s/size", devName, partName)
			data, err := os.ReadFile(sizePath)
			if err != nil {
				continue
			}
			size, err := strconv.ParseInt(strings.TrimSpace(string(data)), 10, 64)
			if err != nil {
				continue
			}
			if size > largestSize {
				largestSize = size
				largest = partition
			}
		}
		return largest, nil
	}

	return "", fmt.Errorf("no suitable partition found")
}

// isBootPartition checks if a partition appears to be a boot partition.
func (m *MountManager) isBootPartition(ctx context.Context, partition string) bool {
	// Check partition label or flags
	cmd := exec.CommandContext(ctx, "blkid", "-o", "value", "-s", "LABEL", partition)
	output, err := cmd.Output()
	if err == nil {
		label := strings.ToLower(strings.TrimSpace(string(output)))
		if strings.Contains(label, "boot") || strings.Contains(label, "efi") {
			return true
		}
	}

	// Also check PARTLABEL for GPT partitions
	cmd = exec.CommandContext(ctx, "blkid", "-o", "value", "-s", "PARTLABEL", partition)
	output, err = cmd.Output()
	if err == nil {
		label := strings.ToLower(strings.TrimSpace(string(output)))
		if strings.Contains(label, "boot") || strings.Contains(label, "efi") {
			return true
		}
	}

	return false
}

// mountPartition mounts a partition read-only at the specified mount point.
func (m *MountManager) mountPartition(ctx context.Context, partition, mountPoint string) error {
	// Mount read-only with common options
	args := []string{
		"-o", "ro,noatime,noexec",
		partition,
		mountPoint,
	}

	cmd := exec.CommandContext(ctx, "mount", args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("mount failed: %w: %s", err, stderr.String())
	}

	return nil
}

// unmount unmounts a filesystem.
func (m *MountManager) unmount(ctx context.Context, mountPoint string) error {
	// First try a regular unmount
	cmd := exec.CommandContext(ctx, "umount", mountPoint)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		// If regular unmount fails, try lazy unmount
		cmd = exec.CommandContext(ctx, "umount", "-l", mountPoint)
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("unmount failed: %w: %s", err, stderr.String())
		}
	}

	return nil
}
