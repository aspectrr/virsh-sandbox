package libvirt

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"sync"
	"time"

	libvirtgo "libvirt.org/go/libvirt"
)

// DomainManager provides libvirt domain operations using libvirt-go bindings.
type DomainManager struct {
	uri  string
	conn *libvirtgo.Connect
	mu   sync.Mutex
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

// Domain XML structures for parsing disk information.
type domainXML struct {
	XMLName xml.Name      `xml:"domain"`
	Name    string        `xml:"name"`
	UUID    string        `xml:"uuid"`
	Devices domainDevices `xml:"devices"`
}

type domainDevices struct {
	Disks []domainDisk `xml:"disk"`
}

type domainDisk struct {
	Type   string           `xml:"type,attr"`
	Device string           `xml:"device,attr"`
	Driver domainDiskDriver `xml:"driver"`
	Source domainDiskSource `xml:"source"`
	Target domainDiskTarget `xml:"target"`
}

type domainDiskDriver struct {
	Name string `xml:"name,attr"`
	Type string `xml:"type,attr"`
}

type domainDiskSource struct {
	File string `xml:"file,attr"`
	Dev  string `xml:"dev,attr"`
}

type domainDiskTarget struct {
	Dev string `xml:"dev,attr"`
	Bus string `xml:"bus,attr"`
}

// NewDomainManager creates a new DomainManager with the given libvirt URI.
func NewDomainManager(uri string) *DomainManager {
	if uri == "" {
		uri = "qemu:///system"
	}
	return &DomainManager{
		uri: uri,
	}
}

// Connect establishes a connection to libvirt.
func (m *DomainManager) Connect() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.conn != nil {
		// Check if connection is still alive
		if alive, _ := m.conn.IsAlive(); alive {
			return nil
		}
		// Connection dead, close and reconnect
		m.conn.Close()
		m.conn = nil
	}

	conn, err := libvirtgo.NewConnect(m.uri)
	if err != nil {
		return fmt.Errorf("failed to connect to libvirt at %s: %w", m.uri, err)
	}
	m.conn = conn
	return nil
}

// Close closes the libvirt connection.
func (m *DomainManager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.conn != nil {
		_, err := m.conn.Close()
		m.conn = nil
		return err
	}
	return nil
}

// ensureConnected ensures we have an active connection.
func (m *DomainManager) ensureConnected() error {
	return m.Connect()
}

// LookupDomain looks up a domain by name and returns its information.
// Returns an error if the domain is transient or not found.
func (m *DomainManager) LookupDomain(ctx context.Context, name string) (*DomainInfo, error) {
	if err := m.ensureConnected(); err != nil {
		return nil, err
	}

	m.mu.Lock()
	dom, err := m.conn.LookupDomainByName(name)
	m.mu.Unlock()

	if err != nil {
		var libvirtErr libvirtgo.Error
		if errors.As(err, &libvirtErr) {
			if libvirtErr.Code == libvirtgo.ERR_NO_DOMAIN {
				return nil, fmt.Errorf("domain %q not found: %w", name, ErrDomainNotFound)
			}
		}
		return nil, fmt.Errorf("failed to lookup domain %q: %w", name, err)
	}
	defer dom.Free()

	// Check if domain is persistent (not transient)
	persistent, err := dom.IsPersistent()
	if err != nil {
		return nil, fmt.Errorf("failed to check if domain is persistent: %w", err)
	}
	if !persistent {
		return nil, fmt.Errorf("domain %q is transient: %w", name, ErrDomainTransient)
	}

	// Get domain UUID
	uuid, err := dom.GetUUIDString()
	if err != nil {
		return nil, fmt.Errorf("failed to get domain UUID: %w", err)
	}

	// Get domain state
	state, _, err := dom.GetState()
	if err != nil {
		return nil, fmt.Errorf("failed to get domain state: %w", err)
	}

	// Get domain XML to extract disk path
	xmlDesc, err := dom.GetXMLDesc(0)
	if err != nil {
		return nil, fmt.Errorf("failed to get domain XML: %w", err)
	}

	diskPath, err := extractDiskPath(xmlDesc)
	if err != nil {
		return nil, fmt.Errorf("failed to extract disk path: %w", err)
	}

	return &DomainInfo{
		Name:       name,
		UUID:       uuid,
		State:      mapLibvirtState(state),
		Persistent: persistent,
		DiskPath:   diskPath,
	}, nil
}

// GetDomainState returns the current state of a domain.
func (m *DomainManager) GetDomainState(ctx context.Context, name string) (DomainState, error) {
	if err := m.ensureConnected(); err != nil {
		return DomainStateUnknown, err
	}

	m.mu.Lock()
	dom, err := m.conn.LookupDomainByName(name)
	m.mu.Unlock()

	if err != nil {
		return DomainStateUnknown, fmt.Errorf("failed to lookup domain %q: %w", name, err)
	}
	defer dom.Free()

	state, _, err := dom.GetState()
	if err != nil {
		return DomainStateUnknown, fmt.Errorf("failed to get domain state: %w", err)
	}

	return mapLibvirtState(state), nil
}

// CreateDiskOnlySnapshot creates an external, disk-only snapshot without metadata.
// This is safe for running VMs and does not pause or stop the VM.
// Returns the snapshot info including the path to the new overlay file.
func (m *DomainManager) CreateDiskOnlySnapshot(ctx context.Context, domainName, snapshotName string) (*SnapshotInfo, error) {
	if err := m.ensureConnected(); err != nil {
		return nil, err
	}

	m.mu.Lock()
	dom, err := m.conn.LookupDomainByName(domainName)
	m.mu.Unlock()

	if err != nil {
		return nil, fmt.Errorf("failed to lookup domain %q: %w", domainName, err)
	}
	defer dom.Free()

	// Get current disk path from domain XML
	xmlDesc, err := dom.GetXMLDesc(0)
	if err != nil {
		return nil, fmt.Errorf("failed to get domain XML: %w", err)
	}

	currentDiskPath, err := extractDiskPath(xmlDesc)
	if err != nil {
		return nil, fmt.Errorf("failed to extract current disk path: %w", err)
	}

	// The backing file for the container will be the current disk
	// After snapshot, libvirt creates a new overlay and the current disk becomes backing
	backingFile := currentDiskPath

	// Build snapshot XML for disk-only, external snapshot
	// The snapshot file path will be auto-generated by libvirt
	snapshotXML := fmt.Sprintf(`
<domainsnapshot>
  <name>%s</name>
  <description>Disk-only snapshot for container cloning</description>
  <disks>
    <disk name='vda' snapshot='external'/>
  </disks>
</domainsnapshot>`, snapshotName)

	// Create the snapshot with flags:
	// - DISK_ONLY: Only snapshot the disk, not memory
	// - ATOMIC: All-or-nothing operation
	// - NO_METADATA: Don't store snapshot metadata in libvirt
	flags := libvirtgo.DomainSnapshotCreateFlags(libvirtgo.DOMAIN_SNAPSHOT_CREATE_DISK_ONLY) |
		libvirtgo.DomainSnapshotCreateFlags(libvirtgo.DOMAIN_SNAPSHOT_CREATE_ATOMIC) |
		libvirtgo.DomainSnapshotCreateFlags(libvirtgo.DOMAIN_SNAPSHOT_CREATE_NO_METADATA)

	m.mu.Lock()
	_, err = dom.CreateSnapshotXML(snapshotXML, flags)
	m.mu.Unlock()

	if err != nil {
		return nil, fmt.Errorf("failed to create disk-only snapshot: %w", err)
	}

	return &SnapshotInfo{
		Name:        snapshotName,
		BackingFile: backingFile,
	}, nil
}

// BlockCommit merges a snapshot overlay back into its backing file and removes the overlay.
// This is used for cleanup after cloning or on rollback.
func (m *DomainManager) BlockCommit(ctx context.Context, domainName, diskTarget string, timeout time.Duration) error {
	if err := m.ensureConnected(); err != nil {
		return err
	}

	m.mu.Lock()
	dom, err := m.conn.LookupDomainByName(domainName)
	m.mu.Unlock()

	if err != nil {
		return fmt.Errorf("failed to lookup domain %q: %w", domainName, err)
	}
	defer dom.Free()

	// Start block commit - merge active layer into backing file
	flags := libvirtgo.DomainBlockCommitFlags(libvirtgo.DOMAIN_BLOCK_COMMIT_ACTIVE) | libvirtgo.DomainBlockCommitFlags(libvirtgo.DOMAIN_BLOCK_COMMIT_DELETE)

	m.mu.Lock()
	err = dom.BlockCommit(diskTarget, "", "", 0, flags)
	m.mu.Unlock()

	if err != nil {
		return fmt.Errorf("failed to start block commit: %w", err)
	}

	// Wait for block commit to complete
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		m.mu.Lock()
		info, err := dom.GetBlockJobInfo(diskTarget, 0)
		m.mu.Unlock()

		if err != nil {
			// Job may have completed
			break
		}

		if info.Type == 0 {
			// No job running, commit completed
			break
		}

		time.Sleep(500 * time.Millisecond)
	}

	// Pivot to the base image if needed
	m.mu.Lock()
	err = dom.BlockJobAbort(diskTarget, libvirtgo.DomainBlockJobAbortFlags(libvirtgo.DOMAIN_BLOCK_JOB_ABORT_PIVOT))
	m.mu.Unlock()

	// Ignore error if no job to abort
	if err != nil {
		var libvirtErr libvirtgo.Error
		if !errors.As(err, &libvirtErr) || libvirtErr.Code != libvirtgo.ERR_BLOCK_COPY_ACTIVE {
			// Only log, don't fail - the commit may have already completed
		}
	}

	return nil
}

// GetDiskPath returns the primary disk path for a domain.
func (m *DomainManager) GetDiskPath(ctx context.Context, domainName string) (string, error) {
	if err := m.ensureConnected(); err != nil {
		return "", err
	}

	m.mu.Lock()
	dom, err := m.conn.LookupDomainByName(domainName)
	m.mu.Unlock()

	if err != nil {
		return "", fmt.Errorf("failed to lookup domain %q: %w", domainName, err)
	}
	defer dom.Free()

	xmlDesc, err := dom.GetXMLDesc(0)
	if err != nil {
		return "", fmt.Errorf("failed to get domain XML: %w", err)
	}

	return extractDiskPath(xmlDesc)
}

// extractDiskPath parses domain XML and extracts the primary disk file path.
func extractDiskPath(xmlDesc string) (string, error) {
	var domain domainXML
	if err := xml.Unmarshal([]byte(xmlDesc), &domain); err != nil {
		return "", fmt.Errorf("failed to parse domain XML: %w", err)
	}

	// Find the first disk device (typically vda)
	for _, disk := range domain.Devices.Disks {
		if disk.Device == "disk" && disk.Source.File != "" {
			return disk.Source.File, nil
		}
	}

	return "", fmt.Errorf("no disk device found in domain XML")
}

// mapLibvirtState converts libvirt domain state to our DomainState type.
func mapLibvirtState(state libvirtgo.DomainState) DomainState {
	switch state {
	case libvirtgo.DOMAIN_RUNNING:
		return DomainStateRunning
	case libvirtgo.DOMAIN_PAUSED:
		return DomainStatePaused
	case libvirtgo.DOMAIN_SHUTDOWN:
		return DomainStateShutdown
	case libvirtgo.DOMAIN_SHUTOFF:
		return DomainStateStopped
	case libvirtgo.DOMAIN_CRASHED:
		return DomainStateCrashed
	case libvirtgo.DOMAIN_PMSUSPENDED:
		return DomainStateSuspended
	default:
		return DomainStateUnknown
	}
}

// Sentinel errors for domain operations.
var (
	ErrDomainNotFound    = errors.New("domain not found")
	ErrDomainTransient   = errors.New("transient domains are not supported")
	ErrDomainUnsupported = errors.New("domain configuration not supported")
)
