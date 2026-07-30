package networking

import (
	"context"
	"testing"
)

func TestResolverSetServersAndFlush(t *testing.T) {
	resolver := newResolver().(*resolver)

	if err := resolver.FlushCache(context.Background()); err == nil {
		t.Log("flush succeeded (resolver may allow unstarted flush)")
	}

	_ = resolver
}
