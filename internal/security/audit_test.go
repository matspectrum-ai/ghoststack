package security

import (
	"context"
	"testing"
)

func TestAuditLoggerLogAndList(t *testing.T) {
	logger := newAuditLogger(10)

	entries, err := logger.List(context.Background(), 10)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected empty list, got %d", len(entries))
	}

	if err := logger.Log(context.Background(), AuditEntry{Action: "start", Source: "test"}); err != nil {
		t.Fatalf("log: %v", err)
	}

	entries, err = logger.List(context.Background(), 10)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(entries) != 1 || entries[0].Action != "start" {
		t.Fatalf("unexpected entries: %v", entries)
	}

	if err := logger.Log(context.Background(), AuditEntry{}); err != nil {
		t.Fatalf("log empty: %v", err)
	}

	entries, err = logger.List(context.Background(), 1)
	if err != nil {
		t.Fatalf("list limit: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected limited list, got %d", len(entries))
	}
}
