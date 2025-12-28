// Package api provides HTTP handlers for the tmux agent API (continued).
package api

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"tmux-client/internal/types"
)

// File handlers (continued)

// @Summary List directory contents
// @Description Lists the contents of a directory
// @Tags File
// @Accept json
// @Produce json
// @Param request body types.ListDirRequest true "List directory request"
// @Success 200 {object} types.ListDirResponse
// @Failure 400 {object} types.APIError
// @Failure 404 {object} types.APIError
// @Failure 500 {object} types.APIError
// @Router /v1/file/list [post]
func (h *Handler) handleListDir(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	var req types.ListDirRequest
	if err := h.decodeJSON(r, &req); err != nil {
		h.writeError(w, r, http.StatusBadRequest, types.ErrCodeValidation, "Invalid request body", err.Error())
		return
	}

	if req.Path == "" {
		h.writeError(w, r, http.StatusBadRequest, types.ErrCodeValidation, "path is required", "")
		return
	}

	resp, err := h.fileTool.ListDir(req)
	if err != nil {
		code := types.ErrCodeInternal
		status := http.StatusInternalServerError
		if err.Error() == "directory not found: "+req.Path {
			code = types.ErrCodeNotFound
			status = http.StatusNotFound
		}
		apiErr := &types.APIError{Code: code, Message: err.Error()}
		h.logAndAudit(r, "file", "list", req, nil, apiErr, time.Since(start))
		h.writeError(w, r, status, code, "Failed to list directory", err.Error())
		return
	}

	h.logAndAudit(r, "file", "list", req, map[string]interface{}{"path": resp.Path, "count": len(resp.Files)}, nil, time.Since(start))
	h.writeResponse(w, r, resp)
}

// @Summary Check if file exists
// @Description Checks if a file or directory exists
// @Tags File
// @Accept json
// @Produce json
// @Param request body object true "File exists request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} types.APIError
// @Failure 500 {object} types.APIError
// @Router /v1/file/exists [post]
func (h *Handler) handleFileExists(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	var req struct {
		Path string `json:"path"`
	}
	if err := h.decodeJSON(r, &req); err != nil {
		h.writeError(w, r, http.StatusBadRequest, types.ErrCodeValidation, "Invalid request body", err.Error())
		return
	}

	if req.Path == "" {
		h.writeError(w, r, http.StatusBadRequest, types.ErrCodeValidation, "path is required", "")
		return
	}

	exists, isDir, err := h.fileTool.Exists(req.Path)
	if err != nil {
		apiErr := &types.APIError{Code: types.ErrCodeInternal, Message: err.Error()}
		h.logAndAudit(r, "file", "exists", req, nil, apiErr, time.Since(start))
		h.writeError(w, r, http.StatusInternalServerError, types.ErrCodeInternal, "Failed to check file", err.Error())
		return
	}

	resp := map[string]interface{}{
		"path":   req.Path,
		"exists": exists,
		"is_dir": isDir,
	}

	h.logAndAudit(r, "file", "exists", req, resp, nil, time.Since(start))
	h.writeResponse(w, r, resp)
}

// @Summary Get file hash
// @Description Computes the SHA256 hash of a file
// @Tags File
// @Accept json
// @Produce json
// @Param request body object true "File hash request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} types.APIError
// @Failure 500 {object} types.APIError
// @Router /v1/file/hash [post]
func (h *Handler) handleFileHash(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	var req struct {
		Path string `json:"path"`
	}
	if err := h.decodeJSON(r, &req); err != nil {
		h.writeError(w, r, http.StatusBadRequest, types.ErrCodeValidation, "Invalid request body", err.Error())
		return
	}

	if req.Path == "" {
		h.writeError(w, r, http.StatusBadRequest, types.ErrCodeValidation, "path is required", "")
		return
	}

	hash, err := h.fileTool.Hash(req.Path)
	if err != nil {
		apiErr := &types.APIError{Code: types.ErrCodeInternal, Message: err.Error()}
		h.logAndAudit(r, "file", "hash", req, nil, apiErr, time.Since(start))
		h.writeError(w, r, http.StatusInternalServerError, types.ErrCodeInternal, "Failed to hash file", err.Error())
		return
	}

	resp := map[string]string{
		"path":      req.Path,
		"hash":      hash,
		"algorithm": "sha256",
	}

	h.logAndAudit(r, "file", "hash", req, resp, nil, time.Since(start))
	h.writeResponse(w, r, resp)
}

// Command handlers

// @Summary Run command
// @Description Executes a shell command
// @Tags Command
// @Accept json
// @Produce json
// @Param request body types.RunCommandRequest true "Run command request"
// @Success 200 {object} types.RunCommandResponse
// @Failure 400 {object} types.APIError
// @Failure 403 {object} types.APIError
// @Failure 500 {object} types.APIError
// @Router /v1/command/run [post]
func (h *Handler) handleRunCommand(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	var req types.RunCommandRequest
	if err := h.decodeJSON(r, &req); err != nil {
		h.writeError(w, r, http.StatusBadRequest, types.ErrCodeValidation, "Invalid request body", err.Error())
		return
	}

	if req.Command == "" {
		h.writeError(w, r, http.StatusBadRequest, types.ErrCodeValidation, "command is required", "")
		return
	}

	resp, err := h.commandTool.RunCommand(r.Context(), req)
	if err != nil {
		code := types.ErrCodeInternal
		status := http.StatusInternalServerError
		// Check for blocked command errors
		if err.Error() == "command is not in the allowed list" ||
			err.Error() == "command is explicitly denied" ||
			err.Error() == "command contains disallowed pattern" {
			code = types.ErrCodeCommandBlocked
			status = http.StatusForbidden
		}
		apiErr := &types.APIError{Code: code, Message: err.Error()}
		// Don't log full args for security
		auditReq := map[string]any{
			"command":  req.Command,
			"args":     req.Args,
			"work_dir": req.WorkDir,
			"dry_run":  req.DryRun,
		}
		h.logAndAudit(r, "command", "run", auditReq, nil, apiErr, time.Since(start))
		h.writeError(w, r, status, code, "Failed to run command", err.Error())
		return
	}

	// Audit with truncated output
	auditResp := map[string]any{
		"command":     resp.Command,
		"args":        resp.Args,
		"exit_code":   resp.ExitCode,
		"dry_run":     resp.DryRun,
		"duration_ms": resp.DurationMs,
		"timed_out":   resp.TimedOut,
		"stdout_len":  len(resp.Stdout),
		"stderr_len":  len(resp.Stderr),
	}
	h.logAndAudit(r, "command", "run", req, auditResp, nil, time.Since(start))
	h.writeResponse(w, r, resp)
}

// @Summary Get allowed commands
// @Description Retrieves the list of allowed and denied commands
// @Tags Command
// @Accept json
// @Produce json
// @Success 200 {object} map[string]any
// @Router /v1/command/allowed [get]
func (h *Handler) handleGetAllowedCommands(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	resp := map[string]any{
		"allowed": h.commandTool.GetAllowedCommands(),
		"denied":  h.commandTool.GetDeniedCommands(),
	}

	h.logAndAudit(r, "command", "get_allowed", nil, resp, nil, time.Since(start))
	h.writeResponse(w, r, resp)
}

// Human approval handlers

// @Summary Request human approval
// @Description Requests approval from a human for an action
// @Tags Human
// @Accept json
// @Produce json
// @Param request body types.AskHumanRequest true "Ask human request"
// @Success 200 {object} types.AskHumanResponse
// @Failure 400 {object} types.APIError
// @Failure 500 {object} types.APIError
// @Router /v1/human/ask [post]
func (h *Handler) handleAskHuman(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	var req types.AskHumanRequest
	if err := h.decodeJSON(r, &req); err != nil {
		h.writeError(w, r, http.StatusBadRequest, types.ErrCodeValidation, "Invalid request body", err.Error())
		return
	}

	if req.Prompt == "" {
		h.writeError(w, r, http.StatusBadRequest, types.ErrCodeValidation, "prompt is required", "")
		return
	}

	if req.ActionType == "" {
		h.writeError(w, r, http.StatusBadRequest, types.ErrCodeValidation, "action_type is required", "")
		return
	}

	resp, err := h.humanTool.AskHuman(r.Context(), req)
	if err != nil {
		apiErr := &types.APIError{Code: types.ErrCodeInternal, Message: err.Error()}
		h.logAndAudit(r, "human", "ask", req, nil, apiErr, time.Since(start))
		h.writeError(w, r, http.StatusInternalServerError, types.ErrCodeInternal, "Failed to request approval", err.Error())
		return
	}

	h.logAndAudit(r, "human", "ask", req, resp, nil, time.Since(start))
	h.writeResponse(w, r, resp)
}

// @Summary Request human approval asynchronously
// @Description Requests approval from a human asynchronously
// @Tags Human
// @Accept json
// @Produce json
// @Param request body types.AskHumanRequest true "Ask human async request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} types.APIError
// @Failure 500 {object} types.APIError
// @Router /v1/human/ask-async [post]
func (h *Handler) handleAskHumanAsync(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	var req types.AskHumanRequest
	if err := h.decodeJSON(r, &req); err != nil {
		h.writeError(w, r, http.StatusBadRequest, types.ErrCodeValidation, "Invalid request body", err.Error())
		return
	}

	if req.Prompt == "" {
		h.writeError(w, r, http.StatusBadRequest, types.ErrCodeValidation, "prompt is required", "")
		return
	}

	if req.ActionType == "" {
		h.writeError(w, r, http.StatusBadRequest, types.ErrCodeValidation, "action_type is required", "")
		return
	}

	requestID, err := h.humanTool.AskHumanAsync(req)
	if err != nil {
		apiErr := &types.APIError{Code: types.ErrCodeInternal, Message: err.Error()}
		h.logAndAudit(r, "human", "ask_async", req, nil, apiErr, time.Since(start))
		h.writeError(w, r, http.StatusInternalServerError, types.ErrCodeInternal, "Failed to request approval", err.Error())
		return
	}

	resp := map[string]string{
		"request_id": requestID,
		"status":     string(types.ApprovalPending),
	}

	h.logAndAudit(r, "human", "ask_async", req, resp, nil, time.Since(start))
	h.writeResponse(w, r, resp)
}

// @Summary List pending approvals
// @Description Lists all pending human approval requests
// @Tags Human
// @Accept json
// @Produce json
// @Success 200 {object} types.ListApprovalsResponse
// @Router /v1/human/pending [get]
func (h *Handler) handleListPendingApprovals(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	pending := h.humanTool.ListPending()

	resp := types.ListApprovalsResponse{
		Pending: pending,
	}

	h.logAndAudit(r, "human", "list_pending", nil, map[string]int{"count": len(pending)}, nil, time.Since(start))
	h.writeResponse(w, r, resp)
}

// @Summary Get pending approval
// @Description Retrieves a specific pending approval request
// @Tags Human
// @Accept json
// @Produce json
// @Param requestID path string true "Request ID"
// @Success 200 {object} types.PendingApproval
// @Failure 400 {object} types.APIError
// @Failure 404 {object} types.APIError
// @Router /v1/human/pending/{requestID} [get]
func (h *Handler) handleGetPendingApproval(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestID := chi.URLParam(r, "requestID")

	if requestID == "" {
		h.writeError(w, r, http.StatusBadRequest, types.ErrCodeValidation, "request_id is required", "")
		return
	}

	approval, err := h.humanTool.GetPending(requestID)
	if err != nil {
		apiErr := &types.APIError{Code: types.ErrCodeNotFound, Message: err.Error()}
		h.logAndAudit(r, "human", "get_pending", map[string]string{"request_id": requestID}, nil, apiErr, time.Since(start))
		h.writeError(w, r, http.StatusNotFound, types.ErrCodeNotFound, "Approval not found", err.Error())
		return
	}

	h.logAndAudit(r, "human", "get_pending", map[string]string{"request_id": requestID}, approval, nil, time.Since(start))
	h.writeResponse(w, r, approval)
}

// @Summary Respond to approval
// @Description Responds to a pending approval request
// @Tags Human
// @Accept json
// @Produce json
// @Param request body types.ApproveRequest true "Approve request"
// @Success 200 {object} types.AskHumanResponse
// @Failure 400 {object} types.APIError
// @Failure 404 {object} types.APIError
// @Router /v1/human/respond [post]
func (h *Handler) handleRespondToApproval(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	var req types.ApproveRequest
	if err := h.decodeJSON(r, &req); err != nil {
		h.writeError(w, r, http.StatusBadRequest, types.ErrCodeValidation, "Invalid request body", err.Error())
		return
	}

	if req.RequestID == "" {
		h.writeError(w, r, http.StatusBadRequest, types.ErrCodeValidation, "request_id is required", "")
		return
	}

	if req.ApprovedBy == "" {
		h.writeError(w, r, http.StatusBadRequest, types.ErrCodeValidation, "approved_by is required", "")
		return
	}

	resp, err := h.humanTool.Respond(req)
	if err != nil {
		apiErr := &types.APIError{Code: types.ErrCodeNotFound, Message: err.Error()}
		h.logAndAudit(r, "human", "respond", req, nil, apiErr, time.Since(start))
		h.writeError(w, r, http.StatusNotFound, types.ErrCodeNotFound, "Approval not found", err.Error())
		return
	}

	h.logAndAudit(r, "human", "respond", req, resp, nil, time.Since(start))
	h.writeResponse(w, r, resp)
}

// @Summary Cancel approval
// @Description Cancels a pending approval request
// @Tags Human
// @Accept json
// @Produce json
// @Param requestID path string true "Request ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} types.APIError
// @Failure 404 {object} types.APIError
// @Router /v1/human/pending/{requestID} [delete]
func (h *Handler) handleCancelApproval(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestID := chi.URLParam(r, "requestID")

	if requestID == "" {
		h.writeError(w, r, http.StatusBadRequest, types.ErrCodeValidation, "request_id is required", "")
		return
	}

	err := h.humanTool.CancelPending(requestID)
	if err != nil {
		apiErr := &types.APIError{Code: types.ErrCodeNotFound, Message: err.Error()}
		h.logAndAudit(r, "human", "cancel", map[string]string{"request_id": requestID}, nil, apiErr, time.Since(start))
		h.writeError(w, r, http.StatusNotFound, types.ErrCodeNotFound, "Approval not found", err.Error())
		return
	}

	resp := map[string]interface{}{
		"request_id": requestID,
		"cancelled":  true,
	}

	h.logAndAudit(r, "human", "cancel", map[string]string{"request_id": requestID}, resp, nil, time.Since(start))
	h.writeResponse(w, r, resp)
}

// Plan handlers

// @Summary Create plan
// @Description Creates a new execution plan
// @Tags Plan
// @Accept json
// @Produce json
// @Param request body types.CreatePlanRequest true "Create plan request"
// @Success 200 {object} types.CreatePlanResponse
// @Failure 400 {object} types.APIError
// @Failure 500 {object} types.APIError
// @Router /v1/plan/create [post]
func (h *Handler) handleCreatePlan(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	var req types.CreatePlanRequest
	if err := h.decodeJSON(r, &req); err != nil {
		h.writeError(w, r, http.StatusBadRequest, types.ErrCodeValidation, "Invalid request body", err.Error())
		return
	}

	if len(req.Steps) == 0 {
		h.writeError(w, r, http.StatusBadRequest, types.ErrCodeValidation, "at least one step is required", "")
		return
	}

	resp, err := h.planTool.CreatePlan(req)
	if err != nil {
		apiErr := &types.APIError{Code: types.ErrCodeInternal, Message: err.Error()}
		h.logAndAudit(r, "plan", "create", req, nil, apiErr, time.Since(start))
		h.writeError(w, r, http.StatusInternalServerError, types.ErrCodeInternal, "Failed to create plan", err.Error())
		return
	}

	h.logAndAudit(r, "plan", "create", req, map[string]interface{}{"plan_id": resp.PlanID, "steps": len(req.Steps)}, nil, time.Since(start))
	h.writeResponse(w, r, resp)
}

// @Summary List plans
// @Description Lists all execution plans
// @Tags Plan
// @Accept json
// @Produce json
// @Success 200 {object} types.ListPlansResponse
// @Failure 500 {object} types.APIError
// @Router /v1/plan/ [get]
func (h *Handler) handleListPlans(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	resp, err := h.planTool.ListPlans()
	if err != nil {
		apiErr := &types.APIError{Code: types.ErrCodeInternal, Message: err.Error()}
		h.logAndAudit(r, "plan", "list", nil, nil, apiErr, time.Since(start))
		h.writeError(w, r, http.StatusInternalServerError, types.ErrCodeInternal, "Failed to list plans", err.Error())
		return
	}

	h.logAndAudit(r, "plan", "list", nil, map[string]int{"count": len(resp.Plans)}, nil, time.Since(start))
	h.writeResponse(w, r, resp)
}

// @Summary Get plan
// @Description Retrieves a specific execution plan
// @Tags Plan
// @Accept json
// @Produce json
// @Param planID path string true "Plan ID"
// @Success 200 {object} types.GetPlanResponse
// @Failure 400 {object} types.APIError
// @Failure 404 {object} types.APIError
// @Router /v1/plan/{planID} [get]
func (h *Handler) handleGetPlan(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	planID := chi.URLParam(r, "planID")

	if planID == "" {
		h.writeError(w, r, http.StatusBadRequest, types.ErrCodeValidation, "plan_id is required", "")
		return
	}

	resp, err := h.planTool.GetPlan(planID)
	if err != nil {
		apiErr := &types.APIError{Code: types.ErrCodeNotFound, Message: err.Error()}
		h.logAndAudit(r, "plan", "get", map[string]string{"plan_id": planID}, nil, apiErr, time.Since(start))
		h.writeError(w, r, http.StatusNotFound, types.ErrCodeNotFound, "Plan not found", err.Error())
		return
	}

	h.logAndAudit(r, "plan", "get", map[string]string{"plan_id": planID}, map[string]interface{}{"status": resp.Plan.Status}, nil, time.Since(start))
	h.writeResponse(w, r, resp)
}

// @Summary Update plan
// @Description Updates an execution plan
// @Tags Plan
// @Accept json
// @Produce json
// @Param request body types.UpdatePlanRequest true "Update plan request"
// @Success 200 {object} types.UpdatePlanResponse
// @Failure 400 {object} types.APIError
// @Failure 404 {object} types.APIError
// @Failure 500 {object} types.APIError
// @Router /v1/plan/update [post]
func (h *Handler) handleUpdatePlan(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	var req types.UpdatePlanRequest
	if err := h.decodeJSON(r, &req); err != nil {
		h.writeError(w, r, http.StatusBadRequest, types.ErrCodeValidation, "Invalid request body", err.Error())
		return
	}

	if req.PlanID == "" {
		h.writeError(w, r, http.StatusBadRequest, types.ErrCodeValidation, "plan_id is required", "")
		return
	}

	resp, err := h.planTool.UpdatePlan(req)
	if err != nil {
		code := types.ErrCodeInternal
		status := http.StatusInternalServerError
		if err.Error() == "plan not found: "+req.PlanID {
			code = types.ErrCodeNotFound
			status = http.StatusNotFound
		}
		apiErr := &types.APIError{Code: code, Message: err.Error()}
		h.logAndAudit(r, "plan", "update", req, nil, apiErr, time.Since(start))
		h.writeError(w, r, status, code, "Failed to update plan", err.Error())
		return
	}

	h.logAndAudit(r, "plan", "update", req, map[string]interface{}{"plan_id": resp.PlanID, "updated": resp.Updated}, nil, time.Since(start))
	h.writeResponse(w, r, resp)
}

// @Summary Advance plan step
// @Description Advances to the next step in a plan
// @Tags Plan
// @Accept json
// @Produce json
// @Param planID path string true "Plan ID"
// @Param request body object false "Advance step request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} types.APIError
// @Failure 404 {object} types.APIError
// @Failure 500 {object} types.APIError
// @Router /v1/plan/{planID}/advance [post]
func (h *Handler) handleAdvanceStep(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	planID := chi.URLParam(r, "planID")

	if planID == "" {
		h.writeError(w, r, http.StatusBadRequest, types.ErrCodeValidation, "plan_id is required", "")
		return
	}

	var req struct {
		Result string `json:"result,omitempty"`
	}
	if err := h.decodeJSON(r, &req); err != nil && err.Error() != "EOF" {
		h.writeError(w, r, http.StatusBadRequest, types.ErrCodeValidation, "Invalid request body", err.Error())
		return
	}

	plan, err := h.planTool.AdvanceStep(planID, req.Result)
	if err != nil {
		code := types.ErrCodeInternal
		status := http.StatusInternalServerError
		if err.Error() == "plan not found: "+planID {
			code = types.ErrCodeNotFound
			status = http.StatusNotFound
		}
		apiErr := &types.APIError{Code: code, Message: err.Error()}
		h.logAndAudit(r, "plan", "advance", map[string]string{"plan_id": planID}, nil, apiErr, time.Since(start))
		h.writeError(w, r, status, code, "Failed to advance plan", err.Error())
		return
	}

	resp := map[string]interface{}{
		"plan_id":      planID,
		"current_step": plan.CurrentStep,
		"status":       plan.Status,
	}

	h.logAndAudit(r, "plan", "advance", map[string]string{"plan_id": planID}, resp, nil, time.Since(start))
	h.writeResponse(w, r, resp)
}

// @Summary Abort plan
// @Description Aborts an execution plan
// @Tags Plan
// @Accept json
// @Produce json
// @Param planID path string true "Plan ID"
// @Param request body object false "Abort plan request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} types.APIError
// @Failure 404 {object} types.APIError
// @Failure 500 {object} types.APIError
// @Router /v1/plan/{planID}/abort [post]
func (h *Handler) handleAbortPlan(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	planID := chi.URLParam(r, "planID")

	if planID == "" {
		h.writeError(w, r, http.StatusBadRequest, types.ErrCodeValidation, "plan_id is required", "")
		return
	}

	var req struct {
		Reason string `json:"reason,omitempty"`
	}
	if err := h.decodeJSON(r, &req); err != nil && err.Error() != "EOF" {
		h.writeError(w, r, http.StatusBadRequest, types.ErrCodeValidation, "Invalid request body", err.Error())
		return
	}

	reason := req.Reason
	if reason == "" {
		reason = "aborted by user"
	}

	err := h.planTool.AbortPlan(planID, reason)
	if err != nil {
		code := types.ErrCodeInternal
		status := http.StatusInternalServerError
		if err.Error() == "plan not found: "+planID {
			code = types.ErrCodeNotFound
			status = http.StatusNotFound
		}
		apiErr := &types.APIError{Code: code, Message: err.Error()}
		h.logAndAudit(r, "plan", "abort", map[string]string{"plan_id": planID, "reason": reason}, nil, apiErr, time.Since(start))
		h.writeError(w, r, status, code, "Failed to abort plan", err.Error())
		return
	}

	resp := map[string]interface{}{
		"plan_id": planID,
		"aborted": true,
		"reason":  reason,
	}

	h.logAndAudit(r, "plan", "abort", map[string]string{"plan_id": planID, "reason": reason}, resp, nil, time.Since(start))
	h.writeResponse(w, r, resp)
}

// @Summary Delete plan
// @Description Deletes an execution plan
// @Tags Plan
// @Accept json
// @Produce json
// @Param planID path string true "Plan ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} types.APIError
// @Failure 404 {object} types.APIError
// @Router /v1/plan/{planID} [delete]
func (h *Handler) handleDeletePlan(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	planID := chi.URLParam(r, "planID")

	if planID == "" {
		h.writeError(w, r, http.StatusBadRequest, types.ErrCodeValidation, "plan_id is required", "")
		return
	}

	err := h.planTool.DeletePlan(planID)
	if err != nil {
		apiErr := &types.APIError{Code: types.ErrCodeNotFound, Message: err.Error()}
		h.logAndAudit(r, "plan", "delete", map[string]string{"plan_id": planID}, nil, apiErr, time.Since(start))
		h.writeError(w, r, http.StatusNotFound, types.ErrCodeNotFound, "Plan not found", err.Error())
		return
	}

	resp := map[string]interface{}{
		"plan_id": planID,
		"deleted": true,
	}

	h.logAndAudit(r, "plan", "delete", map[string]string{"plan_id": planID}, resp, nil, time.Since(start))
	h.writeResponse(w, r, resp)
}

// Audit handlers

// @Summary Query audit log
// @Description Queries the audit log for entries
// @Tags Audit
// @Accept json
// @Produce json
// @Param request body types.AuditQuery false "Audit query"
// @Success 200 {object} types.AuditQueryResponse
// @Failure 400 {object} types.APIError
// @Failure 500 {object} types.APIError
// @Router /v1/audit/query [post]
func (h *Handler) handleQueryAudit(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	var req types.AuditQuery
	if err := h.decodeJSON(r, &req); err != nil && err.Error() != "EOF" {
		h.writeError(w, r, http.StatusBadRequest, types.ErrCodeValidation, "Invalid request body", err.Error())
		return
	}

	entries, count, hasMore, err := h.auditLogger.Query(req)
	if err != nil {
		apiErr := &types.APIError{Code: types.ErrCodeInternal, Message: err.Error()}
		h.logAndAudit(r, "audit", "query", req, nil, apiErr, time.Since(start))
		h.writeError(w, r, http.StatusInternalServerError, types.ErrCodeInternal, "Failed to query audit log", err.Error())
		return
	}

	resp := types.AuditQueryResponse{
		Entries:    entries,
		TotalCount: count,
		HasMore:    hasMore,
	}

	// Don't log audit queries to avoid infinite loops
	h.writeResponse(w, r, resp)
	_ = start // Avoid unused variable warning
}

// @Summary Get audit stats
// @Description Retrieves audit log statistics
// @Tags Audit
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /v1/audit/stats [get]
func (h *Handler) handleAuditStats(w http.ResponseWriter, r *http.Request) {
	stats := h.auditLogger.Stats()
	h.writeResponse(w, r, stats)
}
