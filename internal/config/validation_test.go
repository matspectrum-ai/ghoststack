package config

import (
	"os"
	"testing"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	path := "tmp_test_config.yaml"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write temp config: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Remove(path)
	})
	return path
}

func TestValidateMissingFields(t *testing.T) {
	path := writeTempConfig(t, "profiles:\n  default:\n    providers: []\n")
	result, err := Validate(path)
	if err != nil {
		t.Fatalf("validate: %v", err)
	}
	if result.Valid {
		t.Fatal("expected invalid config")
	}
}

func TestValidateEmptyProviders(t *testing.T) {
	path := writeTempConfig(t, "apiVersion: v1\nkind: Config\nprofiles:\n  default:\n    providers: []\n")
	result, err := Validate(path)
	if err != nil {
		t.Fatalf("validate: %v", err)
	}
	if result.Valid {
		t.Fatal("expected invalid config")
	}
}
