package recovery

import (
	"context"
	"testing"
)

func TestRecoveryManagerRegisterAndExecute(t *testing.T) {
	r := NewRecoveryManager()
	called := false
	r.Register("restart", func(ctx context.Context) error {
		called = true
		return nil
	})

	if err := r.Execute(context.Background(), "restart"); err != nil {
		t.Fatalf("execute: %v", err)
	}
	if !called {
		t.Fatal("expected action to be called")
	}
}

func TestRecoveryManagerMissingAction(t *testing.T) {
	r := NewRecoveryManager()
	if err := r.Execute(context.Background(), "missing"); err == nil {
		t.Fatal("expected error for missing action")
	}
}
