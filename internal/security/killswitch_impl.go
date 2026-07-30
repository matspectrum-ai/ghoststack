package security

import (
	"context"
	"sync"
)

type killSwitch struct {
	mu      sync.RWMutex
	active  bool
	trigger func() error
}

func newKillSwitch(trigger func() error) KillSwitch {
	if trigger == nil {
		trigger = func() error { return nil }
	}
	return &killSwitch{trigger: trigger}
}

func (k *killSwitch) Enable(ctx context.Context) error {
	k.mu.Lock()
	defer k.mu.Unlock()

	if k.active {
		return nil
	}

	if err := k.trigger(); err != nil {
		return err
	}

	k.active = true
	return nil
}

func (k *killSwitch) Disable(ctx context.Context) error {
	k.mu.Lock()
	defer k.mu.Unlock()
	k.active = false
	return nil
}

func (k *killSwitch) Active() bool {
	k.mu.RLock()
	defer k.mu.RUnlock()
	return k.active
}
