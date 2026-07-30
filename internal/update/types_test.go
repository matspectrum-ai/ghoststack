package update

import (
	"testing"
)

func TestParseVersion(t *testing.T) {
	v, err := ParseVersion("1.2.3")
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if v.Major != 1 || v.Minor != 2 || v.Patch != 3 {
		t.Fatalf("unexpected version: %+v", v)
	}
}

func TestParseVersionInvalid(t *testing.T) {
	tests := []string{
		"",
		"1.2",
		"1.2.3.4",
		"1.a.3",
	}
	for _, input := range tests {
		if _, err := ParseVersion(input); err == nil {
			t.Fatalf("expected error for %q", input)
		}
	}
}

func TestVersionLess(t *testing.T) {
	v1 := Version{Major: 1, Minor: 0, Patch: 0}
	v2 := Version{Major: 1, Minor: 1, Patch: 0}
	v3 := Version{Major: 1, Minor: 1, Patch: 0}

	if !v1.Less(v2) {
		t.Fatal("expected v1 < v2")
	}
	if v2.Less(v3) {
		t.Fatal("expected v2 == v3")
	}
}

func TestVersionCompatibleWith(t *testing.T) {
	v1 := Version{Major: 1, Minor: 0, Patch: 0}
	v2 := Version{Major: 1, Minor: 1, Patch: 0}
	v3 := Version{Major: 2, Minor: 0, Patch: 0}

	if !v1.CompatibleWith(v2) {
		t.Fatal("expected compatible")
	}
	if v1.CompatibleWith(v3) {
		t.Fatal("expected incompatible")
	}
}

func TestVersionResolver(t *testing.T) {
	resolver := NewVersionResolver()
	resolver.Register("core", Version{Major: 1, Minor: 0, Patch: 0})

	v, ok := resolver.Resolve("core")
	if !ok || v.String() != "1.0.0" {
		t.Fatalf("unexpected version: %s, ok=%v", v.String(), ok)
	}

	_, ok = resolver.Resolve("missing")
	if ok {
		t.Fatal("expected not found")
	}
}
