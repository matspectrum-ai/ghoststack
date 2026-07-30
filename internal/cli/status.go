package cli

import (
	"fmt"
	"os"

	"github.com/ghoststack/ghoststack/internal/runtime"
	"github.com/spf13/cobra"
)

func newStatusCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show GhostStack status",
		RunE: func(cmd *cobra.Command, args []string) error {
			daemon := runtime.NewDaemon("", nil)
			fmt.Fprintln(os.Stdout, daemon.String())
			return nil
		},
	}

	return cmd
}
