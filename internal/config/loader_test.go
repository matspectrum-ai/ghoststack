package config

import (
	"os"
	"testing"
)

func TestLoadMissingFile(t *testing.T) {
	_, err := Load("missing.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadInvalidYAML(t *testing.T) {
	err := os.WriteFile("tmp_invalid.yaml", []byte(": invalid"), 0o644)
	if err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	defer func() {
		_ = os.Remove("tmp_invalid.yaml")
	}()

	_, err = Load("tmp_invalid.yaml")
	if err == nil {
		t.Fatal("expected error for invalid yaml")
	}
}
