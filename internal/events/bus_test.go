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

func TestEventBusWildcard(t *testing.T) {
	bus := NewEventBus(10)
	received := false

	bus.SubscribeWildcard("", func(ctx context.Context, event Event) error {
		received = true
		return nil
	})

	if err := bus.Publish(context.Background(), Event{Type: "anything"}); err != nil {
		t.Fatalf("publish: %v", err)
	}
	if !received {
		t.Fatal("expected wildcard handler to be called")
	}
}

func TestEventBusPublishWithCorrelation(t *testing.T) {
	bus := NewEventBus(10)
	var capturedCorrelation string

	bus.Subscribe("test", func(ctx context.Context, event Event) error {
		capturedCorrelation = CorrelationID(ctx)
		return nil
	})

	ctx := WithCorrelationID(context.Background(), "abc-123")
	if err := bus.PublishWithCorrelation(ctx, "test", "src", "abc-123", nil); err != nil {
		t.Fatalf("publish with correlation: %v", err)
	}
	if capturedCorrelation != "abc-123" {
		t.Fatalf("expected correlation abc-123, got %s", capturedCorrelation)
	}
}

func TestEventBusMiddleware(t *testing.T) {
	bus := NewEventBus(10)
	middlewareCalled := false

	bus.SetMiddleware(func(ctx context.Context, event Event) error {
		middlewareCalled = true
		return nil
	})

	bus.Subscribe("test", func(ctx context.Context, event Event) error {
		return nil
	})

	if err := bus.Publish(context.Background(), Event{Type: "test"}); err != nil {
		t.Fatalf("publish: %v", err)
	}
	if !middlewareCalled {
		t.Fatal("expected middleware to be called")
	}
}

func TestEventBusUnsubscribe(t *testing.T) {
	bus := NewEventBus(10)
	called := true

	handler := func(ctx context.Context, event Event) error {
		called = false
		return nil
	}
	bus.Subscribe("test", handler)
	bus.Unsubscribe("test", handler)

	if err := bus.Publish(context.Background(), Event{Type: "test"}); err != nil {
		t.Fatalf("publish: %v", err)
	}
	if !called {
		t.Fatal("expected handler to be unsubscribed")
	}
}

func TestCorrelationIDContext(t *testing.T) {
	ctx := WithCorrelationID(context.Background(), "corr-1")
	if got := CorrelationID(ctx); got != "corr-1" {
		t.Fatalf("expected corr-1, got %s", got)
	}
	if got := CorrelationID(context.Background()); got != "" {
		t.Fatalf("expected empty correlation, got %s", got)
	}
}

func TestEventBusBufferOverflow(t *testing.T) {
	bus := NewEventBus(2)
	for i := 0; i < 5; i++ {
		_ = bus.Publish(context.Background(), Event{Type: "x"})
	}
	buffer := bus.Buffer()
	if len(buffer) != 2 {
		t.Fatalf("expected buffer size 2, got %d", len(buffer))
	}
}

func TestEventBusMultipleHandlers(t *testing.T) {
	bus := NewEventBus(10)
	count := 0
	bus.Subscribe("test", func(ctx context.Context, event Event) error {
		count++
		return nil
	})
	bus.Subscribe("test", func(ctx context.Context, event Event) error {
		count++
		return nil
	})
	_ = bus.Publish(context.Background(), Event{Type: "test"})
	if count != 2 {
		t.Fatalf("expected 2 handlers called, got %d", count)
	}
}

func TestPublishWithCorrelationInjectsID(t *testing.T) {
	bus := NewEventBus(10)
	var captured string
	bus.Subscribe("test", func(ctx context.Context, event Event) error {
		if v, ok := event.Payload["correlationId"]; ok {
			captured = v.(string)
		}
		return nil
	})
	_ = bus.PublishWithCorrelation(context.Background(), "test", "src", "", nil)
	if captured != "" {
		t.Fatalf("expected empty correlation when none provided, got %s", captured)
	}
}

func TestEventTypeValidationError(t *testing.T) {
	bus := NewEventBus(10)
	err := bus.Publish(context.Background(), Event{})
	if err == nil {
		t.Fatal("expected error for empty event type")
	}
}
