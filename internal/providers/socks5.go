package providers

import (
	"context"
	"sync"
)

type socks5ProxyProvider struct {
	mu      sync.RWMutex
	running bool
	listen  string
}

func newSocks5ProxyProvider(listen string) Socks5ProxyProvider {
	if listen == "" {
		listen = "127.0.0.1:1080"
	}
	return &socks5ProxyProvider{listen: listen}
}

func (p *socks5ProxyProvider) Name() string {
	return "socks5"
}

func (p *socks5ProxyProvider) Start(ctx context.Context, listen string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.running {
		return nil
	}

	if listen != "" {
		p.listen = listen
	}

	p.running = true
	return nil
}

func (p *socks5ProxyProvider) Stop(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.running {
		return ErrProviderNotStarted
	}

	p.running = false
	return nil
}
