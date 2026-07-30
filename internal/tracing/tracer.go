package tracing

import "context"

type Tracer interface {
	StartSpan(ctx context.Context, name string) (context.Context, Span)
}

type Span interface {
	End()
	SetTag(key, value string)
}

type noopTracer struct{}

func NewNoopTracer() Tracer {
	return &noopTracer{}
}

func (t *noopTracer) StartSpan(ctx context.Context, name string) (context.Context, Span) {
	return ctx, &noopSpan{}
}

type noopSpan struct{}

func (s *noopSpan) End() {}
func (s *noopSpan) SetTag(_, _ string) {}
