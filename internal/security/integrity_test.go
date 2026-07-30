package security

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestIntegrityCheckerVerify(t *testing.T) {
	ic := newIntegrityChecker()

	failures, err := ic.Verify(context.Background(), "some-target")
	if err != nil {
		t.Fatalf("expected no error for new checker, got: %v", err)
	}
	if len(failures) != 0 {
		t.Fatalf("expected no failures, got: %v", failures)
	}
}

func TestIntegrityCheckerEmptyTarget(t *testing.T) {
	ic := newIntegrityChecker()

	_, err := ic.Verify(context.Background(), "")
	if err == nil {
		t.Fatal("expected error for empty target")
	}
}

func TestSecureBootSHA256Verify(t *testing.T) {
	tmpDir := t.TempDir()
	binPath := filepath.Join(tmpDir, "test.bin")

	if err := os.WriteFile(binPath, []byte("hello world"), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}

	boot := NewSecureBoot("")
	failures, err := boot.Verify(context.Background(), binPath)
	if err != nil {
		t.Fatalf("verify without expected hash should not error: %v", err)
	}
	if len(failures) != 1 || failures[0] != "secure boot: no expected hash configured" {
		t.Fatalf("expected 'no hash' failure, got: %v", failures)
	}

	hash, err := boot.ComputeHash(context.Background(), binPath)
	if err != nil {
		t.Fatalf("compute hash: %v", err)
	}
	if hash == "" {
		t.Fatal("expected non-empty hash")
	}

	boot.SetExpectedHash(hash)
	failures, err = boot.Verify(context.Background(), binPath)
	if err != nil {
		t.Fatalf("verify with correct hash: %v", err)
	}
	if len(failures) != 0 {
		t.Fatalf("expected no failures, got: %v", failures)
	}

	boot.SetExpectedHash("badhash")
	failures, err = boot.Verify(context.Background(), binPath)
	if err == nil {
		t.Fatal("expected error for bad hash")
	}
	if len(failures) == 0 {
		t.Fatal("expected failures for bad hash")
	}
}

func TestSecureBootRotateSecret(t *testing.T) {
	boot := NewSecureBoot("abc")

	if err := boot.RotateSecret(context.Background(), "", "new"); err == nil {
		t.Fatal("expected error for empty name")
	}

	if err := boot.RotateSecret(context.Background(), "name", ""); err == nil {
		t.Fatal("expected error for empty value")
	}

	if err := boot.RotateSecret(context.Background(), "wg-key", "new-secret-value"); err != nil {
		t.Fatalf("rotate: %v", err)
	}

	val, ok := boot.GetSecret("wg-key")
	if !ok {
		t.Fatal("expected secret to exist")
	}
	if val != "new-secret-value" {
		t.Fatalf("expected 'new-secret-value', got '%s'", val)
	}
}

func TestSecureBootComputeHash(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "data.bin")

	if err := os.WriteFile(path, []byte("test data for hashing"), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}

	boot := NewSecureBoot("")
	hash, err := boot.ComputeHash(context.Background(), path)
	if err != nil {
		t.Fatalf("compute hash: %v", err)
	}

	if len(hash) != 64 {
		t.Fatalf("expected 64 hex chars, got %d: %s", len(hash), hash)
	}
}
