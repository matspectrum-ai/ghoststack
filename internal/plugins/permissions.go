package plugins

import "fmt"

var knownCapabilities = map[string]struct{}{
	"vpn.provider":      {},
	"proxy.provider":    {},
	"dns.provider":      {},
	"firewall.provider": {},
	"gateway.backend":   {},
	"runtime.provider":  {},
	"dashboard.widget":  {},
	"auth.provider":     {},
	"storage.provider":  {},
	"metrics.provider":  {},
	"notification.provider": {},
}

var knownPermissions = map[string]struct{}{
	"network":    {},
	"filesystem": {},
	"secrets":    {},
	"system":     {},
}

type PermissionChecker struct{}

func NewPermissionChecker() *PermissionChecker {
	return &PermissionChecker{}
}

func (p *PermissionChecker) CheckCapabilities(caps []string) error {
	for _, cap := range caps {
		if _, ok := knownCapabilities[cap]; !ok {
			return fmt.Errorf("unknown capability: %s", cap)
		}
	}
	return nil
}

func (p *PermissionChecker) CheckPermissions(perms []string) error {
	for _, perm := range perms {
		if _, ok := knownPermissions[perm]; !ok {
			return fmt.Errorf("unknown permission: %s", perm)
		}
	}
	return nil
}

func (p *PermissionChecker) AllowedCapabilities() []string {
	out := make([]string, 0, len(knownCapabilities))
	for cap := range knownCapabilities {
		out = append(out, cap)
	}
	return out
}

func (p *PermissionChecker) AllowedPermissions() []string {
	out := make([]string, 0, len(knownPermissions))
	for perm := range knownPermissions {
		out = append(out, perm)
	}
	return out
}
