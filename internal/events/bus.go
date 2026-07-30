package events

import (
	"context"
	"reflect"
	"sync"
	"time"
)

type Event struct {
	Type      string                 `json:"type"`
	Source    string                 `json:"source,omitempty"`
	Timestamp int64                  `json:"timestamp"`
	Payload   map[string]interface{} `json:"payload,omitempty"`
}

type EventHandler func(ctx context.Context, event Event) error

type wildcardEntry struct {
	pattern string
	handler EventHandler
}

type EventBus struct {
	mu         sync.RWMutex
	handlers   map[string][]EventHandler
	wildcards  []wildcardEntry
	buffer     []Event
	maxBuffer  int
	middleware BusMiddleware
}

func NewEventBus(maxBuffer int) *EventBus {
	if maxBuffer <= 0 {
		maxBuffer = 256
	}
	return &EventBus{
		handlers:  make(map[string][]EventHandler),
		buffer:    make([]Event, 0, maxBuffer),
		maxBuffer: maxBuffer,
	}
}

func (b *EventBus) Publish(ctx context.Context, event Event) error {
	if event.Type == "" {
		return ErrEmptyEventType
	}

	if b.middleware != nil {
		if err := b.middleware(ctx, event); err != nil {
			return err
		}
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	if event.Timestamp == 0 {
		event.Timestamp = time.Now().UnixNano()
	}

	b.buffer = append(b.buffer, event)
	if len(b.buffer) > b.maxBuffer {
		b.buffer = b.buffer[len(b.buffer)-b.maxBuffer:]
	}

	for _, handler := range b.handlers[event.Type] {
		if err := handler(ctx, event); err != nil {
			return err
		}
	}

	for _, entry := range b.wildcards {
		if entry.handler != nil {
			_ = entry.handler(ctx, event)
		}
	}

	return nil
}

func (b *EventBus) Subscribe(eventType string, handler EventHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[eventType] = append(b.handlers[eventType], handler)
}

func (b *EventBus) Unsubscribe(eventType string, handler EventHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()

	handlers := b.handlers[eventType]
	for i, h := range handlers {
		if funcPtr(h) == funcPtr(handler) {
			b.handlers[eventType] = append(handlers[:i], handlers[i+1:]...)
			break
		}
	}
}

func funcPtr(f EventHandler) uintptr {
	v := reflect.ValueOf(f)
	if v.IsNil() {
		return 0
	}
	return v.Pointer()
}

func (b *EventBus) Buffer() []Event {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return append([]Event(nil), b.buffer...)
}

func (b *EventBus) SetMiddleware(middleware BusMiddleware) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.middleware = middleware
}
