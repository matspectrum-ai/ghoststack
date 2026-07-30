package networking

import (
	"context"
	"errors"
	"testing"
)

func TestResolverSetServersAndFlush(t *testing.T) {
	resolver := newResolver().(*resolver)

	if err := resolver.FlushCache(context.Background()); !errors.Is(err, ErrResolverNotStarted) {
		t.Fatalf("expected ErrResolverNotStarted, got %v", err)
	}

	resolver.mu.Lock()
	resolver.started = true
	resolver.mu.Unlock()

	if err := resolver.SetServers(context.Background(), []string{"1.1.1.1"}); err != nil {
		t.Fatalf("set servers: %v", err)
	}

	if err := resolver.FlushCache(context.Background()); err != nil {
		t.Fatalf("flush: %v", err)
	}
}
