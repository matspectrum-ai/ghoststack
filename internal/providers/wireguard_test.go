package providers

import (
	"context"
	"os"
	"testing"
)

func TestWireGuardProviderConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		config  WireGuardConfig
		wantErr bool
	}{
		{
			name:    "missing private key",
			config:  WireGuardConfig{PublicKey: "pub", Endpoint: "10.0.0.1:51820"},
			wantErr: true,
		},
		{
			name:    "missing public key",
			config:  WireGuardConfig{PrivateKey: "priv", Endpoint: "10.0.0.1:51820"},
			wantErr: true,
		},
		{
			name:    "missing endpoint",
			config:  WireGuardConfig{PrivateKey: "priv", PublicKey: "pub"},
			wantErr: true,
		},
		{
			name: "valid minimal config",
			config: WireGuardConfig{
				PrivateKey: "priv",
				PublicKey:  "pub",
				Endpoint:   "10.0.0.1:51820",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wp := newWireGuardProvider(tt.config).(*wireGuardProvider)
			err := wp.validateConfig()
			if tt.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestWireGuardProviderDefaults(t *testing.T) {
	wp := newWireGuardProvider(WireGuardConfig{
		PrivateKey: "priv",
		PublicKey:  "pub",
		Endpoint:   "10.0.0.1:51820",
	}).(*wireGuardProvider)

	wp.validateConfig()

	p := wp

	if p.config.MTU != 1420 {
		t.Fatalf("expected MTU 1420, got %d", p.config.MTU)
	}

	if p.config.Interface != "ghost0" {
		t.Fatalf("expected interface ghost0, got %s", p.config.Interface)
	}

	if len(p.config.AllowedIPs) == 0 {
		t.Fatal("expected default allowed IPs")
	}
}

func TestWireGuardProviderBuildConfig(t *testing.T) {
	wp := newWireGuardProvider(WireGuardConfig{
		PrivateKey: "testprivkey",
		PublicKey:  "testpubkey",
		Endpoint:   "203.0.113.1:51820",
		Address:    "10.0.0.2/24",
		MTU:        1420,
		DNS:        []string{"1.1.1.1", "1.0.0.1"},
		AllowedIPs: []string{"0.0.0.0/0"},
	}).(*wireGuardProvider)

	cfg := wp.buildConfig()
	if cfg == "" {
		t.Fatal("expected non-empty config")
	}

	if contains(cfg, "testprivkey") == false {
		t.Fatal("config should contain private key")
	}
	if contains(cfg, "testpubkey") == false {
		t.Fatal("config should contain public key")
	}
	if contains(cfg, "203.0.113.1") == false {
		t.Fatal("config should contain endpoint")
	}
	if contains(cfg, "1.1.1.1") == false {
		t.Fatal("config should contain DNS")
	}
}

func TestWireGuardProviderWriteConfig(t *testing.T) {
	wp := newWireGuardProvider(WireGuardConfig{
		PrivateKey: "testprivkey",
		PublicKey:  "testpubkey",
		Endpoint:   "203.0.113.1:51820",
		Address:    "10.0.0.2/24",
	}).(*wireGuardProvider)

	path, err := wp.writeConfig()
	if err != nil {
		t.Fatalf("write config: %v", err)
	}

	if path == "" {
		t.Fatal("expected non-empty path")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read config: %v", err)
	}

	if len(data) == 0 {
		t.Fatal("expected non-empty config file")
	}

	os.Remove(path)
}

func TestWireGuardProviderStartStop(t *testing.T) {
	if os.Geteuid() != 0 {
		t.Skip("requires root for CAP_NET_ADMIN")
	}

	ctx := context.Background()
	wp := newWireGuardProvider(WireGuardConfig{
		PrivateKey: "testprivkey",
		PublicKey:  "testpubkey",
		Endpoint:   "203.0.113.1:51820",
		Address:    "10.0.0.2/24",
	}).(*wireGuardProvider)

	err := wp.StartWithConfig(ctx, "")
	if err != nil {
		t.Fatalf("start: %v", err)
	}

	if !wp.running {
		t.Fatal("expected running after start")
	}

	err = wp.Stop(ctx)
	if err != nil {
		t.Fatalf("stop: %v", err)
	}

	if wp.running {
		t.Fatal("expected stopped after stop")
	}
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
