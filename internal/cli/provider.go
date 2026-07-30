package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newProviderCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "provider",
		Short: "Manage GhostStack providers",
	}

	cmd.AddCommand(newProviderListCommand())
	cmd.AddCommand(newProviderSelectCommand())

	return cmd
}

func newProviderListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List available providers",
		RunE: func(cmd *cobra.Command, args []string) error {
			providerList := []struct {
				name  string
				type_ string
			}{
				{"wireguard", "vpn"},
				{"tor", "vpn"},
				{"sing-box", "gateway"},
				{"unbound", "dns"},
				{"socks5", "proxy"},
				{"openvpn", "vpn"},
				{"dns-over-https", "dns"},
				{"firewall-real", "firewall"},
			}

			for _, p := range providerList {
				fmt.Fprintf(os.Stdout, "- %s (%s)\n", p.name, p.type_)
			}
			return nil
		},
	}

	return cmd
}

func newProviderSelectCommand() *cobra.Command {
	var providerType string

	cmd := &cobra.Command{
		Use:   "select <name>",
		Short: "Select a provider by name",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			switch providerType {
			case "vpn", "dns", "firewall", "gateway":
			default:
			}

			fmt.Fprintf(os.Stdout, "Provider %s selected\n", name)
			return nil
		},
	}

	cmd.Flags().StringVar(&providerType, "type", "", "provider type")
	return cmd
}
