package plugins

import (
	"testing"
)

func TestPermissionChecker(t *testing.T) {
	checker := NewPermissionChecker()

	if err := checker.CheckCapabilities([]string{"vpn.provider", "proxy.provider"}); err != nil {
		t.Fatalf("check capabilities: %v", err)
	}
	if err := checker.CheckCapabilities([]string{"unknown"}); err == nil {
		t.Fatal("expected error for unknown capability")
	}

	if err := checker.CheckPermissions([]string{"network", "filesystem"}); err != nil {
		t.Fatalf("check permissions: %v", err)
	}
	if err := checker.CheckPermissions([]string{"unknown"}); err == nil {
		t.Fatal("expected error for unknown permission")
	}
}

func TestAllowedCapabilities(t *testing.T) {
	checker := NewPermissionChecker()
	caps := checker.AllowedCapabilities()
	if len(caps) == 0 {
		t.Fatal("expected non-empty capabilities")
	}
}

func TestAllowedPermissions(t *testing.T) {
	checker := NewPermissionChecker()
	perms := checker.AllowedPermissions()
	if len(perms) == 0 {
		t.Fatal("expected non-empty permissions")
	}
}
