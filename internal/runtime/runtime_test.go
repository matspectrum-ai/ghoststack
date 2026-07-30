package runtime

import (
	"context"
	"errors"
	"testing"
)

type recorderHandler struct {
	events []Event
}

func newRecorderHandler() *recorderHandler {
	return &recorderHandler{events: nil}
}

func (r *recorderHandler) Handle(ctx context.Context, event Event) error {
	r.events = append(r.events, event)
	return nil
}

func TestRuntimeStartStop(t *testing.T) {
	handler := newRecorderHandler()
	rt := NewRuntime(handler.Handle)

	if err := rt.Start(context.Background()); err != nil {
		t.Fatalf("start: %v", err)
	}
	if rt.StateNow() != StateRunning {
		t.Fatalf("expected running, got %s", rt.StateNow())
	}

	if err := rt.Stop(context.Background()); err != nil {
		t.Fatalf("stop: %v", err)
	}
	if rt.StateNow() != StateStopped {
		t.Fatalf("expected stopped, got %s", rt.StateNow())
	}

	if len(handler.events) != 3 {
		t.Fatalf("expected 3 events, got %d", len(handler.events))
	}
}

func TestRuntimeDoubleStart(t *testing.T) {
	rt := NewRuntime(nil)
	if err := rt.Start(context.Background()); err != nil {
		t.Fatalf("first start: %v", err)
	}
	if err := rt.Start(context.Background()); !errors.Is(err, ErrAlreadyStarted) {
		t.Fatalf("expected ErrAlreadyStarted, got %v", err)
	}
}

func TestRuntimeStopBeforeStart(t *testing.T) {
	rt := NewRuntime(nil)
	if err := rt.Stop(context.Background()); !errors.Is(err, ErrNotStarted) {
		t.Fatalf("expected ErrNotStarted, got %v", err)
	}
}
