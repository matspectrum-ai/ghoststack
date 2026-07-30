package security

import (
	"context"
	"errors"
	"testing"
)

func TestKillSwitchLifecycle(t *testing.T) {
	triggered := false
	ks := newKillSwitch(func() error {
		triggered = true
		return nil
	})

	if ks.Active() {
		t.Fatal("expected inactive kill switch")
	}

	if err := ks.Enable(context.Background()); err != nil {
		t.Fatalf("enable: %v", err)
	}
	if !ks.Active() {
		t.Fatal("expected active kill switch")
	}
	if !triggered {
		t.Fatal("expected trigger to fire")
	}

	if err := ks.Enable(context.Background()); err != nil {
		t.Fatalf("double enable: %v", err)
	}

	if err := ks.Disable(context.Background()); err != nil {
		t.Fatalf("disable: %v", err)
	}
	if ks.Active() {
		t.Fatal("expected inactive kill switch")
	}
}

func TestKillSwitchTriggerFailure(t *testing.T) {
	ks := newKillSwitch(func() error {
		return errors.New("trigger failed")
	})

	if err := ks.Enable(context.Background()); err == nil {
		t.Fatal("expected error from trigger")
	}
	if ks.Active() {
		t.Fatal("expected inactive kill switch after failed trigger")
	}
}
