package providers

import (
	"context"
	"testing"
	"time"
)

func TestProcessManagerStartStop(t *testing.T) {
	pm := NewProcessManager()

	if pm.Running() {
		t.Fatal("expected not running initially")
	}

	err := pm.Start(context.Background(), ProcessConfig{
		Name: "sleep",
		Args: []string{"10"},
	})
	if err != nil {
		t.Fatalf("start: %v", err)
	}

	if !pm.Running() {
		t.Fatal("expected running after start")
	}

	if pm.PID() <= 0 {
		t.Fatal("expected valid PID")
	}

	if err := pm.Stop(); err != nil {
		t.Fatalf("stop: %v", err)
	}

	if pm.Running() {
		t.Fatal("expected stopped after stop")
	}
}

func TestProcessManagerStateTransitions(t *testing.T) {
	pm := NewProcessManager()

	if pm.State() != ProcessStopped {
		t.Fatalf("expected stopped, got %s", pm.State())
	}

	pm.Start(context.Background(), ProcessConfig{
		Name: "sleep",
		Args: []string{"5"},
	})

	if pm.State() != ProcessRunning {
		t.Fatalf("expected running, got %s", pm.State())
	}

	pm.Stop()
}

func TestProcessManagerContextCancel(t *testing.T) {
	pm := NewProcessManager()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := pm.Start(ctx, ProcessConfig{
		Name: "sleep",
		Args: []string{"10"},
	})
	if err == nil {
		t.Fatal("expected error with cancelled context")
	}
}

func TestProcessManagerStopNotStarted(t *testing.T) {
	pm := NewProcessManager()
	if err := pm.Stop(); err == nil {
		t.Fatal("expected error stopping not started process")
	}
}

func TestProcessManagerPID(t *testing.T) {
	pm := NewProcessManager()

	if pid := pm.PID(); pid != 0 {
		t.Fatalf("expected 0 before start, got %d", pid)
	}

	pm.Start(context.Background(), ProcessConfig{
		Name: "sleep",
		Args: []string{"5"},
	})

	if pid := pm.PID(); pid <= 0 {
		t.Fatalf("expected valid PID, got %d", pid)
	}

	pm.Stop()

	if pid := pm.PID(); pid != 0 {
		t.Fatalf("expected 0 after stop, got %d", pid)
	}
}

func TestProcessManagerProcessExits(t *testing.T) {
	pm := NewProcessManager()

	err := pm.Start(context.Background(), ProcessConfig{
		Name: "true",
	})
	if err != nil {
		t.Fatalf("start: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	if pm.Running() {
		t.Fatal("expected process to have exited")
	}
}

func TestProcessManagerDoubleStart(t *testing.T) {
	pm := NewProcessManager()

	pm.Start(context.Background(), ProcessConfig{
		Name: "sleep",
		Args: []string{"5"},
	})

	err := pm.Start(context.Background(), ProcessConfig{
		Name: "sleep",
		Args: []string{"5"},
	})
	if err == nil {
		t.Fatal("expected error on double start")
	}

	pm.Stop()
}
