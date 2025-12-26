//go:build !libvirt
// +build !libvirt

package libvirt

import (
	"context"
	"errors"
	"time"
)

// Sentinel errors for domain operations.
var (
	ErrDomainNotFound    = errors.New("domain not found")
	ErrDomainTransient   = errors.New("transient domains are not supported")
	ErrDomainUnsupported = errors.New("domain configuration not supported")
)

// DomainManager provides libvirt domain operations using libvirt-go bindings.
// This is a stub implementation that returns errors when libvirt is not available.
type DomainManager struct {
	uri string
}

// DomainInfo contains information about a libvirt domain.
type DomainInfo struct {
	Name       string
	UUID       string
	State      DomainState
	Persistent bool
	DiskPath   string
}

// DomainState represents the state of a domain.
type DomainState int

const (
	DomainStateUnknown DomainState = iota
	DomainStateRunning
	DomainStatePaused
	DomainStateShutdown
	DomainStateStopped
	DomainStateCrashed
	DomainStateSuspended
)

// String returns a human-readable domain state.
func (s DomainState) String() string {
	switch s {
	case DomainStateRunning:
		return "running"
	case DomainStatePaused:
		return "paused"
	case DomainStateShutdown:
		return "shutdown"
	case DomainStateStopped:
		return "stopped"
	case DomainStateCrashed:
		return "crashed"
	case DomainStateSuspended:
		return "suspended"
	default:
		return "unknown"
	}
}

// IsRunning returns true if the domain is in a running state.
func (s DomainState) IsRunning() bool {
	return s == DomainStateRunning || s == DomainStatePaused
}

// SnapshotInfo contains information about a created snapshot.
type SnapshotInfo struct {
	Name        string
	BackingFile string
}

// NewDomainManager creates a new DomainManager with the given libvirt URI.
// Note: This stub implementation will return errors for all operations.
func NewDomainManager(uri string) *DomainManager {
	if uri == "" {
		uri = "qemu:///system"
	}
	return &DomainManager{
		uri: uri,
	}
}

// Connect is a stub that returns an error when libvirt is not available.
func (m *DomainManager) Connect() error {
	return ErrLibvirtNotAvailable
}

// Close is a stub that does nothing when libvirt is not available.
func (m *DomainManager) Close() error {
	return nil
}

// LookupDomain is a stub that returns an error when libvirt is not available.
func (m *DomainManager) LookupDomain(ctx context.Context, name string) (*DomainInfo, error) {
	return nil, ErrLibvirtNotAvailable
}

// GetDomainState is a stub that returns an error when libvirt is not available.
func (m *DomainManager) GetDomainState(ctx context.Context, name string) (DomainState, error) {
	return DomainStateUnknown, ErrLibvirtNotAvailable
}

// CreateDiskOnlySnapshot is a stub that returns an error when libvirt is not available.
func (m *DomainManager) CreateDiskOnlySnapshot(ctx context.Context, domainName, snapshotName string) (*SnapshotInfo, error) {
	return nil, ErrLibvirtNotAvailable
}

// BlockCommit is a stub that returns an error when libvirt is not available.
func (m *DomainManager) BlockCommit(ctx context.Context, domainName, diskTarget string, timeout time.Duration) error {
	return ErrLibvirtNotAvailable
}

// ListDomains is a stub that returns an error when libvirt is not available.
func (m *DomainManager) ListDomains(ctx context.Context) ([]*DomainInfo, error) {
	return nil, ErrLibvirtNotAvailable
}

// GetDiskPath is a stub that returns an error when libvirt is not available.
func (m *DomainManager) GetDiskPath(ctx context.Context, domainName string) (string, error) {
	return "", ErrLibvirtNotAvailable
}
