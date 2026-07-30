package events

import (
	"context"
	"fmt"
	"sync"
)

type Event struct {
	Type      string
	Source    string
	Timestamp int64
	Payload   map[string]interface{}
}

type EventHandler func(ctx context.Context, event Event) error

type EventBus struct {
	mu         sync.RWMutex
	handlers   map[string][]EventHandler
	buffer     []Event
	maxBuffer  int
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
		return fmt.Errorf("event type must not be empty")
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	if event.Timestamp == 0 {
		event.Timestamp = nowUnix()
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
		if &h == &handler {
			b.handlers[eventType] = append(handlers[:i], handlers[i+1:]...)
			break
		}
	}
}

func (b *EventBus) Buffer() []Event {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return append([]Event(nil), b.buffer...)
}

func nowUnix() int64 {
	return 0
}
