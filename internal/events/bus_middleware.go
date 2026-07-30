package events

import "context"

type BusMiddleware func(ctx context.Context, event Event) error

type EventBusOption func(*EventBus)

func WithMiddleware(middleware BusMiddleware) EventBusOption {
	return func(b *EventBus) {
		b.middleware = middleware
	}
}
