package cli

import (
	"os"
	"path/filepath"
	"testing"
)


func TestPluginListCommand(t *testing.T) {
	cmd := newPluginListCommand()
	_, err := executeCommand(cmd)
	if err != nil {
		t.Fatalf("plugin list: %v", err)
	}
}

func TestPluginInstallCommand(t *testing.T) {
	tmpDir := t.TempDir()
	manifest := filepath.Join(tmpDir, "manifest.yaml")
	content := []byte("id: test-plugin\nversion: 0.1.0\nauthor: test\ndescription: test plugin\nentry: main.go\nsdkversion: 1.0.0\n")
	if err := os.WriteFile(manifest, content, 0644); err != nil {
		t.Fatalf("write manifest: %v", err)
	}

	cmd := newPluginInstallCommand()
	_, err := executeCommand(cmd, tmpDir)
	if err != nil {
		t.Fatalf("plugin install: %v", err)
	}
}

func TestConfigValidateCommand(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "ghost-*.yaml")
	if err != nil {
		t.Fatalf("create temp config: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	cmd := newConfigValidateCommand()
	_, err = executeCommand(cmd, "--config", tmpFile.Name())
	if err != nil {
		t.Fatalf("config validate: %v", err)
	}
}

func TestProviderListCommand(t *testing.T) {
	cmd := newProviderListCommand()
	_, err := executeCommand(cmd)
	if err != nil {
		t.Fatalf("provider list: %v", err)
	}
}
