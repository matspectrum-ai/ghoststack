package security

import (
	"context"
	"errors"
	"testing"
)

func TestIntegrityVerify(t *testing.T) {
	ic := newIntegrityChecker()

	results, err := ic.Verify(context.Background(), "")
	if err == nil {
		t.Fatal("expected error for empty target")
	}

	results, err = ic.Verify(context.Background(), "target")
	if err != nil {
		t.Fatalf("verify: %v", err)
	}
	if len(results) != 0 {
		t.Fatalf("expected no failures, got %v", results)
	}
}

func TestIntegrityVerifyWithFailures(t *testing.T) {
	ic := newIntegrityChecker().(*integrityChecker)
	ic.recordFailure("mismatch")

	results, err := ic.Verify(context.Background(), "target")
	if !errors.Is(err, ErrIntegrityFailure) {
		t.Fatalf("expected ErrIntegrityFailure, got %v", err)
	}
	if len(results) != 1 || results[0] != "mismatch" {
		t.Fatalf("unexpected results: %v", results)
	}
}
