package providers

import (
	"context"
	"sync"
)

type singBoxProvider struct {
	mu       sync.RWMutex
	running  bool
	config   string
}

func newSingBoxProvider(config string) SingBoxProvider {
	return &singBoxProvider{config: config}
}

func (p *singBoxProvider) Name() string {
	return "sing-box"
}

func (p *singBoxProvider) Start(ctx context.Context, configPath string) error {
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

func (p *singBoxProvider) Stop(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.running {
		return ErrProviderNotStarted
	}

	p.running = false
	return nil
}
