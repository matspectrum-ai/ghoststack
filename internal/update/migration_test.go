package update

import (
	"context"
	"errors"
	"testing"
)

func TestMigrationEngine(t *testing.T) {
	engine := NewMigrationEngine()

	applied := []string{}
	engine.Register(Migration{
		ID:      "m1",
		Version: "1.0.0",
		Apply: func(ctx context.Context) error {
			applied = append(applied, "m1")
			return nil
		},
	})

	if err := engine.Apply(context.Background(), "1.0.0"); err != nil {
		t.Fatalf("apply: %v", err)
	}

	pending := engine.Pending("1.0.0")
	if len(pending) != 0 {
		t.Fatalf("expected 0 pending, got %d", len(pending))
	}

	if err := engine.Rollback(context.Background(), "0.0.0"); err != nil {
		t.Fatalf("rollback: %v", err)
	}
}

func TestMigrationEngineApplyFailure(t *testing.T) {
	engine := NewMigrationEngine()

	engine.Register(Migration{
		ID:      "m1",
		Version: "1.0.0",
		Apply: func(ctx context.Context) error {
			return errors.New("fail")
		},
	})

	if err := engine.Apply(context.Background(), "1.0.0"); err == nil {
		t.Fatal("expected error")
	}
}
