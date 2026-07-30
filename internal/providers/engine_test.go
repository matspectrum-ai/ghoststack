package providers

import (
	"context"
	"testing"
)

func TestProviderEngineRegister(t *testing.T) {
	e := NewProviderEngine()

	e.Register("socks5", newSocks5Provider)

	if !e.IsRunning("socks5") {
		return
	}
}

func TestProviderEngineStartStop(t *testing.T) {
	e := NewProviderEngine()

	e.Register("socks5", func(config map[string]any) (Provider, error) {
		return newSocks5Provider(config)
	})

	ctx := context.Background()
	err := e.Start(ctx, "socks5", map[string]any{
		"listen": "127.0.0.1:0",
	})
	if err != nil {
		t.Fatalf("start: %v", err)
	}

	if !e.IsRunning("socks5") {
		t.Fatal("expected running")
	}

	err = e.Stop(ctx, "socks5")
	if err != nil {
		t.Fatalf("stop: %v", err)
	}
}

func TestProviderEngineUnknown(t *testing.T) {
	e := NewProviderEngine()

	err := e.Start(context.Background(), "nonexistent", nil)
	if err == nil {
		t.Fatal("expected error for unknown provider")
	}
}

func TestProviderEngineDoubleStart(t *testing.T) {
	e := NewProviderEngine()
	e.Register("socks5", func(config map[string]any) (Provider, error) {
		return newSocks5Provider(config)
	})

	ctx := context.Background()
	e.Start(ctx, "socks5", map[string]any{"listen": "127.0.0.1:0"})

	err := e.Start(ctx, "socks5", nil)
	if err == nil {
		t.Fatal("expected error for double start")
	}

	e.Stop(ctx, "socks5")
}

func TestProviderEngineStopAll(t *testing.T) {
	e := NewProviderEngine()
	e.Register("socks5", func(config map[string]any) (Provider, error) {
		return newSocks5Provider(config)
	})

	ctx := context.Background()
	e.Start(ctx, "socks5", map[string]any{"listen": "127.0.0.1:0"})

	if err := e.StopAll(ctx); err != nil {
		t.Fatalf("stop all: %v", err)
	}
}

func TestProviderEngineList(t *testing.T) {
	e := NewProviderEngine()
	e.Register("socks5", func(config map[string]any) (Provider, error) {
		return newSocks5Provider(config)
	})

	ctx := context.Background()
	e.Start(ctx, "socks5", map[string]any{"listen": "127.0.0.1:0"})

	names := e.List()
	if len(names) != 1 || names[0] != "socks5" {
		t.Fatalf("expected [socks5], got %v", names)
	}

	e.Stop(ctx, "socks5")
}
