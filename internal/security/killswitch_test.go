package security

import (
	"context"
	"testing"
)

func TestKillSwitchLifecycle(t *testing.T) {
	ks := newKillSwitch("wg0")

	if ks.Active() {
		t.Fatal("expected inactive kill switch")
	}

	if ks.Active() {
		t.Fatal("expected inactive initially")
	}

	ks.mu.Lock()
	ks.active = true
	ks.mu.Unlock()

	if !ks.Active() {
		t.Fatal("expected active")
	}

	if err := ks.Disable(context.Background()); err != nil {
		t.Fatalf("disable: %v", err)
	}

	if ks.Active() {
		t.Fatal("expected inactive after disable")
	}
}

func TestKillSwitchEmptyInterface(t *testing.T) {
	ks := newKillSwitch("")

	if err := ks.Enable(context.Background()); err == nil {
		t.Fatal("expected error for empty interface")
	}
}

func TestKillSwitchDoubleEnable(t *testing.T) {
	ks := newKillSwitch("wg0")

	ks.mu.Lock()
	ks.active = true
	ks.mu.Unlock()

	if err := ks.Enable(context.Background()); err != nil {
		t.Fatalf("double enable should no-op: %v", err)
	}
}
