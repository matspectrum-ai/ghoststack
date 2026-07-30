package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

var (
	ErrProfileNotFound = fmt.Errorf("profile not found")
)

type ResolvedProfile struct {
	Name      string
	Providers []string
	Options   map[string]any
}

func ResolveProfile(cfg RawConfig, name string) (ResolvedProfile, error) {
	profile, ok := cfg.Profiles[name]
	if !ok {
		return ResolvedProfile{}, fmt.Errorf("%w: %s", ErrProfileNotFound, name)
	}

	return ResolvedProfile{
		Name:      name,
		Providers: profile.Providers,
		Options:   profile.Options,
	}, nil
}

func LoadSecrets(cfg RawConfig) (map[string]string, error) {
	secrets := make(map[string]string)
	for _, source := range cfg.Secrets.Sources {
		data, err := os.ReadFile(source)
		if err != nil {
			return nil, fmt.Errorf("read secret source %s: %w", source, err)
		}

		var block map[string]string
		if err := yaml.Unmarshal(data, &block); err != nil {
			return nil, fmt.Errorf("parse secret source %s: %w", source, err)
		}
		for k, v := range block {
			secrets[k] = v
		}
	}

	return secrets, nil
}
