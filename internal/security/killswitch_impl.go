package security

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ghoststack/ghoststack/internal/platform/linux"
	"github.com/ghoststack/ghoststack/internal/storage"
)

type killSwitch struct {
	mu        sync.RWMutex
	active    bool
	iface     string
	ifaceIPv6 bool
	firewall  *linux.Firewall
	store     storage.StorageProvider
	dnsForce  bool
	dnsAddrs  []string
}

type KillSwitchOption func(*killSwitch)

func WithStorage(s storage.StorageProvider) KillSwitchOption {
	return func(k *killSwitch) {
		k.store = s
	}
}

func WithDNSForce(addrs []string) KillSwitchOption {
	return func(k *killSwitch) {
		if len(addrs) > 0 {
			k.dnsForce = true
			k.dnsAddrs = addrs
		}
	}
}

func newKillSwitch(iface string, opts ...KillSwitchOption) *killSwitch {
	k := &killSwitch{
		iface:    iface,
		firewall: linux.NewFirewall(),
	}

	for _, opt := range opts {
		opt(k)
	}

	return k
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

	if k.store != nil {
		_ = k.store.AppendAuditLog(ctx, storage.AuditEntry{
			Timestamp: time.Now().Unix(),
			Action:    "killswitch.enabled",
			Source:    "security",
			Detail:    fmt.Sprintf("interface=%s dns_force=%v", k.iface, k.dnsForce),
		})
	}

	return nil
}

func (k *killSwitch) Disable(ctx context.Context) error {
	k.mu.Lock()
	defer k.mu.Unlock()

	if !k.active {
		return nil
	}

	_ = k.firewall.RemoveKillSwitch(ctx)

	k.active = false

	if k.store != nil {
		_ = k.store.AppendAuditLog(ctx, storage.AuditEntry{
			Timestamp: time.Now().Unix(),
			Action:    "killswitch.disabled",
			Source:    "security",
			Detail:    "kill switch removed",
		})
	}

	return nil
}

func (k *killSwitch) Active() bool {
	k.mu.RLock()
	defer k.mu.RUnlock()
	return k.active
}

func (k *killSwitch) Status() map[string]any {
	k.mu.RLock()
	defer k.mu.RUnlock()
	return map[string]any{
		"active":    k.active,
		"interface": k.iface,
		"dns_force": k.dnsForce,
		"dns_addrs": k.dnsAddrs,
	}
}

func (k *killSwitch) RunLeakTest(ctx context.Context) (*LeakCheckResult, error) {
	result, err := CheckLeaks(ctx, k.iface)
	if err != nil {
		return nil, err
	}

	if k.store != nil {
		detail := fmt.Sprintf("dns=%v ipv6=%v icmp=%v public=%s",
			result.DNS, result.IPv6, result.ICMP, result.PublicIP)
		_ = k.store.AppendAuditLog(ctx, storage.AuditEntry{
			Timestamp: time.Now().Unix(),
			Action:    "killswitch.leaktest",
			Source:    "security",
			Detail:    detail,
		})
	}

	return result, nil
}
