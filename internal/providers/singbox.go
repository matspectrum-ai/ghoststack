package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type singBoxProvider struct {
	mu        sync.RWMutex
	state     ProviderState
	proc      *ProcessManager
	configDir string
	addr      string
	port      int
}

func newSingBoxProvider(config map[string]any) (Provider, error) {
	return &singBoxProvider{
		state: ProviderStopped,
		proc:  NewProcessManager(),
		addr:  "0.0.0.0",
		port:  1080,
	}, nil
}

func (p *singBoxProvider) Name() string {
	return "sing-box"
}

func (p *singBoxProvider) State() ProviderState {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.state
}

func (p *singBoxProvider) Start(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.state == ProviderRunning {
		return nil
	}

	dir, err := os.MkdirTemp("", "ghoststack-singbox-*")
	if err != nil {
		return fmt.Errorf("temp dir: %w", err)
	}

	cfg := map[string]any{
		"log": map[string]any{
			"level":  "info",
			"output": filepath.Join(dir, "sing-box.log"),
		},
		"inbounds": []map[string]any{
			{
				"type":        "socks",
				"tag":         "socks-in",
				"listen":      p.addr,
				"listen_port": p.port,
			},
		},
	}

	configPath := filepath.Join(dir, "config.json")
	data, _ := json.MarshalIndent(cfg, "", "  ")
	if err := os.WriteFile(configPath, data, 0600); err != nil {
		os.RemoveAll(dir)
		return fmt.Errorf("write sing-box config: %w", err)
	}

	if err := p.proc.Start(ctx, ProcessConfig{
		Name: "sing-box",
		Args: []string{"run", "-c", configPath},
	}); err != nil {
		os.RemoveAll(dir)
		return fmt.Errorf("start sing-box: %w", err)
	}

	p.configDir = dir
	p.state = ProviderRunning
	return nil
}

func (p *singBoxProvider) Stop(ctx context.Context) error {
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
