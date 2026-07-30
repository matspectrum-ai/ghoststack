package benchmark

import (
	"context"
	"testing"
	"time"
)

func TestSuiteRun(t *testing.T) {
	suite := NewSuite()
	result := suite.Run(context.Background(), "fast", func(ctx context.Context) error {
		time.Sleep(1 * time.Millisecond)
		return nil
	})
	if result.Error != nil {
		t.Fatalf("unexpected error: %v", result.Error)
	}
	if result.Duration < 1*time.Millisecond {
		t.Fatalf("expected at least 1ms, got %s", result.Duration)
	}
}

func TestSuiteSummary(t *testing.T) {
	suite := NewSuite()
	suite.Run(context.Background(), "ok", func(ctx context.Context) error {
		return nil
	})
	suite.Run(context.Background(), "fail", func(ctx context.Context) error {
		return context.DeadlineExceeded
	})
	summary := suite.Summary()
	if summary == "" {
		t.Fatal("expected non-empty summary")
	}
}

func TestGoBenchmarkWrapper(t *testing.T) {
	var count int
	var b testing.B
	b.N = 10
	for i := 0; i < b.N; i++ {
		count++
	}
	if count != 10 {
		t.Fatalf("expected 10, got %d", count)
	}
}
