package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ghoststack/ghoststack/internal/security"
	"github.com/spf13/cobra"
)

func newKillSwitchCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "killswitch",
		Short: "Manage kill switch and leak protection",
	}

	cmd.AddCommand(newKillSwitchEnableCommand())
	cmd.AddCommand(newKillSwitchDisableCommand())
	cmd.AddCommand(newKillSwitchStatusCommand())
	cmd.AddCommand(newKillSwitchLeakTestCommand())

	return cmd
}

func newKillSwitchEnableCommand() *cobra.Command {
	var iface string
	var dnsForce bool
	var dnsAddrs []string

	cmd := &cobra.Command{
		Use:   "enable",
		Short: "Enable kill switch firewall rules",
		RunE: func(cmd *cobra.Command, args []string) error {
			ks := security.NewKillSwitch(iface,
				security.WithDNSForce(dnsAddrs),
			)
			if err := ks.Enable(cmd.Context()); err != nil {
				return fmt.Errorf("enable kill switch: %w", err)
			}
			fmt.Fprintf(os.Stdout, "Kill switch enabled (interface: %s)\n", iface)
			return nil
		},
	}

	cmd.Flags().StringVar(&iface, "interface", "ghost0", "VPN interface to allow")
	cmd.Flags().BoolVar(&dnsForce, "dns-force", false, "force DNS through tunnel only")
	cmd.Flags().StringSliceVar(&dnsAddrs, "dns", []string{"1.1.1.1", "9.9.9.9"}, "allowed DNS servers")
	return cmd
}

func newKillSwitchDisableCommand() *cobra.Command {
	var iface string

	cmd := &cobra.Command{
		Use:   "disable",
		Short: "Disable kill switch firewall rules",
		RunE: func(cmd *cobra.Command, args []string) error {
			ks := security.NewKillSwitch(iface)
			if err := ks.Disable(cmd.Context()); err != nil {
				return fmt.Errorf("disable kill switch: %w", err)
			}
			fmt.Fprintln(os.Stdout, "Kill switch disabled")
			return nil
		},
	}

	cmd.Flags().StringVar(&iface, "interface", "ghost0", "VPN interface")
	return cmd
}

func newKillSwitchStatusCommand() *cobra.Command {
	var iface string
	var jsonOut bool

	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show kill switch status",
		RunE: func(cmd *cobra.Command, args []string) error {
			ks := security.NewKillSwitch(iface)
			status := ks.Status()

			if jsonOut {
				enc := json.NewEncoder(os.Stdout)
				enc.SetIndent("", "  ")
				return enc.Encode(status)
			}

			fmt.Fprintf(os.Stdout, "Active:     %v\n", status["active"])
			fmt.Fprintf(os.Stdout, "Interface:  %s\n", status["interface"])
			fmt.Fprintf(os.Stdout, "DNS Force:  %v\n", status["dns_force"])
			fmt.Fprintf(os.Stdout, "DNS Addrs:  %v\n", status["dns_addrs"])
			return nil
		},
	}

	cmd.Flags().StringVar(&iface, "interface", "ghost0", "VPN interface")
	cmd.Flags().BoolVar(&jsonOut, "json", false, "output as JSON")
	return cmd
}

func newKillSwitchLeakTestCommand() *cobra.Command {
	var iface string
	var jsonOut bool

	cmd := &cobra.Command{
		Use:   "leak-test",
		Short: "Run DNS/IPv6/ICMP leak test",
		RunE: func(cmd *cobra.Command, args []string) error {
			ks := security.NewKillSwitch(iface)
			result, err := ks.RunLeakTest(cmd.Context())
			if err != nil {
				return fmt.Errorf("leak test: %w", err)
			}

			if jsonOut {
				enc := json.NewEncoder(os.Stdout)
				enc.SetIndent("", "  ")
				return enc.Encode(result)
			}

			fmt.Fprintf(os.Stdout, "DNS Leak:     %v\n", result.DNS)
			fmt.Fprintf(os.Stdout, "IPv6 Leak:    %v\n", result.IPv6)
			fmt.Fprintf(os.Stdout, "ICMP Leak:    %v\n", result.ICMP)
			if result.PublicIP != "" {
				fmt.Fprintf(os.Stdout, "Public IP:    %s\n", result.PublicIP)
			}
			if result.Error != "" {
				fmt.Fprintf(os.Stdout, "Error:        %s\n", result.Error)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&iface, "interface", "ghost0", "VPN interface")
	cmd.Flags().BoolVar(&jsonOut, "json", false, "output as JSON")
	return cmd
}
