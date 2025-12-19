// Package api provides HTTP handlers for the tmux agent API (continued).
package api

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"tmux-client/internal/types"
)

// File handlers (continued)

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
		auditReq := map[string]interface{}{
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
	auditResp := map[string]interface{}{
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

func (h *Handler) handleGetAllowedCommands(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	resp := map[string]interface{}{
		"allowed": h.commandTool.GetAllowedCommands(),
		"denied":  h.commandTool.GetDeniedCommands(),
	}

	h.logAndAudit(r, "command", "get_allowed", nil, resp, nil, time.Since(start))
	h.writeResponse(w, r, resp)
}

// Human approval handlers

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

func (h *Handler) handleListPendingApprovals(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	pending := h.humanTool.ListPending()

	resp := types.ListApprovalsResponse{
		Pending: pending,
	}

	h.logAndAudit(r, "human", "list_pending", nil, map[string]int{"count": len(pending)}, nil, time.Since(start))
	h.writeResponse(w, r, resp)
}

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

func (h *Handler) handleAuditStats(w http.ResponseWriter, r *http.Request) {
	stats := h.auditLogger.Stats()
	h.writeResponse(w, r, stats)
}
