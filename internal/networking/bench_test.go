package networking

import (
	"context"
	"testing"
)

func BenchmarkGatewayStartStop(b *testing.B) {
	g := NewGateway("bench")
	ctx := context.Background()

	for i := 0; i < b.N; i++ {
		_ = g.Start(ctx)
		_ = g.Stop(ctx)
	}
}

func BenchmarkRouteTableAdd(b *testing.B) {
	table := newRouteTable()
	ctx := context.Background()

	for i := 0; i < b.N; i++ {
		_ = table.Add(ctx, "10.0.0.0/24", "10.0.0.1")
	}
}
