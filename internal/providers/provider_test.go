package providers

import (
	"context"
	"testing"
)

func TestWireGuardProviderName(t *testing.T) {
	p := newWireGuardProvider(WireGuardConfig{
		PrivateKey: "priv",
		PublicKey:  "pub",
		Endpoint:   "10.0.0.1:51820",
	})

	if p.Name() != "wireguard" {
		t.Fatalf("unexpected name: %s", p.Name())
	}

	if p.State() != ProviderStopped {
		t.Fatalf("expected stopped, got %s", p.State())
	}
}

func TestTorProviderName(t *testing.T) {
	p, err := newTorProvider(nil)
	if err != nil {
		t.Fatalf("new: %v", err)
	}

	if p.Name() != "tor" {
		t.Fatalf("unexpected name: %s", p.Name())
	}
}

func TestSingBoxProviderName(t *testing.T) {
	p, err := newSingBoxProvider(nil)
	if err != nil {
		t.Fatalf("new: %v", err)
	}

	if p.Name() != "sing-box" {
		t.Fatalf("unexpected name: %s", p.Name())
	}
}

func TestUnboundProviderName(t *testing.T) {
	p, err := newUnboundProvider(nil)
	if err != nil {
		t.Fatalf("new: %v", err)
	}

	if p.Name() != "unbound" {
		t.Fatalf("unexpected name: %s", p.Name())
	}
}

func TestSocks5ProviderName(t *testing.T) {
	p, err := newSocks5Provider(map[string]any{
		"listen": "127.0.0.1:0",
	})
	if err != nil {
		t.Fatalf("new: %v", err)
	}

	if p.Name() != "socks5" {
		t.Fatalf("unexpected name: %s", p.Name())
	}
}

func TestSocks5ProviderLifecycle(t *testing.T) {
	p, err := newSocks5Provider(map[string]any{
		"listen": "127.0.0.1:0",
	})
	if err != nil {
		t.Fatalf("new: %v", err)
	}

	if err := p.Start(context.Background()); err != nil {
		t.Fatalf("start: %v", err)
	}

	if p.State() != ProviderRunning {
		t.Fatalf("expected running, got %s", p.State())
	}

	if err := p.Stop(context.Background()); err != nil {
		t.Fatalf("stop: %v", err)
	}

	if p.State() != ProviderStopped {
		t.Fatalf("expected stopped, got %s", p.State())
	}
}
