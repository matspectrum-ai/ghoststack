package security

import (
	"context"
	"testing"
)

func TestStructuredAuditLoggerLogAndList(t *testing.T) {
	logger := NewStructuredAuditLogger("")
	entry := StructuredAuditEntry{Action: "start", Source: "test", Detail: "ok"}
	if err := logger.Log(context.Background(), entry); err != nil {
		t.Fatalf("log: %v", err)
	}
	list, err := logger.List(context.Background(), 10)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(list))
	}
	if list[0].Action != "start" {
		t.Fatalf("unexpected action: %s", list[0].Action)
	}
}

func TestNoopSandboxApply(t *testing.T) {
	sandbox := NewNoopSandbox()
	if err := sandbox.Apply(context.Background(), SandboxPolicy{AllowNetwork: true}); err != nil {
		t.Fatalf("apply: %v", err)
	}
}

func TestSecureBootRotateSecretValidation(t *testing.T) {
	boot := NewSecureBoot("abc")
	if err := boot.RotateSecret(context.Background(), "", "new"); err == nil {
		t.Fatal("expected error for empty name")
	}
	if err := boot.RotateSecret(context.Background(), "name", ""); err == nil {
		t.Fatal("expected error for empty value")
	}
}
