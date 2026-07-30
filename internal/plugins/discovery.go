package plugins

import (
	"fmt"
	"os"
	"path/filepath"
)

var officialDirectories = []string{
	"/usr/lib/ghoststack/plugins",
	"~/.local/share/ghoststack/plugins",
	"./plugins",
}

func DiscoverPlugins(dirs []string) ([]string, error) {
	if len(dirs) == 0 {
		dirs = officialDirectories
	}

	var manifests []string
	for _, dir := range dirs {
		expanded := expandPath(dir)
		entries, err := os.ReadDir(expanded)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}

			manifestPath := filepath.Join(expanded, entry.Name(), "manifest.yaml")
			if _, err := os.Stat(manifestPath); err == nil {
				manifests = append(manifests, manifestPath)
			}
		}
	}

	return manifests, nil
}

func expandPath(path string) string {
	if len(path) > 0 && path[0] == '~' {
		home, err := os.UserHomeDir()
		if err == nil {
			return filepath.Join(home, path[1:])
		}
	}
	return path
}

func ValidateManifestPath(path string) error {
	if path == "" {
		return fmt.Errorf("manifest path must not be empty")
	}

	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("manifest not found: %w", err)
	}

	if info.IsDir() {
		return fmt.Errorf("manifest path is a directory")
	}

	return nil
}
