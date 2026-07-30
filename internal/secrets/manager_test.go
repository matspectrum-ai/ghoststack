package secrets

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInitAndLoad(t *testing.T) {
	tmpDir := t.TempDir()
	homeDir := filepath.Join(tmpDir, ".ghoststack")

	sm := NewSecretsManager(homeDir)
	passphrase := "test-master-password-123"

	if err := sm.Init(passphrase); err != nil {
		t.Fatalf("init: %v", err)
	}

	if _, err := os.Stat(filepath.Join(homeDir, "secrets.enc")); os.IsNotExist(err) {
		t.Fatal("secrets.enc not created")
	}

	sm2 := NewSecretsManager(homeDir)
	if err := sm2.Load(passphrase); err != nil {
		t.Fatalf("load: %v", err)
	}
}

func TestInitTwiceFails(t *testing.T) {
	tmpDir := t.TempDir()
	sm := NewSecretsManager(filepath.Join(tmpDir, ".ghoststack"))

	if err := sm.Init("pass"); err != nil {
		t.Fatalf("init: %v", err)
	}

	if err := sm.Init("pass"); err == nil {
		t.Fatal("expected error on second init")
	}
}

func TestSetGetDelete(t *testing.T) {
	tmpDir := t.TempDir()
	sm := NewSecretsManager(filepath.Join(tmpDir, ".ghoststack"))

	if err := sm.Init("mypass"); err != nil {
		t.Fatalf("init: %v", err)
	}

	if err := sm.Set("wg-key", "supersecretkey123"); err != nil {
		t.Fatalf("set: %v", err)
	}

	val, err := sm.Get("wg-key")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if val != "supersecretkey123" {
		t.Fatalf("expected supersecretkey123, got %s", val)
	}

	if err := sm.Set("api-token", "token456"); err != nil {
		t.Fatalf("set2: %v", err)
	}

	keys := sm.List()
	if len(keys) != 2 {
		t.Fatalf("expected 2 keys, got %d: %v", len(keys), keys)
	}

	if err := sm.Delete("wg-key"); err != nil {
		t.Fatalf("delete: %v", err)
	}

	if sm.Exists("wg-key") {
		t.Fatal("wg-key should not exist after delete")
	}
	if !sm.Exists("api-token") {
		t.Fatal("api-token should exist")
	}
}

func TestPersistAcrossSessions(t *testing.T) {
	tmpDir := t.TempDir()
	homeDir := filepath.Join(tmpDir, ".ghoststack")
	passphrase := "persist-test-pass"

	sm := NewSecretsManager(homeDir)
	if err := sm.Init(passphrase); err != nil {
		t.Fatalf("init: %v", err)
	}

	if err := sm.Set("key1", "value1"); err != nil {
		t.Fatalf("set: %v", err)
	}

	if err := sm.Save(passphrase); err != nil {
		t.Fatalf("save: %v", err)
	}

	sm2 := NewSecretsManager(homeDir)
	if err := sm2.Load(passphrase); err != nil {
		t.Fatalf("load: %v", err)
	}

	val, err := sm2.Get("key1")
	if err != nil {
		t.Fatalf("get after reload: %v", err)
	}
	if val != "value1" {
		t.Fatalf("expected value1, got %s", val)
	}
}

func TestWrongPassphrase(t *testing.T) {
	tmpDir := t.TempDir()
	homeDir := filepath.Join(tmpDir, ".ghoststack")

	sm := NewSecretsManager(homeDir)
	if err := sm.Init("correct-pass"); err != nil {
		t.Fatalf("init: %v", err)
	}

	sm2 := NewSecretsManager(homeDir)
	if err := sm2.Load("wrong-pass"); err == nil {
		t.Fatal("expected error with wrong passphrase")
	}
}

func TestEmptyValue(t *testing.T) {
	sm := NewSecretsManager(t.TempDir())
	if err := sm.Set("key", ""); err == nil {
		t.Fatal("expected error for empty value")
	}
}

func TestKeyHash(t *testing.T) {
	h1 := KeyHash("pass1")
	h2 := KeyHash("pass2")
	if h1 == h2 {
		t.Fatal("expected different hashes for different passwords")
	}
}
