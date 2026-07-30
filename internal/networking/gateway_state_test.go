package networking

import (
	"errors"
	"testing"
)

func TestGatewayStateTransitions(t *testing.T) {
	g := newGatewayState()
	running, _, started, _ := g.snapshot()
	if running {
		t.Fatal("expected idle gateway")
	}
	if started != 0 {
		t.Fatal("expected zero start time")
	}

	g.markStarted("test")
	var config string
	running, config, started, _ = g.snapshot()
	if !running {
		t.Fatal("expected running gateway")
	}
	if config != "test" {
		t.Fatal("expected preserved config")
	}
	if started == 0 {
		t.Fatal("expected non-zero start time")
	}

	g.markStopped()
	running, _, _, stopped := g.snapshot()
	if running {
		t.Fatal("expected stopped gateway")
	}
	if stopped == 0 {
		t.Fatal("expected non-zero stop time")
	}
}

func TestGatewayStateErrors(t *testing.T) {
	g := newGatewayState()
	g.addError(nil)
	g.addError(errors.New("boom"))

	if len(g.errors) != 1 || g.errors[0] != "boom" {
		t.Fatalf("unexpected errors: %v", g.errors)
	}
}
