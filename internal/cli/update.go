package cli

import "github.com/spf13/cobra"

func newUpdateCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "update",
		Short: "Update GhostStack",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, _ = cmd, args
			return nil
		},
	}
}
