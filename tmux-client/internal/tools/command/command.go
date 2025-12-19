// Package command provides safe command execution for the agent API.
package command

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"tmux-client/internal/config"
	"tmux-client/internal/types"
)

// Tool provides command execution with strict safety constraints.
type Tool struct {
	config *config.CommandConfig
}

// NewTool creates a new command tool.
func NewTool(cfg *config.CommandConfig) (*Tool, error) {
	return &Tool{config: cfg}, nil
}

// dangerousPatterns contains patterns that are not allowed in commands or arguments.
var dangerousPatterns = []string{
	"|",   // Pipe
	">",   // Redirect stdout
	"<",   // Redirect stdin
	">>",  // Append redirect
	"2>",  // Redirect stderr
	"&>",  // Redirect all
	"&&",  // AND chaining
	"||",  // OR chaining
	";",   // Command separator
	"`",   // Backtick command substitution
	"$(",  // Command substitution
	"${",  // Variable expansion (potentially dangerous)
	"\n",  // Newline (could be used for injection)
	"\r",  // Carriage return
	"\\!", // History expansion
}

// reservedCommands are commands that should never be allowed.
var reservedCommands = []string{
	"eval",
	"exec",
	"source",
	".",
	"bash",
	"sh",
	"zsh",
	"fish",
	"csh",
	"tcsh",
	"ksh",
}

// RunCommand executes a command with the given arguments.
func (t *Tool) RunCommand(ctx context.Context, req types.RunCommandRequest) (*types.RunCommandResponse, error) {
	startTime := time.Now()

	// Validate command
	if err := t.validateCommand(req.Command); err != nil {
		return nil, err
	}

	// Validate arguments
	if err := t.validateArguments(req.Args); err != nil {
		return nil, err
	}

	// Validate environment variables if provided
	if len(req.Env) > 0 {
		if err := t.validateEnvVars(req.Env); err != nil {
			return nil, err
		}
	}

	// Validate working directory if provided
	if req.WorkDir != "" {
		if err := t.validateWorkDir(req.WorkDir); err != nil {
			return nil, err
		}
	}

	// Handle dry-run mode
	if req.DryRun {
		return &types.RunCommandResponse{
			Command:    req.Command,
			Args:       req.Args,
			Stdout:     "",
			Stderr:     "",
			ExitCode:   0,
			DryRun:     true,
			DurationMs: 0,
			TimedOut:   false,
		}, nil
	}

	// Determine timeout
	timeout := t.config.DefaultTimeout
	if req.Timeout > 0 {
		requestedTimeout := time.Duration(req.Timeout) * time.Second
		if requestedTimeout > t.config.MaxTimeout {
			return nil, fmt.Errorf("requested timeout %v exceeds maximum allowed %v", requestedTimeout, t.config.MaxTimeout)
		}
		timeout = requestedTimeout
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Find the command executable
	cmdPath, err := t.resolveCommand(req.Command)
	if err != nil {
		return nil, err
	}

	// Create command
	cmd := exec.CommandContext(ctx, cmdPath, req.Args...)

	// Set working directory
	if req.WorkDir != "" {
		cmd.Dir = req.WorkDir
	} else if t.config.WorkingDirectory != "" {
		cmd.Dir = t.config.WorkingDirectory
	}

	// Set environment
	if t.config.InheritEnv {
		cmd.Env = os.Environ()
	} else {
		cmd.Env = []string{}
	}

	// Add allowed environment variables
	if len(req.Env) > 0 && t.config.AllowEnvVars {
		cmd.Env = append(cmd.Env, req.Env...)
	}

	// Capture output
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &limitedWriter{w: &stdout, limit: t.config.MaxOutputSize}
	cmd.Stderr = &limitedWriter{w: &stderr, limit: t.config.MaxOutputSize}

	// Run command
	err = cmd.Run()

	duration := time.Since(startTime)

	// Prepare response
	resp := &types.RunCommandResponse{
		Command:    req.Command,
		Args:       req.Args,
		Stdout:     stdout.String(),
		Stderr:     stderr.String(),
		ExitCode:   0,
		DryRun:     false,
		DurationMs: duration.Milliseconds(),
		TimedOut:   false,
	}

	// Check for timeout
	if ctx.Err() == context.DeadlineExceeded {
		resp.TimedOut = true
		resp.ExitCode = -1
		resp.Stderr = fmt.Sprintf("command timed out after %v\n%s", timeout, resp.Stderr)
		return resp, nil
	}

	// Get exit code
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			resp.ExitCode = exitErr.ExitCode()
		} else {
			return nil, fmt.Errorf("command execution failed: %w", err)
		}
	}

	return resp, nil
}

// validateCommand validates the command name.
func (t *Tool) validateCommand(cmd string) error {
	if cmd == "" {
		return fmt.Errorf("command cannot be empty")
	}

	// Check for path traversal
	if strings.Contains(cmd, "..") {
		return fmt.Errorf("command cannot contain path traversal")
	}

	// Check for dangerous patterns
	for _, pattern := range dangerousPatterns {
		if strings.Contains(cmd, pattern) {
			return fmt.Errorf("command contains disallowed pattern: %q", pattern)
		}
	}

	// Check reserved commands
	baseName := filepath.Base(cmd)
	for _, reserved := range reservedCommands {
		if baseName == reserved {
			return fmt.Errorf("command %q is not allowed", cmd)
		}
	}

	// Check denied list
	for _, denied := range t.config.DeniedCommands {
		if baseName == denied || cmd == denied {
			return fmt.Errorf("command %q is explicitly denied", cmd)
		}
	}

	// Check allowed list (if configured)
	if len(t.config.AllowedCommands) > 0 {
		allowed := false
		for _, allowedCmd := range t.config.AllowedCommands {
			if baseName == allowedCmd || cmd == allowedCmd {
				allowed = true
				break
			}
		}
		if !allowed {
			return fmt.Errorf("command %q is not in the allowed list", cmd)
		}
	}

	return nil
}

// validateArguments validates command arguments.
func (t *Tool) validateArguments(args []string) error {
	for i, arg := range args {
		// Check for dangerous patterns in arguments
		for _, pattern := range dangerousPatterns {
			if strings.Contains(arg, pattern) {
				return fmt.Errorf("argument %d contains disallowed pattern: %q", i, pattern)
			}
		}
	}

	return nil
}

// validateEnvVars validates environment variables.
func (t *Tool) validateEnvVars(envVars []string) error {
	if !t.config.AllowEnvVars {
		return fmt.Errorf("environment variables are not allowed")
	}

	for _, env := range envVars {
		// Must be in KEY=VALUE format
		parts := strings.SplitN(env, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid environment variable format: %q (expected KEY=VALUE)", env)
		}

		key := parts[0]
		value := parts[1]

		// Validate key format
		if key == "" {
			return fmt.Errorf("environment variable key cannot be empty")
		}

		for _, c := range key {
			if !((c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') ||
				(c >= '0' && c <= '9') || c == '_') {
				return fmt.Errorf("invalid character in environment variable key: %q", key)
			}
		}

		// Check against allowed env vars list if configured
		if len(t.config.AllowedEnvVars) > 0 {
			allowed := false
			for _, allowedKey := range t.config.AllowedEnvVars {
				if key == allowedKey {
					allowed = true
					break
				}
			}
			if !allowed {
				return fmt.Errorf("environment variable %q is not in the allowed list", key)
			}
		}

		// Check for dangerous patterns in value
		for _, pattern := range dangerousPatterns {
			if strings.Contains(value, pattern) {
				return fmt.Errorf("environment variable value contains disallowed pattern")
			}
		}
	}

	return nil
}

// validateWorkDir validates the working directory.
func (t *Tool) validateWorkDir(workDir string) error {
	// Make absolute
	absPath, err := filepath.Abs(workDir)
	if err != nil {
		return fmt.Errorf("invalid working directory: %w", err)
	}

	// Check it exists and is a directory
	info, err := os.Stat(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("working directory does not exist: %s", workDir)
		}
		return fmt.Errorf("failed to stat working directory: %w", err)
	}

	if !info.IsDir() {
		return fmt.Errorf("working directory is not a directory: %s", workDir)
	}

	// Check against allowed paths if configured
	if len(t.config.AllowedPaths) > 0 {
		allowed := false
		for _, allowedPath := range t.config.AllowedPaths {
			allowedAbs, err := filepath.Abs(allowedPath)
			if err != nil {
				continue
			}
			if strings.HasPrefix(absPath, allowedAbs) {
				allowed = true
				break
			}
		}
		if !allowed {
			return fmt.Errorf("working directory is not in allowed paths")
		}
	}

	return nil
}

// resolveCommand resolves a command name to its full path.
func (t *Tool) resolveCommand(cmd string) (string, error) {
	// If it's an absolute path, validate and return
	if filepath.IsAbs(cmd) {
		if _, err := os.Stat(cmd); err != nil {
			return "", fmt.Errorf("command not found: %s", cmd)
		}
		return cmd, nil
	}

	// Look up command in PATH
	cmdPath, err := exec.LookPath(cmd)
	if err != nil {
		return "", fmt.Errorf("command not found in PATH: %s", cmd)
	}

	return cmdPath, nil
}

// limitedWriter limits the amount of data written to prevent memory exhaustion.
type limitedWriter struct {
	w       *bytes.Buffer
	limit   int
	written int
}

func (lw *limitedWriter) Write(p []byte) (n int, err error) {
	remaining := lw.limit - lw.written
	if remaining <= 0 {
		return len(p), nil // Silently discard
	}

	if len(p) > remaining {
		p = p[:remaining]
	}

	n, err = lw.w.Write(p)
	lw.written += n
	return len(p), err // Report full write to avoid errors
}

// IsCommandAllowed checks if a command would be allowed to run.
func (t *Tool) IsCommandAllowed(cmd string) bool {
	return t.validateCommand(cmd) == nil
}

// GetAllowedCommands returns the list of allowed commands.
func (t *Tool) GetAllowedCommands() []string {
	return t.config.AllowedCommands
}

// GetDeniedCommands returns the list of denied commands.
func (t *Tool) GetDeniedCommands() []string {
	return t.config.DeniedCommands
}
