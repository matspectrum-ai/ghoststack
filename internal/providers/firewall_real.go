package providers

import (
	"context"
	"fmt"
	"sync"
)

type FirewallReal struct {
	mu      sync.RWMutex
	rules   []string
	started bool
}

func NewFirewallReal() *FirewallReal {
	return &FirewallReal{rules: make([]string, 0)}
}

func (f *FirewallReal) Name() string {
	return "firewall-real"
}

func (f *FirewallReal) Start(ctx context.Context) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.started {
		return nil
	}

	f.started = true
	return nil
}

func (f *FirewallReal) Stop(ctx context.Context) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if !f.started {
		return ErrProviderNotStarted
	}

	f.started = false
	f.rules = nil
	return nil
}

func (f *FirewallReal) AddRule(ctx context.Context, rule string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if !f.started {
		return ErrProviderNotStarted
	}
	if rule == "" {
		return fmt.Errorf("rule must not be empty")
	}

	f.rules = append(f.rules, rule)
	return nil
}

func (f *FirewallReal) RemoveRule(ctx context.Context, rule string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	for i, r := range f.rules {
		if r == rule {
			f.rules = append(f.rules[:i], f.rules[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("rule not found: %s", rule)
}

func (f *FirewallReal) Rules() []string {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return append([]string(nil), f.rules...)
}
