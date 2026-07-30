package events

import (
	"context"
	"testing"
)

func TestMetricsMiddleware(t *testing.T) {
	var count int
	recorder := &testMetricsRecorder{increment: func(eventType string) {
		count++
	}}
	bus := NewEventBus(10)
	EventBusWithMetrics(bus, recorder)
	_ = bus.Publish(context.Background(), Event{Type: "test"})
	if count != 1 {
		t.Fatalf("expected 1 metric, got %d", count)
	}
}

type testMetricsRecorder struct {
	increment func(string)
}

func (r *testMetricsRecorder) Increment(eventType string) {
	if r.increment != nil {
		r.increment(eventType)
	}
}

func (r *testMetricsRecorder) ObserveLatency(_ string, _ float64) {}
