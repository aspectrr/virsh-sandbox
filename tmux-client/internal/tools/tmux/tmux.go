// Package tmux provides safe tmux interaction capabilities for the agent API.
package tmux

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"tmux-client/internal/config"
	"tmux-client/internal/tools/sandbox"
	"tmux-client/internal/types"
)

// Tool provides tmux operations with safety constraints.
type Tool struct {
	config      *config.TmuxConfig
	tmuxPath    string
	socketPath  string
	sandboxTool *sandbox.Tool
}

// NewTool creates a new tmux tool.
func NewTool(cfg *config.TmuxConfig) (*Tool, error) {
	// Find tmux binary
	tmuxPath, err := exec.LookPath("tmux")
	if err != nil {
		return nil, fmt.Errorf("tmux not found in PATH: %w", err)
	}

	tool := &Tool{
		config:     cfg,
		tmuxPath:   tmuxPath,
		socketPath: cfg.SocketPath,
	}

	return tool, nil
}

// SetSandboxTool sets the sandbox tool for SSH certificate-based access.
func (t *Tool) SetSandboxTool(st *sandbox.Tool) {
	t.sandboxTool = st
}

// CheckSandboxAPIHealth checks if the virsh-sandbox API is reachable.
// Returns nil if healthy, error otherwise. Returns nil if sandbox tool is not configured.
func (t *Tool) CheckSandboxAPIHealth(ctx context.Context) error {
	if t.sandboxTool == nil {
		return nil // Not configured, skip check
	}
	return t.sandboxTool.CheckAPIHealth(ctx)
}

// IsSandboxEnabled returns true if the sandbox tool is configured.
func (t *Tool) IsSandboxEnabled() bool {
	return t.sandboxTool != nil
}

// execTmux executes a tmux command and returns its output.
func (t *Tool) execTmux(ctx context.Context, args ...string) (string, string, error) {
	// Add socket path if configured
	if t.socketPath != "" {
		args = append([]string{"-S", t.socketPath}, args...)
	}

	cmd := exec.CommandContext(ctx, t.tmuxPath, args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	return stdout.String(), stderr.String(), err
}

// IsTmuxRunning checks if tmux server is running.
func (t *Tool) IsTmuxRunning(ctx context.Context) bool {
	_, _, err := t.execTmux(ctx, "list-sessions")
	return err == nil
}

// ListSessions returns all tmux sessions.
func (t *Tool) ListSessions(ctx context.Context) ([]types.SessionInfo, error) {
	format := "#{session_name}\t#{session_id}\t#{session_windows}\t#{session_created}\t#{session_attached}"
	stdout, stderr, err := t.execTmux(ctx, "list-sessions", "-F", format)
	if err != nil {
		if strings.Contains(stderr, "no server running") {
			return []types.SessionInfo{}, nil
		}
		return nil, fmt.Errorf("failed to list sessions: %s", stderr)
	}

	var sessions []types.SessionInfo
	scanner := bufio.NewScanner(strings.NewReader(stdout))
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "\t")
		if len(parts) < 5 {
			continue
		}

		windows, _ := strconv.Atoi(parts[2])
		createdUnix, _ := strconv.ParseInt(parts[3], 10, 64)
		attached := parts[4] == "1"

		sessions = append(sessions, types.SessionInfo{
			Name:     parts[0],
			ID:       parts[1],
			Windows:  windows,
			Created:  time.Unix(createdUnix, 0),
			Attached: attached,
		})
	}

	return sessions, nil
}

// ListWindows returns all windows in a session.
func (t *Tool) ListWindows(ctx context.Context, sessionName string) ([]types.WindowInfo, error) {
	args := []string{"list-windows", "-F", "#{session_name}\t#{window_index}\t#{window_name}\t#{window_active}\t#{window_panes}\t#{window_width}\t#{window_height}"}
	if sessionName != "" {
		args = append(args, "-t", sessionName)
	} else {
		args = append(args, "-a")
	}

	stdout, stderr, err := t.execTmux(ctx, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list windows: %s", stderr)
	}

	var windows []types.WindowInfo
	scanner := bufio.NewScanner(strings.NewReader(stdout))
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "\t")
		if len(parts) < 7 {
			continue
		}

		index, _ := strconv.Atoi(parts[1])
		active := parts[3] == "1"
		panes, _ := strconv.Atoi(parts[4])
		width, _ := strconv.Atoi(parts[5])
		height, _ := strconv.Atoi(parts[6])

		windows = append(windows, types.WindowInfo{
			SessionName: parts[0],
			Index:       index,
			Name:        parts[2],
			Active:      active,
			Panes:       panes,
			Width:       width,
			Height:      height,
		})
	}

	return windows, nil
}

// ListPanes returns all panes, optionally filtered by session.
func (t *Tool) ListPanes(ctx context.Context, sessionName string) ([]types.PaneInfo, error) {
	format := "#{session_name}\t#{window_index}\t#{window_name}\t#{pane_index}\t#{pane_id}\t#{pane_pid}\t#{pane_title}\t#{pane_width}\t#{pane_height}\t#{pane_active}\t#{pane_current_path}"

	args := []string{"list-panes", "-F", format}
	if sessionName != "" {
		args = append(args, "-t", sessionName, "-s")
	} else {
		args = append(args, "-a")
	}

	stdout, stderr, err := t.execTmux(ctx, args...)
	if err != nil {
		if strings.Contains(stderr, "no server running") {
			return []types.PaneInfo{}, nil
		}
		return nil, fmt.Errorf("failed to list panes: %s", stderr)
	}

	var panes []types.PaneInfo
	scanner := bufio.NewScanner(strings.NewReader(stdout))
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "\t")
		if len(parts) < 11 {
			continue
		}

		windowIndex, _ := strconv.Atoi(parts[1])
		paneIndex, _ := strconv.Atoi(parts[3])
		panePID, _ := strconv.Atoi(parts[5])
		paneWidth, _ := strconv.Atoi(parts[7])
		paneHeight, _ := strconv.Atoi(parts[8])
		active := parts[9] == "1"

		panes = append(panes, types.PaneInfo{
			SessionName: parts[0],
			WindowIndex: windowIndex,
			WindowName:  parts[2],
			PaneIndex:   paneIndex,
			PaneID:      parts[4],
			PanePID:     panePID,
			PaneTitle:   parts[6],
			PaneWidth:   paneWidth,
			PaneHeight:  paneHeight,
			Active:      active,
			CurrentPath: parts[10],
		})
	}

	return panes, nil
}

// ReadPane reads the content of a pane.
func (t *Tool) ReadPane(ctx context.Context, paneID string, lastNLines int) (string, int, error) {
	// Validate pane ID format (e.g., %0, %1, etc.)
	if !isValidPaneID(paneID) {
		return "", 0, fmt.Errorf("invalid pane ID format: %s", paneID)
	}

	// Limit the number of lines that can be read
	if lastNLines <= 0 || lastNLines > t.config.MaxPaneReadLines {
		lastNLines = t.config.MaxPaneReadLines
	}

	// Capture pane content
	// -p prints to stdout, -t specifies target pane
	// We use capture-pane with history
	args := []string{"capture-pane", "-t", paneID, "-p", "-S", fmt.Sprintf("-%d", lastNLines)}

	stdout, stderr, err := t.execTmux(ctx, args...)
	if err != nil {
		return "", 0, fmt.Errorf("failed to read pane: %s", stderr)
	}

	// Count lines
	lines := strings.Count(stdout, "\n")
	if len(stdout) > 0 && !strings.HasSuffix(stdout, "\n") {
		lines++
	}

	// Trim trailing empty lines but preserve content
	content := strings.TrimRight(stdout, "\n")

	return content, lines, nil
}

// ReadPaneWithHistory reads the pane content including scrollback history.
func (t *Tool) ReadPaneWithHistory(ctx context.Context, paneID string, startLine, endLine int) (string, error) {
	if !isValidPaneID(paneID) {
		return "", fmt.Errorf("invalid pane ID format: %s", paneID)
	}

	args := []string{"capture-pane", "-t", paneID, "-p"}

	if startLine != 0 {
		args = append(args, "-S", strconv.Itoa(startLine))
	}
	if endLine != 0 {
		args = append(args, "-E", strconv.Itoa(endLine))
	}

	stdout, stderr, err := t.execTmux(ctx, args...)
	if err != nil {
		return "", fmt.Errorf("failed to read pane history: %s", stderr)
	}

	return strings.TrimRight(stdout, "\n"), nil
}

// SwitchPane switches focus to a specific pane.
func (t *Tool) SwitchPane(ctx context.Context, paneID string) error {
	if !isValidPaneID(paneID) {
		return fmt.Errorf("invalid pane ID format: %s", paneID)
	}

	_, stderr, err := t.execTmux(ctx, "select-pane", "-t", paneID)
	if err != nil {
		return fmt.Errorf("failed to switch pane: %s", stderr)
	}

	return nil
}

// SwitchWindow switches focus to a specific window.
func (t *Tool) SwitchWindow(ctx context.Context, target string) error {
	if target == "" {
		return fmt.Errorf("target window is required")
	}

	_, stderr, err := t.execTmux(ctx, "select-window", "-t", target)
	if err != nil {
		return fmt.Errorf("failed to switch window: %s", stderr)
	}

	return nil
}

// CreatePane creates a new pane by splitting an existing one.
func (t *Tool) CreatePane(ctx context.Context, req types.CreatePaneRequest) (*types.CreatePaneResponse, error) {
	var args []string

	if req.NewWindow {
		// Create a new window instead of splitting
		args = []string{"new-window"}

		if req.SessionName != "" {
			args = append(args, "-t", req.SessionName+":")
		}

		if req.WindowName != "" {
			args = append(args, "-n", req.WindowName)
		}

		// Print the pane ID of the new window's pane
		args = append(args, "-P", "-F", "#{pane_id}\t#{session_name}\t#{window_index}\t#{pane_index}")
	} else {
		// Split the current pane
		args = []string{"split-window"}

		if req.SessionName != "" {
			args = append(args, "-t", req.SessionName+":")
		}

		if req.Horizontal {
			args = append(args, "-h")
		} else {
			args = append(args, "-v")
		}

		// Print the pane ID
		args = append(args, "-P", "-F", "#{pane_id}\t#{session_name}\t#{window_index}\t#{pane_index}")
	}

	// Add command if specified
	if req.Command != "" {
		// Validate the command doesn't contain dangerous patterns
		if containsDangerousPatterns(req.Command) {
			return nil, fmt.Errorf("command contains disallowed patterns")
		}
		args = append(args, req.Command)
	}

	stdout, stderr, err := t.execTmux(ctx, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to create pane: %s", stderr)
	}

	// Parse the output
	stdout = strings.TrimSpace(stdout)
	parts := strings.Split(stdout, "\t")
	if len(parts) < 4 {
		return nil, fmt.Errorf("unexpected output format from tmux")
	}

	windowIndex, _ := strconv.Atoi(parts[2])
	paneIndex, _ := strconv.Atoi(parts[3])

	return &types.CreatePaneResponse{
		PaneID:      parts[0],
		SessionName: parts[1],
		WindowIndex: windowIndex,
		PaneIndex:   paneIndex,
	}, nil
}

// CreateSession creates a new tmux session.
func (t *Tool) CreateSession(ctx context.Context, name string, command string) (string, error) {
	if name == "" {
		return "", fmt.Errorf("session name is required")
	}

	// Validate session name (alphanumeric, underscore, hyphen only)
	if !isValidSessionName(name) {
		return "", fmt.Errorf("invalid session name: must be alphanumeric with underscores or hyphens")
	}

	args := []string{"new-session", "-d", "-s", name, "-P", "-F", "#{session_id}"}

	if command != "" {
		if containsDangerousPatterns(command) {
			return "", fmt.Errorf("command contains disallowed patterns")
		}
		args = append(args, command)
	}

	stdout, stderr, err := t.execTmux(ctx, args...)
	if err != nil {
		return "", fmt.Errorf("failed to create session: %s", stderr)
	}

	return strings.TrimSpace(stdout), nil
}

// SandboxSessionInfo contains information about a sandbox-connected session.
type SandboxSessionInfo struct {
	// SessionID is the tmux session ID
	SessionID string

	// SessionName is the tmux session name
	SessionName string

	// SandboxID is the sandbox being accessed
	SandboxID string

	// VMIPAddress is the IP of the sandbox VM
	VMIPAddress string

	// Username is the SSH username
	Username string

	// ValidUntil is when the certificate expires
	ValidUntil time.Time

	// ConnectionInfo contains the full connection details
	ConnectionInfo *sandbox.ConnectionInfo
}

// CreateSandboxSession creates a new tmux session that SSHs into a sandbox VM.
// It requests an SSH certificate from the virsh-sandbox API and creates a session
// that automatically connects to the sandbox.
func (t *Tool) CreateSandboxSession(ctx context.Context, sandboxID string, sessionName string, ttlMinutes int) (*SandboxSessionInfo, error) {
	if t.sandboxTool == nil {
		return nil, fmt.Errorf("sandbox tool not configured")
	}

	if sandboxID == "" {
		return nil, fmt.Errorf("sandbox_id is required")
	}

	// Generate session name if not provided
	if sessionName == "" {
		sessionName = fmt.Sprintf("sandbox_%s", sandboxID)
	}

	// Validate session name
	if !isValidSessionName(sessionName) {
		return nil, fmt.Errorf("invalid session name: must be alphanumeric with underscores or hyphens")
	}

	// Request access to the sandbox (generates keys and gets certificate)
	connInfo, err := t.sandboxTool.RequestAccess(ctx, sandboxID, ttlMinutes)
	if err != nil {
		return nil, fmt.Errorf("request sandbox access: %w", err)
	}

	// Build SSH command
	sshCommand := connInfo.SSHCommand

	// Create tmux session with SSH command
	args := []string{"new-session", "-d", "-s", sessionName, "-P", "-F", "#{session_id}", sshCommand}

	stdout, stderr, err := t.execTmux(ctx, args...)
	if err != nil {
		// Clean up keys on failure
		connInfo.Cleanup()
		return nil, fmt.Errorf("failed to create sandbox session: %s", stderr)
	}

	sessionID := strings.TrimSpace(stdout)

	return &SandboxSessionInfo{
		SessionID:      sessionID,
		SessionName:    sessionName,
		SandboxID:      sandboxID,
		VMIPAddress:    connInfo.VMIPAddress,
		Username:       connInfo.Username,
		ValidUntil:     connInfo.ValidUntil,
		ConnectionInfo: connInfo,
	}, nil
}

// KillSandboxSession kills a sandbox session and cleans up its credentials.
func (t *Tool) KillSandboxSession(ctx context.Context, sessionName string, connInfo *sandbox.ConnectionInfo) error {
	// Kill the tmux session
	if err := t.KillSession(ctx, sessionName); err != nil {
		// Log but continue with cleanup
		_ = err
	}

	// Clean up credentials
	if connInfo != nil {
		if err := connInfo.Cleanup(); err != nil {
			return fmt.Errorf("cleanup credentials: %w", err)
		}
	}

	return nil
}

// KillPane kills a specific pane.
func (t *Tool) KillPane(ctx context.Context, paneID string) error {
	if !isValidPaneID(paneID) {
		return fmt.Errorf("invalid pane ID format: %s", paneID)
	}

	_, stderr, err := t.execTmux(ctx, "kill-pane", "-t", paneID)
	if err != nil {
		return fmt.Errorf("failed to kill pane: %s", stderr)
	}

	return nil
}

// KillSession kills an entire session.
func (t *Tool) KillSession(ctx context.Context, sessionName string) error {
	if sessionName == "" {
		return fmt.Errorf("session name is required")
	}

	_, stderr, err := t.execTmux(ctx, "kill-session", "-t", sessionName)
	if err != nil {
		return fmt.Errorf("failed to kill session: %s", stderr)
	}

	return nil
}

// SendKeys sends limited, approved keystrokes to a pane.
// Only keys in the allowed list can be sent for safety.
func (t *Tool) SendKeys(ctx context.Context, paneID string, key string) error {
	if !isValidPaneID(paneID) {
		return fmt.Errorf("invalid pane ID format: %s", paneID)
	}

	// Check if the key is allowed
	if !t.isKeyAllowed(key) {
		return fmt.Errorf("key '%s' is not in the allowed list", key)
	}

	_, stderr, err := t.execTmux(ctx, "send-keys", "-t", paneID, key)
	if err != nil {
		return fmt.Errorf("failed to send key: %s", stderr)
	}

	return nil
}

// SendText sends text to a pane followed by Enter (for running commands).
// This is more restrictive than SendKeys and validates the input.
func (t *Tool) SendText(ctx context.Context, paneID string, text string, pressEnter bool) error {
	if !isValidPaneID(paneID) {
		return fmt.Errorf("invalid pane ID format: %s", paneID)
	}

	// Validate text length
	if len(text) > 4096 {
		return fmt.Errorf("text too long: maximum 4096 characters")
	}

	// Check for dangerous patterns
	if containsDangerousPatterns(text) {
		return fmt.Errorf("text contains disallowed patterns")
	}

	args := []string{"send-keys", "-t", paneID, text}
	if pressEnter {
		args = append(args, "Enter")
	}

	_, stderr, err := t.execTmux(ctx, args...)
	if err != nil {
		return fmt.Errorf("failed to send text: %s", stderr)
	}

	return nil
}

// GetPaneInfo gets detailed information about a specific pane.
func (t *Tool) GetPaneInfo(ctx context.Context, paneID string) (*types.PaneInfo, error) {
	if !isValidPaneID(paneID) {
		return nil, fmt.Errorf("invalid pane ID format: %s", paneID)
	}

	format := "#{session_name}\t#{window_index}\t#{window_name}\t#{pane_index}\t#{pane_id}\t#{pane_pid}\t#{pane_title}\t#{pane_width}\t#{pane_height}\t#{pane_active}\t#{pane_current_path}"

	stdout, stderr, err := t.execTmux(ctx, "display-message", "-t", paneID, "-p", format)
	if err != nil {
		return nil, fmt.Errorf("failed to get pane info: %s", stderr)
	}

	stdout = strings.TrimSpace(stdout)
	parts := strings.Split(stdout, "\t")
	if len(parts) < 11 {
		return nil, fmt.Errorf("unexpected output format from tmux")
	}

	windowIndex, _ := strconv.Atoi(parts[1])
	paneIndex, _ := strconv.Atoi(parts[3])
	panePID, _ := strconv.Atoi(parts[5])
	paneWidth, _ := strconv.Atoi(parts[7])
	paneHeight, _ := strconv.Atoi(parts[8])
	active := parts[9] == "1"

	return &types.PaneInfo{
		SessionName: parts[0],
		WindowIndex: windowIndex,
		WindowName:  parts[2],
		PaneIndex:   paneIndex,
		PaneID:      parts[4],
		PanePID:     panePID,
		PaneTitle:   parts[6],
		PaneWidth:   paneWidth,
		PaneHeight:  paneHeight,
		Active:      active,
		CurrentPath: parts[10],
	}, nil
}

// RenameSession renames a tmux session.
func (t *Tool) RenameSession(ctx context.Context, oldName, newName string) error {
	if !isValidSessionName(newName) {
		return fmt.Errorf("invalid session name: must be alphanumeric with underscores or hyphens")
	}

	_, stderr, err := t.execTmux(ctx, "rename-session", "-t", oldName, newName)
	if err != nil {
		return fmt.Errorf("failed to rename session: %s", stderr)
	}

	return nil
}

// RenameWindow renames a tmux window.
func (t *Tool) RenameWindow(ctx context.Context, target, newName string) error {
	if newName == "" {
		return fmt.Errorf("new window name is required")
	}

	_, stderr, err := t.execTmux(ctx, "rename-window", "-t", target, newName)
	if err != nil {
		return fmt.Errorf("failed to rename window: %s", stderr)
	}

	return nil
}

// ResizePane resizes a pane.
func (t *Tool) ResizePane(ctx context.Context, paneID string, width, height int) error {
	if !isValidPaneID(paneID) {
		return fmt.Errorf("invalid pane ID format: %s", paneID)
	}

	if width > 0 {
		_, stderr, err := t.execTmux(ctx, "resize-pane", "-t", paneID, "-x", strconv.Itoa(width))
		if err != nil {
			return fmt.Errorf("failed to resize pane width: %s", stderr)
		}
	}

	if height > 0 {
		_, stderr, err := t.execTmux(ctx, "resize-pane", "-t", paneID, "-y", strconv.Itoa(height))
		if err != nil {
			return fmt.Errorf("failed to resize pane height: %s", stderr)
		}
	}

	return nil
}

// WaitForOutput waits for specific output to appear in a pane.
func (t *Tool) WaitForOutput(ctx context.Context, paneID string, pattern string, timeout time.Duration) (bool, string, error) {
	if !isValidPaneID(paneID) {
		return false, "", fmt.Errorf("invalid pane ID format: %s", paneID)
	}

	deadline := time.Now().Add(timeout)
	pollInterval := 100 * time.Millisecond

	for time.Now().Before(deadline) {
		content, _, err := t.ReadPane(ctx, paneID, 100)
		if err != nil {
			return false, "", err
		}

		if strings.Contains(content, pattern) {
			return true, content, nil
		}

		select {
		case <-ctx.Done():
			return false, "", ctx.Err()
		case <-time.After(pollInterval):
			// Continue polling
		}
	}

	return false, "", nil
}

// Validation helpers

// isValidPaneID validates a tmux pane ID format.
func isValidPaneID(paneID string) bool {
	if paneID == "" {
		return false
	}

	// Pane IDs start with % followed by digits
	if strings.HasPrefix(paneID, "%") {
		_, err := strconv.Atoi(paneID[1:])
		return err == nil
	}

	// Also allow session:window.pane format
	// e.g., "mysession:0.1" or "0:0.0"
	if strings.Contains(paneID, ":") || strings.Contains(paneID, ".") {
		// Basic validation - contains only allowed characters
		for _, c := range paneID {
			if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') ||
				(c >= '0' && c <= '9') || c == ':' || c == '.' || c == '-' || c == '_') {
				return false
			}
		}
		return true
	}

	return false
}

// isValidSessionName validates a tmux session name.
func isValidSessionName(name string) bool {
	if name == "" || len(name) > 256 {
		return false
	}

	for _, c := range name {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') ||
			(c >= '0' && c <= '9') || c == '-' || c == '_') {
			return false
		}
	}

	return true
}

// containsDangerousPatterns checks for potentially dangerous shell patterns.
func containsDangerousPatterns(s string) bool {
	dangerous := []string{
		"$(",   // Command substitution
		"`",    // Backtick command substitution
		"&&",   // Command chaining
		"||",   // Command chaining
		";",    // Command separator
		"|",    // Pipe (could be used for command chaining)
		">",    // Redirect
		"<",    // Redirect
		"eval", // eval command
		"exec", // exec command
	}

	for _, pattern := range dangerous {
		if strings.Contains(s, pattern) {
			return true
		}
	}

	return false
}

// isKeyAllowed checks if a key is in the allowed list.
func (t *Tool) isKeyAllowed(key string) bool {
	for _, allowed := range t.config.AllowedKeys {
		if key == allowed {
			return true
		}
	}
	return false
}
