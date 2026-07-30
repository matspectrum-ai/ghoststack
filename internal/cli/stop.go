package cli

import "github.com/spf13/cobra"

func newStopCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "stop",
		Short: "Stop GhostStack runtime",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, _ = cmd, args
			return nil
		},
	}
}
