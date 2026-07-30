package update

import (
	"context"
	"testing"
)

func TestUpdateManifestValidate(t *testing.T) {
	m := UpdateManifest{Version: "1.0.0", Checksum: "abc"}
	if err := m.Validate(); err != nil {
		t.Fatalf("validate: %v", err)
	}

	mEmpty := UpdateManifest{}
	if err := mEmpty.Validate(); err == nil {
		t.Fatal("expected error for empty manifest")
	}
}

func TestUpdateManagerCheck(t *testing.T) {
	manager := NewUpdateManager(nil)
	result, err := manager.Check(context.Background())
	if err != nil {
		t.Fatalf("check: %v", err)
	}
	if result.Available {
		t.Fatal("expected no available update")
	}
}

func TestUpdateManagerPrepare(t *testing.T) {
	manager := NewUpdateManager(nil)
	m := UpdateManifest{Version: "1.0.0", Checksum: "abc"}

	if err := manager.Prepare(context.Background(), m); err != nil {
		t.Fatalf("prepare: %v", err)
	}

	invalid := UpdateManifest{Version: "", Checksum: ""}
	if err := manager.Prepare(context.Background(), invalid); err == nil {
		t.Fatal("expected error for invalid manifest")
	}
}

func TestUpdateManagerRollback(t *testing.T) {
	manager := NewUpdateManager(nil)
	if err := manager.Rollback(context.Background()); err == nil {
		t.Fatal("expected error when no update in progress")
	}
}
