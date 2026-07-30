package plugins

import (
	"os"
	"testing"
)

func TestDiscoverPlugins(t *testing.T) {
	manifests, err := DiscoverPlugins(nil)
	if err != nil {
		t.Fatalf("discover: %v", err)
	}
	if len(manifests) != 0 {
		t.Fatalf("expected 0 manifests in empty dirs, got %d", len(manifests))
	}
}

func TestDiscoverPluginsWithDir(t *testing.T) {
	tmpDir := t.TempDir()
	pluginDir := pluginDir(tmpDir, "example")
	if err := os.MkdirAll(pluginDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(manifestPath(pluginDir), []byte("id: example\nversion: 1.0.0\nsdk: 1.0\nentry: plugin.so\nauthor: test\ndescription: test\n"), 0644); err != nil {
		t.Fatalf("write manifest: %v", err)
	}

	manifests, err := DiscoverPlugins([]string{tmpDir})
	if err != nil {
		t.Fatalf("discover: %v", err)
	}
	if len(manifests) != 1 {
		t.Fatalf("expected 1 manifest, got %d", len(manifests))
	}
}

func TestExpandPath(t *testing.T) {
	result := expandPath("~/plugins")
	if result == "~/plugins" {
		t.Fatal("expected expanded path")
	}
}

func TestValidateManifestPath(t *testing.T) {
	if err := ValidateManifestPath(""); err == nil {
		t.Fatal("expected error for empty path")
	}
	if err := ValidateManifestPath("/nonexistent/manifest.yaml"); err == nil {
		t.Fatal("expected error for missing path")
	}
}

func pluginDir(base, name string) string {
	return base + "/" + name
}

func manifestPath(dir string) string {
	return dir + "/manifest.yaml"
}
