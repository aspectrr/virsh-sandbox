package sshca

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestGenerateCA(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "sshca-test-")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	keyPath := filepath.Join(tempDir, "test_ca")
	comment := "test-ssh-ca"

	// Generate CA
	err = GenerateCA(keyPath, comment)
	if err != nil {
		t.Fatalf("GenerateCA failed: %v", err)
	}

	// Check private key exists
	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		t.Error("private key was not created")
	}

	// Check public key exists
	pubKeyPath := keyPath + ".pub"
	if _, err := os.Stat(pubKeyPath); os.IsNotExist(err) {
		t.Error("public key was not created")
	}

	// Check private key permissions
	info, err := os.Stat(keyPath)
	if err != nil {
		t.Fatalf("failed to stat private key: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("private key has wrong permissions: %o, expected 0600", info.Mode().Perm())
	}

	// Check public key content
	pubKeyBytes, err := os.ReadFile(pubKeyPath)
	if err != nil {
		t.Fatalf("failed to read public key: %v", err)
	}
	pubKey := string(pubKeyBytes)
	if !strings.HasPrefix(pubKey, "ssh-ed25519 ") {
		t.Errorf("public key has wrong format: %s", pubKey[:min(len(pubKey), 50)])
	}
	if !strings.Contains(pubKey, comment) {
		t.Errorf("public key does not contain comment: %s", pubKey)
	}
}

func TestGenerateUserKeyPair(t *testing.T) {
	comment := "test-user-key"

	privateKey, publicKey, err := GenerateUserKeyPair(comment)
	if err != nil {
		t.Fatalf("GenerateUserKeyPair failed: %v", err)
	}

	// Check private key format
	if !strings.Contains(privateKey, "OPENSSH PRIVATE KEY") {
		t.Error("private key is not in OpenSSH format")
	}

	// Check public key format
	if !strings.HasPrefix(publicKey, "ssh-ed25519 ") {
		t.Errorf("public key has wrong format: %s", publicKey[:min(len(publicKey), 50)])
	}
	if !strings.Contains(publicKey, comment) {
		t.Errorf("public key does not contain comment")
	}
}

func TestNewCA(t *testing.T) {
	cfg := DefaultConfig()
	cfg.CAKeyPath = "/nonexistent/path"
	cfg.EnforceKeyPermissions = false

	ca, err := NewCA(cfg)
	if err != nil {
		t.Fatalf("NewCA failed: %v", err)
	}

	if ca == nil {
		t.Error("NewCA returned nil")
	}

	if ca.sshKeygen == "" {
		t.Error("ssh-keygen path not set")
	}
}

func TestCAInitialize(t *testing.T) {
	// Create temp directory with CA keys
	tempDir, err := os.MkdirTemp("", "sshca-test-")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	keyPath := filepath.Join(tempDir, "test_ca")
	err = GenerateCA(keyPath, "test-ca")
	if err != nil {
		t.Fatalf("failed to generate CA: %v", err)
	}

	// Create CA instance
	cfg := Config{
		CAKeyPath:             keyPath,
		CAPubKeyPath:          keyPath + ".pub",
		WorkDir:               tempDir,
		DefaultTTL:            5 * time.Minute,
		MaxTTL:                10 * time.Minute,
		DefaultPrincipals:     []string{"sandbox"},
		EnforceKeyPermissions: true,
	}

	ca, err := NewCA(cfg)
	if err != nil {
		t.Fatalf("NewCA failed: %v", err)
	}

	// Initialize CA
	err = ca.Initialize(context.Background())
	if err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	if !ca.initialized {
		t.Error("CA not marked as initialized")
	}

	// Get public key
	pubKey, err := ca.GetPublicKey()
	if err != nil {
		t.Fatalf("GetPublicKey failed: %v", err)
	}
	if !strings.HasPrefix(pubKey, "ssh-ed25519 ") {
		t.Errorf("public key has wrong format: %s", pubKey[:min(len(pubKey), 50)])
	}
}

func TestCAInitializeNotFound(t *testing.T) {
	cfg := Config{
		CAKeyPath:             "/nonexistent/path/ssh_ca",
		CAPubKeyPath:          "/nonexistent/path/ssh_ca.pub",
		WorkDir:               "/tmp",
		EnforceKeyPermissions: false,
	}

	ca, err := NewCA(cfg)
	if err != nil {
		t.Fatalf("NewCA failed: %v", err)
	}

	err = ca.Initialize(context.Background())
	if err == nil {
		t.Error("Initialize should have failed with nonexistent key")
	}
}

func TestCAIssueCertificate(t *testing.T) {
	// Create temp directory with CA keys
	tempDir, err := os.MkdirTemp("", "sshca-test-")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	keyPath := filepath.Join(tempDir, "test_ca")
	err = GenerateCA(keyPath, "test-ca")
	if err != nil {
		t.Fatalf("failed to generate CA: %v", err)
	}

	// Generate user key
	_, userPubKey, err := GenerateUserKeyPair("test-user")
	if err != nil {
		t.Fatalf("failed to generate user key: %v", err)
	}

	// Create and initialize CA
	cfg := Config{
		CAKeyPath:             keyPath,
		CAPubKeyPath:          keyPath + ".pub",
		WorkDir:               tempDir,
		DefaultTTL:            5 * time.Minute,
		MaxTTL:                10 * time.Minute,
		DefaultPrincipals:     []string{"sandbox"},
		EnforceKeyPermissions: true,
	}

	ca, err := NewCA(cfg)
	if err != nil {
		t.Fatalf("NewCA failed: %v", err)
	}

	err = ca.Initialize(context.Background())
	if err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	// Issue certificate
	req := &CertificateRequest{
		UserID:      "test-user",
		VMID:        "test-vm",
		SandboxID:   "SBX-123",
		PublicKey:   userPubKey,
		TTL:         5 * time.Minute,
		Principals:  []string{"sandbox"},
		SourceIP:    "127.0.0.1",
		RequestTime: time.Now(),
	}

	cert, err := ca.IssueCertificate(context.Background(), req)
	if err != nil {
		t.Fatalf("IssueCertificate failed: %v", err)
	}

	// Validate certificate
	if cert.ID == "" {
		t.Error("certificate ID is empty")
	}
	if cert.Identity == "" {
		t.Error("certificate identity is empty")
	}
	if !strings.Contains(cert.Identity, "test-user") {
		t.Errorf("identity should contain user ID: %s", cert.Identity)
	}
	if !strings.Contains(cert.Identity, "test-vm") {
		t.Errorf("identity should contain VM ID: %s", cert.Identity)
	}
	if cert.Certificate == "" {
		t.Error("certificate content is empty")
	}
	if !strings.Contains(cert.Certificate, "cert-v01@openssh.com") {
		t.Error("certificate is not in OpenSSH certificate format")
	}
	if cert.SerialNumber == 0 {
		t.Error("serial number should not be zero")
	}
	if len(cert.Principals) == 0 {
		t.Error("principals should not be empty")
	}
	if cert.ValidBefore.Before(cert.ValidAfter) {
		t.Error("ValidBefore should be after ValidAfter")
	}

	// Check certificate info
	info := cert.GetCertInfo()
	if info.IsExpired {
		t.Error("certificate should not be expired immediately after issuance")
	}
	if info.TimeToExpiry <= 0 {
		t.Error("time to expiry should be positive")
	}
}

func TestCAIssueCertificateNotInitialized(t *testing.T) {
	cfg := DefaultConfig()
	cfg.EnforceKeyPermissions = false

	ca, err := NewCA(cfg)
	if err != nil {
		t.Fatalf("NewCA failed: %v", err)
	}

	req := &CertificateRequest{
		UserID:    "test-user",
		VMID:      "test-vm",
		SandboxID: "SBX-123",
		PublicKey: "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAITest test",
	}

	_, err = ca.IssueCertificate(context.Background(), req)
	if err != ErrCANotInitialized {
		t.Errorf("expected ErrCANotInitialized, got: %v", err)
	}
}

func TestCAValidateRequest(t *testing.T) {
	// Create temp directory with CA keys
	tempDir, err := os.MkdirTemp("", "sshca-test-")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	keyPath := filepath.Join(tempDir, "test_ca")
	err = GenerateCA(keyPath, "test-ca")
	if err != nil {
		t.Fatalf("failed to generate CA: %v", err)
	}

	cfg := Config{
		CAKeyPath:             keyPath,
		CAPubKeyPath:          keyPath + ".pub",
		WorkDir:               tempDir,
		DefaultTTL:            5 * time.Minute,
		MaxTTL:                10 * time.Minute,
		DefaultPrincipals:     []string{"sandbox"},
		EnforceKeyPermissions: true,
	}

	ca, err := NewCA(cfg)
	if err != nil {
		t.Fatalf("NewCA failed: %v", err)
	}

	err = ca.Initialize(context.Background())
	if err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	tests := []struct {
		name    string
		req     *CertificateRequest
		wantErr bool
	}{
		{
			name: "missing UserID",
			req: &CertificateRequest{
				VMID:      "test-vm",
				SandboxID: "SBX-123",
				PublicKey: "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAITest test",
			},
			wantErr: true,
		},
		{
			name: "missing VMID",
			req: &CertificateRequest{
				UserID:    "test-user",
				SandboxID: "SBX-123",
				PublicKey: "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAITest test",
			},
			wantErr: true,
		},
		{
			name: "missing PublicKey",
			req: &CertificateRequest{
				UserID:    "test-user",
				VMID:      "test-vm",
				SandboxID: "SBX-123",
			},
			wantErr: true,
		},
		{
			name: "invalid PublicKey format",
			req: &CertificateRequest{
				UserID:    "test-user",
				VMID:      "test-vm",
				SandboxID: "SBX-123",
				PublicKey: "not-a-valid-key",
			},
			wantErr: true,
		},
		{
			name: "unsupported key type",
			req: &CertificateRequest{
				UserID:    "test-user",
				VMID:      "test-vm",
				SandboxID: "SBX-123",
				PublicKey: "ssh-dss AAAAB3NzaC1kc3MAAACBA test",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ca.IssueCertificate(context.Background(), tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("IssueCertificate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCATTLCapping(t *testing.T) {
	// Create temp directory with CA keys
	tempDir, err := os.MkdirTemp("", "sshca-test-")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	keyPath := filepath.Join(tempDir, "test_ca")
	err = GenerateCA(keyPath, "test-ca")
	if err != nil {
		t.Fatalf("failed to generate CA: %v", err)
	}

	_, userPubKey, err := GenerateUserKeyPair("test-user")
	if err != nil {
		t.Fatalf("failed to generate user key: %v", err)
	}

	cfg := Config{
		CAKeyPath:             keyPath,
		CAPubKeyPath:          keyPath + ".pub",
		WorkDir:               tempDir,
		DefaultTTL:            5 * time.Minute,
		MaxTTL:                10 * time.Minute,
		DefaultPrincipals:     []string{"sandbox"},
		EnforceKeyPermissions: true,
	}

	ca, err := NewCA(cfg)
	if err != nil {
		t.Fatalf("NewCA failed: %v", err)
	}

	err = ca.Initialize(context.Background())
	if err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	// Request certificate with TTL exceeding max
	req := &CertificateRequest{
		UserID:    "test-user",
		VMID:      "test-vm",
		SandboxID: "SBX-123",
		PublicKey: userPubKey,
		TTL:       30 * time.Minute, // Exceeds MaxTTL
	}

	cert, err := ca.IssueCertificate(context.Background(), req)
	if err != nil {
		t.Fatalf("IssueCertificate failed: %v", err)
	}

	// Check that TTL was capped
	actualTTL := cert.ValidBefore.Sub(cert.IssuedAt)
	// Allow for some clock skew (the cert adds 1 minute before valid_after)
	if actualTTL > 11*time.Minute {
		t.Errorf("TTL should be capped to MaxTTL (10m), got: %v", actualTTL)
	}
}

func TestCertificateConnectCommand(t *testing.T) {
	cert := &Certificate{
		Principals: []string{"sandbox"},
	}

	cmd := cert.SSHConnectCommand("/path/to/key", "/path/to/key-cert.pub", "192.168.1.100", 22)

	if !strings.Contains(cmd, "-i /path/to/key") {
		t.Error("command should contain private key path")
	}
	if !strings.Contains(cmd, "CertificateFile=/path/to/key-cert.pub") {
		t.Error("command should contain certificate path")
	}
	if !strings.Contains(cmd, "sandbox@192.168.1.100") {
		t.Error("command should contain user@host")
	}
	if !strings.Contains(cmd, "-p 22") {
		t.Error("command should contain port")
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
