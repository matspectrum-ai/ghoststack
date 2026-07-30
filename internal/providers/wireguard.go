package providers

import (
	"context"
	"fmt"
	"sync"
)

var (
	ErrProviderNotStarted = fmt.Errorf("provider not started")
)

type wireGuardProvider struct {
	mu       sync.RWMutex
	running  bool
	config   string
}

func newWireGuardProvider(config string) WireGuardProvider {
	return &wireGuardProvider{config: config}
}

func (p *wireGuardProvider) Name() string {
	return "wireguard"
}

func (p *wireGuardProvider) Start(ctx context.Context, configPath string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.running {
		return nil
	}

	if configPath != "" {
		p.config = configPath
	}

	p.running = true
	return nil
}

func (p *wireGuardProvider) Stop(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.running {
		return ErrProviderNotStarted
	}

	p.running = false
	return nil
}
