package plugins

import (
	"testing"
)

func TestParseManifest(t *testing.T) {
	manifest := PluginManifest{
		ID:          "test",
		Version:     "1.0.0",
		SDKVersion:  "1.0",
		Entry:       "plugin.so",
		Author:      "test",
		Description: "test plugin",
	}

	if err := validateManifestFields(manifest); err != nil {
		t.Fatalf("validate: %v", err)
	}
}

func TestParseManifestInvalid(t *testing.T) {
	tests := []PluginManifest{
		{Version: "1.0.0", SDKVersion: "1.0", Entry: "plugin.so", Author: "test", Description: "test"},
		{ID: "test", SDKVersion: "1.0", Entry: "plugin.so", Author: "test", Description: "test"},
		{ID: "test", Version: "1.0.0", Entry: "plugin.so", Author: "test", Description: "test"},
		{ID: "test", Version: "1.0.0", SDKVersion: "1.0", Author: "test", Description: "test"},
		{ID: "test", Version: "1.0.0", SDKVersion: "1.0", Entry: "plugin.so", Description: "test"},
	}

	for i, m := range tests {
		if err := validateManifestFields(m); err == nil {
			t.Fatalf("test %d: expected error", i)
		}
	}
}

func TestValidateCapabilities(t *testing.T) {
	manifest := PluginManifest{Capabilities: []string{"vpn.provider", "unknown"}}
	if err := manifest.ValidateCapabilities(knownCapabilities); err == nil {
		t.Fatal("expected error for unknown capability")
	}
}

func TestValidatePermissions(t *testing.T) {
	manifest := PluginManifest{Permissions: []string{"network", "unknown"}}
	if err := manifest.ValidatePermissions(knownPermissions); err == nil {
		t.Fatal("expected error for unknown permission")
	}
}

func TestPluginDir(t *testing.T) {
	manifest := PluginManifest{Entry: "/path/to/plugin.so"}
	dir := manifest.PluginDir()
	if dir != "/path/to" {
		t.Fatalf("unexpected dir: %s", dir)
	}
}
