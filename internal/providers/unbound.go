package providers

import (
	"context"
	"sync"
)

type unboundProvider struct {
	mu      sync.RWMutex
	running bool
}

func newUnboundProvider() UnboundProvider {
	return &unboundProvider{}
}

func (p *unboundProvider) Name() string {
	return "unbound"
}

func (p *unboundProvider) Start(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.running {
		return nil
	}

	p.running = true
	return nil
}

func (p *unboundProvider) Stop(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.running {
		return ErrProviderNotStarted
	}

	p.running = false
	return nil
}
