package update

import (
	"context"
	"fmt"
	"sort"
	"sync"
)

type Migration struct {
	ID      string
	Version string
	Apply   func(ctx context.Context) error
	Rollback func(ctx context.Context) error
}

type MigrationEngine struct {
	mu         sync.RWMutex
	migrations []Migration
	applied    []string
}

func NewMigrationEngine() *MigrationEngine {
	return &MigrationEngine{
		migrations: make([]Migration, 0),
		applied:    make([]string, 0),
	}
}

func (e *MigrationEngine) Register(m Migration) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.migrations = append(e.migrations, m)
}

func (e *MigrationEngine) Apply(ctx context.Context, targetVersion string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	var pending []Migration
	for _, m := range e.migrations {
		if !contains(e.applied, m.ID) && m.Version <= targetVersion {
			pending = append(pending, m)
		}
	}

	for _, m := range pending {
		if err := m.Apply(ctx); err != nil {
			return fmt.Errorf("apply migration %s: %w", m.ID, err)
		}
		e.applied = append(e.applied, m.ID)
	}

	return nil
}

func (e *MigrationEngine) Rollback(ctx context.Context, targetVersion string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	var toRollback []Migration
	for i := len(e.applied) - 1; i >= 0; i-- {
		id := e.applied[i]
		for _, m := range e.migrations {
			if m.ID == id && m.Version > targetVersion && m.Rollback != nil {
				toRollback = append(toRollback, m)
				break
			}
		}
	}

	for _, m := range toRollback {
		if err := m.Rollback(ctx); err != nil {
			return fmt.Errorf("rollback migration %s: %w", m.ID, err)
		}
	}

	return nil
}

func (e *MigrationEngine) Applied() []string {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return append([]string(nil), e.applied...)
}

func (e *MigrationEngine) Pending(targetVersion string) []string {
	e.mu.RLock()
	defer e.mu.RUnlock()

	var pending []string
	for _, m := range e.migrations {
		if !contains(e.applied, m.ID) && m.Version <= targetVersion {
			pending = append(pending, m.ID)
		}
	}
	sort.Strings(pending)
	return pending
}

func contains(items []string, item string) bool {
	for _, i := range items {
		if i == item {
			return true
		}
	}
	return false
}
