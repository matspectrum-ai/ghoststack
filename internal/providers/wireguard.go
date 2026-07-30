package providers

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/ghoststack/ghoststack/internal/networking"
	"github.com/ghoststack/ghoststack/internal/security"
)

var (
	ErrProviderNotStarted = fmt.Errorf("provider not started")
)

type WireGuardConfig struct {
	Interface          string   `yaml:"interface"`
	Endpoint           string   `yaml:"endpoint"`
	PrivateKey         string   `yaml:"private_key"`
	PublicKey          string   `yaml:"public_key"`
	AllowedIPs         []string `yaml:"allowed_ips"`
	DNS                []string `yaml:"dns"`
	PersistentKeepalive int     `yaml:"persistent_keepalive"`
	MTU                int      `yaml:"mtu"`
	Address            string   `yaml:"address"`
}

type wireGuardProvider struct {
	mu         sync.RWMutex
	running    bool
	config     WireGuardConfig
	proc       *ProcessManager
	tun        networking.TunDevice
	iface      string
	configPath string
	killSwitch security.KillSwitch
}

func newWireGuardProviderFromConfig(config map[string]any) (Provider, error) {
	cfg := WireGuardConfig{}
	if v, ok := config["interface"]; ok {
		cfg.Interface = fmt.Sprintf("%v", v)
	}
	if v, ok := config["endpoint"]; ok {
		cfg.Endpoint = fmt.Sprintf("%v", v)
	}
	if v, ok := config["private_key"]; ok {
		cfg.PrivateKey = fmt.Sprintf("%v", v)
	}
	if v, ok := config["public_key"]; ok {
		cfg.PublicKey = fmt.Sprintf("%v", v)
	}
	if v, ok := config["address"]; ok {
		cfg.Address = fmt.Sprintf("%v", v)
	}
	if v, ok := config["mtu"]; ok {
		if m, ok := v.(int); ok {
			cfg.MTU = m
		}
	}
	if v, ok := config["persistent_keepalive"]; ok {
		if pk, ok := v.(int); ok {
			cfg.PersistentKeepalive = pk
		}
	}
	return newWireGuardProvider(cfg), nil
}

func newWireGuardProvider(cfg WireGuardConfig) WireGuardProvider {
	return &wireGuardProvider{
		config:     cfg,
		proc:       NewProcessManager(),
		iface:      "ghost0",
		killSwitch: security.NewKillSwitch("ghost0"),
	}
}

func (p *wireGuardProvider) Name() string {
	return "wireguard"
}

func (p *wireGuardProvider) State() ProviderState {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if p.running {
		return ProviderRunning
	}
	return ProviderStopped
}

func (p *wireGuardProvider) Start(ctx context.Context) error {
	return p.StartWithConfig(ctx, "")
}

func (p *wireGuardProvider) StartWithConfig(ctx context.Context, configPath string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.running {
		return nil
	}

	if configPath != "" {
		p.configPath = configPath
	}

	if err := p.validateConfig(); err != nil {
		return fmt.Errorf("wireguard config: %w", err)
	}

	if p.configPath == "" {
		path, err := p.writeConfig()
		if err != nil {
			return fmt.Errorf("write wireguard config: %w", err)
		}
		p.configPath = path
	}

	p.tun = networking.NewTUN()
	if err := p.tun.Create(ctx, p.iface, p.config.MTU); err != nil {
		return fmt.Errorf("create tun: %w", err)
	}

	if p.config.Address != "" {
		if err := p.tun.SetIP(ctx, p.config.Address); err != nil {
			p.tun.Delete(ctx)
			return fmt.Errorf("set tun ip: %w", err)
		}
	}

	if err := p.tun.Up(ctx); err != nil {
		p.tun.Delete(ctx)
		return fmt.Errorf("up tun: %w", err)
	}

	if err := p.proc.Start(ctx, ProcessConfig{
		Name: "wg-quick",
		Args: []string{"up", p.configPath},
	}); err != nil {
		p.tun.Delete(ctx)
		return fmt.Errorf("wg-quick up: %w", err)
	}

	if err := p.killSwitch.Enable(ctx); err != nil {
		p.tun.Delete(ctx)
		return fmt.Errorf("enable kill switch: %w", err)
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

	if p.proc != nil {
		p.proc.Stop()
	}

	if p.configPath != "" {
		exec.CommandContext(ctx, "wg-quick", "down", p.configPath).Run()
	}

	if p.tun != nil {
		p.tun.Down(ctx)
		p.tun.Delete(ctx)
	}

	if p.configPath != "" {
		os.Remove(p.configPath)
	}

	p.killSwitch.Disable(ctx)

	p.running = false
	return nil
}

func (p *wireGuardProvider) Status(ctx context.Context) (map[string]any, error) {
	out, err := exec.CommandContext(ctx, "wg", "show", p.iface, "transfer").Output()
	if err != nil {
		return map[string]any{
			"interface": p.iface,
			"running":   p.running,
			"error":     err.Error(),
		}, nil
	}

	status := map[string]any{
		"interface": p.iface,
		"running":   p.running,
		"transfer":  string(out),
	}

	return status, nil
}

func (p *wireGuardProvider) validateConfig() error {
	if p.config.PrivateKey == "" {
		return fmt.Errorf("private_key is required")
	}
	if p.config.PublicKey == "" {
		return fmt.Errorf("public_key is required")
	}
	if p.config.Endpoint == "" {
		return fmt.Errorf("endpoint is required")
	}
	if len(p.config.AllowedIPs) == 0 {
		p.config.AllowedIPs = []string{"0.0.0.0/0", "::/0"}
	}
	if p.config.MTU <= 0 {
		p.config.MTU = 1420
	}
	if p.config.Interface == "" {
		p.config.Interface = "ghost0"
	}
	return nil
}

func (p *wireGuardProvider) writeConfig() (string, error) {
	dir, err := os.MkdirTemp("", "ghoststack-wg-*")
	if err != nil {
		return "", fmt.Errorf("temp dir: %w", err)
	}

	path := filepath.Join(dir, "wg.conf")
	content := p.buildConfig()
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		os.RemoveAll(dir)
		return "", fmt.Errorf("write config: %w", err)
	}

	return path, nil
}

func (p *wireGuardProvider) buildConfig() string {
	cfg := fmt.Sprintf(`[Interface]
PrivateKey = %s
Address = %s
MTU = %d
`, p.config.PrivateKey, p.config.Address, p.config.MTU)

	if len(p.config.DNS) > 0 {
		cfg += fmt.Sprintf("DNS = %s\n", joinStrings(p.config.DNS, ", "))
	}

	cfg += fmt.Sprintf(`
[Peer]
PublicKey = %s
Endpoint = %s
`, p.config.PublicKey, p.config.Endpoint)

	if p.config.PersistentKeepalive > 0 {
		cfg += fmt.Sprintf("PersistentKeepalive = %d\n", p.config.PersistentKeepalive)
	}

	if len(p.config.AllowedIPs) > 0 {
		cfg += fmt.Sprintf("AllowedIPs = %s\n", joinStrings(p.config.AllowedIPs, ", "))
	}

	return cfg
}

func joinStrings(strs []string, sep string) string {
	result := ""
	for i, s := range strs {
		if i > 0 {
			result += sep
		}
		result += s
	}
	return result
}
