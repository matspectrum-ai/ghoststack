package profiling

import (
	"context"
	"testing"
)

func TestProfilerEnableDisable(t *testing.T) {
	p := NewProfiler("/tmp/ghost.prof")
	if err := p.Enable(context.Background()); err != nil {
		t.Fatalf("enable: %v", err)
	}
	if err := p.Enable(context.Background()); err != nil {
		t.Fatalf("enable idempotent: %v", err)
	}
	if err := p.Disable(context.Background()); err != nil {
		t.Fatalf("disable: %v", err)
	}
	if err := p.Disable(context.Background()); err != nil {
		t.Fatalf("disable idempotent: %v", err)
	}
}

func TestProfilerSnapshot(t *testing.T) {
	p := NewProfiler("/tmp/ghost.prof")
	if err := p.Enable(context.Background()); err != nil {
		t.Fatalf("enable: %v", err)
	}

	snap, err := p.Snapshot()
	if err != nil {
		t.Fatalf("snapshot: %v", err)
	}
	if snap.Timestamp == 0 {
		t.Fatal("expected non-zero timestamp")
	}
	if snap.AllocMB < 0 {
		t.Fatal("expected non-negative alloc")
	}
}

func TestProfilerEmptyOutput(t *testing.T) {
	p := NewProfiler("")
	if err := p.Enable(context.Background()); err == nil {
		t.Fatal("expected error for empty output")
	}
}

func TestMemoryPoolStats(t *testing.T) {
	pool := newMemoryPool(1024)
	stats := pool.Stats()
	if stats.AllocatedBytes() != 1024 {
		t.Fatalf("expected 1024, got %d", stats.AllocatedBytes())
	}
	pool.Reset()
	stats = pool.Stats()
	if stats.AllocatedBytes() != 0 {
		t.Fatalf("expected 0 after reset, got %d", stats.AllocatedBytes())
	}
}
