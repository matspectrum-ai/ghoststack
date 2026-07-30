package linux

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

type FirewallMode string

const (
	FirewallNftables    FirewallMode = "nftables"
	FirewallIptables    FirewallMode = "iptables"
	FirewallIptablesNft FirewallMode = "iptables-nft"
)

type Firewall struct {
	mode  FirewallMode
	iface string
}

func NewFirewall() *Firewall {
	return &Firewall{mode: detectMode()}
}

func detectMode() FirewallMode {
	if err := exec.Command("nft", "--version").Run(); err == nil {
		return FirewallNftables
	}
	if err := exec.Command("iptables-nft", "--version").Run(); err == nil {
		return FirewallIptablesNft
	}
	if err := exec.Command("iptables-legacy", "--version").Run(); err == nil {
		return FirewallIptables
	}
	return FirewallIptables
}

func (f *Firewall) Apply(ctx context.Context, rules []string) error {
	switch f.mode {
	case FirewallNftables:
		return f.applyNftables(ctx, rules)
	default:
		return f.applyIptables(ctx, rules)
	}
}

func (f *Firewall) applyNftables(ctx context.Context, rules []string) error {
	setup := `
table inet ghoststack {
	delete table inet ghoststack
}
table inet ghoststack {
	chain output {
		type filter hook output priority 0; policy accept;
	}
}
`
	cmd := exec.CommandContext(ctx, "nft", "-f", "-")
	cmd.Stdin = strings.NewReader(setup)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("nft setup: %w", err)
	}

	for _, rule := range rules {
		parts := strings.Fields(rule)
		if len(parts) < 2 {
			continue
		}
		action := parts[0]
		dest := parts[1]

		switch action {
		case "allow":
			if err := exec.CommandContext(ctx,
				"nft", "add", "rule", "inet", "ghoststack", "output",
				"ip", "daddr", dest, "accept",
			).Run(); err != nil {
				return fmt.Errorf("nft rule allow %s: %w", dest, err)
			}
		case "drop":
			if err := exec.CommandContext(ctx,
				"nft", "add", "rule", "inet", "ghoststack", "output",
				"ip", "daddr", dest, "drop",
			).Run(); err != nil {
				return fmt.Errorf("nft rule drop %s: %w", dest, err)
			}
		case "allow_port":
			if err := exec.CommandContext(ctx,
				"nft", "add", "rule", "inet", "ghoststack", "output",
				dest, "accept",
			).Run(); err != nil {
				return fmt.Errorf("nft rule allow_port %s: %w", dest, err)
			}
		}
	}

	return nil
}

func (f *Firewall) applyIptables(ctx context.Context, rules []string) error {
	ipt := string(f.mode)
	if f.mode == FirewallIptablesNft {
		ipt = "iptables-nft"
	}

	exec.CommandContext(ctx, ipt, "-N", "GHOSTSTACK").Run()
	exec.CommandContext(ctx, ipt, "-I", "OUTPUT", "-j", "GHOSTSTACK").Run()

	for _, rule := range rules {
		parts := strings.Fields(rule)
		if len(parts) < 2 {
			continue
		}
		action := parts[0]
		dest := parts[1]

		var args []string
		switch action {
		case "allow":
			args = []string{ipt, "-A", "GHOSTSTACK", "-d", dest, "-j", "ACCEPT"}
		case "drop":
			args = []string{ipt, "-A", "GHOSTSTACK", "-d", dest, "-j", "DROP"}
		case "allow_port":
			args = []string{ipt, "-A", "GHOSTSTACK", dest, "-j", "ACCEPT"}
		}

		if args != nil {
			if err := exec.CommandContext(ctx, args[0], args[1:]...).Run(); err != nil {
				return fmt.Errorf("iptables rule: %w", err)
			}
		}
	}

	return nil
}

func (f *Firewall) Flush(ctx context.Context) error {
	switch f.mode {
	case FirewallNftables:
		return f.flushNftables(ctx)
	default:
		return f.flushIptables(ctx)
	}
}

func (f *Firewall) flushNftables(ctx context.Context) error {
	return exec.CommandContext(ctx, "nft", "delete", "table", "inet", "ghoststack").Run()
}

func (f *Firewall) flushIptables(ctx context.Context) error {
	ipt := string(f.mode)
	if f.mode == FirewallIptablesNft {
		ipt = "iptables-nft"
	}

	exec.CommandContext(ctx, ipt, "-F", "GHOSTSTACK").Run()
	exec.CommandContext(ctx, ipt, "-D", "OUTPUT", "-j", "GHOSTSTACK").Run()
	exec.CommandContext(ctx, ipt, "-X", "GHOSTSTACK").Run()
	return nil
}

func (f *Firewall) ApplyKillSwitch(ctx context.Context, iface string) error {
	f.iface = iface

	switch f.mode {
	case FirewallNftables:
		return f.applyKillSwitchNftables(ctx)
	default:
		return f.applyKillSwitchIptables(ctx)
	}
}

func (f *Firewall) applyKillSwitchNftables(ctx context.Context) error {
	if err := exec.CommandContext(ctx,
		"nft", "add", "rule", "inet", "ghoststack", "output",
		"oif", "!=", f.iface, "drop",
	).Run(); err != nil {
		return fmt.Errorf("kill switch nft: %w", err)
	}

	if err := exec.CommandContext(ctx,
		"nft", "add", "rule", "inet", "ghoststack", "output",
		"oif", "lo", "accept",
	).Run(); err != nil {
		return fmt.Errorf("kill switch nft lo: %w", err)
	}

	if err := exec.CommandContext(ctx,
		"nft", "add", "rule", "inet", "ghoststack", "output",
		"udp", "dport", "53", "accept",
	).Run(); err != nil {
		return fmt.Errorf("kill switch nft dns: %w", err)
	}

	if err := exec.CommandContext(ctx,
		"nft", "add", "rule", "inet", "ghoststack", "output",
		"udp", "dport", "67-68", "accept",
	).Run(); err != nil {
		return fmt.Errorf("kill switch nft dhcp: %w", err)
	}

	return nil
}

func (f *Firewall) applyKillSwitchIptables(ctx context.Context) error {
	ipt := string(f.mode)
	if f.mode == FirewallIptablesNft {
		ipt = "iptables-nft"
	}

	rules := [][]string{
		{ipt, "-A", "GHOSTSTACK", "-o", f.iface, "-j", "ACCEPT"},
		{ipt, "-A", "GHOSTSTACK", "-o", "lo", "-j", "ACCEPT"},
		{ipt, "-A", "GHOSTSTACK", "-p", "udp", "--dport", "53", "-j", "ACCEPT"},
		{ipt, "-A", "GHOSTSTACK", "-p", "udp", "--dport", "67:68", "-j", "ACCEPT"},
		{ipt, "-A", "GHOSTSTACK", "-j", "DROP"},
	}

	for _, rule := range rules {
		if err := exec.CommandContext(ctx, rule[0], rule[1:]...).Run(); err != nil {
			return fmt.Errorf("kill switch iptables: %w", err)
		}
	}
	return nil
}

func (f *Firewall) RemoveKillSwitch(ctx context.Context) error {
	switch f.mode {
	case FirewallNftables:
		return exec.CommandContext(ctx, "nft", "flush", "rule", "inet", "ghoststack", "output").Run()
	default:
		return f.Flush(ctx)
	}
}

func (f *Firewall) Mode() FirewallMode {
	return f.mode
}

func Supported() bool {
	return detectMode() != FirewallIptables
}
