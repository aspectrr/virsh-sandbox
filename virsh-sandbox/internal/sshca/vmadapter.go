// Package sshca provides SSH Certificate Authority management for ephemeral sandbox access.
package sshca

import (
	"context"
	"fmt"

	"virsh-sandbox/internal/store"
)

// SandboxStore defines the minimal interface needed to look up sandbox information.
type SandboxStore interface {
	GetSandbox(ctx context.Context, id string) (*store.Sandbox, error)
}

// VMAdapter implements VMInfoProvider by delegating to the sandbox store.
type VMAdapter struct {
	store SandboxStore
}

// NewVMAdapter creates a new VM adapter.
func NewVMAdapter(st SandboxStore) *VMAdapter {
	return &VMAdapter{
		store: st,
	}
}

// GetSandboxIP returns the IP address of a sandbox.
func (a *VMAdapter) GetSandboxIP(ctx context.Context, sandboxID string) (string, error) {
	sb, err := a.store.GetSandbox(ctx, sandboxID)
	if err != nil {
		return "", fmt.Errorf("get sandbox: %w", err)
	}

	if sb.IPAddress == nil || *sb.IPAddress == "" {
		return "", fmt.Errorf("sandbox %s has no IP address", sandboxID)
	}

	return *sb.IPAddress, nil
}

// GetSandboxVMName returns the VM name for a sandbox.
func (a *VMAdapter) GetSandboxVMName(ctx context.Context, sandboxID string) (string, error) {
	sb, err := a.store.GetSandbox(ctx, sandboxID)
	if err != nil {
		return "", fmt.Errorf("get sandbox: %w", err)
	}

	return sb.SandboxName, nil
}

// IsSandboxRunning checks if the sandbox is in a running state.
func (a *VMAdapter) IsSandboxRunning(ctx context.Context, sandboxID string) (bool, error) {
	sb, err := a.store.GetSandbox(ctx, sandboxID)
	if err != nil {
		return false, fmt.Errorf("get sandbox: %w", err)
	}

	return sb.State == store.SandboxStateRunning, nil
}

// Verify VMAdapter implements VMInfoProvider at compile time.
var _ VMInfoProvider = (*VMAdapter)(nil)
