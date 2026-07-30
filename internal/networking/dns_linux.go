package networking

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func setResolvectl(ctx context.Context, iface string, servers []string) error {
	if err := exec.CommandContext(ctx, "resolvectl", "dns", iface, servers[0]).Run(); err != nil {
		return fmt.Errorf("resolvectl dns: %w", err)
	}

	for _, server := range servers[1:] {
		exec.CommandContext(ctx, "resolvectl", "dns", iface, server).Run()
	}

	if err := exec.CommandContext(ctx, "resolvectl", "domain", iface, "~.").Run(); err != nil {
		return fmt.Errorf("resolvectl domain: %w", err)
	}

	return nil
}

func flushResolvectl(ctx context.Context) error {
	if err := exec.CommandContext(ctx, "resolvectl", "flush-caches").Run(); err != nil {
		return fmt.Errorf("resolvectl flush-caches: %w", err)
	}
	return nil
}

type resolvConfBackup struct {
	path    string
	content []byte
}

var dnsBackup *resolvConfBackup

func setResolvConf(ctx context.Context, servers []string) error {
	content, err := os.ReadFile("/etc/resolv.conf")
	if err == nil {
		dnsBackup = &resolvConfBackup{path: "/etc/resolv.conf", content: content}
	}

	var sb strings.Builder
	for _, s := range servers {
		sb.WriteString(fmt.Sprintf("nameserver %s\n", s))
	}

	tmp := filepath.Dir("/etc/resolv.conf")
	if err := os.WriteFile(filepath.Join(tmp, "resolv.conf.ghoststack"), []byte(sb.String()), 0644); err != nil {
		return fmt.Errorf("write resolv.conf: %w", err)
	}

	if err := os.Rename(filepath.Join(tmp, "resolv.conf.ghoststack"), "/etc/resolv.conf"); err != nil {
		return fmt.Errorf("replace resolv.conf: %w", err)
	}

	return nil
}

func restoreResolvConf(ctx context.Context) error {
	if dnsBackup == nil {
		return nil
	}

	if err := os.WriteFile("/etc/resolv.conf", dnsBackup.content, 0644); err != nil {
		return fmt.Errorf("restore resolv.conf: %w", err)
	}

	dnsBackup = nil
	return nil
}
