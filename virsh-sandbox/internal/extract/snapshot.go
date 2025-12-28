package extract

import (
	"context"
	"fmt"
	"time"

	"virsh-sandbox/internal/libvirt"
	"virsh-sandbox/internal/model"
	"virsh-sandbox/internal/workflow"
)

// SnapshotManager handles snapshot creation and extraction mode detection.
type SnapshotManager struct {
	domainMgr *libvirt.DomainManager
}

// NewSnapshotManager creates a new SnapshotManager.
func NewSnapshotManager(domainMgr *libvirt.DomainManager) *SnapshotManager {
	return &SnapshotManager{
		domainMgr: domainMgr,
	}
}

// ExtractionPlan describes how to extract a VM's filesystem.
type ExtractionPlan struct {
	// VMName is the name of the source VM.
	VMName string

	// Mode is the extraction mode: "snapshot" for running VMs, "offline" for stopped VMs.
	Mode string

	// DiskPath is the path to the disk image to extract from.
	// For snapshot mode, this is the backing file of the snapshot.
	// For offline mode, this is the VM's primary disk.
	DiskPath string

	// SnapshotName is the name of the created snapshot (empty for offline mode).
	SnapshotName string

	// Cleanup is a function to call to clean up the snapshot (nil for offline mode).
	Cleanup workflow.CleanupFunc
}

// DetermineExtractionMode determines whether to use snapshot or offline mode
// based on the VM's current state.
func (m *SnapshotManager) DetermineExtractionMode(ctx context.Context, vmName string) (string, error) {
	state, err := m.domainMgr.GetDomainState(ctx, vmName)
	if err != nil {
		return "", fmt.Errorf("failed to get domain state: %w", err)
	}

	if state.IsRunning() {
		return model.ModeSnapshot, nil
	}
	return model.ModeOffline, nil
}

// PrepareExtraction prepares the extraction plan for a VM.
// For running VMs, it creates a disk-only snapshot.
// For stopped VMs, it returns the disk path directly.
func (m *SnapshotManager) PrepareExtraction(ctx context.Context, vmName string) (*ExtractionPlan, error) {
	// Get domain info
	domainInfo, err := m.domainMgr.LookupDomain(ctx, vmName)
	if err != nil {
		return nil, fmt.Errorf("failed to lookup domain: %w", err)
	}

	// Determine extraction mode
	mode, err := m.DetermineExtractionMode(ctx, vmName)
	if err != nil {
		return nil, err
	}

	plan := &ExtractionPlan{
		VMName: vmName,
		Mode:   mode,
	}

	if mode == model.ModeOffline {
		// For offline mode, use the disk directly
		plan.DiskPath = domainInfo.DiskPath
		return plan, nil
	}

	// For snapshot mode, create a disk-only snapshot
	snapshotName := generateSnapshotName(vmName)

	snapshotInfo, err := m.domainMgr.CreateDiskOnlySnapshot(ctx, vmName, snapshotName)
	if err != nil {
		return nil, workflow.NewWorkflowError(
			workflow.StageCreateSnapshot,
			workflow.ErrSnapshotFailed,
			fmt.Sprintf("unable to create disk-only snapshot: %v", err),
		)
	}

	plan.SnapshotName = snapshotName
	plan.DiskPath = snapshotInfo.BackingFile

	// Create cleanup function that commits the snapshot back
	plan.Cleanup = func() error {
		commitCtx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()
		return m.domainMgr.BlockCommit(commitCtx, vmName, "vda", 5*time.Minute)
	}

	return plan, nil
}

// generateSnapshotName generates a unique snapshot name based on VM name and timestamp.
func generateSnapshotName(vmName string) string {
	return fmt.Sprintf("clone-%s-%d", vmName, time.Now().UnixNano())
}

// CleanupSnapshot removes a snapshot created during extraction.
// This is typically called on successful completion to clean up resources.
func (m *SnapshotManager) CleanupSnapshot(ctx context.Context, vmName string, plan *ExtractionPlan) error {
	if plan == nil || plan.Mode == model.ModeOffline {
		// Nothing to clean up for offline mode
		return nil
	}

	if plan.Cleanup != nil {
		return plan.Cleanup()
	}

	return nil
}
