package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	Version   = "0.0.0"
	Commit    = "unknown"
	BuildTime = "unknown"
)

func newVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintf(os.Stdout, "GhostStack %s (%s)\n", Version, Commit)
			fmt.Fprintf(os.Stdout, "Built: %s\n", BuildTime)
			return nil
		},
	}
}
