package providers

import (
	"context"
	"sync"
)

type DNSOverHTTPSProvider struct {
	mu       sync.RWMutex
	running  bool
	endpoint string
}

func NewDNSOverHTTPSProvider(endpoint string) *DNSOverHTTPSProvider {
	if endpoint == "" {
		endpoint = "https://dns.google/dns-query"
	}
	return &DNSOverHTTPSProvider{endpoint: endpoint}
}

func (p *DNSOverHTTPSProvider) Name() string {
	return "dns-over-https"
}

func (p *DNSOverHTTPSProvider) Start(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.running {
		return nil
	}

	p.running = true
	return nil
}

func (p *DNSOverHTTPSProvider) Stop(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.running {
		return ErrProviderNotStarted
	}

	p.running = false
	return nil
}

func (p *DNSOverHTTPSProvider) Status() string {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.running {
		return "running"
	}
	return "stopped"
}
