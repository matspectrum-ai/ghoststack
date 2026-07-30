package providers

import (
	"context"
	"sync"
)

type torProvider struct {
	mu      sync.RWMutex
	running bool
}

func newTorProvider() TorProvider {
	return &torProvider{}
}

func (p *torProvider) Name() string {
	return "tor"
}

func (p *torProvider) Start(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.running {
		return nil
	}

	p.running = true
	return nil
}

func (p *torProvider) Stop(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.running {
		return ErrProviderNotStarted
	}

	p.running = false
	return nil
}
