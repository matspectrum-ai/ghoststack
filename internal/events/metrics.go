package events

import "context"

type MetricsRecorder interface {
	Increment(eventType string)
	ObserveLatency(eventType string, latencyMs float64)
}

type NoopMetrics struct{}

func (NoopMetrics) Increment(_ string)          {}
func (NoopMetrics) ObserveLatency(_ string, _ float64) {}

type EventMetrics struct {
	recorder MetricsRecorder
}

func NewEventMetrics(recorder MetricsRecorder) *EventMetrics {
	if recorder == nil {
		recorder = NoopMetrics{}
	}
	return &EventMetrics{recorder: recorder}
}

func (m *EventMetrics) Middleware() BusMiddleware {
	return func(ctx context.Context, event Event) error {
		m.recorder.Increment(event.Type)
		return nil
	}
}

func EventBusWithMetrics(bus *EventBus, recorder MetricsRecorder) {
	metrics := NewEventMetrics(recorder)
	bus.SetMiddleware(metrics.Middleware())
}
