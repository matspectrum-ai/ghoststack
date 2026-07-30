package networking

import (
	"context"
	"errors"
	"testing"
)

func TestTUNLifecycle(t *testing.T) {
	tun := newTUN()

	if err := tun.Up(context.Background()); !errors.Is(err, ErrTUNNotCreated) {
		t.Fatalf("expected ErrTUNNotCreated, got %v", err)
	}

	if err := tun.Create(context.Background(), "gtun0", 0); err != nil {
		t.Fatalf("create: %v", err)
	}

	if err := tun.Up(context.Background()); err != nil {
		t.Fatalf("up: %v", err)
	}

	addrs, err := tun.Addresses(context.Background())
	if err != nil {
		t.Fatalf("addresses: %v", err)
	}
	if len(addrs) != 0 {
		t.Fatalf("expected empty addresses, got %v", addrs)
	}

	if err := tun.Down(context.Background()); err != nil {
		t.Fatalf("down: %v", err)
	}
}
