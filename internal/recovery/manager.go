package recovery

import (
	"context"
	"fmt"
	"sync"
)

type RecoveryAction func(ctx context.Context) error

type RecoveryManager struct {
	mu      sync.RWMutex
	actions map[string]RecoveryAction
}

func NewRecoveryManager() *RecoveryManager {
	return &RecoveryManager{actions: make(map[string]RecoveryAction)}
}

func (r *RecoveryManager) Register(name string, action RecoveryAction) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.actions[name] = action
}

func (r *RecoveryManager) Execute(ctx context.Context, name string) error {
	r.mu.RLock()
	action, ok := r.actions[name]
	r.mu.RUnlock()

	if !ok {
		return fmt.Errorf("recovery action not found: %s", name)
	}
	return action(ctx)
}
