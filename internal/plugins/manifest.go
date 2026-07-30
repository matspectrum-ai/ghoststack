package plugins

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type PluginManifestPayload struct {
	Capabilities []string `yaml:"capabilities"`
	Permissions  []string `yaml:"permissions"`
	Dependencies  map[string]string `yaml:"dependencies"`
	Resources     PluginResources `yaml:"resources"`
}

func ParseManifest(path string) (PluginManifest, error) {
	var manifest PluginManifest

	b, err := os.ReadFile(path)
	if err != nil {
		return manifest, fmt.Errorf("read manifest: %w", err)
	}

	if err := yaml.Unmarshal(b, &manifest); err != nil {
		return manifest, fmt.Errorf("parse manifest: %w", err)
	}

	if err := validateManifestFields(manifest); err != nil {
		return manifest, err
	}

	return manifest, nil
}

func validateManifestFields(m PluginManifest) error {
	if strings.TrimSpace(m.ID) == "" {
		return fmt.Errorf("manifest.id must not be empty")
	}
	if strings.TrimSpace(m.Version) == "" {
		return fmt.Errorf("manifest.version must not be empty")
	}
	if strings.TrimSpace(m.SDKVersion) == "" {
		return fmt.Errorf("manifest.sdk must not be empty")
	}
	if strings.TrimSpace(m.Entry) == "" {
		return fmt.Errorf("manifest.entry must not be empty")
	}
	if strings.TrimSpace(m.Author) == "" {
		return fmt.Errorf("manifest.author must not be empty")
	}
	if strings.TrimSpace(m.Description) == "" {
		return fmt.Errorf("manifest.description must not be empty")
	}
	return nil
}

func (m PluginManifest) ValidateCapabilities(known map[string]struct{}) error {
	if len(m.Capabilities) == 0 {
		return nil
	}

	for _, cap := range m.Capabilities {
		if _, ok := known[cap]; !ok {
			return fmt.Errorf("unknown capability: %s", cap)
		}
	}
	return nil
}

func (m PluginManifest) ValidatePermissions(known map[string]struct{}) error {
	if len(m.Permissions) == 0 {
		return nil
	}

	for _, perm := range m.Permissions {
		if _, ok := known[perm]; !ok {
			return fmt.Errorf("unknown permission: %s", perm)
		}
	}
	return nil
}

func (m PluginManifest) PluginDir() string {
	return filepath.Dir(m.Entry)
}
