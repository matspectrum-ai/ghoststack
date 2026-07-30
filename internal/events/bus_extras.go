package events

import "context"

var ErrEmptyEventType = &eventError{"event type must not be empty"}

type eventError struct {
	msg string
}

func (e *eventError) Error() string {
	return e.msg
}

func (b *EventBus) SubscribeWildcard(pattern string, handler EventHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.wildcards = append(b.wildcards, wildcardEntry{pattern: pattern, handler: handler})
}

func (b *EventBus) PublishWithCorrelation(ctx context.Context, eventType, source, correlationID string, payload map[string]interface{}) error {
	if correlationID == "" {
		correlationID = CorrelationID(ctx)
	}
	if correlationID != "" {
		ctx = WithCorrelationID(ctx, correlationID)
	}
	if payload == nil {
		payload = map[string]interface{}{}
	}
	if source != "" {
		payload["source"] = source
	}
	if correlationID != "" {
		payload["correlationId"] = correlationID
	}
	return b.Publish(ctx, Event{Type: eventType, Source: source, Payload: payload})
}
