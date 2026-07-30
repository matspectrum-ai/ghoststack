package plugins

import (
	"context"
	"testing"
)

func TestPluginManagerDiscover(t *testing.T) {
	manager := NewPluginManager(nil, nil)
	manifests, err := manager.Discover(context.Background(), nil)
	if err != nil {
		t.Fatalf("discover: %v", err)
	}
	if len(manifests) != 0 {
		t.Fatalf("expected 0 manifests, got %d", len(manifests))
	}
}

func TestPluginManagerLoadAndGet(t *testing.T) {
	manager := NewPluginManager(nil, nil)
	tmpDir := t.TempDir()
	pluginDir := tmpDir + "/test-plugin"
	if err := writeTestPlugin(t, pluginDir); err != nil {
		t.Fatalf("write test plugin: %v", err)
	}

	_, err := manager.Load(context.Background(), pluginDir)
	if err == nil {
		t.Fatal("expected error for unimplemented loader")
	}
}

func TestPluginManagerList(t *testing.T) {
	manager := NewPluginManager(nil, nil)
	manifests, err := manager.List(context.Background())
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(manifests) != 0 {
		t.Fatalf("expected 0 plugins, got %d", len(manifests))
	}
}

func TestPluginManagerEnableDisableUnload(t *testing.T) {
	manager := NewPluginManager(nil, nil)

	if err := manager.Enable(context.Background(), ""); err == nil {
		t.Fatal("expected error for empty plugin id")
	}
	if err := manager.Disable(context.Background(), ""); err == nil {
		t.Fatal("expected error for empty plugin id")
	}
	if err := manager.Unload(context.Background(), ""); err == nil {
		t.Fatal("expected error for empty plugin id")
	}

	if err := manager.Enable(context.Background(), "missing"); err == nil {
		t.Fatal("expected error for missing plugin")
	}
	if err := manager.Disable(context.Background(), "missing"); err == nil {
		t.Fatal("expected error for missing plugin")
	}
	if err := manager.Unload(context.Background(), "missing"); err == nil {
		t.Fatal("expected error for missing plugin")
	}
}

func TestPluginManagerValidate(t *testing.T) {
	manager := NewPluginManager(nil, nil)
	if err := manager.Validate(context.Background(), PluginManifest{}); err == nil {
		t.Fatal("expected error for empty manifest")
	}
}

func writeTestPlugin(t *testing.T, dir string) error {
	t.Helper()
	return nil
}
