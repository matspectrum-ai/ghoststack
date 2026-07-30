package cli

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/ghoststack/ghoststack/internal/security"
	"github.com/spf13/cobra"
	"golang.org/x/sys/unix"
)

type diagnoseCheck struct {
	name  string
	check func() error
	fix   string
}

var diagnoseChecks = []diagnoseCheck{
	{
		name: "CAP_NET_ADMIN",
		check: func() error {
			f, err := os.OpenFile("/dev/net/tun", os.O_RDONLY, 0)
			if err != nil {
				return fmt.Errorf("cannot access /dev/net/tun: %w", err)
			}
			f.Close()
			return nil
		},
		fix: "Run with sudo or add CAP_NET_ADMIN: sudo setcap cap_net_admin+ep $(which ghost)",
	},
	{
		name: "wg binary",
		check: func() error {
			if _, err := exec.LookPath("wg"); err != nil {
				return fmt.Errorf("wg not found in PATH")
			}
			return nil
		},
		fix: "Install wireguard-tools: sudo apt install wireguard-tools",
	},
	{
		name: "wg-quick binary",
		check: func() error {
			if _, err := exec.LookPath("wg-quick"); err != nil {
				return fmt.Errorf("wg-quick not found in PATH")
			}
			return nil
		},
		fix: "Install wireguard-tools: sudo apt install wireguard-tools",
	},
	{
		name: "nftables or iptables",
		check: func() error {
			if err := exec.Command("nft", "--version").Run(); err == nil {
				return nil
			}
			if err := exec.Command("iptables-nft", "--version").Run(); err == nil {
				return nil
			}
			if err := exec.Command("iptables-legacy", "--version").Run(); err == nil {
				return nil
			}
			return fmt.Errorf("no firewall tool found (nft, iptables-nft, or iptables-legacy)")
		},
		fix: "Install nftables: sudo apt install nftables",
	},
	{
		name: "/dev/net/tun",
		check: func() error {
			info, err := os.Stat("/dev/net/tun")
			if err != nil {
				return fmt.Errorf("cannot stat /dev/net/tun: %w", err)
			}
			if info.Mode()&os.ModeDevice == 0 {
				return fmt.Errorf("/dev/net/tun is not a device")
			}
			return nil
		},
		fix: "Enable TUN module: sudo modprobe tun",
	},
	{
		name: "resolvectl",
		check: func() error {
			if _, err := exec.LookPath("resolvectl"); err != nil {
				return fmt.Errorf("resolvectl not found in PATH")
			}
			return nil
		},
		fix: "Install systemd-resolved: sudo apt install systemd-resolved",
	},
	{
		name: "WireGuard kernel module",
		check: func() error {
			if _, err := os.Stat("/sys/module/wireguard"); os.IsNotExist(err) {
				if err := exec.Command("modprobe", "wireguard").Run(); err != nil {
					return fmt.Errorf("wireguard kernel module not loaded: %w", err)
				}
			}
			return nil
		},
		fix: "Load module: sudo modprobe wireguard",
	},
	{
		name: "seccomp syscall filter",
		check: func() error {
			if err := unix.Prctl(unix.PR_GET_SECCOMP, 0, 0, 0, 0); err != nil {
				return fmt.Errorf("seccomp not available: %w", err)
			}
			return nil
		},
		fix: "Ensure kernel 3.17+ with CONFIG_SECCOMP=y",
	},
	{
		name: "binary integrity",
		check: func() error {
			exe, err := os.Executable()
			if err != nil {
				return fmt.Errorf("get executable path: %w", err)
			}
			failures, err := security.NewSecureBoot("").Verify(context.Background(), exe)
			if err != nil && len(failures) > 0 {
				if failures[0] == "secure boot: no expected hash configured" {
					return nil
				}
				return fmt.Errorf("integrity: %v", failures)
			}
			return nil
		},
		fix: "Set GHOSTSTACK_EXPECTED_HASH env var or check binary permissions",
	},
}

func newDiagnoseCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "diagnose",
		Short: "Run system diagnostics",
		Long: `Check system prerequisites for GhostStack:
  - CAP_NET_ADMIN capability
  - WireGuard tools (wg, wg-quick)
  - Firewall (nftables/iptables)
  - TUN device (/dev/net/tun)
  - DNS resolver (resolvectl)
  - WireGuard kernel module`,
		RunE: func(cmd *cobra.Command, args []string) error {
			passed := 0
			failed := 0

			fmt.Fprintln(os.Stdout, "GhostStack System Diagnostics")
			fmt.Fprintln(os.Stdout, "=============================")
			fmt.Fprintln(os.Stdout)

			for _, check := range diagnoseChecks {
				fmt.Fprintf(os.Stdout, "  [ ] %s... ", check.name)
				if err := check.check(); err != nil {
					fmt.Fprintln(os.Stdout, "FAIL")
					fmt.Fprintf(os.Stdout, "       Error: %v\n", err)
					fmt.Fprintf(os.Stdout, "       Fix: %s\n", check.fix)
					failed++
				} else {
					fmt.Fprintln(os.Stdout, "OK")
					passed++
				}
			}

			fmt.Fprintln(os.Stdout)
			fmt.Fprintf(os.Stdout, "Result: %d/%d checks passed\n", passed, passed+failed)

			if failed > 0 {
				fmt.Fprintln(os.Stdout, "Some checks failed. Run 'ghost start' with --force to skip.")
				fmt.Fprintln(os.Stdout, "Install missing dependencies and re-run 'ghost diagnose'.")
			} else {
				fmt.Fprintln(os.Stdout, "All checks passed. Ready to start GhostStack.")
			}

			if failed > 0 {
				return fmt.Errorf("%d diagnostics failed", failed)
			}
			return nil
		},
	}
}

func runDiagnostics() error {
	for _, check := range diagnoseChecks {
		if err := check.check(); err != nil {
			return fmt.Errorf("%s: %w. Fix: %s", check.name, err, check.fix)
		}
	}
	return nil
}
