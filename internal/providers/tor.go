package providers

import (
	"context"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type torProvider struct {
	mu        sync.RWMutex
	state     ProviderState
	proc      *ProcessManager
	configDir string
	socksPort int
	ctrlPort  int
}

func newTorProvider(config map[string]any) (Provider, error) {
	return &torProvider{
		state:     ProviderStopped,
		proc:      NewProcessManager(),
		socksPort: 9050,
		ctrlPort:  9051,
	}, nil
}

func (p *torProvider) Name() string {
	return "tor"
}

func (p *torProvider) State() ProviderState {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.state
}

func (p *torProvider) Start(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.state == ProviderRunning {
		return nil
	}

	dir, err := os.MkdirTemp("", "ghoststack-tor-*")
	if err != nil {
		return fmt.Errorf("temp dir: %w", err)
	}

	torrc := fmt.Sprintf(`SOCKSPort 127.0.0.1:%d
ControlPort 127.0.0.1:%d
DataDirectory %s
Log notice stdout
`, p.socksPort, p.ctrlPort, dir)

	configPath := filepath.Join(dir, "torrc")
	if err := os.WriteFile(configPath, []byte(torrc), 0600); err != nil {
		os.RemoveAll(dir)
		return fmt.Errorf("write torrc: %w", err)
	}

	if err := p.proc.Start(ctx, ProcessConfig{
		Name: "tor",
		Args: []string{"-f", configPath},
	}); err != nil {
		os.RemoveAll(dir)
		return fmt.Errorf("start tor: %w", err)
	}

	p.configDir = dir
	p.state = ProviderRunning

	if err := p.waitForPort(ctx, "127.0.0.1", p.socksPort, 10*time.Second); err != nil {
		p.proc.Stop()
		os.RemoveAll(dir)
		p.state = ProviderFailed
		return fmt.Errorf("tor not ready: %w", err)
	}

	return nil
}

func (p *torProvider) Stop(ctx context.Context) error {
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

func (p *torProvider) waitForPort(ctx context.Context, host string, port int, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), time.Second)
		if err == nil {
			conn.Close()
			return nil
		}

		time.Sleep(200 * time.Millisecond)
	}
	return fmt.Errorf("timeout waiting for %s:%d", host, port)
}
