package cli

import "github.com/spf13/cobra"

func newStatusCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show GhostStack status",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, _ = cmd, args
			return nil
		},
	}
}
