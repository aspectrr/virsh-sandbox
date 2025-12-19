// Package human provides human-in-the-loop approval functionality for sensitive actions.
package human

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"

	"tmux-client/internal/config"
	"tmux-client/internal/types"
)

// Tool provides human approval operations.
type Tool struct {
	config   *config.HumanConfig
	mu       sync.RWMutex
	pending  map[string]*pendingRequest
	approved map[string]*types.AskHumanResponse
}

// pendingRequest represents a pending approval request with its channel.
type pendingRequest struct {
	approval *types.PendingApproval
	resultCh chan *types.AskHumanResponse
}

// NewTool creates a new human approval tool.
func NewTool(cfg *config.HumanConfig) (*Tool, error) {
	return &Tool{
		config:   cfg,
		pending:  make(map[string]*pendingRequest),
		approved: make(map[string]*types.AskHumanResponse),
	}, nil
}

// AskHuman requests human approval for an action.
// This method blocks until the human approves, rejects, or the request times out.
func (t *Tool) AskHuman(ctx context.Context, req types.AskHumanRequest) (*types.AskHumanResponse, error) {
	// Validate request
	if req.Prompt == "" {
		return nil, fmt.Errorf("prompt is required")
	}

	if req.ActionType == "" {
		return nil, fmt.Errorf("action_type is required")
	}

	// Check if we've hit the pending limit
	t.mu.RLock()
	pendingCount := len(t.pending)
	t.mu.RUnlock()

	if pendingCount >= t.config.MaxPending {
		return nil, fmt.Errorf("maximum pending approvals reached (%d)", t.config.MaxPending)
	}

	// Determine timeout
	timeout := t.config.DefaultTimeout
	if req.TimeoutSecs > 0 {
		timeout = time.Duration(req.TimeoutSecs) * time.Second
	}

	// Generate request ID
	requestID := uuid.New().String()

	// Calculate expiry
	var expiresAt *time.Time
	if timeout > 0 {
		exp := time.Now().Add(timeout)
		expiresAt = &exp
	}

	// Set default urgency
	urgency := req.Urgency
	if urgency == "" {
		urgency = "medium"
	}

	// Create pending approval
	approval := &types.PendingApproval{
		RequestID:  requestID,
		Prompt:     req.Prompt,
		Context:    req.Context,
		ActionType: req.ActionType,
		Urgency:    urgency,
		Status:     types.ApprovalPending,
		CreatedAt:  time.Now().UTC(),
		ExpiresAt:  expiresAt,
	}

	// Create result channel
	resultCh := make(chan *types.AskHumanResponse, 1)

	// Store pending request
	t.mu.Lock()
	t.pending[requestID] = &pendingRequest{
		approval: approval,
		resultCh: resultCh,
	}
	t.mu.Unlock()

	// Send notification if configured
	t.sendNotification(approval)

	// Create context with timeout if specified
	var cancel context.CancelFunc
	if timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	// Wait for approval, rejection, or timeout
	select {
	case result := <-resultCh:
		return result, nil

	case <-ctx.Done():
		// Request timed out or was cancelled
		t.mu.Lock()
		delete(t.pending, requestID)
		t.mu.Unlock()

		return &types.AskHumanResponse{
			RequestID: requestID,
			Approved:  false,
			Status:    types.ApprovalExpired,
			Comment:   "request timed out",
			ExpiresAt: expiresAt,
		}, nil
	}
}

// AskHumanAsync requests human approval without blocking.
// Returns the request ID which can be used to check status or respond.
func (t *Tool) AskHumanAsync(req types.AskHumanRequest) (string, error) {
	// Validate request
	if req.Prompt == "" {
		return "", fmt.Errorf("prompt is required")
	}

	if req.ActionType == "" {
		return "", fmt.Errorf("action_type is required")
	}

	// Check if we've hit the pending limit
	t.mu.RLock()
	pendingCount := len(t.pending)
	t.mu.RUnlock()

	if pendingCount >= t.config.MaxPending {
		return "", fmt.Errorf("maximum pending approvals reached (%d)", t.config.MaxPending)
	}

	// Determine timeout
	timeout := t.config.DefaultTimeout
	if req.TimeoutSecs > 0 {
		timeout = time.Duration(req.TimeoutSecs) * time.Second
	}

	// Generate request ID
	requestID := uuid.New().String()

	// Calculate expiry
	var expiresAt *time.Time
	if timeout > 0 {
		exp := time.Now().Add(timeout)
		expiresAt = &exp
	}

	// Set default urgency
	urgency := req.Urgency
	if urgency == "" {
		urgency = "medium"
	}

	// Create pending approval
	approval := &types.PendingApproval{
		RequestID:  requestID,
		Prompt:     req.Prompt,
		Context:    req.Context,
		ActionType: req.ActionType,
		Urgency:    urgency,
		Status:     types.ApprovalPending,
		CreatedAt:  time.Now().UTC(),
		ExpiresAt:  expiresAt,
	}

	// Create result channel
	resultCh := make(chan *types.AskHumanResponse, 1)

	// Store pending request
	t.mu.Lock()
	t.pending[requestID] = &pendingRequest{
		approval: approval,
		resultCh: resultCh,
	}
	t.mu.Unlock()

	// Start expiry goroutine if timeout is set
	if timeout > 0 {
		go t.handleExpiry(requestID, timeout)
	}

	// Send notification if configured
	t.sendNotification(approval)

	return requestID, nil
}

// handleExpiry handles the expiry of a pending request.
func (t *Tool) handleExpiry(requestID string, timeout time.Duration) {
	time.Sleep(timeout)

	t.mu.Lock()
	defer t.mu.Unlock()

	if req, exists := t.pending[requestID]; exists {
		// Send expired result
		response := &types.AskHumanResponse{
			RequestID: requestID,
			Approved:  false,
			Status:    types.ApprovalExpired,
			Comment:   "request expired",
		}

		select {
		case req.resultCh <- response:
		default:
		}

		delete(t.pending, requestID)
	}
}

// Respond handles a human's response to an approval request.
func (t *Tool) Respond(req types.ApproveRequest) (*types.AskHumanResponse, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	// Find the pending request
	pending, exists := t.pending[req.RequestID]
	if !exists {
		return nil, fmt.Errorf("approval request not found: %s", req.RequestID)
	}

	// Check if already expired
	if pending.approval.ExpiresAt != nil && time.Now().After(*pending.approval.ExpiresAt) {
		delete(t.pending, req.RequestID)
		return nil, fmt.Errorf("approval request has expired")
	}

	now := time.Now().UTC()

	// Determine status
	var status types.ApprovalStatus
	if req.Approved {
		status = types.ApprovalApproved
	} else {
		status = types.ApprovalRejected
	}

	// Create response
	response := &types.AskHumanResponse{
		RequestID:  req.RequestID,
		Approved:   req.Approved,
		Status:     status,
		Comment:    req.Comment,
		ApprovedBy: req.ApprovedBy,
		ApprovedAt: &now,
		ExpiresAt:  pending.approval.ExpiresAt,
	}

	// Send response to waiting goroutine
	select {
	case pending.resultCh <- response:
	default:
		// Channel might already have been read
	}

	// Store in approved map for later reference
	t.approved[req.RequestID] = response

	// Remove from pending
	delete(t.pending, req.RequestID)

	return response, nil
}

// GetPending returns a pending approval request by ID.
func (t *Tool) GetPending(requestID string) (*types.PendingApproval, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	pending, exists := t.pending[requestID]
	if !exists {
		return nil, fmt.Errorf("approval request not found: %s", requestID)
	}

	// Check if expired
	if pending.approval.ExpiresAt != nil && time.Now().After(*pending.approval.ExpiresAt) {
		return nil, fmt.Errorf("approval request has expired")
	}

	return pending.approval, nil
}

// ListPending returns all pending approval requests.
func (t *Tool) ListPending() []types.PendingApproval {
	t.mu.RLock()
	defer t.mu.RUnlock()

	var approvals []types.PendingApproval
	now := time.Now()

	for _, req := range t.pending {
		// Skip expired requests
		if req.approval.ExpiresAt != nil && now.After(*req.approval.ExpiresAt) {
			continue
		}
		approvals = append(approvals, *req.approval)
	}

	return approvals
}

// GetResponse returns the response for a completed approval request.
func (t *Tool) GetResponse(requestID string) (*types.AskHumanResponse, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	response, exists := t.approved[requestID]
	if !exists {
		return nil, fmt.Errorf("approval response not found: %s", requestID)
	}

	return response, nil
}

// CancelPending cancels a pending approval request.
func (t *Tool) CancelPending(requestID string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	pending, exists := t.pending[requestID]
	if !exists {
		return fmt.Errorf("approval request not found: %s", requestID)
	}

	// Send cancelled response
	now := time.Now().UTC()
	response := &types.AskHumanResponse{
		RequestID:  requestID,
		Approved:   false,
		Status:     types.ApprovalRejected,
		Comment:    "request cancelled",
		ApprovedAt: &now,
	}

	select {
	case pending.resultCh <- response:
	default:
	}

	delete(t.pending, requestID)

	return nil
}

// CleanupExpired removes expired pending requests.
func (t *Tool) CleanupExpired() int {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now()
	cleaned := 0

	for id, req := range t.pending {
		if req.approval.ExpiresAt != nil && now.After(*req.approval.ExpiresAt) {
			// Send expired response
			response := &types.AskHumanResponse{
				RequestID: id,
				Approved:  false,
				Status:    types.ApprovalExpired,
				Comment:   "request expired during cleanup",
			}

			select {
			case req.resultCh <- response:
			default:
			}

			delete(t.pending, id)
			cleaned++
		}
	}

	return cleaned
}

// sendNotification sends a notification about a pending approval.
func (t *Tool) sendNotification(approval *types.PendingApproval) {
	// Send via command if configured
	if t.config.NotifyCommand != "" {
		go t.sendCommandNotification(approval)
	}

	// Send via webhook if configured
	if t.config.WebhookURL != "" {
		go t.sendWebhookNotification(approval)
	}
}

// sendCommandNotification executes the notification command.
func (t *Tool) sendCommandNotification(approval *types.PendingApproval) {
	args := make([]string, len(t.config.NotifyArgs))
	copy(args, t.config.NotifyArgs)

	// Replace placeholders in args
	for i, arg := range args {
		arg = strings.ReplaceAll(arg, "{request_id}", approval.RequestID)
		arg = strings.ReplaceAll(arg, "{prompt}", approval.Prompt)
		arg = strings.ReplaceAll(arg, "{action_type}", approval.ActionType)
		arg = strings.ReplaceAll(arg, "{urgency}", approval.Urgency)
		args[i] = arg
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, t.config.NotifyCommand, args...)
	cmd.Run() // Ignore errors for notifications
}

// sendWebhookNotification sends a webhook notification.
func (t *Tool) sendWebhookNotification(approval *types.PendingApproval) {
	// This is a placeholder for webhook notification
	// In a real implementation, you would use http.Client to POST to the webhook URL
	// For now, we'll leave this as a TODO
}

// RequiresApproval checks if an action type requires human approval based on config.
func (t *Tool) RequiresApproval(actionType string) bool {
	// Check auto-approve list first
	for _, autoApprove := range t.config.AutoApproveTypes {
		if actionType == autoApprove {
			return false
		}
	}

	// Check require approval list
	for _, require := range t.config.RequireApproval {
		if actionType == require {
			return true
		}
	}

	return false
}

// Stats returns statistics about the approval system.
func (t *Tool) Stats() ApprovalStats {
	t.mu.RLock()
	defer t.mu.RUnlock()

	stats := ApprovalStats{
		PendingCount:  len(t.pending),
		ApprovedCount: len(t.approved),
		ByActionType:  make(map[string]int),
		ByStatus:      make(map[string]int),
	}

	for _, req := range t.pending {
		stats.ByActionType[req.approval.ActionType]++
		stats.ByStatus["pending"]++
	}

	for _, resp := range t.approved {
		stats.ByStatus[string(resp.Status)]++
	}

	return stats
}

// ApprovalStats contains statistics about the approval system.
type ApprovalStats struct {
	PendingCount  int            `json:"pending_count"`
	ApprovedCount int            `json:"approved_count"`
	ByActionType  map[string]int `json:"by_action_type"`
	ByStatus      map[string]int `json:"by_status"`
}

// Close cleans up the approval tool.
func (t *Tool) Close() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	// Cancel all pending requests
	for id, req := range t.pending {
		response := &types.AskHumanResponse{
			RequestID: id,
			Approved:  false,
			Status:    types.ApprovalRejected,
			Comment:   "system shutdown",
		}

		select {
		case req.resultCh <- response:
		default:
		}

		close(req.resultCh)
	}

	t.pending = make(map[string]*pendingRequest)

	return nil
}
