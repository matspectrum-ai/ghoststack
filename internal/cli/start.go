package cli

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ghoststack/ghoststack/internal/api"
	"github.com/ghoststack/ghoststack/internal/runtime"
	"github.com/spf13/cobra"
)

func newStartCommand() *cobra.Command {
	var configPath string
	var apiAddr string

	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start GhostStack runtime",
		RunE: func(cmd *cobra.Command, args []string) error {
			daemon := runtime.NewDaemon(configPath, nil)
			if err := daemon.Start(cmd.Context()); err != nil {
				return fmt.Errorf("start daemon: %w", err)
			}

			ctx, cancel := context.WithTimeout(cmd.Context(), 2*time.Second)
			defer cancel()

			server := api.NewServer(daemon)
			httpServer, err := server.Start(ctx, apiAddr)
			if err != nil {
				return fmt.Errorf("start api: %w", err)
			}
			defer httpServer.Close()

			fmt.Fprintln(os.Stdout, daemon.String())
			fmt.Fprintf(os.Stdout, "API listening on http://%s\n", apiAddr)
			return nil
		},
	}

	cmd.Flags().StringVar(&configPath, "config", "", "path to GhostStack config")
	cmd.Flags().StringVar(&apiAddr, "api-addr", "127.0.0.1:8080", "API listen address")
	return cmd
}
