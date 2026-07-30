package networking

import (
	"context"
	"fmt"
	"sync"
)

type Firewall interface {
	Allow(ctx context.Context, rule string) error
	Drop(ctx context.Context, rule string) error
	List(ctx context.Context) ([]string, error)
}

var (
	ErrFirewallNotStarted = fmt.Errorf("firewall not started")
)

type firewall struct {
	mu      sync.RWMutex
	rules   map[string]struct{}
	started bool
}

func newFirewall() Firewall {
	return &firewall{rules: make(map[string]struct{})}
}

func (f *firewall) Allow(ctx context.Context, rule string) error {
	if !f.started {
		return ErrFirewallNotStarted
	}
	if rule == "" {
		return fmt.Errorf("rule must not be empty")
	}

	f.mu.Lock()
	defer f.mu.Unlock()
	f.rules[rule] = struct{}{}
	return nil
}

func (f *firewall) Drop(ctx context.Context, rule string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	delete(f.rules, rule)
	return nil
}

func (f *firewall) List(ctx context.Context) ([]string, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	out := make([]string, 0, len(f.rules))
	for rule := range f.rules {
		out = append(out, rule)
	}
	return out, nil
}
