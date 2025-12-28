// Package api provides HTTP handlers for the tmux agent API.
package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/MarceloPetrucio/go-scalar-api-reference"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"tmux-client/internal/audit"
	"tmux-client/internal/config"
	"tmux-client/internal/tools/command"
	"tmux-client/internal/tools/file"
	"tmux-client/internal/tools/human"
	"tmux-client/internal/tools/plan"
	"tmux-client/internal/tools/tmux"
	"tmux-client/internal/types"
)

// Handler provides HTTP handlers for the API.
type Handler struct {
	config      *config.Config
	logger      zerolog.Logger
	auditLogger *audit.Logger
	tmuxTool    *tmux.Tool
	fileTool    *file.Tool
	commandTool *command.Tool
	humanTool   *human.Tool
	planTool    *plan.Tool
	startTime   time.Time
	version     string
}

// NewHandler creates a new API handler.
func NewHandler(
	cfg *config.Config,
	logger zerolog.Logger,
	auditLogger *audit.Logger,
	tmuxTool *tmux.Tool,
	fileTool *file.Tool,
	commandTool *command.Tool,
	humanTool *human.Tool,
	planTool *plan.Tool,
	version string,
) *Handler {
	return &Handler{
		config:      cfg,
		logger:      logger,
		auditLogger: auditLogger,
		tmuxTool:    tmuxTool,
		fileTool:    fileTool,
		commandTool: commandTool,
		humanTool:   humanTool,
		planTool:    planTool,
		startTime:   time.Now(),
		version:     version,
	}
}

// RegisterRoutes registers all API routes.
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Get("/docs", func(w http.ResponseWriter, r *http.Request) {
		htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
			// SpecURL: "https://generator3.swagger.io/openapi.json",// allow external URL or local path file
			SpecURL: "./docs/swagger.json",
			CustomOptions: scalar.CustomOptions{
				PageTitle: "tmux-client API",
			},
			DarkMode: true,
		})
		if err != nil {
			fmt.Printf("%v", err)
		}

		fmt.Fprintln(w, htmlContent)
	})

	// API v1
	r.Route("/v1", func(r chi.Router) {
		// Health check
		r.Get("/health", h.handleHealth)

		// Tmux endpoints
		r.Route("/tmux", func(r chi.Router) {
			r.Get("/sessions", h.handleListSessions)
			r.Get("/windows", h.handleListWindows)
			r.Get("/panes", h.handleListPanes)
			r.Post("/panes/read", h.handleReadPane)
			r.Post("/panes/switch", h.handleSwitchPane)
			r.Post("/panes/create", h.handleCreatePane)
			r.Post("/panes/send-keys", h.handleSendKeys)
			r.Post("/sessions/create", h.handleCreateSession)
			r.Delete("/panes/{paneID}", h.handleKillPane)
			r.Post("/sessions/{sessionId}/release", h.handleReleaseSession)
		})

		// File endpoints
		r.Route("/file", func(r chi.Router) {
			r.Post("/read", h.handleReadFile)
			r.Post("/write", h.handleWriteFile)
			r.Post("/edit", h.handleEditFile)
			r.Post("/copy", h.handleCopyFile)
			r.Post("/delete", h.handleDeleteFile)
			r.Post("/list", h.handleListDir)
			r.Post("/exists", h.handleFileExists)
			r.Post("/hash", h.handleFileHash)
		})

		// Command endpoints
		r.Route("/command", func(r chi.Router) {
			r.Post("/run", h.handleRunCommand)
			r.Get("/allowed", h.handleGetAllowedCommands)
		})

		// Human approval endpoints
		r.Route("/human", func(r chi.Router) {
			r.Post("/ask", h.handleAskHuman)
			r.Post("/ask-async", h.handleAskHumanAsync)
			r.Get("/pending", h.handleListPendingApprovals)
			r.Get("/pending/{requestID}", h.handleGetPendingApproval)
			r.Post("/respond", h.handleRespondToApproval)
			r.Delete("/pending/{requestID}", h.handleCancelApproval)
		})

		// Plan endpoints
		r.Route("/plan", func(r chi.Router) {
			r.Post("/create", h.handleCreatePlan)
			r.Get("/", h.handleListPlans)
			r.Get("/{planID}", h.handleGetPlan)
			r.Post("/update", h.handleUpdatePlan)
			r.Post("/{planID}/advance", h.handleAdvanceStep)
			r.Post("/{planID}/abort", h.handleAbortPlan)
			r.Delete("/{planID}", h.handleDeletePlan)
		})

		// Audit endpoints
		r.Route("/audit", func(r chi.Router) {
			r.Post("/query", h.handleQueryAudit)
			r.Get("/stats", h.handleAuditStats)
		})
	})
}

// Helper functions for response handling

func (h *Handler) requestID(r *http.Request) string {
	if id := r.Header.Get("X-Request-ID"); id != "" {
		return id
	}
	return uuid.New().String()
}

func (h *Handler) writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *Handler) writeResponse(w http.ResponseWriter, r *http.Request, data any) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		h.writeError(w, r, http.StatusInternalServerError, types.ErrCodeInternal, "Failed to marshal response", err.Error())
		return
	}

	resp := types.APIResponse{
		Success:   true,
		Data:      jsonData,
		Timestamp: time.Now().UTC(),
		RequestID: h.requestID(r),
	}

	h.writeJSON(w, http.StatusOK, resp)
}

func (h *Handler) writeError(w http.ResponseWriter, r *http.Request, status int, code, message, details string) {
	resp := types.APIResponse{
		Success: false,
		Error: &types.APIError{
			Code:    code,
			Message: message,
			Details: details,
		},
		Timestamp: time.Now().UTC(),
		RequestID: h.requestID(r),
	}

	h.writeJSON(w, status, resp)
}

func (h *Handler) decodeJSON(r *http.Request, v any) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func (h *Handler) logAndAudit(r *http.Request, tool, action string, args, result any, apiErr *types.APIError, duration time.Duration) {
	if h.auditLogger != nil {
		h.auditLogger.LogToolCall(
			h.requestID(r),
			tool,
			action,
			args,
			result,
			apiErr,
			duration,
			r.RemoteAddr,
			r.UserAgent(),
		)
	}
}

// @Summary Get health status
// @Description Retrieves the health status of the API server and its components
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} types.HealthResponse
// @Router /v1/health [get]
// Health check handler
func (h *Handler) handleHealth(w http.ResponseWriter, r *http.Request) {
	components := []types.ComponentHealth{
		{Name: "api", Status: types.HealthStatusHealthy},
	}

	// Check tmux
	if h.tmuxTool != nil {
		if h.tmuxTool.IsTmuxRunning(r.Context()) {
			components = append(components, types.ComponentHealth{
				Name:   "tmux",
				Status: types.HealthStatusHealthy,
			})
		} else {
			components = append(components, types.ComponentHealth{
				Name:    "tmux",
				Status:  types.HealthStatusDegraded,
				Message: "tmux server not running",
			})
		}
	}

	// Check file tool
	if h.fileTool != nil {
		components = append(components, types.ComponentHealth{
			Name:   "file",
			Status: types.HealthStatusHealthy,
		})
	}

	// Check command tool
	if h.commandTool != nil {
		components = append(components, types.ComponentHealth{
			Name:   "command",
			Status: types.HealthStatusHealthy,
		})
	}

	// Determine overall status
	overallStatus := types.HealthStatusHealthy
	for _, c := range components {
		if c.Status == types.HealthStatusUnhealthy {
			overallStatus = types.HealthStatusUnhealthy
			break
		}
		if c.Status == types.HealthStatusDegraded {
			overallStatus = types.HealthStatusDegraded
		}
	}

	uptime := time.Since(h.startTime).Round(time.Second).String()

	resp := types.HealthResponse{
		Status:     overallStatus,
		Version:    h.version,
		Uptime:     uptime,
		Components: components,
	}

	h.writeJSON(w, http.StatusOK, resp)
}

// Tmux handlers

// @Summary List tmux sessions
// @Description Get a list of all active tmux sessions
// @Tags Tmux
// @Accept json
// @Produce json
// @Success 200 {array} types.SessionInfo
// @Failure 500 {object} types.APIError
// @Router /v1/tmux/sessions [get]
func (h *Handler) handleListSessions(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	sessions, err := h.tmuxTool.ListSessions(r.Context())
	if err != nil {
		apiErr := &types.APIError{Code: types.ErrCodeInternal, Message: err.Error()}
		h.logAndAudit(r, "tmux", "list_sessions", nil, nil, apiErr, time.Since(start))
		h.writeError(w, r, http.StatusInternalServerError, types.ErrCodeInternal, "Failed to list sessions", err.Error())
		return
	}

	h.logAndAudit(r, "tmux", "list_sessions", nil, sessions, nil, time.Since(start))
	h.writeResponse(w, r, sessions)
}

// @Summary Release tmux session
// @Description Releases (kills) a tmux session by ID
// @Tags Tmux
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID"
// @Success 200 {object} types.KillSessionResponse
// @Failure 500 {object} types.APIError
// @Router /v1/tmux/sessions/{sessionId}/release [post]
func (h *Handler) handleReleaseSession(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	sessionName := r.URL.Query().Get("session")

	err := h.tmuxTool.KillSession(r.Context(), sessionName)
	if err != nil {
		apiErr := &types.APIError{Code: types.ErrCodeInternal, Message: err.Error()}
		h.logAndAudit(r, "tmux", "release_session", map[string]string{"session": sessionName}, nil, apiErr, time.Since(start))
		h.writeError(w, r, http.StatusInternalServerError, types.ErrCodeInternal, "Failed to release session", err.Error())
		return
	}

	h.logAndAudit(r, "tmux", "release_session", map[string]string{"session": sessionName}, nil, nil, time.Since(start))
	resp := types.KillSessionResponse{
		SessionName: sessionName,
		Success:     true,
	}
	h.writeResponse(w, r, resp)
}

// @Summary List tmux windows
// @Description Get a list of windows in a tmux session
// @Tags Tmux
// @Accept json
// @Produce json
// @Param session query string false "Session name"
// @Success 200 {array} types.WindowInfo
// @Failure 500 {object} types.APIError
// @Router /v1/tmux/windows [get]
func (h *Handler) handleListWindows(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	sessionName := r.URL.Query().Get("session")

	windows, err := h.tmuxTool.ListWindows(r.Context(), sessionName)
	if err != nil {
		apiErr := &types.APIError{Code: types.ErrCodeInternal, Message: err.Error()}
		h.logAndAudit(r, "tmux", "list_windows", map[string]string{"session": sessionName}, nil, apiErr, time.Since(start))
		h.writeError(w, r, http.StatusInternalServerError, types.ErrCodeInternal, "Failed to list windows", err.Error())
		return
	}

	h.logAndAudit(r, "tmux", "list_windows", map[string]string{"session": sessionName}, windows, nil, time.Since(start))
	h.writeResponse(w, r, windows)
}

// @Summary List tmux panes
// @Description Get a list of panes in a tmux session
// @Tags Tmux
// @Accept json
// @Produce json
// @Param session query string false "Session name"
// @Success 200 {object} types.ListPanesResponse
// @Failure 400 {object} types.APIError
// @Failure 500 {object} types.APIError
// @Router /v1/tmux/panes [get]
func (h *Handler) handleListPanes(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	var req types.ListPanesRequest
	if r.Method == http.MethodPost {
		if err := h.decodeJSON(r, &req); err != nil && err.Error() != "EOF" {
			h.writeError(w, r, http.StatusBadRequest, types.ErrCodeValidation, "Invalid request body", err.Error())
			return
		}
	} else {
		req.SessionName = r.URL.Query().Get("session")
	}

	panes, err := h.tmuxTool.ListPanes(r.Context(), req.SessionName)
	if err != nil {
		apiErr := &types.APIError{Code: types.ErrCodeInternal, Message: err.Error()}
		h.logAndAudit(r, "tmux", "list_panes", req, nil, apiErr, time.Since(start))
		h.writeError(w, r, http.StatusInternalServerError, types.ErrCodeInternal, "Failed to list panes", err.Error())
		return
	}

	resp := types.ListPanesResponse{Panes: panes}
	h.logAndAudit(r, "tmux", "list_panes", req, resp, nil, time.Since(start))
	h.writeResponse(w, r, resp)
}

// @Summary Read tmux pane
// @Description Reads the content of a tmux pane
// @Tags Tmux
// @Accept json
// @Produce json
// @Param request body types.ReadPaneRequest true "Read pane request"
// @Success 200 {object} types.ReadPaneResponse
// @Failure 400 {object} types.APIError
// @Failure 500 {object} types.APIError
// @Router /v1/tmux/panes/read [post]
func (h *Handler) handleReadPane(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	var req types.ReadPaneRequest
	if err := h.decodeJSON(r, &req); err != nil {
		h.writeError(w, r, http.StatusBadRequest, types.ErrCodeValidation, "Invalid request body", err.Error())
		return
	}

	if req.PaneID == "" {
		h.writeError(w, r, http.StatusBadRequest, types.ErrCodeValidation, "pane_id is required", "")
		return
	}

	content, lines, err := h.tmuxTool.ReadPane(r.Context(), req.PaneID, req.LastNLines)
	if err != nil {
		apiErr := &types.APIError{Code: types.ErrCodeInternal, Message: err.Error()}
		h.logAndAudit(r, "tmux", "read_pane", req, nil, apiErr, time.Since(start))
		h.writeError(w, r, http.StatusInternalServerError, types.ErrCodeInternal, "Failed to read pane", err.Error())
		return
	}

	resp := types.ReadPaneResponse{
		PaneID:  req.PaneID,
		Content: content,
		Lines:   lines,
	}

	h.logAndAudit(r, "tmux", "read_pane", req, map[string]interface{}{"pane_id": req.PaneID, "lines": lines}, nil, time.Since(start))
	h.writeResponse(w, r, resp)
}

// @Summary Switch tmux pane
// @Description Switches to a specific tmux pane
// @Tags Tmux
// @Accept json
// @Produce json
// @Param request body types.SwitchPaneRequest true "Switch pane request"
// @Success 200 {object} types.SwitchPaneResponse
// @Failure 400 {object} types.APIError
// @Failure 500 {object} types.APIError
// @Router /v1/tmux/panes/switch [post]
func (h *Handler) handleSwitchPane(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	var req types.SwitchPaneRequest
	if err := h.decodeJSON(r, &req); err != nil {
		h.writeError(w, r, http.StatusBadRequest, types.ErrCodeValidation, "Invalid request body", err.Error())
		return
	}

	if req.PaneID == "" {
		h.writeError(w, r, http.StatusBadRequest, types.ErrCodeValidation, "pane_id is required", "")
		return
	}

	err := h.tmuxTool.SwitchPane(r.Context(), req.PaneID)
	if err != nil {
		apiErr := &types.APIError{Code: types.ErrCodeInternal, Message: err.Error()}
		h.logAndAudit(r, "tmux", "switch_pane", req, nil, apiErr, time.Since(start))
		h.writeError(w, r, http.StatusInternalServerError, types.ErrCodeInternal, "Failed to switch pane", err.Error())
		return
	}

	resp := types.SwitchPaneResponse{
		PaneID:   req.PaneID,
		Switched: true,
	}

	h.logAndAudit(r, "tmux", "switch_pane", req, resp, nil, time.Since(start))
	h.writeResponse(w, r, resp)
}

// @Summary Create tmux pane
// @Description Creates a new tmux pane
// @Tags Tmux
// @Accept json
// @Produce json
// @Param request body types.CreatePaneRequest true "Create pane request"
// @Success 200 {object} types.CreatePaneResponse
// @Failure 400 {object} types.APIError
// @Failure 500 {object} types.APIError
// @Router /v1/tmux/panes/create [post]
func (h *Handler) handleCreatePane(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	var req types.CreatePaneRequest
	if err := h.decodeJSON(r, &req); err != nil {
		h.writeError(w, r, http.StatusBadRequest, types.ErrCodeValidation, "Invalid request body", err.Error())
		return
	}

	resp, err := h.tmuxTool.CreatePane(r.Context(), req)
	if err != nil {
		apiErr := &types.APIError{Code: types.ErrCodeInternal, Message: err.Error()}
		h.logAndAudit(r, "tmux", "create_pane", req, nil, apiErr, time.Since(start))
		h.writeError(w, r, http.StatusInternalServerError, types.ErrCodeInternal, "Failed to create pane", err.Error())
		return
	}

	h.logAndAudit(r, "tmux", "create_pane", req, resp, nil, time.Since(start))
	h.writeResponse(w, r, resp)
}

// @Summary Send keys to tmux pane
// @Description Sends keystrokes to a tmux pane
// @Tags Tmux
// @Accept json
// @Produce json
// @Param request body types.SendKeysRequest true "Send keys request"
// @Success 200 {object} types.SendKeysResponse
// @Failure 400 {object} types.APIError
// @Failure 403 {object} types.APIError
// @Failure 500 {object} types.APIError
// @Router /v1/tmux/panes/send-keys [post]
func (h *Handler) handleSendKeys(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	var req types.SendKeysRequest
	if err := h.decodeJSON(r, &req); err != nil {
		h.writeError(w, r, http.StatusBadRequest, types.ErrCodeValidation, "Invalid request body", err.Error())
		return
	}

	if req.PaneID == "" {
		h.writeError(w, r, http.StatusBadRequest, types.ErrCodeValidation, "pane_id is required", "")
		return
	}

	if req.Key == "" {
		h.writeError(w, r, http.StatusBadRequest, types.ErrCodeValidation, "key is required", "")
		return
	}

	err := h.tmuxTool.SendKeys(r.Context(), req.PaneID, req.Key)
	if err != nil {
		apiErr := &types.APIError{Code: types.ErrCodeForbidden, Message: err.Error()}
		h.logAndAudit(r, "tmux", "send_keys", req, nil, apiErr, time.Since(start))
		h.writeError(w, r, http.StatusForbidden, types.ErrCodeForbidden, "Failed to send keys", err.Error())
		return
	}

	resp := types.SendKeysResponse{
		PaneID: req.PaneID,
		Sent:   true,
	}

	h.logAndAudit(r, "tmux", "send_keys", req, resp, nil, time.Since(start))
	h.writeResponse(w, r, resp)
}

// @Summary Create tmux session
// @Description Creates a new tmux session
// @Tags Tmux
// @Accept json
// @Produce json
// @Param request body object true "Create session request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} types.APIError
// @Failure 500 {object} types.APIError
// @Router /v1/tmux/sessions/create [post]
func (h *Handler) handleCreateSession(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	var req struct {
		Name    string `json:"name"`
		Command string `json:"command,omitempty"`
	}

	if err := h.decodeJSON(r, &req); err != nil {
		h.writeError(w, r, http.StatusBadRequest, types.ErrCodeValidation, "Invalid request body", err.Error())
		return
	}

	if req.Name == "" {
		h.writeError(w, r, http.StatusBadRequest, types.ErrCodeValidation, "name is required", "")
		return
	}

	sessionID, err := h.tmuxTool.CreateSession(r.Context(), req.Name, req.Command)
	if err != nil {
		apiErr := &types.APIError{Code: types.ErrCodeInternal, Message: err.Error()}
		h.logAndAudit(r, "tmux", "create_session", req, nil, apiErr, time.Since(start))
		h.writeError(w, r, http.StatusInternalServerError, types.ErrCodeInternal, "Failed to create session", err.Error())
		return
	}

	resp := map[string]string{
		"session_id":   sessionID,
		"session_name": req.Name,
	}

	h.logAndAudit(r, "tmux", "create_session", req, resp, nil, time.Since(start))
	h.writeResponse(w, r, resp)
}

// @Summary Kill tmux pane
// @Description Kills a tmux pane
// @Tags Tmux
// @Accept json
// @Produce json
// @Param paneID path string true "Pane ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} types.APIError
// @Failure 500 {object} types.APIError
// @Router /v1/tmux/panes/{paneID} [delete]
func (h *Handler) handleKillPane(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	paneID := chi.URLParam(r, "paneID")

	if paneID == "" {
		h.writeError(w, r, http.StatusBadRequest, types.ErrCodeValidation, "pane_id is required", "")
		return
	}

	err := h.tmuxTool.KillPane(r.Context(), paneID)
	if err != nil {
		apiErr := &types.APIError{Code: types.ErrCodeInternal, Message: err.Error()}
		h.logAndAudit(r, "tmux", "kill_pane", map[string]string{"pane_id": paneID}, nil, apiErr, time.Since(start))
		h.writeError(w, r, http.StatusInternalServerError, types.ErrCodeInternal, "Failed to kill pane", err.Error())
		return
	}

	resp := map[string]interface{}{
		"pane_id": paneID,
		"killed":  true,
	}

	h.logAndAudit(r, "tmux", "kill_pane", map[string]string{"pane_id": paneID}, resp, nil, time.Since(start))
	h.writeResponse(w, r, resp)
}

// @Summary Kill tmux session
// @Description Kills a tmux session
// @Tags Tmux
// @Accept json
// @Produce json
// @Param sessionName path string true "Session name"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} types.APIError
// @Failure 500 {object} types.APIError
// @Router /v1/tmux/sessions/{sessionName} [delete]
func (h *Handler) handleKillSession(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	sessionName := chi.URLParam(r, "sessionName")

	if sessionName == "" {
		h.writeError(w, r, http.StatusBadRequest, types.ErrCodeValidation, "session_name is required", "")
		return
	}

	err := h.tmuxTool.KillSession(r.Context(), sessionName)
	if err != nil {
		apiErr := &types.APIError{Code: types.ErrCodeInternal, Message: err.Error()}
		h.logAndAudit(r, "tmux", "kill_session", map[string]string{"session_name": sessionName}, nil, apiErr, time.Since(start))
		h.writeError(w, r, http.StatusInternalServerError, types.ErrCodeInternal, "Failed to kill session", err.Error())
		return
	}

	resp := map[string]interface{}{
		"session_name": sessionName,
		"killed":       true,
	}

	h.logAndAudit(r, "tmux", "kill_session", map[string]string{"session_name": sessionName}, resp, nil, time.Since(start))
	h.writeResponse(w, r, resp)
}

// File handlers

// @Summary Read file
// @Description Reads the content of a file
// @Tags File
// @Accept json
// @Produce json
// @Param request body types.ReadFileRequest true "Read file request"
// @Success 200 {object} types.ReadFileResponse
// @Failure 400 {object} types.APIError
// @Failure 403 {object} types.APIError
// @Failure 404 {object} types.APIError
// @Failure 500 {object} types.APIError
// @Router /v1/file/read [post]
func (h *Handler) handleReadFile(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	var req types.ReadFileRequest
	if err := h.decodeJSON(r, &req); err != nil {
		h.writeError(w, r, http.StatusBadRequest, types.ErrCodeValidation, "Invalid request body", err.Error())
		return
	}

	if req.Path == "" {
		h.writeError(w, r, http.StatusBadRequest, types.ErrCodeValidation, "path is required", "")
		return
	}

	resp, err := h.fileTool.ReadFile(req)
	if err != nil {
		code := types.ErrCodeInternal
		status := http.StatusInternalServerError
		if err.Error() == "file not found: "+req.Path {
			code = types.ErrCodeNotFound
			status = http.StatusNotFound
		} else if err.Error() == "access denied: path is in denied list" || err.Error() == "access denied: file extension is not allowed" {
			code = types.ErrCodeForbidden
			status = http.StatusForbidden
		}
		apiErr := &types.APIError{Code: code, Message: err.Error()}
		h.logAndAudit(r, "file", "read", req, nil, apiErr, time.Since(start))
		h.writeError(w, r, status, code, "Failed to read file", err.Error())
		return
	}

	// Don't log file content in audit for privacy/size reasons
	auditResp := map[string]any{
		"path":        resp.Path,
		"total_lines": resp.TotalLines,
		"size":        resp.Size,
		"truncated":   resp.Truncated,
	}
	h.logAndAudit(r, "file", "read", req, auditResp, nil, time.Since(start))
	h.writeResponse(w, r, resp)
}

// @Summary Write file
// @Description Writes content to a file
// @Tags File
// @Accept json
// @Produce json
// @Param request body types.WriteFileRequest true "Write file request"
// @Success 200 {object} types.WriteFileResponse
// @Failure 400 {object} types.APIError
// @Failure 403 {object} types.APIError
// @Failure 500 {object} types.APIError
// @Router /v1/file/write [post]
func (h *Handler) handleWriteFile(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	var req types.WriteFileRequest
	if err := h.decodeJSON(r, &req); err != nil {
		h.writeError(w, r, http.StatusBadRequest, types.ErrCodeValidation, "Invalid request body", err.Error())
		return
	}

	if req.Path == "" {
		h.writeError(w, r, http.StatusBadRequest, types.ErrCodeValidation, "path is required", "")
		return
	}

	resp, err := h.fileTool.WriteFile(req)
	if err != nil {
		code := types.ErrCodeInternal
		status := http.StatusInternalServerError
		if err.Error() == "access denied: path is in denied list" || err.Error() == "access denied: file extension is not allowed" {
			code = types.ErrCodeForbidden
			status = http.StatusForbidden
		}
		apiErr := &types.APIError{Code: code, Message: err.Error()}
		// Don't log content in audit
		auditReq := map[string]interface{}{
			"path":       req.Path,
			"mode":       req.Mode,
			"overwrite":  req.Overwrite,
			"create_dir": req.CreateDir,
			"size":       len(req.Content),
		}
		h.logAndAudit(r, "file", "write", auditReq, nil, apiErr, time.Since(start))
		h.writeError(w, r, status, code, "Failed to write file", err.Error())
		return
	}

	auditReq := map[string]interface{}{
		"path":       req.Path,
		"mode":       req.Mode,
		"overwrite":  req.Overwrite,
		"create_dir": req.CreateDir,
		"size":       len(req.Content),
	}
	h.logAndAudit(r, "file", "write", auditReq, resp, nil, time.Since(start))
	h.writeResponse(w, r, resp)
}

// @Summary Edit file
// @Description Edits the content of a file
// @Tags File
// @Accept json
// @Produce json
// @Param request body types.EditFileRequest true "Edit file request"
// @Success 200 {object} types.EditFileResponse
// @Failure 400 {object} types.APIError
// @Failure 404 {object} types.APIError
// @Failure 500 {object} types.APIError
// @Router /v1/file/edit [post]
func (h *Handler) handleEditFile(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	var req types.EditFileRequest
	if err := h.decodeJSON(r, &req); err != nil {
		h.writeError(w, r, http.StatusBadRequest, types.ErrCodeValidation, "Invalid request body", err.Error())
		return
	}

	if req.Path == "" {
		h.writeError(w, r, http.StatusBadRequest, types.ErrCodeValidation, "path is required", "")
		return
	}

	resp, err := h.fileTool.EditFile(req)
	if err != nil {
		code := types.ErrCodeInternal
		status := http.StatusInternalServerError
		if err.Error() == "file not found: "+req.Path {
			code = types.ErrCodeNotFound
			status = http.StatusNotFound
		} else if err.Error() == "old text not found in file" {
			code = types.ErrCodeValidation
			status = http.StatusBadRequest
		}
		apiErr := &types.APIError{Code: code, Message: err.Error()}
		h.logAndAudit(r, "file", "edit", req, nil, apiErr, time.Since(start))
		h.writeError(w, r, status, code, "Failed to edit file", err.Error())
		return
	}

	// Log with diff but without full content
	auditResp := map[string]any{
		"path":         resp.Path,
		"edited":       resp.Edited,
		"replacements": resp.Replacements,
		"diff":         resp.Diff,
	}
	h.logAndAudit(r, "file", "edit", req, auditResp, nil, time.Since(start))
	h.writeResponse(w, r, resp)
}

// @Summary Copy file
// @Description Copies a file from source to destination
// @Tags File
// @Accept json
// @Produce json
// @Param request body types.CopyFileRequest true "Copy file request"
// @Success 200 {object} types.CopyFileResponse
// @Failure 400 {object} types.APIError
// @Failure 500 {object} types.APIError
// @Router /v1/file/copy [post]
func (h *Handler) handleCopyFile(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	var req types.CopyFileRequest
	if err := h.decodeJSON(r, &req); err != nil {
		h.writeError(w, r, http.StatusBadRequest, types.ErrCodeValidation, "Invalid request body", err.Error())
		return
	}

	if req.Source == "" || req.Destination == "" {
		h.writeError(w, r, http.StatusBadRequest, types.ErrCodeValidation, "source and destination are required", "")
		return
	}

	resp, err := h.fileTool.CopyFile(req)
	if err != nil {
		apiErr := &types.APIError{Code: types.ErrCodeInternal, Message: err.Error()}
		h.logAndAudit(r, "file", "copy", req, nil, apiErr, time.Since(start))
		h.writeError(w, r, http.StatusInternalServerError, types.ErrCodeInternal, "Failed to copy file", err.Error())
		return
	}

	h.logAndAudit(r, "file", "copy", req, resp, nil, time.Since(start))
	h.writeResponse(w, r, resp)
}

// @Summary Delete file
// @Description Deletes a file or directory
// @Tags File
// @Accept json
// @Produce json
// @Param request body types.DeleteFileRequest true "Delete file request"
// @Success 200 {object} types.DeleteFileResponse
// @Failure 400 {object} types.APIError
// @Failure 403 {object} types.APIError
// @Failure 500 {object} types.APIError
// @Router /v1/file/delete [post]
func (h *Handler) handleDeleteFile(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	var req types.DeleteFileRequest
	if err := h.decodeJSON(r, &req); err != nil {
		h.writeError(w, r, http.StatusBadRequest, types.ErrCodeValidation, "Invalid request body", err.Error())
		return
	}

	if req.Path == "" {
		h.writeError(w, r, http.StatusBadRequest, types.ErrCodeValidation, "path is required", "")
		return
	}

	resp, err := h.fileTool.DeleteFile(req)
	if err != nil {
		code := types.ErrCodeInternal
		status := http.StatusInternalServerError
		if err.Error() == "delete operations are disabled in config" {
			code = types.ErrCodeForbidden
			status = http.StatusForbidden
		}
		apiErr := &types.APIError{Code: code, Message: err.Error()}
		h.logAndAudit(r, "file", "delete", req, nil, apiErr, time.Since(start))
		h.writeError(w, r, status, code, "Failed to delete file", err.Error())
		return
	}

	h.logAndAudit(r, "file", "delete", req, resp, nil, time.Since(start))
	h.writeResponse(w, r, resp)
}
