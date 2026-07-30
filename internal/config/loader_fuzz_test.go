package config

import (
	"testing"
)

func FuzzLoad(f *testing.F) {
	f.Add("apiVersion: v1\nkind: GhostStack\n")
	f.Add("profiles:\n  default:\n    providers: [wireguard]\n")
	f.Add("secrets:\n  sources: []\n")

	for _, seed := range []string{
		"",
		"invalid: yaml: [",
	} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, data string) {
		_, err := LoadFromString(data)
		if err != nil {
			return
		}
	})
}
