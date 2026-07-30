package providers

import (
	"context"
	"testing"
)

func TestWireGuardProviderLifecycle(t *testing.T) {
	p := newWireGuardProvider("test")

	if p.Name() != "wireguard" {
		t.Fatalf("unexpected name: %s", p.Name())
	}

	if err := p.Start(context.Background(), ""); err != nil {
		t.Fatalf("start: %v", err)
	}

	if err := p.Stop(context.Background()); err != nil {
		t.Fatalf("stop: %v", err)
	}
}

func TestTorProviderLifecycle(t *testing.T) {
	p := newTorProvider()

	if p.Name() != "tor" {
		t.Fatalf("unexpected name: %s", p.Name())
	}

	if err := p.Start(context.Background()); err != nil {
		t.Fatalf("start: %v", err)
	}

	if err := p.Stop(context.Background()); err != nil {
		t.Fatalf("stop: %v", err)
	}
}

func TestSingBoxProviderLifecycle(t *testing.T) {
	p := newSingBoxProvider("test")

	if p.Name() != "sing-box" {
		t.Fatalf("unexpected name: %s", p.Name())
	}

	if err := p.Start(context.Background(), ""); err != nil {
		t.Fatalf("start: %v", err)
	}

	if err := p.Stop(context.Background()); err != nil {
		t.Fatalf("stop: %v", err)
	}
}

func TestUnboundProviderLifecycle(t *testing.T) {
	p := newUnboundProvider()

	if p.Name() != "unbound" {
		t.Fatalf("unexpected name: %s", p.Name())
	}

	if err := p.Start(context.Background()); err != nil {
		t.Fatalf("start: %v", err)
	}

	if err := p.Stop(context.Background()); err != nil {
		t.Fatalf("stop: %v", err)
	}
}

func TestSocks5ProviderLifecycle(t *testing.T) {
	p := newSocks5ProxyProvider("")

	if p.Name() != "socks5" {
		t.Fatalf("unexpected name: %s", p.Name())
	}

	if err := p.Start(context.Background(), ""); err != nil {
		t.Fatalf("start: %v", err)
	}

	if err := p.Stop(context.Background()); err != nil {
		t.Fatalf("stop: %v", err)
	}
}
