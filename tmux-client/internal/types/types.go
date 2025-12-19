// Package types defines shared types and response structures for the tmux agent API.
package types

import (
	"encoding/json"
	"time"
)

// APIResponse wraps all API responses with a consistent structure.
type APIResponse struct {
	Success   bool            `json:"success"`
	Data      json.RawMessage `json:"data,omitempty"`
	Error     *APIError       `json:"error,omitempty"`
	Timestamp time.Time       `json:"timestamp"`
	RequestID string          `json:"request_id"`
}

// APIError represents a structured error response.
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// Common error codes
const (
	ErrCodeValidation     = "VALIDATION_ERROR"
	ErrCodeNotFound       = "NOT_FOUND"
	ErrCodeForbidden      = "FORBIDDEN"
	ErrCodeInternal       = "INTERNAL_ERROR"
	ErrCodeTimeout        = "TIMEOUT"
	ErrCodeDenied         = "DENIED"
	ErrCodePathViolation  = "PATH_VIOLATION"
	ErrCodeCommandBlocked = "COMMAND_BLOCKED"
)

// ============================================================================
// TMUX Types
// ============================================================================

// PaneInfo represents information about a tmux pane.
type PaneInfo struct {
	SessionName string `json:"session_name"`
	WindowIndex int    `json:"window_index"`
	WindowName  string `json:"window_name"`
	PaneIndex   int    `json:"pane_index"`
	PaneID      string `json:"pane_id"`
	PanePID     int    `json:"pane_pid"`
	PaneTitle   string `json:"pane_title"`
	PaneWidth   int    `json:"pane_width"`
	PaneHeight  int    `json:"pane_height"`
	Active      bool   `json:"active"`
	CurrentPath string `json:"current_path"`
}

// SessionInfo represents information about a tmux session.
type SessionInfo struct {
	Name      string    `json:"name"`
	ID        string    `json:"id"`
	Windows   int       `json:"windows"`
	Created   time.Time `json:"created"`
	Attached  bool      `json:"attached"`
	LastPaneX int       `json:"last_pane_x,omitempty"`
	LastPaneY int       `json:"last_pane_y,omitempty"`
}

// WindowInfo represents information about a tmux window.
type WindowInfo struct {
	SessionName string `json:"session_name"`
	Index       int    `json:"index"`
	Name        string `json:"name"`
	Active      bool   `json:"active"`
	Panes       int    `json:"panes"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
}

// ListPanesRequest is the request for listing panes.
type ListPanesRequest struct {
	SessionName string `json:"session_name,omitempty"`
}

// ListPanesResponse is the response for listing panes.
type ListPanesResponse struct {
	Panes []PaneInfo `json:"panes"`
}

// ReadPaneRequest is the request for reading pane content.
type ReadPaneRequest struct {
	PaneID     string `json:"pane_id"`
	LastNLines int    `json:"last_n_lines,omitempty"` // 0 means all visible content
}

// ReadPaneResponse is the response for reading pane content.
type ReadPaneResponse struct {
	PaneID  string `json:"pane_id"`
	Content string `json:"content"`
	Lines   int    `json:"lines"`
}

// SwitchPaneRequest is the request for switching focus to a pane.
type SwitchPaneRequest struct {
	PaneID string `json:"pane_id"`
}

// SwitchPaneResponse is the response for switching panes.
type SwitchPaneResponse struct {
	PaneID   string `json:"pane_id"`
	Switched bool   `json:"switched"`
}

// CreatePaneRequest is the request for creating a new pane.
type CreatePaneRequest struct {
	SessionName string `json:"session_name"`
	WindowName  string `json:"window_name,omitempty"`
	Command     string `json:"command,omitempty"`
	Horizontal  bool   `json:"horizontal,omitempty"` // false = vertical split
	NewWindow   bool   `json:"new_window,omitempty"` // true = create new window instead of split
}

// CreatePaneResponse is the response for creating a pane.
type CreatePaneResponse struct {
	PaneID      string `json:"pane_id"`
	SessionName string `json:"session_name"`
	WindowIndex int    `json:"window_index"`
	PaneIndex   int    `json:"pane_index"`
}

// SendKeysRequest allows sending limited, approved keystrokes to a pane.
// This is intentionally restrictive - only special keys like Enter, Ctrl+C are allowed.
type SendKeysRequest struct {
	PaneID string `json:"pane_id"`
	Key    string `json:"key"` // Must be from approved list: "Enter", "C-c", "C-d", "Escape"
}

// SendKeysResponse is the response for sending keys.
type SendKeysResponse struct {
	PaneID string `json:"pane_id"`
	Sent   bool   `json:"sent"`
}

// ============================================================================
// File Types
// ============================================================================

// ReadFileRequest is the request for reading a file.
type ReadFileRequest struct {
	Path     string `json:"path"`
	MaxLines int    `json:"max_lines,omitempty"` // 0 = no limit
	FromLine int    `json:"from_line,omitempty"` // 1-indexed, 0 = start
	ToLine   int    `json:"to_line,omitempty"`   // 1-indexed, 0 = end
}

// ReadFileResponse is the response for reading a file.
type ReadFileResponse struct {
	Path       string `json:"path"`
	Content    string `json:"content"`
	TotalLines int    `json:"total_lines"`
	FromLine   int    `json:"from_line"`
	ToLine     int    `json:"to_line"`
	Truncated  bool   `json:"truncated"`
	Size       int64  `json:"size"`
	Mode       string `json:"mode"`
	ModTime    string `json:"mod_time"`
}

// WriteFileRequest is the request for writing a file.
type WriteFileRequest struct {
	Path      string `json:"path"`
	Content   string `json:"content"`
	Mode      string `json:"mode,omitempty"` // e.g., "0644"
	Overwrite bool   `json:"overwrite"`      // Must be true to overwrite existing
	CreateDir bool   `json:"create_dir"`     // Create parent directories if needed
}

// WriteFileResponse is the response for writing a file.
type WriteFileResponse struct {
	Path         string `json:"path"`
	Written      bool   `json:"written"`
	BytesWritten int64  `json:"bytes_written"`
	Created      bool   `json:"created"` // true if file was created, false if overwritten
}

// EditFileRequest is the request for editing a file using patch semantics.
type EditFileRequest struct {
	Path    string `json:"path"`
	OldText string `json:"old_text"` // Text to find and replace
	NewText string `json:"new_text"` // Replacement text
	All     bool   `json:"all"`      // Replace all occurrences (default: first only)
}

// EditFileResponse is the response for editing a file.
type EditFileResponse struct {
	Path          string `json:"path"`
	Edited        bool   `json:"edited"`
	Replacements  int    `json:"replacements"`
	Diff          string `json:"diff"`           // Unified diff format
	ContentBefore string `json:"content_before"` // For audit trail
	ContentAfter  string `json:"content_after"`  // For audit trail
}

// CopyFileRequest is the request for copying a file.
type CopyFileRequest struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Overwrite   bool   `json:"overwrite"`
}

// CopyFileResponse is the response for copying a file.
type CopyFileResponse struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Copied      bool   `json:"copied"`
	BytesCopied int64  `json:"bytes_copied"`
}

// DeleteFileRequest is the request for deleting a file.
type DeleteFileRequest struct {
	Path      string `json:"path"`
	Recursive bool   `json:"recursive"` // For directories
}

// DeleteFileResponse is the response for deleting a file.
type DeleteFileResponse struct {
	Path    string `json:"path"`
	Deleted bool   `json:"deleted"`
	WasDir  bool   `json:"was_dir"`
}

// ListDirRequest is the request for listing directory contents.
type ListDirRequest struct {
	Path      string `json:"path"`
	Recursive bool   `json:"recursive,omitempty"`
	MaxDepth  int    `json:"max_depth,omitempty"`
}

// FileInfo represents basic file information.
type FileInfo struct {
	Name    string    `json:"name"`
	Path    string    `json:"path"`
	Size    int64     `json:"size"`
	Mode    string    `json:"mode"`
	ModTime time.Time `json:"mod_time"`
	IsDir   bool      `json:"is_dir"`
}

// ListDirResponse is the response for listing directory contents.
type ListDirResponse struct {
	Path  string     `json:"path"`
	Files []FileInfo `json:"files"`
}

// ============================================================================
// Command Types
// ============================================================================

// RunCommandRequest is the request for running a command.
type RunCommandRequest struct {
	Command string   `json:"command"`            // Executable name only
	Args    []string `json:"args,omitempty"`     // Arguments as separate items
	WorkDir string   `json:"work_dir,omitempty"` // Working directory
	Timeout int      `json:"timeout,omitempty"`  // Seconds, 0 = default (30s)
	DryRun  bool     `json:"dry_run,omitempty"`  // If true, don't actually run
	Env     []string `json:"env,omitempty"`      // Additional env vars (KEY=VALUE)
}

// RunCommandResponse is the response for running a command.
type RunCommandResponse struct {
	Command    string   `json:"command"`
	Args       []string `json:"args"`
	Stdout     string   `json:"stdout"`
	Stderr     string   `json:"stderr"`
	ExitCode   int      `json:"exit_code"`
	DryRun     bool     `json:"dry_run"`
	DurationMs int64    `json:"duration_ms"`
	TimedOut   bool     `json:"timed_out"`
}

// ============================================================================
// Human Approval Types
// ============================================================================

// ApprovalStatus represents the status of a human approval request.
type ApprovalStatus string

const (
	ApprovalPending  ApprovalStatus = "pending"
	ApprovalApproved ApprovalStatus = "approved"
	ApprovalRejected ApprovalStatus = "rejected"
	ApprovalExpired  ApprovalStatus = "expired"
)

// AskHumanRequest is the request for asking human approval.
type AskHumanRequest struct {
	Prompt       string   `json:"prompt"`                 // Human-readable description
	Context      string   `json:"context,omitempty"`      // Additional context
	ActionType   string   `json:"action_type"`            // Category: "destructive", "sensitive", "irreversible"
	Urgency      string   `json:"urgency,omitempty"`      // "low", "medium", "high"
	TimeoutSecs  int      `json:"timeout_secs,omitempty"` // Auto-reject after timeout, 0 = no timeout
	Alternatives []string `json:"alternatives,omitempty"` // Suggested alternative actions
}

// AskHumanResponse is the response for human approval.
type AskHumanResponse struct {
	RequestID  string         `json:"request_id"`
	Approved   bool           `json:"approved"`
	Status     ApprovalStatus `json:"status"`
	Comment    string         `json:"comment,omitempty"`
	ApprovedBy string         `json:"approved_by,omitempty"`
	ApprovedAt *time.Time     `json:"approved_at,omitempty"`
	ExpiresAt  *time.Time     `json:"expires_at,omitempty"`
}

// PendingApproval represents a pending human approval request.
type PendingApproval struct {
	RequestID  string         `json:"request_id"`
	Prompt     string         `json:"prompt"`
	Context    string         `json:"context,omitempty"`
	ActionType string         `json:"action_type"`
	Urgency    string         `json:"urgency"`
	Status     ApprovalStatus `json:"status"`
	CreatedAt  time.Time      `json:"created_at"`
	ExpiresAt  *time.Time     `json:"expires_at,omitempty"`
}

// ApproveRequest is sent by a human to approve a pending request.
type ApproveRequest struct {
	RequestID  string `json:"request_id"`
	Approved   bool   `json:"approved"`
	Comment    string `json:"comment,omitempty"`
	ApprovedBy string `json:"approved_by"`
}

// ListApprovalsResponse is the response for listing pending approvals.
type ListApprovalsResponse struct {
	Pending []PendingApproval `json:"pending"`
}

// ============================================================================
// Plan Types
// ============================================================================

// PlanStatus represents the status of a plan.
type PlanStatus string

const (
	PlanStatusPending    PlanStatus = "pending"
	PlanStatusInProgress PlanStatus = "in_progress"
	PlanStatusCompleted  PlanStatus = "completed"
	PlanStatusFailed     PlanStatus = "failed"
	PlanStatusAborted    PlanStatus = "aborted"
)

// StepStatus represents the status of a plan step.
type StepStatus string

const (
	StepStatusPending   StepStatus = "pending"
	StepStatusActive    StepStatus = "active"
	StepStatusCompleted StepStatus = "completed"
	StepStatusFailed    StepStatus = "failed"
	StepStatusSkipped   StepStatus = "skipped"
)

// PlanStep represents a single step in a plan.
type PlanStep struct {
	Index       int        `json:"index"`
	Description string     `json:"description"`
	Status      StepStatus `json:"status"`
	Result      string     `json:"result,omitempty"`
	Error       string     `json:"error,omitempty"`
	StartedAt   *time.Time `json:"started_at,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

// Plan represents a multi-step plan.
type Plan struct {
	ID          string     `json:"id"`
	Name        string     `json:"name,omitempty"`
	Description string     `json:"description,omitempty"`
	Steps       []PlanStep `json:"steps"`
	Status      PlanStatus `json:"status"`
	CurrentStep int        `json:"current_step"` // -1 if not started
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

// CreatePlanRequest is the request for creating a plan.
type CreatePlanRequest struct {
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description,omitempty"`
	Steps       []string `json:"steps"` // Step descriptions
}

// CreatePlanResponse is the response for creating a plan.
type CreatePlanResponse struct {
	PlanID string `json:"plan_id"`
	Plan   Plan   `json:"plan"`
}

// UpdatePlanRequest is the request for updating a plan step.
type UpdatePlanRequest struct {
	PlanID    string     `json:"plan_id"`
	StepIndex int        `json:"step_index"`
	Status    StepStatus `json:"status"`
	Result    string     `json:"result,omitempty"`
	Error     string     `json:"error,omitempty"`
}

// UpdatePlanResponse is the response for updating a plan.
type UpdatePlanResponse struct {
	PlanID  string `json:"plan_id"`
	Updated bool   `json:"updated"`
	Plan    Plan   `json:"plan"`
}

// GetPlanRequest is the request for getting a plan.
type GetPlanRequest struct {
	PlanID string `json:"plan_id"`
}

// GetPlanResponse is the response for getting a plan.
type GetPlanResponse struct {
	Plan Plan `json:"plan"`
}

// ListPlansResponse is the response for listing plans.
type ListPlansResponse struct {
	Plans []Plan `json:"plans"`
}

// ============================================================================
// Audit Types
// ============================================================================

// AuditEntry represents a single audit log entry.
type AuditEntry struct {
	Timestamp  time.Time       `json:"timestamp"`
	RequestID  string          `json:"request_id"`
	Tool       string          `json:"tool"`
	Action     string          `json:"action"`
	Arguments  json.RawMessage `json:"arguments"`
	Result     json.RawMessage `json:"result,omitempty"`
	Error      *APIError       `json:"error,omitempty"`
	DurationMs int64           `json:"duration_ms"`
	ClientIP   string          `json:"client_ip,omitempty"`
	UserAgent  string          `json:"user_agent,omitempty"`
}

// AuditQuery represents a query for audit logs.
type AuditQuery struct {
	Tool      string     `json:"tool,omitempty"`
	Action    string     `json:"action,omitempty"`
	Since     *time.Time `json:"since,omitempty"`
	Until     *time.Time `json:"until,omitempty"`
	Limit     int        `json:"limit,omitempty"`
	RequestID string     `json:"request_id,omitempty"`
}

// AuditQueryResponse is the response for querying audit logs.
type AuditQueryResponse struct {
	Entries    []AuditEntry `json:"entries"`
	TotalCount int          `json:"total_count"`
	HasMore    bool         `json:"has_more"`
}

// ============================================================================
// Health Types
// ============================================================================

// HealthStatus represents the health of a component.
type HealthStatus string

const (
	HealthStatusHealthy   HealthStatus = "healthy"
	HealthStatusDegraded  HealthStatus = "degraded"
	HealthStatusUnhealthy HealthStatus = "unhealthy"
)

// ComponentHealth represents the health of a single component.
type ComponentHealth struct {
	Name    string       `json:"name"`
	Status  HealthStatus `json:"status"`
	Message string       `json:"message,omitempty"`
}

// HealthResponse is the response for health checks.
type HealthResponse struct {
	Status     HealthStatus      `json:"status"`
	Version    string            `json:"version"`
	Uptime     string            `json:"uptime"`
	Components []ComponentHealth `json:"components"`
}
