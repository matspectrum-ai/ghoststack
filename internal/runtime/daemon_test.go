package runtime

import (
	"context"
	"errors"
	"testing"
)

type eventRecorder struct {
	events []Event
}

func newEventRecorder() *eventRecorder {
	return &eventRecorder{}
}

func (r *eventRecorder) Handle(ctx context.Context, event Event) error {
	r.events = append(r.events, event)
	return nil
}

func TestDaemonStartStopLifecycle(t *testing.T) {
	rec := newEventRecorder()
	d := NewDaemon("", rec.Handle)

	if err := d.Start(context.Background()); err != nil {
		t.Fatalf("start: %v", err)
	}
	if d.State() != StateRunning {
		t.Fatalf("expected running, got %s", d.State())
	}

	if d.Uptime() <= 0 {
		t.Fatalf("expected positive uptime, got %v", d.Uptime())
	}

	if err := d.Stop(context.Background()); err != nil {
		t.Fatalf("stop: %v", err)
	}
	if d.State() != StateStopped {
		t.Fatalf("expected stopped, got %s", d.State())
	}

	if len(rec.events) < 3 {
		t.Fatalf("expected at least 3 events, got %d", len(rec.events))
	}
}

func TestDaemonDoubleStart(t *testing.T) {
	d := NewDaemon("", nil)
	if err := d.Start(context.Background()); err != nil {
		t.Fatalf("first start: %v", err)
	}
	if err := d.Start(context.Background()); !errors.Is(err, ErrAlreadyStarted) {
		t.Fatalf("expected ErrAlreadyStarted, got %v", err)
	}
}

func TestDaemonStringStates(t *testing.T) {
	d := NewDaemon("", nil)
	if d.State() != StateIdle {
		t.Fatalf("expected idle, got %s", d.State())
	}
	if d.String() != "GhostStack idle" {
		t.Fatalf("unexpected string: %s", d.String())
	}
}
