package cli

import (
	"fmt"
	"os"

	"github.com/ghoststack/ghoststack/internal/runtime"
	"github.com/spf13/cobra"
)

func newStopCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop GhostStack runtime",
		RunE: func(cmd *cobra.Command, args []string) error {
			daemon := runtime.NewDaemon(nil)
			if err := daemon.Stop(cmd.Context()); err != nil {
				return fmt.Errorf("stop daemon: %w", err)
			}
			fmt.Fprintln(os.Stdout, daemon.String())
			return nil
		},
	}

	return cmd
}
