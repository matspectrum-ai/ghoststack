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
	mode     FirewallMode
	iface    string
	dnsForce bool
	dnsAddrs []string
}

type FirewallOption func(*Firewall)

func WithDNSForce(addrs []string) FirewallOption {
	return func(f *Firewall) {
		f.dnsForce = true
		f.dnsAddrs = addrs
	}
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
	chain input {
		type filter hook input priority 0; policy accept;
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
		case "allow6":
			if err := exec.CommandContext(ctx,
				"nft", "add", "rule", "inet", "ghoststack", "output",
				"ip6", "daddr", dest, "accept",
			).Run(); err != nil {
				return fmt.Errorf("nft rule allow6 %s: %w", dest, err)
			}
		case "drop":
			if err := exec.CommandContext(ctx,
				"nft", "add", "rule", "inet", "ghoststack", "output",
				"ip", "daddr", dest, "drop",
			).Run(); err != nil {
				return fmt.Errorf("nft rule drop %s: %w", dest, err)
			}
		case "drop6":
			if err := exec.CommandContext(ctx,
				"nft", "add", "rule", "inet", "ghoststack", "output",
				"ip6", "daddr", dest, "drop",
			).Run(); err != nil {
				return fmt.Errorf("nft rule drop6 %s: %w", dest, err)
			}
		case "allow_port":
			if err := exec.CommandContext(ctx,
				"nft", "add", "rule", "inet", "ghoststack", "output",
				dest, "accept",
			).Run(); err != nil {
				return fmt.Errorf("nft rule allow_port %s: %w", dest, err)
			}
		case "block_port":
			if err := exec.CommandContext(ctx,
				"nft", "add", "rule", "inet", "ghoststack", "output",
				dest, "drop",
			).Run(); err != nil {
				return fmt.Errorf("nft rule block_port %s: %w", dest, err)
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
		case "block_port":
			args = []string{ipt, "-A", "GHOSTSTACK", dest, "-j", "DROP"}
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
	rules := []struct {
		args []string
		desc string
	}{
		{[]string{"nft", "add", "rule", "inet", "ghoststack", "output", "oif", f.iface, "accept"}, "vpn iface"},
		{[]string{"nft", "add", "rule", "inet", "ghoststack", "output", "oif", "lo", "accept"}, "loopback"},
		{[]string{"nft", "add", "rule", "inet", "ghoststack", "input", "iif", "lo", "accept"}, "loopback input"},

		{[]string{"nft", "add", "rule", "inet", "ghoststack", "output", "ip", "protocol", "icmp", "drop"}, "icmp4 block"},
		{[]string{"nft", "add", "rule", "inet", "ghoststack", "output", "ip6", "nexthdr", "icmpv6", "drop"}, "icmp6 block"},

		{[]string{"nft", "add", "rule", "inet", "ghoststack", "output", "ip6", "daddr", "fe80::/10", "accept"}, "link-local v6"},
		{[]string{"nft", "add", "rule", "inet", "ghoststack", "output", "ip6", "daddr", "ff00::/8", "accept"}, "multicast v6"},
		{[]string{"nft", "add", "rule", "inet", "ghoststack", "output", "udp", "dport", "67-68", "accept"}, "dhcp"},

		{[]string{"nft", "add", "rule", "inet", "ghoststack", "output", "meta", "l4proto", "udp", "th", "dport", "53", "accept"}, "dns"},
		{[]string{"nft", "add", "rule", "inet", "ghoststack", "output", "meta", "l4proto", "tcp", "th", "dport", "53", "accept"}, "dns tcp"},

		{[]string{"nft", "add", "rule", "inet", "ghoststack", "output", "oif", "!=", f.iface, "drop"}, "global block"},
	}

	for _, r := range rules {
		if err := exec.CommandContext(ctx, r.args[0], r.args[1:]...).Run(); err != nil {
			return fmt.Errorf("kill switch nft %s: %w", r.desc, err)
		}
	}

	if f.dnsForce && len(f.dnsAddrs) > 0 {
		for _, dns := range f.dnsAddrs {
			if err := exec.CommandContext(ctx,
				"nft", "add", "rule", "inet", "ghoststack", "output",
				"udp", "dport", "53", "ip", "daddr", "!=", dns, "drop",
			).Run(); err != nil {
				return fmt.Errorf("dns force %s: %w", dns, err)
			}
		}
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
		{ipt, "-A", "GHOSTSTACK", "-p", "tcp", "--dport", "53", "-j", "ACCEPT"},
		{ipt, "-A", "GHOSTSTACK", "-p", "udp", "--dport", "67:68", "-j", "ACCEPT"},
		{ipt, "-A", "GHOSTSTACK", "-p", "icmp", "-j", "DROP"},
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
		err := exec.CommandContext(ctx, "nft", "delete", "table", "inet", "ghoststack").Run()
		if err != nil {
			return nil
		}
		return nil
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
