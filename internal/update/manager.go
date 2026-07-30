package update

import (
	"context"
	"fmt"
	"sync"
)

type UpdateManager struct {
	mu           sync.RWMutex
	state        UpdateState
	manifest     *UpdateManifest
	backupPath   string
	migrations   *MigrationEngine
}

func NewUpdateManager(migrations *MigrationEngine) *UpdateManager {
	if migrations == nil {
		migrations = NewMigrationEngine()
	}
	return &UpdateManager{
		state:      UpdateStateIdle,
		migrations: migrations,
	}
}

func (m *UpdateManager) State() UpdateState {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.state
}

func (m *UpdateManager) SetState(state UpdateState) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.state = state
}

func (m *UpdateManager) Check(ctx context.Context) (*UpdateCheckResult, error) {
	m.SetState(UpdateStateChecking)

	result := &UpdateCheckResult{
		State:   UpdateStateChecking,
		Available: false,
	}

	if m.manifest != nil {
		result.Available = true
		result.Manifest = m.manifest
	}

	m.SetState(UpdateStateIdle)
	result.State = UpdateStateIdle

	return result, nil
}

func (m *UpdateManager) Prepare(ctx context.Context, manifest UpdateManifest) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if err := manifest.Validate(); err != nil {
		return fmt.Errorf("invalid manifest: %w", err)
	}

	m.manifest = &manifest
	m.backupPath = fmt.Sprintf("backup/%s", manifest.Version)
	return nil
}

func (m *UpdateManager) Migrate(ctx context.Context, targetVersion string) error {
	if m.migrations == nil {
		return fmt.Errorf("migration engine not configured")
	}

	m.SetState(UpdateStateMigrating)
	defer m.SetState(UpdateStateIdle)

	if err := m.migrations.Apply(ctx, targetVersion); err != nil {
		m.SetState(UpdateStateFailed)
		return err
	}

	return nil
}

func (m *UpdateManager) Rollback(ctx context.Context) error {
	if m.manifest == nil {
		return fmt.Errorf("no update in progress")
	}

	m.SetState(UpdateStateFailed)
	defer m.SetState(UpdateStateIdle)

	previousVersion := "0.0.0"
	if err := m.migrations.Rollback(ctx, previousVersion); err != nil {
		return fmt.Errorf("rollback failed: %w", err)
	}

	return nil
}

func (m *UpdateManager) Complete() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.manifest == nil {
		return fmt.Errorf("no update in progress")
	}

	m.state = UpdateStateCompleted
	m.manifest = nil
	m.backupPath = ""
	return nil
}
