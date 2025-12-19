// Package plan provides multi-step task planning and tracking for the agent API.
package plan

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"

	"tmux-client/internal/config"
	"tmux-client/internal/types"
)

// Tool provides plan management operations.
type Tool struct {
	config *config.PlanConfig
	mu     sync.RWMutex
	plans  map[string]*types.Plan
}

// NewTool creates a new plan tool.
func NewTool(cfg *config.PlanConfig) (*Tool, error) {
	tool := &Tool{
		config: cfg,
		plans:  make(map[string]*types.Plan),
	}

	// Create plan directory if persistence is enabled
	if cfg.PersistPlans && cfg.PlanDirectory != "" {
		if err := os.MkdirAll(cfg.PlanDirectory, 0750); err != nil {
			return nil, fmt.Errorf("failed to create plan directory: %w", err)
		}

		// Load existing plans
		if err := tool.loadPlans(); err != nil {
			// Log but don't fail - plans might be corrupted
			fmt.Printf("Warning: failed to load existing plans: %v\n", err)
		}
	}

	return tool, nil
}

// CreatePlan creates a new multi-step plan.
func (t *Tool) CreatePlan(req types.CreatePlanRequest) (*types.CreatePlanResponse, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	// Check plan limit
	if len(t.plans) >= t.config.MaxPlans {
		return nil, fmt.Errorf("maximum number of plans reached (%d)", t.config.MaxPlans)
	}

	// Validate steps
	if len(req.Steps) == 0 {
		return nil, fmt.Errorf("at least one step is required")
	}

	if len(req.Steps) > t.config.MaxStepsPerPlan {
		return nil, fmt.Errorf("too many steps: %d (max: %d)", len(req.Steps), t.config.MaxStepsPerPlan)
	}

	// Generate plan ID
	planID := uuid.New().String()

	// Create steps
	steps := make([]types.PlanStep, len(req.Steps))
	for i, desc := range req.Steps {
		steps[i] = types.PlanStep{
			Index:       i,
			Description: desc,
			Status:      types.StepStatusPending,
		}
	}

	now := time.Now().UTC()

	// Create plan
	plan := &types.Plan{
		ID:          planID,
		Name:        req.Name,
		Description: req.Description,
		Steps:       steps,
		Status:      types.PlanStatusPending,
		CurrentStep: -1, // Not started
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Store plan
	t.plans[planID] = plan

	// Persist if enabled
	if t.config.PersistPlans {
		if err := t.savePlan(plan); err != nil {
			// Log but don't fail
			fmt.Printf("Warning: failed to persist plan: %v\n", err)
		}
	}

	return &types.CreatePlanResponse{
		PlanID: planID,
		Plan:   *plan,
	}, nil
}

// UpdatePlan updates the status of a plan step.
func (t *Tool) UpdatePlan(req types.UpdatePlanRequest) (*types.UpdatePlanResponse, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	// Find plan
	plan, exists := t.plans[req.PlanID]
	if !exists {
		return nil, fmt.Errorf("plan not found: %s", req.PlanID)
	}

	// Validate step index
	if req.StepIndex < 0 || req.StepIndex >= len(plan.Steps) {
		return nil, fmt.Errorf("invalid step index: %d (plan has %d steps)", req.StepIndex, len(plan.Steps))
	}

	now := time.Now().UTC()

	// Update step
	step := &plan.Steps[req.StepIndex]

	// Handle status transitions
	switch req.Status {
	case types.StepStatusActive:
		if step.StartedAt == nil {
			step.StartedAt = &now
		}
		// Update current step if this is the first active step
		if plan.CurrentStep < req.StepIndex {
			plan.CurrentStep = req.StepIndex
		}
		// Update plan status
		if plan.Status == types.PlanStatusPending {
			plan.Status = types.PlanStatusInProgress
		}

	case types.StepStatusCompleted:
		step.CompletedAt = &now
		if step.StartedAt == nil {
			step.StartedAt = &now
		}

	case types.StepStatusFailed:
		step.CompletedAt = &now
		if step.StartedAt == nil {
			step.StartedAt = &now
		}

	case types.StepStatusSkipped:
		step.CompletedAt = &now
	}

	step.Status = req.Status
	step.Result = req.Result
	step.Error = req.Error

	// Update plan metadata
	plan.UpdatedAt = now

	// Check if plan is complete
	t.updatePlanStatus(plan)

	// Persist if enabled
	if t.config.PersistPlans {
		if err := t.savePlan(plan); err != nil {
			fmt.Printf("Warning: failed to persist plan: %v\n", err)
		}
	}

	return &types.UpdatePlanResponse{
		PlanID:  req.PlanID,
		Updated: true,
		Plan:    *plan,
	}, nil
}

// updatePlanStatus updates the overall plan status based on step statuses.
func (t *Tool) updatePlanStatus(plan *types.Plan) {
	completedCount := 0
	failedCount := 0
	skippedCount := 0
	activeCount := 0

	for _, step := range plan.Steps {
		switch step.Status {
		case types.StepStatusCompleted:
			completedCount++
		case types.StepStatusFailed:
			failedCount++
		case types.StepStatusSkipped:
			skippedCount++
		case types.StepStatusActive:
			activeCount++
		}
	}

	totalSteps := len(plan.Steps)
	finishedSteps := completedCount + failedCount + skippedCount

	// Determine plan status
	if failedCount > 0 && activeCount == 0 {
		// Plan failed if any step failed and no steps are active
		plan.Status = types.PlanStatusFailed
		now := time.Now().UTC()
		plan.CompletedAt = &now
	} else if finishedSteps == totalSteps {
		// All steps are done
		if failedCount > 0 {
			plan.Status = types.PlanStatusFailed
		} else {
			plan.Status = types.PlanStatusCompleted
		}
		now := time.Now().UTC()
		plan.CompletedAt = &now
	} else if activeCount > 0 || completedCount > 0 {
		plan.Status = types.PlanStatusInProgress
	}
}

// GetPlan retrieves a plan by ID.
func (t *Tool) GetPlan(planID string) (*types.GetPlanResponse, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	plan, exists := t.plans[planID]
	if !exists {
		return nil, fmt.Errorf("plan not found: %s", planID)
	}

	return &types.GetPlanResponse{
		Plan: *plan,
	}, nil
}

// ListPlans returns all plans.
func (t *Tool) ListPlans() (*types.ListPlansResponse, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	plans := make([]types.Plan, 0, len(t.plans))
	for _, plan := range t.plans {
		plans = append(plans, *plan)
	}

	// Sort by creation time (newest first)
	sort.Slice(plans, func(i, j int) bool {
		return plans[i].CreatedAt.After(plans[j].CreatedAt)
	})

	return &types.ListPlansResponse{
		Plans: plans,
	}, nil
}

// DeletePlan deletes a plan by ID.
func (t *Tool) DeletePlan(planID string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if _, exists := t.plans[planID]; !exists {
		return fmt.Errorf("plan not found: %s", planID)
	}

	delete(t.plans, planID)

	// Delete persisted plan if enabled
	if t.config.PersistPlans && t.config.PlanDirectory != "" {
		planFile := filepath.Join(t.config.PlanDirectory, planID+".json")
		os.Remove(planFile) // Ignore errors
	}

	return nil
}

// AbortPlan aborts a plan and marks all pending steps as skipped.
func (t *Tool) AbortPlan(planID string, reason string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	plan, exists := t.plans[planID]
	if !exists {
		return fmt.Errorf("plan not found: %s", planID)
	}

	if plan.Status == types.PlanStatusCompleted || plan.Status == types.PlanStatusAborted {
		return fmt.Errorf("plan is already %s", plan.Status)
	}

	now := time.Now().UTC()

	// Mark all pending/active steps as skipped
	for i := range plan.Steps {
		step := &plan.Steps[i]
		if step.Status == types.StepStatusPending || step.Status == types.StepStatusActive {
			step.Status = types.StepStatusSkipped
			step.Result = reason
			step.CompletedAt = &now
		}
	}

	plan.Status = types.PlanStatusAborted
	plan.UpdatedAt = now
	plan.CompletedAt = &now

	// Persist if enabled
	if t.config.PersistPlans {
		if err := t.savePlan(plan); err != nil {
			fmt.Printf("Warning: failed to persist plan: %v\n", err)
		}
	}

	return nil
}

// AdvanceStep marks the current step as completed and activates the next step.
func (t *Tool) AdvanceStep(planID string, result string) (*types.Plan, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	plan, exists := t.plans[planID]
	if !exists {
		return nil, fmt.Errorf("plan not found: %s", planID)
	}

	now := time.Now().UTC()

	// If plan hasn't started, start it
	if plan.CurrentStep < 0 {
		plan.CurrentStep = 0
		plan.Status = types.PlanStatusInProgress
		plan.Steps[0].Status = types.StepStatusActive
		plan.Steps[0].StartedAt = &now
		plan.UpdatedAt = now

		if t.config.PersistPlans {
			t.savePlan(plan)
		}

		return plan, nil
	}

	// Complete current step
	currentIdx := plan.CurrentStep
	if currentIdx >= 0 && currentIdx < len(plan.Steps) {
		step := &plan.Steps[currentIdx]
		step.Status = types.StepStatusCompleted
		step.Result = result
		step.CompletedAt = &now
	}

	// Advance to next step
	nextIdx := currentIdx + 1
	if nextIdx < len(plan.Steps) {
		plan.CurrentStep = nextIdx
		plan.Steps[nextIdx].Status = types.StepStatusActive
		plan.Steps[nextIdx].StartedAt = &now
	} else {
		// Plan is complete
		plan.Status = types.PlanStatusCompleted
		plan.CompletedAt = &now
	}

	plan.UpdatedAt = now

	if t.config.PersistPlans {
		t.savePlan(plan)
	}

	return plan, nil
}

// GetCurrentStep returns the current step of a plan.
func (t *Tool) GetCurrentStep(planID string) (*types.PlanStep, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	plan, exists := t.plans[planID]
	if !exists {
		return nil, fmt.Errorf("plan not found: %s", planID)
	}

	if plan.CurrentStep < 0 {
		return nil, fmt.Errorf("plan has not started")
	}

	if plan.CurrentStep >= len(plan.Steps) {
		return nil, fmt.Errorf("plan is complete")
	}

	step := plan.Steps[plan.CurrentStep]
	return &step, nil
}

// GetProgress returns the progress of a plan as a percentage.
func (t *Tool) GetProgress(planID string) (float64, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	plan, exists := t.plans[planID]
	if !exists {
		return 0, fmt.Errorf("plan not found: %s", planID)
	}

	if len(plan.Steps) == 0 {
		return 100, nil
	}

	completedCount := 0
	for _, step := range plan.Steps {
		if step.Status == types.StepStatusCompleted || step.Status == types.StepStatusSkipped {
			completedCount++
		}
	}

	return float64(completedCount) / float64(len(plan.Steps)) * 100, nil
}

// savePlan saves a plan to disk.
func (t *Tool) savePlan(plan *types.Plan) error {
	if t.config.PlanDirectory == "" {
		return nil
	}

	planFile := filepath.Join(t.config.PlanDirectory, plan.ID+".json")

	data, err := json.MarshalIndent(plan, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal plan: %w", err)
	}

	if err := os.WriteFile(planFile, data, 0600); err != nil {
		return fmt.Errorf("failed to write plan file: %w", err)
	}

	return nil
}

// loadPlans loads all plans from disk.
func (t *Tool) loadPlans() error {
	if t.config.PlanDirectory == "" {
		return nil
	}

	entries, err := os.ReadDir(t.config.PlanDirectory)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to read plan directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		planFile := filepath.Join(t.config.PlanDirectory, entry.Name())
		data, err := os.ReadFile(planFile)
		if err != nil {
			continue // Skip unreadable files
		}

		var plan types.Plan
		if err := json.Unmarshal(data, &plan); err != nil {
			continue // Skip malformed files
		}

		// Check retention period
		if t.config.RetentionPeriod > 0 {
			if time.Since(plan.UpdatedAt) > t.config.RetentionPeriod {
				// Plan has expired, delete it
				os.Remove(planFile)
				continue
			}
		}

		t.plans[plan.ID] = &plan
	}

	return nil
}

// CleanupExpired removes plans that have exceeded the retention period.
func (t *Tool) CleanupExpired() int {
	if t.config.RetentionPeriod <= 0 {
		return 0
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	cleaned := 0
	cutoff := time.Now().Add(-t.config.RetentionPeriod)

	for id, plan := range t.plans {
		if plan.UpdatedAt.Before(cutoff) {
			delete(t.plans, id)

			if t.config.PersistPlans && t.config.PlanDirectory != "" {
				planFile := filepath.Join(t.config.PlanDirectory, id+".json")
				os.Remove(planFile)
			}

			cleaned++
		}
	}

	return cleaned
}

// Stats returns statistics about the plan system.
func (t *Tool) Stats() PlanStats {
	t.mu.RLock()
	defer t.mu.RUnlock()

	stats := PlanStats{
		TotalPlans: len(t.plans),
		ByStatus:   make(map[string]int),
	}

	for _, plan := range t.plans {
		stats.ByStatus[string(plan.Status)]++

		totalSteps := len(plan.Steps)
		stats.TotalSteps += totalSteps

		for _, step := range plan.Steps {
			if step.Status == types.StepStatusCompleted {
				stats.CompletedSteps++
			}
		}
	}

	return stats
}

// PlanStats contains statistics about the plan system.
type PlanStats struct {
	TotalPlans     int            `json:"total_plans"`
	TotalSteps     int            `json:"total_steps"`
	CompletedSteps int            `json:"completed_steps"`
	ByStatus       map[string]int `json:"by_status"`
}

// Close cleans up the plan tool.
func (t *Tool) Close() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	// Persist all plans if enabled
	if t.config.PersistPlans {
		for _, plan := range t.plans {
			t.savePlan(plan)
		}
	}

	return nil
}
