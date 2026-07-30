package events

import (
	"context"
	"testing"
)

func TestEventBusPublishSubscribe(t *testing.T) {
	bus := NewEventBus(10)
	received := false

	bus.Subscribe("test", func(ctx context.Context, event Event) error {
		received = true
		return nil
	})

	if err := bus.Publish(context.Background(), Event{Type: "test"}); err != nil {
		t.Fatalf("publish: %v", err)
	}
	if !received {
		t.Fatal("expected handler to be called")
	}
}

func TestEventBusEmptyType(t *testing.T) {
	bus := NewEventBus(10)
	if err := bus.Publish(context.Background(), Event{}); err == nil {
		t.Fatal("expected error for empty event type")
	}
}

func TestEventBusBuffer(t *testing.T) {
	bus := NewEventBus(2)

	_ = bus.Publish(context.Background(), Event{Type: "a"})
	_ = bus.Publish(context.Background(), Event{Type: "b"})
	_ = bus.Publish(context.Background(), Event{Type: "c"})

	buffer := bus.Buffer()
	if len(buffer) != 2 {
		t.Fatalf("expected buffer size 2, got %d", len(buffer))
	}
}
