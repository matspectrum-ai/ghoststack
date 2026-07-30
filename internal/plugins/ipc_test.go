package plugins

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestSubprocessPluginLoader(t *testing.T) {
	tmpDir := t.TempDir()

	manifestContent := `id: test-plugin
name: Test Plugin
version: 0.1.0
author: test
description: A test plugin
entry: /bin/true
sdkversion: 1.0.0
capabilities:
  - runtime.provider
`

	if err := os.WriteFile(filepath.Join(tmpDir, "manifest.yaml"), []byte(manifestContent), 0644); err != nil {
		t.Fatalf("write manifest: %v", err)
	}

	loader := &subprocessPluginLoader{}
	plugin, err := loader.Load(tmpDir)
	if err != nil {
		t.Fatalf("load: %v", err)
	}

	if plugin == nil {
		t.Fatal("expected non-nil plugin")
	}

	manifest := plugin.Manifest()
	if manifest.ID != "test-plugin" {
		t.Fatalf("expected id test-plugin, got %s", manifest.ID)
	}
	if manifest.Version != "0.1.0" {
		t.Fatalf("expected version 0.1.0, got %s", manifest.Version)
	}
}

func TestSubprocessPluginManifest(t *testing.T) {
	tmpDir := t.TempDir()
	entryPath := "/bin/echo"

	manifestContent := `id: echo-plugin
name: Echo Plugin
version: 1.0.0
author: ghoststack
description: Simple echo plugin
entry: ` + entryPath + `
sdkversion: 1.0.0
capabilities:
  - runtime.provider
permissions:
  - network
`

	if err := os.WriteFile(filepath.Join(tmpDir, "manifest.yaml"), []byte(manifestContent), 0644); err != nil {
		t.Fatalf("write manifest: %v", err)
	}

	manifest, err := ParseManifest(filepath.Join(tmpDir, "manifest.yaml"))
	if err != nil {
		t.Fatalf("parse manifest: %v", err)
	}

	if manifest.ID != "echo-plugin" {
		t.Fatalf("expected id echo-plugin, got %s", manifest.ID)
	}

	if len(manifest.Capabilities) != 1 || manifest.Capabilities[0] != "runtime.provider" {
		t.Fatalf("unexpected capabilities: %v", manifest.Capabilities)
	}

	if len(manifest.Permissions) != 1 || manifest.Permissions[0] != "network" {
		t.Fatalf("unexpected permissions: %v", manifest.Permissions)
	}
}

func TestSubprocessPluginInitialize(t *testing.T) {
	tmpDir := t.TempDir()
	entryPath := "/bin/sleep"

	manifestContent := `id: sleeper
name: Sleeper Plugin
version: 0.1.0
author: test
description: Sleeps for a bit
entry: ` + entryPath + `
sdkversion: 1.0.0
`

	if err := os.WriteFile(filepath.Join(tmpDir, "manifest.yaml"), []byte(manifestContent), 0644); err != nil {
		t.Fatalf("write manifest: %v", err)
	}

	loader := &subprocessPluginLoader{}
	plugin, err := loader.Load(tmpDir)
	if err != nil {
		t.Fatalf("load: %v", err)
	}

	sp, ok := plugin.(*subprocessPlugin)
	if !ok {
		t.Fatalf("expected *subprocessPlugin, got %T", plugin)
	}

	pc := PluginContext{}
	err = sp.Initialize(context.Background(), pc)

	if err == nil {
		t.Log("Initialize succeeded (plugin binary may or may not be running)")
	} else {
		t.Logf("Initialize expected for non-plugin binary: %v", err)
	}

	sp.Disable(context.Background())
}

func TestSubprocessPluginStates(t *testing.T) {
	manifest := PluginManifest{
		ID:      "test",
		Version: "0.1.0",
		Entry:   "/bin/true",
	}

	plugin := newSubprocessPlugin(manifest)

	if plugin.Manifest().ID != "test" {
		t.Fatalf("expected id test")
	}

	ctx := context.Background()
	pc := PluginContext{}

	if err := plugin.Initialize(ctx, pc); err != nil {
		t.Logf("Initialize (expected to fail without binary): %v", err)
	}

	if err := plugin.Unload(ctx); err != nil {
		t.Fatalf("unload: %v", err)
	}
}
