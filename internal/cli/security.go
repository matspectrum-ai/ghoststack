package cli

import "github.com/spf13/cobra"

func newSecurityCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "security",
		Short: "Security operations",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, _ = cmd, args
			return nil
		},
	}
}
