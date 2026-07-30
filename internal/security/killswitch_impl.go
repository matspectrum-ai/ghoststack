package security

import (
	"context"
	"fmt"
	"sync"

	"github.com/ghoststack/ghoststack/internal/platform/linux"
)

type killSwitch struct {
	mu      sync.RWMutex
	active  bool
	iface   string
	firewall *linux.Firewall
}

func newKillSwitch(iface string) *killSwitch {
	return &killSwitch{
		iface:    iface,
		firewall: linux.NewFirewall(),
	}
}

func (k *killSwitch) Enable(ctx context.Context) error {
	k.mu.Lock()
	defer k.mu.Unlock()

	if k.active {
		return nil
	}

	if k.iface == "" {
		return fmt.Errorf("kill switch: no interface set")
	}

	if err := k.firewall.ApplyKillSwitch(ctx, k.iface); err != nil {
		return fmt.Errorf("kill switch enable: %w", err)
	}

	k.active = true
	return nil
}

func (k *killSwitch) Disable(ctx context.Context) error {
	k.mu.Lock()
	defer k.mu.Unlock()

	if !k.active {
		return nil
	}

	if err := k.firewall.RemoveKillSwitch(ctx); err != nil {
		return fmt.Errorf("kill switch disable: %w", err)
	}

	k.active = false
	return nil
}

func (k *killSwitch) Active() bool {
	k.mu.RLock()
	defer k.mu.RUnlock()
	return k.active
}
