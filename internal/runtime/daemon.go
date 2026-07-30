package runtime

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ghoststack/ghoststack/internal/storage"
)

type Daemon struct {
	mu        sync.RWMutex
	runtime   *Runtime
	store     storage.StorageProvider
	started   time.Time
	stopped   time.Time
	events    []Event
	maxEvents int
	config    string
}

func NewDaemon(config string, handler EventHandler) *Daemon {
	return &Daemon{
		runtime:   NewRuntime(handler),
		maxEvents: 256,
		config:    config,
	}
}

func (d *Daemon) SetStorage(s storage.StorageProvider) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.store = s
}

func (d *Daemon) Start(ctx context.Context) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if err := d.runtime.Start(ctx); err != nil {
		return err
	}

	d.started = time.Now().UTC()

	if d.store != nil {
		_ = d.store.SaveRuntimeState(ctx, storage.RuntimeState{
			Status:    "running",
			Mode:      d.config,
			StartedAt: d.started.Unix(),
			UpdatedAt: time.Now().Unix(),
		})
	}

	return d.record(ctx, Event{Type: "daemon.started", Source: "daemon"})
}

func (d *Daemon) Stop(ctx context.Context) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if err := d.runtime.Stop(ctx); err != nil {
		return err
	}

	d.stopped = time.Now().UTC()

	if d.store != nil {
		_ = d.store.SaveRuntimeState(ctx, storage.RuntimeState{
			Status:    "stopped",
			Mode:      d.config,
			UpdatedAt: time.Now().Unix(),
		})
	}

	return d.record(ctx, Event{Type: "daemon.stopped", Source: "daemon"})
}

func (d *Daemon) State() State {
	return d.runtime.StateNow()
}

func (d *Daemon) Uptime() time.Duration {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if d.runtime.StateNow() != StateRunning {
		return 0
	}
	return time.Since(d.started)
}

func (d *Daemon) Events() []Event {
	d.mu.RLock()
	defer d.mu.RUnlock()
	out := make([]Event, len(d.events))
	copy(out, d.events)
	return out
}

func (d *Daemon) ConfigString() string {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.config
}

func (d *Daemon) Storage() storage.StorageProvider {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.store
}

func (d *Daemon) record(ctx context.Context, event Event) error {
	if d.maxEvents <= 0 {
		return nil
	}

	event.Timestamp = time.Now().UTC()
	if event.ID == "" {
		event.ID = generateEventID()
	}

	d.events = append(d.events, event)
	if len(d.events) > d.maxEvents {
		d.events = d.events[len(d.events)-d.maxEvents:]
	}

	if d.store != nil {
		_ = d.store.AppendAuditLog(ctx, storage.AuditEntry{
			Timestamp: time.Now().Unix(),
			Action:    event.Type,
			Source:    event.Source,
			Detail:    fmt.Sprintf("%v", event.Payload),
		})
	}

	return d.runtime.emit(ctx, event)
}

func (d *Daemon) String() string {
	d.mu.RLock()
	defer d.mu.RUnlock()

	switch d.runtime.StateNow() {
	case StateRunning:
		return fmt.Sprintf("GhostStack running since %s", d.started.Format(time.RFC3339))
	case StateStopped:
		return fmt.Sprintf("GhostStack stopped at %s", d.stopped.Format(time.RFC3339))
	case StateStopping:
		return "GhostStack stopping"
	case StateFailed:
		return "GhostStack failed"
	default:
		return "GhostStack idle"
	}
}
