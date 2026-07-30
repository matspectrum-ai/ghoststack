package runtime

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	ErrAlreadyStarted = errors.New("runtime already started")
	ErrNotStarted     = errors.New("runtime not started")
)

type State string

const (
	StateIdle    State = "idle"
	StateRunning State = "running"
	StateStopping State = "stopping"
	StateStopped State = "stopped"
	StateFailed  State = "failed"
)

type Event struct {
	ID        string
	Type      string
	Source    string
	Timestamp time.Time
	Payload   map[string]any
}

type EventHandler func(ctx context.Context, event Event) error

type Runtime struct {
	mu      sync.RWMutex
	state   State
	started time.Time
	handler EventHandler
}

func NewRuntime(handler EventHandler) *Runtime {
	if handler == nil {
		handler = func(ctx context.Context, event Event) error {
			return nil
		}
	}
	return &Runtime{
		state:   StateIdle,
		handler: handler,
	}
}

func (r *Runtime) Start(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.state == StateRunning {
		return ErrAlreadyStarted
	}

	r.state = StateRunning
	r.started = time.Now().UTC()
	return r.emit(ctx, Event{Type: "runtime.started", Source: "runtime"})
}

func (r *Runtime) Stop(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.state != StateRunning {
		return ErrNotStarted
	}

	r.state = StateStopping
	if err := r.emit(ctx, Event{Type: "runtime.stopping", Source: "runtime"}); err != nil {
		r.state = StateFailed
		return err
	}

	r.state = StateStopped
	return r.emit(ctx, Event{Type: "runtime.stopped", Source: "runtime"})
}

func (r *Runtime) StateNow() State {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.state
}

func (r *Runtime) emit(ctx context.Context, event Event) error {
	if r.handler == nil {
		return nil
	}

	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now().UTC()
	}
	if event.ID == "" {
		event.ID = generateEventID()
	}

	return r.handler(ctx, event)
}

func generateEventID() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
