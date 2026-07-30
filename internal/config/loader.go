package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type RawConfig struct {
	APIVersion string                 `yaml:"apiVersion"`
	Kind       string                 `yaml:"kind"`
	Meta       map[string]string      `yaml:"meta"`
	Profiles   map[string]ProfileSpec `yaml:"profiles"`
	Secrets    SecretSpec             `yaml:"secrets"`
}

type ProfileSpec struct {
	Providers []string         `yaml:"providers"`
	Options   map[string]any   `yaml:"options"`
}

type SecretSpec struct {
	Sources []string `yaml:"sources"`
}

func Load(path string) (RawConfig, error) {
	var cfg RawConfig
	b, err := os.ReadFile(path)
	if err != nil {
		return cfg, fmt.Errorf("read config: %w", err)
	}
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return cfg, fmt.Errorf("parse yaml: %w", err)
	}
	return cfg, nil
}

func LoadFromString(data string) (RawConfig, error) {
	var cfg RawConfig
	if err := yaml.Unmarshal([]byte(data), &cfg); err != nil {
		return cfg, fmt.Errorf("parse yaml: %w", err)
	}
	return cfg, nil
}
