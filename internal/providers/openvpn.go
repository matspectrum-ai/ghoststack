package providers

import (
	"context"
	"sync"
)

type OpenVPNProvider struct {
	mu      sync.RWMutex
	running bool
	config  string
}

func NewOpenVPNProvider(config string) *OpenVPNProvider {
	return &OpenVPNProvider{config: config}
}

func (p *OpenVPNProvider) Name() string {
	return "openvpn"
}

func (p *OpenVPNProvider) Start(ctx context.Context, configPath string) error {
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

func (p *OpenVPNProvider) Stop(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.running {
		return ErrProviderNotStarted
	}

	p.running = false
	return nil
}

func (p *OpenVPNProvider) Status() string {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.running {
		return "running"
	}
	return "stopped"
}
