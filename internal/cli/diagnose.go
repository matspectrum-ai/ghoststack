package cli

import "github.com/spf13/cobra"

func newDiagnoseCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "diagnose",
		Short: "Run diagnostics",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, _ = cmd, args
			return nil
		},
	}
}
