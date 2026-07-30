package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func writeTempFile(t *testing.T, name string, content []byte) string {
	t.Helper()
	path := "tmp_" + name
	if err := os.WriteFile(path, content, 0o644); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Remove(path)
	})
	return path
}

func TestResolveProfile(t *testing.T) {
	cfg := RawConfig{
		APIVersion: "v1",
		Kind:       "Config",
		Profiles: map[string]ProfileSpec{
			"default": {Providers: []string{"wg"}},
		},
	}

	profile, err := ResolveProfile(cfg, "default")
	require.NoError(t, err)
	require.Equal(t, []string{"wg"}, profile.Providers)
}

func TestResolveProfileMissing(t *testing.T) {
	_, err := ResolveProfile(RawConfig{}, "default")
	require.ErrorIs(t, err, ErrProfileNotFound)
}

func TestLoadSecrets(t *testing.T) {
	content, err := yaml.Marshal(map[string]string{"token": "abc"})
	require.NoError(t, err)
	path := writeTempFile(t, "secrets.yaml", content)

	cfg := RawConfig{
		Secrets: SecretSpec{Sources: []string{path}},
	}
	secrets, err := LoadSecrets(cfg)
	require.NoError(t, err)
	require.Equal(t, "abc", secrets["token"])
}
