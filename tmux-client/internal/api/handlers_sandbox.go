// Package api provides HTTP handlers for the tmux agent API.
package api

import (
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"

	"tmux-client/internal/tools/sandbox"
	"tmux-client/internal/tools/tmux"
	"tmux-client/internal/types"
)

// SandboxSessionTracker tracks active sandbox sessions for cleanup.
type SandboxSessionTracker struct {
	mu       sync.RWMutex
	sessions map[string]*tmux.SandboxSessionInfo
}

// NewSandboxSessionTracker creates a new session tracker.
func NewSandboxSessionTracker() *SandboxSessionTracker {
	return &SandboxSessionTracker{
		sessions: make(map[string]*tmux.SandboxSessionInfo),
	}
}

// Add adds a session to the tracker.
func (t *SandboxSessionTracker) Add(info *tmux.SandboxSessionInfo) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.sessions[info.SessionName] = info
}

// Remove removes a session from the tracker.
func (t *SandboxSessionTracker) Remove(sessionName string) *tmux.SandboxSessionInfo {
	t.mu.Lock()
	defer t.mu.Unlock()
	info := t.sessions[sessionName]
	delete(t.sessions, sessionName)
	return info
}

// Get gets a session from the tracker.
func (t *SandboxSessionTracker) Get(sessionName string) *tmux.SandboxSessionInfo {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.sessions[sessionName]
}

// List returns all tracked sessions.
func (t *SandboxSessionTracker) List() []*tmux.SandboxSessionInfo {
	t.mu.RLock()
	defer t.mu.RUnlock()
	result := make([]*tmux.SandboxSessionInfo, 0, len(t.sessions))
	for _, info := range t.sessions {
		result = append(result, info)
	}
	return result
}

// sandboxSessionTracker is the global tracker for sandbox sessions.
var sandboxSessionTracker = NewSandboxSessionTracker()

// RegisterSandboxRoutes registers sandbox-related routes.
func (h *Handler) RegisterSandboxRoutes(r chi.Router) {
	r.Route("/sandbox", func(r chi.Router) {
		r.Post("/sessions/create", h.handleCreateSandboxSession)
		r.Get("/sessions", h.handleListSandboxSessions)
		r.Get("/sessions/{sessionName}", h.handleGetSandboxSession)
		r.Delete("/sessions/{sessionName}", h.handleKillSandboxSession)
		r.Get("/health", h.handleSandboxAPIHealth)
	})
}

// CreateSandboxSessionRequest is the request body for creating a sandbox session.
type CreateSandboxSessionRequest struct {
	// SandboxID is the ID of the sandbox to connect to
	SandboxID string `json:"sandbox_id"`

	// SessionName is the optional tmux session name (auto-generated if empty)
	SessionName string `json:"session_name,omitempty"`

	// TTLMinutes is the certificate TTL in minutes (1-10, default 5)
	TTLMinutes int `json:"ttl_minutes,omitempty"`
}

// CreateSandboxSessionResponse is the response for creating a sandbox session.
type CreateSandboxSessionResponse struct {
	// SessionID is the tmux session ID
	SessionID string `json:"session_id"`

	// SessionName is the tmux session name
	SessionName string `json:"session_name"`

	// SandboxID is the sandbox being accessed
	SandboxID string `json:"sandbox_id"`

	// VMIPAddress is the IP of the sandbox VM
	VMIPAddress string `json:"vm_ip_address"`

	// Username is the SSH username
	Username string `json:"username"`

	// ValidUntil is when the certificate expires (RFC3339)
	ValidUntil string `json:"valid_until"`

	// TTLSeconds is the remaining certificate validity in seconds
	TTLSeconds int `json:"ttl_seconds"`

	// Message provides additional information
	Message string `json:"message,omitempty"`
}

// SandboxSessionInfo is the response format for sandbox session queries.
type SandboxSessionInfo struct {
	SessionID   string `json:"session_id"`
	SessionName string `json:"session_name"`
	SandboxID   string `json:"sandbox_id"`
	VMIPAddress string `json:"vm_ip_address"`
	Username    string `json:"username"`
	ValidUntil  string `json:"valid_until"`
	TTLSeconds  int    `json:"ttl_seconds"`
	IsExpired   bool   `json:"is_expired"`
}

// ListSandboxSessionsResponse is the response for listing sandbox sessions.
type ListSandboxSessionsResponse struct {
	Sessions []SandboxSessionInfo `json:"sessions"`
	Total    int                  `json:"total"`
}

// @Summary Create sandbox session
// @Description Creates a new tmux session connected to a sandbox VM via SSH certificate
// @Tags Sandbox
// @Accept json
// @Produce json
// @Param request body CreateSandboxSessionRequest true "Create sandbox session request"
// @Success 200 {object} CreateSandboxSessionResponse
// @Failure 400 {object} types.APIError
// @Failure 500 {object} types.APIError
// @Id createSandboxSession
// @Router /v1/sandbox/sessions/create [post]
func (h *Handler) handleCreateSandboxSession(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	var req CreateSandboxSessionRequest
	if err := h.decodeJSON(r, &req); err != nil {
		h.writeError(w, r, http.StatusBadRequest, types.ErrCodeValidation, "Invalid request body", err.Error())
		return
	}

	if req.SandboxID == "" {
		h.writeError(w, r, http.StatusBadRequest, types.ErrCodeValidation, "sandbox_id is required", "")
		return
	}

	// Create the sandbox session
	sessionInfo, err := h.tmuxTool.CreateSandboxSession(r.Context(), req.SandboxID, req.SessionName, req.TTLMinutes)
	if err != nil {
		apiErr := &types.APIError{Code: types.ErrCodeInternal, Message: err.Error()}
		h.logAndAudit(r, "sandbox", "create_session", req, nil, apiErr, time.Since(start))
		h.writeError(w, r, http.StatusInternalServerError, types.ErrCodeInternal, "Failed to create sandbox session", err.Error())
		return
	}

	// Track the session for cleanup
	sandboxSessionTracker.Add(sessionInfo)

	// Calculate TTL
	ttlSeconds := int(time.Until(sessionInfo.ValidUntil).Seconds())
	if ttlSeconds < 0 {
		ttlSeconds = 0
	}

	resp := CreateSandboxSessionResponse{
		SessionID:   sessionInfo.SessionID,
		SessionName: sessionInfo.SessionName,
		SandboxID:   sessionInfo.SandboxID,
		VMIPAddress: sessionInfo.VMIPAddress,
		Username:    sessionInfo.Username,
		ValidUntil:  sessionInfo.ValidUntil.Format(time.RFC3339),
		TTLSeconds:  ttlSeconds,
		Message:     "Sandbox session created. Certificate valid for " + time.Until(sessionInfo.ValidUntil).Round(time.Second).String(),
	}

	h.logAndAudit(r, "sandbox", "create_session", req, resp, nil, time.Since(start))
	h.writeResponse(w, r, resp)
}

// @Summary List sandbox sessions
// @Description Lists all active sandbox sessions
// @Tags Sandbox
// @Produce json
// @Success 200 {object} ListSandboxSessionsResponse
// @Failure 500 {object} types.APIError
// @Id listSandboxSessions
// @Router /v1/sandbox/sessions [get]
func (h *Handler) handleListSandboxSessions(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	sessions := sandboxSessionTracker.List()
	result := make([]SandboxSessionInfo, len(sessions))

	for i, s := range sessions {
		ttlSeconds := int(time.Until(s.ValidUntil).Seconds())
		if ttlSeconds < 0 {
			ttlSeconds = 0
		}

		result[i] = SandboxSessionInfo{
			SessionID:   s.SessionID,
			SessionName: s.SessionName,
			SandboxID:   s.SandboxID,
			VMIPAddress: s.VMIPAddress,
			Username:    s.Username,
			ValidUntil:  s.ValidUntil.Format(time.RFC3339),
			TTLSeconds:  ttlSeconds,
			IsExpired:   time.Now().After(s.ValidUntil),
		}
	}

	resp := ListSandboxSessionsResponse{
		Sessions: result,
		Total:    len(result),
	}

	h.logAndAudit(r, "sandbox", "list_sessions", nil, resp, nil, time.Since(start))
	h.writeResponse(w, r, resp)
}

// @Summary Get sandbox session
// @Description Gets details of a specific sandbox session
// @Tags Sandbox
// @Produce json
// @Param sessionName path string true "Session name"
// @Success 200 {object} SandboxSessionInfo
// @Failure 404 {object} types.APIError
// @Id getSandboxSession
// @Router /v1/sandbox/sessions/{sessionName} [get]
func (h *Handler) handleGetSandboxSession(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	sessionName := chi.URLParam(r, "sessionName")

	if sessionName == "" {
		h.writeError(w, r, http.StatusBadRequest, types.ErrCodeValidation, "session_name is required", "")
		return
	}

	session := sandboxSessionTracker.Get(sessionName)
	if session == nil {
		h.writeError(w, r, http.StatusNotFound, types.ErrCodeNotFound, "Sandbox session not found", sessionName)
		return
	}

	ttlSeconds := int(time.Until(session.ValidUntil).Seconds())
	if ttlSeconds < 0 {
		ttlSeconds = 0
	}

	resp := SandboxSessionInfo{
		SessionID:   session.SessionID,
		SessionName: session.SessionName,
		SandboxID:   session.SandboxID,
		VMIPAddress: session.VMIPAddress,
		Username:    session.Username,
		ValidUntil:  session.ValidUntil.Format(time.RFC3339),
		TTLSeconds:  ttlSeconds,
		IsExpired:   time.Now().After(session.ValidUntil),
	}

	h.logAndAudit(r, "sandbox", "get_session", map[string]string{"session_name": sessionName}, resp, nil, time.Since(start))
	h.writeResponse(w, r, resp)
}

// @Summary Kill sandbox session
// @Description Kills a sandbox session and cleans up its credentials
// @Tags Sandbox
// @Produce json
// @Param sessionName path string true "Session name"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} types.APIError
// @Failure 404 {object} types.APIError
// @Failure 500 {object} types.APIError
// @Id killSandboxSession
// @Router /v1/sandbox/sessions/{sessionName} [delete]
func (h *Handler) handleKillSandboxSession(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	sessionName := chi.URLParam(r, "sessionName")

	if sessionName == "" {
		h.writeError(w, r, http.StatusBadRequest, types.ErrCodeValidation, "session_name is required", "")
		return
	}

	// Get and remove from tracker
	session := sandboxSessionTracker.Remove(sessionName)

	// Get connection info for cleanup
	var connInfo *sandbox.ConnectionInfo
	if session != nil {
		connInfo = session.ConnectionInfo
	}

	// Kill session and cleanup credentials
	err := h.tmuxTool.KillSandboxSession(r.Context(), sessionName, connInfo)
	if err != nil {
		apiErr := &types.APIError{Code: types.ErrCodeInternal, Message: err.Error()}
		h.logAndAudit(r, "sandbox", "kill_session", map[string]string{"session_name": sessionName}, nil, apiErr, time.Since(start))
		h.writeError(w, r, http.StatusInternalServerError, types.ErrCodeInternal, "Failed to kill sandbox session", err.Error())
		return
	}

	resp := map[string]interface{}{
		"session_name":        sessionName,
		"killed":              true,
		"credentials_cleaned": connInfo != nil,
	}

	h.logAndAudit(r, "sandbox", "kill_session", map[string]string{"session_name": sessionName}, resp, nil, time.Since(start))
	h.writeResponse(w, r, resp)
}

// @Summary Check sandbox API health
// @Description Checks if the virsh-sandbox API is reachable
// @Tags Sandbox
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} types.APIError
// @Id sandboxAPIHealth
// @Router /v1/sandbox/health [get]
func (h *Handler) handleSandboxAPIHealth(w http.ResponseWriter, r *http.Request) {
	// This would need access to the sandbox tool
	// For now, return a placeholder response
	resp := map[string]interface{}{
		"status":  "sandbox_api_check_not_implemented",
		"message": "Use the virsh-sandbox API directly to check its health",
	}
	h.writeResponse(w, r, resp)
}
