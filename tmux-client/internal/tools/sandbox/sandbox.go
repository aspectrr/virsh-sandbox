// Package sandbox provides SSH certificate-based access to virsh-sandbox VMs.
// It handles certificate fetching from the virsh-sandbox API and SSH connection setup.
package sandbox

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

// Config holds configuration for the sandbox tool.
type Config struct {
	// APIBaseURL is the base URL of the virsh-sandbox API (e.g., "http://localhost:8080")
	APIBaseURL string `yaml:"api_base_url"`

	// KeyDir is the directory to store ephemeral SSH keys and certificates
	KeyDir string `yaml:"key_dir"`

	// DefaultTTLMinutes is the default certificate TTL to request
	DefaultTTLMinutes int `yaml:"default_ttl_minutes"`

	// UserID is the user identifier for certificate requests
	UserID string `yaml:"user_id"`

	// HTTPTimeout is the timeout for API requests
	HTTPTimeout time.Duration `yaml:"http_timeout"`
}

// DefaultConfig returns sensible defaults for the sandbox tool.
func DefaultConfig() *Config {
	return &Config{
		APIBaseURL:        "http://localhost:8080",
		KeyDir:            "/tmp/sandbox-keys",
		DefaultTTLMinutes: 5,
		UserID:            "",
		HTTPTimeout:       30 * time.Second,
	}
}

// Tool provides sandbox access operations.
type Tool struct {
	config     *Config
	httpClient *http.Client
}

// NewTool creates a new sandbox tool.
func NewTool(cfg *Config) (*Tool, error) {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	// Create key directory if it doesn't exist
	if err := os.MkdirAll(cfg.KeyDir, 0o700); err != nil {
		return nil, fmt.Errorf("create key directory: %w", err)
	}

	return &Tool{
		config: cfg,
		httpClient: &http.Client{
			Timeout: cfg.HTTPTimeout,
		},
	}, nil
}

// AccessRequest represents a request for sandbox access.
type AccessRequest struct {
	SandboxID  string `json:"sandbox_id"`
	UserID     string `json:"user_id"`
	PublicKey  string `json:"public_key"`
	TTLMinutes int    `json:"ttl_minutes,omitempty"`
}

// AccessResponse represents the response from a successful access request.
type AccessResponse struct {
	CertificateID  string    `json:"certificate_id"`
	Certificate    string    `json:"certificate"`
	VMIPAddress    string    `json:"vm_ip_address"`
	SSHPort        int       `json:"ssh_port"`
	Username       string    `json:"username"`
	ValidUntil     time.Time `json:"valid_until"`
	TTLSeconds     int       `json:"ttl_seconds"`
	ConnectCommand string    `json:"connect_command"`
	Instructions   string    `json:"instructions"`
}

// ConnectionInfo contains everything needed to SSH into a sandbox.
type ConnectionInfo struct {
	// SandboxID is the sandbox being accessed
	SandboxID string

	// PrivateKeyPath is the path to the ephemeral private key
	PrivateKeyPath string

	// CertificatePath is the path to the SSH certificate
	CertificatePath string

	// VMIPAddress is the IP of the sandbox VM
	VMIPAddress string

	// SSHPort is the SSH port (usually 22)
	SSHPort int

	// Username is the SSH username
	Username string

	// ValidUntil is when the certificate expires
	ValidUntil time.Time

	// SSHCommand is the full SSH command to connect
	SSHCommand string
}

// RequestAccess requests SSH certificate access to a sandbox.
// It generates an ephemeral key pair, requests a certificate, and returns connection info.
func (t *Tool) RequestAccess(ctx context.Context, sandboxID string, ttlMinutes int) (*ConnectionInfo, error) {
	if sandboxID == "" {
		return nil, fmt.Errorf("sandbox_id is required")
	}

	if ttlMinutes <= 0 {
		ttlMinutes = t.config.DefaultTTLMinutes
	}

	userID := t.config.UserID
	if userID == "" {
		userID = os.Getenv("USER")
		if userID == "" {
			userID = "agent"
		}
	}

	// Generate ephemeral key pair
	privateKeyPath, publicKey, err := t.generateEphemeralKeyPair(sandboxID)
	if err != nil {
		return nil, fmt.Errorf("generate key pair: %w", err)
	}

	// Request certificate from API
	req := AccessRequest{
		SandboxID:  sandboxID,
		UserID:     userID,
		PublicKey:  publicKey,
		TTLMinutes: ttlMinutes,
	}

	resp, err := t.requestCertificate(ctx, req)
	if err != nil {
		// Clean up key on failure
		os.Remove(privateKeyPath)
		return nil, fmt.Errorf("request certificate: %w", err)
	}

	// Save certificate
	certPath := privateKeyPath + "-cert.pub"
	if err := os.WriteFile(certPath, []byte(resp.Certificate), 0o600); err != nil {
		os.Remove(privateKeyPath)
		return nil, fmt.Errorf("save certificate: %w", err)
	}

	// Build SSH command
	sshCommand := fmt.Sprintf("ssh -i %s -o CertificateFile=%s -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -p %d %s@%s",
		privateKeyPath, certPath, resp.SSHPort, resp.Username, resp.VMIPAddress)

	return &ConnectionInfo{
		SandboxID:       sandboxID,
		PrivateKeyPath:  privateKeyPath,
		CertificatePath: certPath,
		VMIPAddress:     resp.VMIPAddress,
		SSHPort:         resp.SSHPort,
		Username:        resp.Username,
		ValidUntil:      resp.ValidUntil,
		SSHCommand:      sshCommand,
	}, nil
}

// GetSSHArgs returns the arguments needed for ssh command to connect to the sandbox.
func (c *ConnectionInfo) GetSSHArgs() []string {
	return []string{
		"-i", c.PrivateKeyPath,
		"-o", "CertificateFile=" + c.CertificatePath,
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null",
		"-p", fmt.Sprintf("%d", c.SSHPort),
		fmt.Sprintf("%s@%s", c.Username, c.VMIPAddress),
	}
}

// Cleanup removes the ephemeral keys and certificates.
func (c *ConnectionInfo) Cleanup() error {
	var errs []error

	if c.PrivateKeyPath != "" {
		if err := os.Remove(c.PrivateKeyPath); err != nil && !os.IsNotExist(err) {
			errs = append(errs, fmt.Errorf("remove private key: %w", err))
		}
		// Also remove the .pub file if it exists
		pubPath := c.PrivateKeyPath + ".pub"
		os.Remove(pubPath) // Ignore errors for .pub
	}

	if c.CertificatePath != "" {
		if err := os.Remove(c.CertificatePath); err != nil && !os.IsNotExist(err) {
			errs = append(errs, fmt.Errorf("remove certificate: %w", err))
		}
	}

	if len(errs) > 0 {
		return errs[0]
	}
	return nil
}

// IsExpired returns true if the certificate has expired.
func (c *ConnectionInfo) IsExpired() bool {
	return time.Now().After(c.ValidUntil)
}

// TimeToExpiry returns the duration until the certificate expires.
func (c *ConnectionInfo) TimeToExpiry() time.Duration {
	remaining := time.Until(c.ValidUntil)
	if remaining < 0 {
		return 0
	}
	return remaining
}

// generateEphemeralKeyPair generates an ed25519 SSH key pair and saves the private key.
// Returns the path to the private key and the public key string.
func (t *Tool) generateEphemeralKeyPair(sandboxID string) (string, string, error) {
	// Generate ed25519 key pair
	pubKey, privKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return "", "", fmt.Errorf("generate ed25519 key: %w", err)
	}

	// Convert to SSH key format
	sshPubKey, err := ssh.NewPublicKey(pubKey)
	if err != nil {
		return "", "", fmt.Errorf("convert to ssh public key: %w", err)
	}

	// Generate filename with timestamp to avoid collisions
	timestamp := time.Now().UnixNano()
	keyName := fmt.Sprintf("sandbox_%s_%d", sanitizeFilename(sandboxID), timestamp)
	privateKeyPath := filepath.Join(t.config.KeyDir, keyName)

	// Encode private key to OpenSSH format
	privateKeyPEM, err := encodeED25519PrivateKey(privKey)
	if err != nil {
		return "", "", fmt.Errorf("encode private key: %w", err)
	}

	// Save private key with secure permissions
	if err := os.WriteFile(privateKeyPath, privateKeyPEM, 0o600); err != nil {
		return "", "", fmt.Errorf("save private key: %w", err)
	}

	// Format public key as authorized_keys format
	publicKeyStr := strings.TrimSpace(string(ssh.MarshalAuthorizedKey(sshPubKey)))

	return privateKeyPath, publicKeyStr, nil
}

// requestCertificate makes the API call to request a certificate.
func (t *Tool) requestCertificate(ctx context.Context, req AccessRequest) (*AccessResponse, error) {
	url := fmt.Sprintf("%s/v1/access/request", strings.TrimSuffix(t.config.APIBaseURL, "/"))

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := t.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		// Try to parse error response
		var errResp struct {
			Error   string `json:"error"`
			Details string `json:"details"`
		}
		if json.Unmarshal(respBody, &errResp) == nil && errResp.Error != "" {
			return nil, fmt.Errorf("API error (%d): %s - %s", resp.StatusCode, errResp.Error, errResp.Details)
		}
		return nil, fmt.Errorf("API error (%d): %s", resp.StatusCode, string(respBody))
	}

	var accessResp AccessResponse
	if err := json.Unmarshal(respBody, &accessResp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	return &accessResp, nil
}

// encodeED25519PrivateKey encodes an ed25519 private key to OpenSSH format.
func encodeED25519PrivateKey(key ed25519.PrivateKey) ([]byte, error) {
	// OpenSSH private key format
	// This is a simplified version - for production, consider using a library

	pubKey := key.Public().(ed25519.PublicKey)

	// Generate random check bytes
	checkBytes := make([]byte, 4)
	if _, err := rand.Read(checkBytes); err != nil {
		return nil, err
	}

	// Build the private key section
	var privSection bytes.Buffer

	// Check bytes (repeated)
	privSection.Write(checkBytes)
	privSection.Write(checkBytes)

	// Key type
	writeString(&privSection, "ssh-ed25519")

	// Public key
	writeBytes(&privSection, pubKey)

	// Private key (64 bytes = 32 byte seed + 32 byte public)
	writeBytes(&privSection, key)

	// Comment (empty)
	writeString(&privSection, "")

	// Padding
	for i := 1; (privSection.Len() % 8) != 0; i++ {
		privSection.WriteByte(byte(i))
	}

	// Build the full key
	var fullKey bytes.Buffer

	// Auth magic
	fullKey.WriteString("openssh-key-v1\x00")

	// Cipher name (none)
	writeString(&fullKey, "none")

	// KDF name (none)
	writeString(&fullKey, "none")

	// KDF options (empty)
	writeString(&fullKey, "")

	// Number of keys
	writeUint32(&fullKey, 1)

	// Public key section
	var pubSection bytes.Buffer
	writeString(&pubSection, "ssh-ed25519")
	writeBytes(&pubSection, pubKey)
	writeBytes(&fullKey, pubSection.Bytes())

	// Private key section
	writeBytes(&fullKey, privSection.Bytes())

	// PEM encode
	block := &pem.Block{
		Type:  "OPENSSH PRIVATE KEY",
		Bytes: fullKey.Bytes(),
	}

	return pem.EncodeToMemory(block), nil
}

func writeString(w *bytes.Buffer, s string) {
	writeBytes(w, []byte(s))
}

func writeBytes(w *bytes.Buffer, b []byte) {
	writeUint32(w, uint32(len(b)))
	w.Write(b)
}

func writeUint32(w *bytes.Buffer, n uint32) {
	w.WriteByte(byte(n >> 24))
	w.WriteByte(byte(n >> 16))
	w.WriteByte(byte(n >> 8))
	w.WriteByte(byte(n))
}

// sanitizeFilename removes characters that aren't safe for filenames.
func sanitizeFilename(s string) string {
	// Replace unsafe characters with underscores
	result := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-' || c == '_' {
			result = append(result, c)
		} else {
			result = append(result, '_')
		}
	}
	return string(result)
}

// CheckAPIHealth checks if the virsh-sandbox API is reachable.
func (t *Tool) CheckAPIHealth(ctx context.Context) error {
	url := fmt.Sprintf("%s/v1/health", strings.TrimSuffix(t.config.APIBaseURL, "/"))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API unhealthy: status %d", resp.StatusCode)
	}

	return nil
}

// GetCAPublicKey retrieves the CA public key from the API.
func (t *Tool) GetCAPublicKey(ctx context.Context) (string, error) {
	url := fmt.Sprintf("%s/v1/access/ca-pubkey", strings.TrimSuffix(t.config.APIBaseURL, "/"))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API error: status %d", resp.StatusCode)
	}

	var caResp struct {
		PublicKey string `json:"public_key"`
		Usage     string `json:"usage"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&caResp); err != nil {
		return "", fmt.Errorf("parse response: %w", err)
	}

	return caResp.PublicKey, nil
}
