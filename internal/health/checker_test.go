package health

import (
	"context"
	"testing"
)

func TestCheckerRegisterAndRun(t *testing.T) {
	c := NewChecker()
	c.Register("db", func(ctx context.Context) (HealthStatus, string) {
		return HealthUp, "ok"
	})
	c.Register("cache", func(ctx context.Context) (HealthStatus, string) {
		return HealthDown, "timeout"
	})

	checks, err := c.Run(context.Background())
	if err == nil {
		t.Fatal("expected error for down check")
	}
	if len(checks) != 2 {
		t.Fatalf("expected 2 checks, got %d", len(checks))
	}
}
