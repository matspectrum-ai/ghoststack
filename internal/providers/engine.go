package providers

import (
	"context"
	"fmt"
	"sync"
)

type ProviderEngine struct {
	mu        sync.RWMutex
	factories map[string]ProviderFactory
	active    map[string]Provider
}

func NewProviderEngine() *ProviderEngine {
	return &ProviderEngine{
		factories: make(map[string]ProviderFactory),
		active:    make(map[string]Provider),
	}
}

func (e *ProviderEngine) Register(name string, factory ProviderFactory) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.factories[name] = factory
}

func (e *ProviderEngine) Start(ctx context.Context, name string, config map[string]any) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if _, exists := e.active[name]; exists {
		return fmt.Errorf("provider already active: %s", name)
	}

	factory, exists := e.factories[name]
	if !exists {
		return fmt.Errorf("unknown provider: %s", name)
	}

	provider, err := factory(config)
	if err != nil {
		return fmt.Errorf("create provider %s: %w", name, err)
	}

	if err := provider.Start(ctx); err != nil {
		return fmt.Errorf("start provider %s: %w", name, err)
	}

	e.active[name] = provider
	return nil
}

func (e *ProviderEngine) Stop(ctx context.Context, name string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	provider, exists := e.active[name]
	if !exists {
		return fmt.Errorf("provider not active: %s", name)
	}

	if err := provider.Stop(ctx); err != nil {
		return fmt.Errorf("stop provider %s: %w", name, err)
	}

	delete(e.active, name)
	return nil
}

func (e *ProviderEngine) StopAll(ctx context.Context) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	var errs []error
	for name, provider := range e.active {
		if err := provider.Stop(ctx); err != nil {
			errs = append(errs, fmt.Errorf("stop provider %s: %w", name, err))
		}
		delete(e.active, name)
	}
	if len(errs) > 0 {
		return fmt.Errorf("%d provider(s) failed to stop: %v", len(errs), errs[0])
	}
	return nil
}

func (e *ProviderEngine) Get(ctx context.Context, name string) (Provider, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	provider, exists := e.active[name]
	if !exists {
		return nil, fmt.Errorf("provider not active: %s", name)
	}
	return provider, nil
}

func (e *ProviderEngine) List() []string {
	e.mu.RLock()
	defer e.mu.RUnlock()

	names := make([]string, 0, len(e.active))
	for name := range e.active {
		names = append(names, name)
	}
	return names
}

func (e *ProviderEngine) IsRunning(name string) bool {
	e.mu.RLock()
	defer e.mu.RUnlock()

	provider, exists := e.active[name]
	if !exists {
		return false
	}
	return provider.State() == ProviderRunning
}

func RegisterBuiltins(e *ProviderEngine) {
	e.Register("wireguard", newWireGuardProviderFromConfig)
	e.Register("tor", newTorProvider)
	e.Register("sing-box", newSingBoxProvider)
	e.Register("unbound", newUnboundProvider)
	e.Register("socks5", newSocks5Provider)
}
