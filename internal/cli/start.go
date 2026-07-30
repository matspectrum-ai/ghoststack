package cli

import (
	"fmt"
	"os"

	"github.com/ghoststack/ghoststack/internal/runtime"
	"github.com/spf13/cobra"
)

func newStartCommand() *cobra.Command {
	var configPath string

	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start GhostStack runtime",
		RunE: func(cmd *cobra.Command, args []string) error {
			daemon := runtime.NewDaemon(nil)
			if err := daemon.Start(cmd.Context()); err != nil {
				return fmt.Errorf("start daemon: %w", err)
			}

			fmt.Fprintln(os.Stdout, daemon.String())
			return nil
		},
	}

	cmd.Flags().StringVar(&configPath, "config", "", "path to GhostStack config")
	return cmd
}
