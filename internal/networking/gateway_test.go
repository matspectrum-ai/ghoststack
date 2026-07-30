package networking

import (
	"context"
	"errors"
	"testing"
)

func TestGatewayStartStop(t *testing.T) {
	g := NewGateway("test")

	if g.Status() != "stopped" {
		t.Fatalf("expected stopped, got %s", g.Status())
	}

	if err := g.Start(context.Background()); err != nil {
		t.Fatalf("start: %v", err)
	}
	if g.Status() != "running" {
		t.Fatalf("expected running, got %s", g.Status())
	}

	if err := g.Start(context.Background()); !errors.Is(err, ErrGatewayAlreadyStarted) {
		t.Fatalf("expected ErrGatewayAlreadyStarted, got %v", err)
	}

	if err := g.Stop(context.Background()); err != nil {
		t.Fatalf("stop: %v", err)
	}
	if g.Status() != "stopped" {
		t.Fatalf("expected stopped, got %s", g.Status())
	}

	if err := g.Stop(context.Background()); !errors.Is(err, ErrGatewayNotStarted) {
		t.Fatalf("expected ErrGatewayNotStarted, got %v", err)
	}
}
