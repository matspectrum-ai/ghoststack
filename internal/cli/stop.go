package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/ghoststack/ghoststack/internal/providers"
	"github.com/ghoststack/ghoststack/internal/runtime"
	"github.com/spf13/cobra"
)

func newStopCommand() *cobra.Command {
	var configPath string

	cmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop GhostStack runtime",
		RunE: func(cmd *cobra.Command, args []string) error {
			engine := providers.NewProviderEngine()
			providers.RegisterBuiltins(engine)
			engine.StopAll(cmd.Context())

			daemon := runtime.NewDaemon(configPath, nil)
			if err := daemon.Stop(cmd.Context()); err != nil {
				if errors.Is(err, runtime.ErrNotStarted) {
					fmt.Fprintln(os.Stdout, "GhostStack idle")
					return nil
				}
				return fmt.Errorf("stop daemon: %w", err)
			}
			fmt.Fprintln(os.Stdout, daemon.String())
			return nil
		},
	}

	cmd.Flags().StringVar(&configPath, "config", "", "path to GhostStack config")
	return cmd
}

func newEmergencyStopCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "emergency-stop",
		Short: "Force stop all GhostStack processes",
		Long: `Force stop all GhostStack processes and cleanup.
Kills any ghoststack-related processes and flushes firewall rules.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			engine := providers.NewProviderEngine()
			providers.RegisterBuiltins(engine)
			engine.StopAll(cmd.Context())

			fmt.Fprintln(os.Stdout, "All GhostStack processes stopped")
			return nil
		},
	}
}
