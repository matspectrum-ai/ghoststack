package cli

import "github.com/spf13/cobra"

func newStartCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start GhostStack runtime",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, _ = cmd, args
			return nil
		},
	}
	return cmd
}
