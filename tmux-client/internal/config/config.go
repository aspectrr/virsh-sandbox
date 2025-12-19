// Package config provides configuration management for the tmux agent API.
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Config is the top-level configuration structure.
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Tmux     TmuxConfig     `yaml:"tmux"`
	File     FileConfig     `yaml:"file"`
	Command  CommandConfig  `yaml:"command"`
	Human    HumanConfig    `yaml:"human"`
	Plan     PlanConfig     `yaml:"plan"`
	Audit    AuditConfig    `yaml:"audit"`
	Security SecurityConfig `yaml:"security"`
}

// ServerConfig contains HTTP server configuration.
type ServerConfig struct {
	Host            string        `yaml:"host"`
	Port            int           `yaml:"port"`
	ReadTimeout     time.Duration `yaml:"read_timeout"`
	WriteTimeout    time.Duration `yaml:"write_timeout"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
	TLSCertFile     string        `yaml:"tls_cert_file"`
	TLSKeyFile      string        `yaml:"tls_key_file"`
	EnableTLS       bool          `yaml:"enable_tls"`
}

// TmuxConfig contains tmux-specific configuration.
type TmuxConfig struct {
	SocketPath       string   `yaml:"socket_path"`
	DefaultSession   string   `yaml:"default_session"`
	MaxPaneReadLines int      `yaml:"max_pane_read_lines"`
	AllowedKeys      []string `yaml:"allowed_keys"`
}

// FileConfig contains file operation configuration.
type FileConfig struct {
	RootDirectory    string   `yaml:"root_directory"`
	AllowedPaths     []string `yaml:"allowed_paths"`
	DeniedPaths      []string `yaml:"denied_paths"`
	MaxFileSize      int64    `yaml:"max_file_size"`
	MaxReadLines     int      `yaml:"max_read_lines"`
	AllowDelete      bool     `yaml:"allow_delete"`
	AllowOverwrite   bool     `yaml:"allow_overwrite"`
	BackupOnEdit     bool     `yaml:"backup_on_edit"`
	BackupDirectory  string   `yaml:"backup_directory"`
	DeniedExtensions []string `yaml:"denied_extensions"`
}

// CommandConfig contains command execution configuration.
type CommandConfig struct {
	AllowedCommands  []string      `yaml:"allowed_commands"`
	DeniedCommands   []string      `yaml:"denied_commands"`
	AllowedPaths     []string      `yaml:"allowed_paths"`
	DefaultTimeout   time.Duration `yaml:"default_timeout"`
	MaxTimeout       time.Duration `yaml:"max_timeout"`
	MaxOutputSize    int           `yaml:"max_output_size"`
	AllowEnvVars     bool          `yaml:"allow_env_vars"`
	AllowedEnvVars   []string      `yaml:"allowed_env_vars"`
	InheritEnv       bool          `yaml:"inherit_env"`
	WorkingDirectory string        `yaml:"working_directory"`
}

// HumanConfig contains human approval configuration.
type HumanConfig struct {
	DefaultTimeout   time.Duration `yaml:"default_timeout"`
	MaxPending       int           `yaml:"max_pending"`
	NotifyCommand    string        `yaml:"notify_command"`
	NotifyArgs       []string      `yaml:"notify_args"`
	WebhookURL       string        `yaml:"webhook_url"`
	RequireApproval  []string      `yaml:"require_approval"`
	AutoApproveTypes []string      `yaml:"auto_approve_types"`
}

// PlanConfig contains plan management configuration.
type PlanConfig struct {
	MaxPlans        int           `yaml:"max_plans"`
	MaxStepsPerPlan int           `yaml:"max_steps_per_plan"`
	RetentionPeriod time.Duration `yaml:"retention_period"`
	PersistPlans    bool          `yaml:"persist_plans"`
	PlanDirectory   string        `yaml:"plan_directory"`
}

// AuditConfig contains audit logging configuration.
type AuditConfig struct {
	Enabled         bool          `yaml:"enabled"`
	LogFile         string        `yaml:"log_file"`
	RotateSize      int64         `yaml:"rotate_size"`
	RotateCount     int           `yaml:"rotate_count"`
	RetentionPeriod time.Duration `yaml:"retention_period"`
	LogArguments    bool          `yaml:"log_arguments"`
	LogResults      bool          `yaml:"log_results"`
	SensitiveFields []string      `yaml:"sensitive_fields"`
}

// SecurityConfig contains security-related configuration.
type SecurityConfig struct {
	EnableAuth      bool     `yaml:"enable_auth"`
	APIKeys         []string `yaml:"api_keys"`
	AllowedIPs      []string `yaml:"allowed_ips"`
	RateLimitPerMin int      `yaml:"rate_limit_per_min"`
	MaxRequestSize  int64    `yaml:"max_request_size"`
}

// DefaultConfig returns a configuration with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Host:            "127.0.0.1",
			Port:            8080,
			ReadTimeout:     30 * time.Second,
			WriteTimeout:    30 * time.Second,
			ShutdownTimeout: 10 * time.Second,
			EnableTLS:       false,
		},
		Tmux: TmuxConfig{
			SocketPath:       "",
			DefaultSession:   "",
			MaxPaneReadLines: 1000,
			AllowedKeys: []string{
				"Enter",
				"C-c",
				"C-d",
				"C-z",
				"Escape",
				"Up",
				"Down",
				"Left",
				"Right",
				"Tab",
			},
		},
		File: FileConfig{
			RootDirectory: "",
			AllowedPaths:  []string{},
			DeniedPaths: []string{
				"/etc/passwd",
				"/etc/shadow",
				"/etc/sudoers",
				"~/.ssh",
				"~/.gnupg",
			},
			MaxFileSize:     10 * 1024 * 1024, // 10MB
			MaxReadLines:    10000,
			AllowDelete:     false,
			AllowOverwrite:  true,
			BackupOnEdit:    true,
			BackupDirectory: "",
			DeniedExtensions: []string{
				".key",
				".pem",
				".p12",
				".pfx",
			},
		},
		Command: CommandConfig{
			AllowedCommands: []string{
				"ls", "cat", "head", "tail", "grep", "find", "wc",
				"echo", "date", "whoami", "pwd", "env", "which",
				"git", "go", "make", "npm", "node", "python", "python3",
				"docker", "kubectl", "curl", "wget",
				"ps", "top", "df", "du", "free", "uptime",
				"mkdir", "cp", "mv", "rm", "touch", "chmod",
				"tar", "gzip", "gunzip", "zip", "unzip",
				"sed", "awk", "sort", "uniq", "cut", "tr",
				"diff", "patch", "file", "stat",
			},
			DeniedCommands: []string{
				"sudo", "su", "passwd", "chown", "chgrp",
				"mount", "umount", "fdisk", "mkfs",
				"iptables", "firewall-cmd", "ufw",
				"systemctl", "service", "init",
				"reboot", "shutdown", "halt", "poweroff",
				"dd", "shred", "wipefs",
				"ssh", "scp", "sftp", "rsync",
				"nc", "netcat", "nmap", "tcpdump",
				"eval", "exec", "source",
			},
			AllowedPaths:     []string{},
			DefaultTimeout:   30 * time.Second,
			MaxTimeout:       5 * time.Minute,
			MaxOutputSize:    1024 * 1024, // 1MB
			AllowEnvVars:     false,
			AllowedEnvVars:   []string{},
			InheritEnv:       true,
			WorkingDirectory: "",
		},
		Human: HumanConfig{
			DefaultTimeout: 5 * time.Minute,
			MaxPending:     100,
			NotifyCommand:  "",
			NotifyArgs:     []string{},
			WebhookURL:     "",
			RequireApproval: []string{
				"destructive",
				"irreversible",
				"sensitive",
			},
			AutoApproveTypes: []string{},
		},
		Plan: PlanConfig{
			MaxPlans:        100,
			MaxStepsPerPlan: 50,
			RetentionPeriod: 24 * time.Hour,
			PersistPlans:    true,
			PlanDirectory:   "",
		},
		Audit: AuditConfig{
			Enabled:         true,
			LogFile:         "audit.log",
			RotateSize:      100 * 1024 * 1024, // 100MB
			RotateCount:     5,
			RetentionPeriod: 30 * 24 * time.Hour, // 30 days
			LogArguments:    true,
			LogResults:      true,
			SensitiveFields: []string{
				"password",
				"secret",
				"token",
				"api_key",
				"private_key",
			},
		},
		Security: SecurityConfig{
			EnableAuth:      false,
			APIKeys:         []string{},
			AllowedIPs:      []string{"127.0.0.1", "::1"},
			RateLimitPerMin: 60,
			MaxRequestSize:  10 * 1024 * 1024, // 10MB
		},
	}
}

// Load loads configuration from a YAML file.
func Load(path string) (*Config, error) {
	cfg := DefaultConfig()

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Apply environment variable overrides
	cfg.applyEnvOverrides()

	// Expand paths
	cfg.expandPaths()

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

// LoadOrDefault loads configuration from a file or returns defaults.
func LoadOrDefault(path string) (*Config, error) {
	if path == "" {
		cfg := DefaultConfig()
		cfg.applyEnvOverrides()
		cfg.expandPaths()
		return cfg, nil
	}

	return Load(path)
}

// applyEnvOverrides applies environment variable overrides to the configuration.
func (c *Config) applyEnvOverrides() {
	// Server
	if v := os.Getenv("TMUX_AGENT_HOST"); v != "" {
		c.Server.Host = v
	}
	if v := os.Getenv("TMUX_AGENT_PORT"); v != "" {
		if port, err := strconv.Atoi(v); err == nil {
			c.Server.Port = port
		}
	}
	if v := os.Getenv("TMUX_AGENT_TLS_CERT"); v != "" {
		c.Server.TLSCertFile = v
		c.Server.EnableTLS = true
	}
	if v := os.Getenv("TMUX_AGENT_TLS_KEY"); v != "" {
		c.Server.TLSKeyFile = v
	}

	// File
	if v := os.Getenv("TMUX_AGENT_ROOT_DIR"); v != "" {
		c.File.RootDirectory = v
	}
	if v := os.Getenv("TMUX_AGENT_BACKUP_DIR"); v != "" {
		c.File.BackupDirectory = v
	}

	// Command
	if v := os.Getenv("TMUX_AGENT_WORK_DIR"); v != "" {
		c.Command.WorkingDirectory = v
	}

	// Audit
	if v := os.Getenv("TMUX_AGENT_AUDIT_FILE"); v != "" {
		c.Audit.LogFile = v
	}
	if v := os.Getenv("TMUX_AGENT_AUDIT_ENABLED"); v != "" {
		c.Audit.Enabled = v == "true" || v == "1" || v == "yes"
	}

	// Security
	if v := os.Getenv("TMUX_AGENT_API_KEYS"); v != "" {
		c.Security.APIKeys = strings.Split(v, ",")
		c.Security.EnableAuth = true
	}

	// Plan
	if v := os.Getenv("TMUX_AGENT_PLAN_DIR"); v != "" {
		c.Plan.PlanDirectory = v
	}
}

// expandPaths expands ~ and environment variables in paths.
func (c *Config) expandPaths() {
	c.File.RootDirectory = expandPath(c.File.RootDirectory)
	c.File.BackupDirectory = expandPath(c.File.BackupDirectory)
	c.Command.WorkingDirectory = expandPath(c.Command.WorkingDirectory)
	c.Audit.LogFile = expandPath(c.Audit.LogFile)
	c.Plan.PlanDirectory = expandPath(c.Plan.PlanDirectory)
	c.Server.TLSCertFile = expandPath(c.Server.TLSCertFile)
	c.Server.TLSKeyFile = expandPath(c.Server.TLSKeyFile)

	for i, p := range c.File.AllowedPaths {
		c.File.AllowedPaths[i] = expandPath(p)
	}
	for i, p := range c.File.DeniedPaths {
		c.File.DeniedPaths[i] = expandPath(p)
	}
	for i, p := range c.Command.AllowedPaths {
		c.Command.AllowedPaths[i] = expandPath(p)
	}
}

// expandPath expands ~ and environment variables in a path.
func expandPath(path string) string {
	if path == "" {
		return path
	}

	// Expand ~
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err == nil {
			path = filepath.Join(home, path[2:])
		}
	} else if path == "~" {
		home, err := os.UserHomeDir()
		if err == nil {
			path = home
		}
	}

	// Expand environment variables
	path = os.ExpandEnv(path)

	return path
}

// Validate validates the configuration.
func (c *Config) Validate() error {
	var errs []string

	// Server validation
	if c.Server.Port < 1 || c.Server.Port > 65535 {
		errs = append(errs, "server.port must be between 1 and 65535")
	}
	if c.Server.EnableTLS {
		if c.Server.TLSCertFile == "" {
			errs = append(errs, "server.tls_cert_file is required when TLS is enabled")
		}
		if c.Server.TLSKeyFile == "" {
			errs = append(errs, "server.tls_key_file is required when TLS is enabled")
		}
	}

	// File validation
	if c.File.MaxFileSize <= 0 {
		errs = append(errs, "file.max_file_size must be positive")
	}
	if c.File.MaxReadLines <= 0 {
		errs = append(errs, "file.max_read_lines must be positive")
	}

	// Command validation
	if c.Command.DefaultTimeout <= 0 {
		errs = append(errs, "command.default_timeout must be positive")
	}
	if c.Command.MaxTimeout <= 0 {
		errs = append(errs, "command.max_timeout must be positive")
	}
	if c.Command.MaxTimeout < c.Command.DefaultTimeout {
		errs = append(errs, "command.max_timeout must be >= command.default_timeout")
	}
	if c.Command.MaxOutputSize <= 0 {
		errs = append(errs, "command.max_output_size must be positive")
	}

	// Tmux validation
	if c.Tmux.MaxPaneReadLines <= 0 {
		errs = append(errs, "tmux.max_pane_read_lines must be positive")
	}

	// Plan validation
	if c.Plan.MaxPlans <= 0 {
		errs = append(errs, "plan.max_plans must be positive")
	}
	if c.Plan.MaxStepsPerPlan <= 0 {
		errs = append(errs, "plan.max_steps_per_plan must be positive")
	}

	// Human validation
	if c.Human.MaxPending <= 0 {
		errs = append(errs, "human.max_pending must be positive")
	}

	// Security validation
	if c.Security.EnableAuth && len(c.Security.APIKeys) == 0 {
		errs = append(errs, "security.api_keys is required when auth is enabled")
	}
	if c.Security.RateLimitPerMin < 0 {
		errs = append(errs, "security.rate_limit_per_min must be non-negative")
	}

	if len(errs) > 0 {
		return fmt.Errorf("validation errors:\n  - %s", strings.Join(errs, "\n  - "))
	}

	return nil
}

// Save saves the configuration to a YAML file.
func (c *Config) Save(path string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// IsCommandAllowed checks if a command is allowed to run.
func (c *Config) IsCommandAllowed(cmd string) bool {
	// Check denied list first
	for _, denied := range c.Command.DeniedCommands {
		if cmd == denied {
			return false
		}
	}

	// If allowed list is empty, allow all non-denied commands
	if len(c.Command.AllowedCommands) == 0 {
		return true
	}

	// Check allowed list
	for _, allowed := range c.Command.AllowedCommands {
		if cmd == allowed {
			return true
		}
	}

	return false
}

// IsPathAllowed checks if a file path is allowed.
func (c *Config) IsPathAllowed(path string) bool {
	// Clean and make path absolute
	absPath, err := filepath.Abs(path)
	if err != nil {
		return false
	}

	// Check denied paths first
	for _, denied := range c.File.DeniedPaths {
		deniedAbs, err := filepath.Abs(denied)
		if err != nil {
			continue
		}
		if strings.HasPrefix(absPath, deniedAbs) {
			return false
		}
	}

	// Check denied extensions
	ext := strings.ToLower(filepath.Ext(absPath))
	for _, deniedExt := range c.File.DeniedExtensions {
		if ext == strings.ToLower(deniedExt) {
			return false
		}
	}

	// If root directory is set, path must be under it
	if c.File.RootDirectory != "" {
		rootAbs, err := filepath.Abs(c.File.RootDirectory)
		if err != nil {
			return false
		}
		if !strings.HasPrefix(absPath, rootAbs) {
			return false
		}
	}

	// If allowed paths are empty, allow all non-denied paths
	if len(c.File.AllowedPaths) == 0 {
		return true
	}

	// Check allowed paths
	for _, allowed := range c.File.AllowedPaths {
		allowedAbs, err := filepath.Abs(allowed)
		if err != nil {
			continue
		}
		if strings.HasPrefix(absPath, allowedAbs) {
			return true
		}
	}

	return false
}

// IsKeyAllowed checks if a tmux key is in the allowed list.
func (c *Config) IsKeyAllowed(key string) bool {
	for _, allowed := range c.Tmux.AllowedKeys {
		if key == allowed {
			return true
		}
	}
	return false
}

// RequiresApproval checks if an action type requires human approval.
func (c *Config) RequiresApproval(actionType string) bool {
	// Check auto-approve first
	for _, autoApprove := range c.Human.AutoApproveTypes {
		if actionType == autoApprove {
			return false
		}
	}

	// Check require approval list
	for _, require := range c.Human.RequireApproval {
		if actionType == require {
			return true
		}
	}

	return false
}

// Address returns the server address in host:port format.
func (c *ServerConfig) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
