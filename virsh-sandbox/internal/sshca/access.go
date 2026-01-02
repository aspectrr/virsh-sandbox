// Package sshca provides SSH Certificate Authority management for ephemeral sandbox access.
package sshca

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

// AccessService orchestrates SSH certificate-based access to sandboxes.
// It handles certificate issuance, session tracking, and cleanup.
type AccessService struct {
	ca        *CA
	store     CertificateStore
	vmLookup  VMInfoProvider
	timeNowFn func() time.Time
	mu        sync.RWMutex

	// Configuration
	defaultTTL time.Duration
	maxTTL     time.Duration
	sshPort    int
	username   string
}

// VMInfoProvider defines the interface for looking up VM/sandbox information.
type VMInfoProvider interface {
	// GetSandboxIP returns the IP address of a sandbox.
	GetSandboxIP(ctx context.Context, sandboxID string) (string, error)

	// GetSandboxVMName returns the VM name for a sandbox.
	GetSandboxVMName(ctx context.Context, sandboxID string) (string, error)

	// IsSandboxRunning checks if the sandbox is in a running state.
	IsSandboxRunning(ctx context.Context, sandboxID string) (bool, error)
}

// AccessServiceConfig configures the access service.
type AccessServiceConfig struct {
	// DefaultTTL is the default certificate lifetime.
	DefaultTTL time.Duration

	// MaxTTL is the maximum allowed certificate lifetime.
	MaxTTL time.Duration

	// SSHPort is the SSH port on VMs (default 22).
	SSHPort int

	// Username is the SSH username (default "sandbox").
	Username string
}

// DefaultAccessServiceConfig returns sensible defaults.
func DefaultAccessServiceConfig() AccessServiceConfig {
	return AccessServiceConfig{
		DefaultTTL: 5 * time.Minute,
		MaxTTL:     10 * time.Minute,
		SSHPort:    22,
		Username:   "sandbox",
	}
}

// AccessServiceOption configures the AccessService.
type AccessServiceOption func(*AccessService)

// WithAccessTimeNow overrides the clock (useful for tests).
func WithAccessTimeNow(fn func() time.Time) AccessServiceOption {
	return func(s *AccessService) { s.timeNowFn = fn }
}

// NewAccessService creates a new access service.
func NewAccessService(ca *CA, store CertificateStore, vmLookup VMInfoProvider, cfg AccessServiceConfig, opts ...AccessServiceOption) *AccessService {
	if cfg.DefaultTTL == 0 {
		cfg.DefaultTTL = 5 * time.Minute
	}
	if cfg.MaxTTL == 0 {
		cfg.MaxTTL = 10 * time.Minute
	}
	if cfg.SSHPort == 0 {
		cfg.SSHPort = 22
	}
	if cfg.Username == "" {
		cfg.Username = "sandbox"
	}

	s := &AccessService{
		ca:         ca,
		store:      store,
		vmLookup:   vmLookup,
		timeNowFn:  time.Now,
		defaultTTL: cfg.DefaultTTL,
		maxTTL:     cfg.MaxTTL,
		sshPort:    cfg.SSHPort,
		username:   cfg.Username,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// RequestAccess issues a short-lived SSH certificate for sandbox access.
func (s *AccessService) RequestAccess(ctx context.Context, req *AccessRequest) (*AccessResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Validate request
	if req.SandboxID == "" {
		return nil, fmt.Errorf("sandbox_id is required")
	}
	if req.UserID == "" {
		return nil, fmt.Errorf("user_id is required")
	}
	if req.PublicKey == "" {
		return nil, fmt.Errorf("public_key is required")
	}

	// Determine TTL
	ttl := time.Duration(req.TTLMinutes) * time.Minute
	if ttl == 0 {
		ttl = s.defaultTTL
	}
	if ttl < time.Minute {
		ttl = time.Minute
	}
	if ttl > s.maxTTL {
		ttl = s.maxTTL
	}

	// Check if sandbox is running
	running, err := s.vmLookup.IsSandboxRunning(ctx, req.SandboxID)
	if err != nil {
		return nil, fmt.Errorf("check sandbox status: %w", err)
	}
	if !running {
		return nil, fmt.Errorf("sandbox %s is not running", req.SandboxID)
	}

	// Get sandbox IP
	vmIP, err := s.vmLookup.GetSandboxIP(ctx, req.SandboxID)
	if err != nil {
		return nil, fmt.Errorf("get sandbox IP: %w", err)
	}
	if vmIP == "" {
		return nil, fmt.Errorf("sandbox %s has no IP address", req.SandboxID)
	}

	// Get VM name for certificate identity
	vmName, err := s.vmLookup.GetSandboxVMName(ctx, req.SandboxID)
	if err != nil {
		return nil, fmt.Errorf("get sandbox VM name: %w", err)
	}

	now := s.timeNowFn()
	if req.RequestTime.IsZero() {
		req.RequestTime = now
	}

	// Issue certificate
	certReq := &CertificateRequest{
		UserID:      req.UserID,
		VMID:        vmName,
		SandboxID:   req.SandboxID,
		PublicKey:   req.PublicKey,
		TTL:         ttl,
		Principals:  []string{s.username},
		SourceIP:    req.SourceIP,
		RequestTime: req.RequestTime,
	}

	cert, err := s.ca.IssueCertificate(ctx, certReq)
	if err != nil {
		return nil, fmt.Errorf("issue certificate: %w", err)
	}

	// Calculate public key fingerprint
	fingerprint := s.calculateFingerprint(req.PublicKey)

	// Persist certificate record
	record := &CertificateRecord{
		ID:                   cert.ID,
		SandboxID:            req.SandboxID,
		UserID:               req.UserID,
		VMID:                 vmName,
		Identity:             cert.Identity,
		SerialNumber:         cert.SerialNumber,
		Principals:           cert.Principals,
		PublicKeyFingerprint: fingerprint,
		ValidAfter:           cert.ValidAfter,
		ValidBefore:          cert.ValidBefore,
		SourceIP:             req.SourceIP,
		Status:               CertStatusActive,
		IssuedAt:             now,
	}

	if s.store != nil {
		if err := s.store.CreateCertificate(ctx, record); err != nil {
			return nil, fmt.Errorf("persist certificate: %w", err)
		}
	}

	// Build response
	validUntil := cert.ValidBefore
	ttlSeconds := int(validUntil.Sub(now).Seconds())

	connectCmd := fmt.Sprintf("ssh -i /path/to/key -o CertificateFile=/path/to/key-cert.pub -o StrictHostKeyChecking=no %s@%s",
		s.username, vmIP)

	return &AccessResponse{
		CertificateID:  cert.ID,
		Certificate:    cert.Certificate,
		VMIPAddress:    vmIP,
		SSHPort:        s.sshPort,
		Username:       s.username,
		ValidUntil:     validUntil,
		TTLSeconds:     ttlSeconds,
		ConnectCommand: connectCmd,
	}, nil
}

// RevokeAccess revokes a certificate, immediately terminating access.
func (s *AccessService) RevokeAccess(ctx context.Context, certificateID, reason string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.store == nil {
		return fmt.Errorf("no certificate store configured")
	}

	// Get certificate
	cert, err := s.store.GetCertificate(ctx, certificateID)
	if err != nil {
		return fmt.Errorf("get certificate: %w", err)
	}

	if cert.Status == CertStatusRevoked {
		return ErrCertAlreadyRevoked
	}

	// Revoke certificate
	if err := s.store.RevokeCertificate(ctx, certificateID, reason); err != nil {
		return fmt.Errorf("revoke certificate: %w", err)
	}

	// End any active sessions for this certificate
	sessions, err := s.store.GetSessionsByCertificate(ctx, certificateID)
	if err != nil {
		return fmt.Errorf("get sessions: %w", err)
	}

	for _, session := range sessions {
		if session.Status == SessionStatusActive || session.Status == SessionStatusPending {
			now := s.timeNowFn()
			if err := s.store.EndSession(ctx, session.ID, now, "certificate revoked: "+reason); err != nil {
				// Log but continue
				continue
			}
		}
	}

	return nil
}

// RecordSessionStart records the start of an SSH session.
func (s *AccessService) RecordSessionStart(ctx context.Context, certificateID, sourceIP string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.store == nil {
		return "", fmt.Errorf("no certificate store configured")
	}

	// Get certificate
	cert, err := s.store.GetCertificate(ctx, certificateID)
	if err != nil {
		return "", fmt.Errorf("get certificate: %w", err)
	}

	// Validate certificate is still valid
	if cert.Status != CertStatusActive {
		return "", fmt.Errorf("certificate status is %s, not active", cert.Status)
	}
	if cert.IsExpired() {
		return "", fmt.Errorf("certificate has expired")
	}

	// Get VM IP
	vmIP, err := s.vmLookup.GetSandboxIP(ctx, cert.SandboxID)
	if err != nil {
		vmIP = "" // Non-fatal
	}

	// Create session record
	sessionID := s.generateSessionID()
	now := s.timeNowFn()

	session := &AccessSession{
		ID:            sessionID,
		CertificateID: certificateID,
		SandboxID:     cert.SandboxID,
		UserID:        cert.UserID,
		VMID:          cert.VMID,
		VMIPAddress:   vmIP,
		SourceIP:      sourceIP,
		Status:        SessionStatusActive,
		StartedAt:     now,
	}

	if err := s.store.CreateSession(ctx, session); err != nil {
		return "", fmt.Errorf("create session: %w", err)
	}

	// Update certificate last used
	if err := s.store.UpdateCertificateLastUsed(ctx, certificateID, now); err != nil {
		// Non-fatal
	}

	return sessionID, nil
}

// RecordSessionEnd records the end of an SSH session.
func (s *AccessService) RecordSessionEnd(ctx context.Context, sessionID, reason string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.store == nil {
		return fmt.Errorf("no certificate store configured")
	}

	now := s.timeNowFn()
	return s.store.EndSession(ctx, sessionID, now, reason)
}

// GetCertificate retrieves certificate information.
func (s *AccessService) GetCertificate(ctx context.Context, id string) (*CertificateRecord, error) {
	if s.store == nil {
		return nil, fmt.Errorf("no certificate store configured")
	}
	return s.store.GetCertificate(ctx, id)
}

// ListCertificates lists certificates with optional filtering.
func (s *AccessService) ListCertificates(ctx context.Context, filter CertificateFilter, opts *ListOptions) ([]*CertificateRecord, error) {
	if s.store == nil {
		return nil, fmt.Errorf("no certificate store configured")
	}
	return s.store.ListCertificates(ctx, filter, opts)
}

// GetActiveCertificatesForSandbox returns all active certificates for a sandbox.
func (s *AccessService) GetActiveCertificatesForSandbox(ctx context.Context, sandboxID string) ([]*CertificateRecord, error) {
	if s.store == nil {
		return nil, fmt.Errorf("no certificate store configured")
	}
	filter := CertificateFilter{
		SandboxID:  &sandboxID,
		ActiveOnly: true,
	}
	return s.store.ListCertificates(ctx, filter, nil)
}

// GetActiveSessionsForSandbox returns all active sessions for a sandbox.
func (s *AccessService) GetActiveSessionsForSandbox(ctx context.Context, sandboxID string) ([]*AccessSession, error) {
	if s.store == nil {
		return nil, fmt.Errorf("no certificate store configured")
	}
	filter := SessionFilter{
		SandboxID:  &sandboxID,
		ActiveOnly: true,
	}
	return s.store.ListSessions(ctx, filter, nil)
}

// CleanupExpiredCertificates marks expired certificates and ends associated sessions.
func (s *AccessService) CleanupExpiredCertificates(ctx context.Context) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.store == nil {
		return 0, fmt.Errorf("no certificate store configured")
	}

	// Mark expired certificates
	count, err := s.store.ExpireCertificates(ctx)
	if err != nil {
		return 0, fmt.Errorf("expire certificates: %w", err)
	}

	// End sessions for expired certificates
	filter := SessionFilter{
		ActiveOnly: true,
	}
	sessions, err := s.store.ListSessions(ctx, filter, nil)
	if err != nil {
		return count, fmt.Errorf("list sessions: %w", err)
	}

	now := s.timeNowFn()
	for _, session := range sessions {
		// Check if certificate is expired
		cert, err := s.store.GetCertificate(ctx, session.CertificateID)
		if err != nil {
			continue
		}
		if cert.IsExpired() || cert.Status == CertStatusExpired {
			if err := s.store.EndSession(ctx, session.ID, now, "certificate expired"); err != nil {
				// Log but continue
				continue
			}
		}
	}

	return count, nil
}

// RevokeAllForSandbox revokes all certificates for a sandbox.
// This is typically called when destroying a sandbox.
func (s *AccessService) RevokeAllForSandbox(ctx context.Context, sandboxID, reason string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.store == nil {
		return nil // No store configured, nothing to revoke
	}

	// Get all active certificates for the sandbox
	filter := CertificateFilter{
		SandboxID:  &sandboxID,
		ActiveOnly: true,
	}
	certs, err := s.store.ListCertificates(ctx, filter, nil)
	if err != nil {
		return fmt.Errorf("list certificates: %w", err)
	}

	// Revoke each certificate
	for _, cert := range certs {
		if err := s.store.RevokeCertificate(ctx, cert.ID, reason); err != nil {
			// Log but continue
			continue
		}
	}

	// End all active sessions
	sessionFilter := SessionFilter{
		SandboxID:  &sandboxID,
		ActiveOnly: true,
	}
	sessions, err := s.store.ListSessions(ctx, sessionFilter, nil)
	if err != nil {
		return fmt.Errorf("list sessions: %w", err)
	}

	now := s.timeNowFn()
	for _, session := range sessions {
		if err := s.store.EndSession(ctx, session.ID, now, reason); err != nil {
			// Log but continue
			continue
		}
	}

	return nil
}

// GetCAPublicKey returns the CA public key for VM configuration.
func (s *AccessService) GetCAPublicKey() (string, error) {
	return s.ca.GetPublicKey()
}

// calculateFingerprint computes the SHA256 fingerprint of a public key.
func (s *AccessService) calculateFingerprint(publicKey string) string {
	parts := strings.SplitN(publicKey, " ", 3)
	if len(parts) < 2 {
		return ""
	}

	keyData, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return ""
	}

	hash := sha256.Sum256(keyData)
	return fmt.Sprintf("SHA256:%s", base64.StdEncoding.EncodeToString(hash[:]))
}

// generateSessionID generates a unique session identifier.
func (s *AccessService) generateSessionID() string {
	id := uuid.NewString()
	return fmt.Sprintf("SESS-%s", strings.ToUpper(id[:8]))
}

// StartCleanupRoutine starts a background goroutine to periodically clean up expired certificates.
func (s *AccessService) StartCleanupRoutine(ctx context.Context, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if _, err := s.CleanupExpiredCertificates(ctx); err != nil {
					// Log error but continue
				}
			}
		}
	}()
}
