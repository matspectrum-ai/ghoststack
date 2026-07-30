package providers

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type unboundProvider struct {
	mu        sync.RWMutex
	state     ProviderState
	proc      *ProcessManager
	configDir string
	dnsPort   int
}

func newUnboundProvider(config map[string]any) (Provider, error) {
	return &unboundProvider{
		state:   ProviderStopped,
		proc:    NewProcessManager(),
		dnsPort: 53,
	}, nil
}

func (p *unboundProvider) Name() string {
	return "unbound"
}

func (p *unboundProvider) State() ProviderState {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.state
}

func (p *unboundProvider) Start(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.state == ProviderRunning {
		return nil
	}

	dir, err := os.MkdirTemp("", "ghoststack-unbound-*")
	if err != nil {
		return fmt.Errorf("temp dir: %w", err)
	}

	conf := fmt.Sprintf(`server:
    interface: 127.0.0.1
    port: %d
    access-control: 127.0.0.0/8 allow
    access-control: 10.0.0.0/8 allow
    do-ip4: yes
    do-ip6: yes
    do-udp: yes
    do-tcp: yes
    hide-identity: yes
    hide-version: yes
    prefetch: yes
    cache-min-ttl: 300
    cache-max-ttl: 86400
`, p.dnsPort)

	configPath := filepath.Join(dir, "unbound.conf")
	if err := os.WriteFile(configPath, []byte(conf), 0600); err != nil {
		os.RemoveAll(dir)
		return fmt.Errorf("write unbound config: %w", err)
	}

	if err := p.proc.Start(ctx, ProcessConfig{
		Name: "unbound",
		Args: []string{"-d", "-c", configPath},
	}); err != nil {
		os.RemoveAll(dir)
		return fmt.Errorf("start unbound: %w", err)
	}

	p.configDir = dir
	p.state = ProviderRunning
	return nil
}

func (p *unboundProvider) Stop(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.state != ProviderRunning {
		return ErrProviderNotStarted
	}

	p.proc.Stop()

	if p.configDir != "" {
		os.RemoveAll(p.configDir)
	}

	p.state = ProviderStopped
	return nil
}
