package cli

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ghoststack/ghoststack/internal/agent"
	"github.com/spf13/cobra"
)

func newAgentCommand() *cobra.Command {
	var controllerURL string
	var apiKey string
	var name string

	cmd := &cobra.Command{
		Use:   "agent",
		Short: "Run GhostStack in agent mode for remote coordination",
	}

	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Start agent and connect to controller",
		RunE: func(cmd *cobra.Command, args []string) error {
			if controllerURL == "" {
				return fmt.Errorf("--controller is required")
			}
			if apiKey == "" {
				return fmt.Errorf("--api-key is required")
			}
			if name == "" {
				hostname, _ := os.Hostname()
				name = hostname
			}

			agentRuntime := agent.NewRuntime(controllerURL, apiKey, name, Version)

			ctx, cancel := context.WithCancel(cmd.Context())
			defer cancel()

			sigCh := make(chan os.Signal, 1)
			signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

			go func() {
				<-sigCh
				cancel()
			}()

			fmt.Fprintf(os.Stdout, "Agent %q connecting to %s...\n", name, controllerURL)
			if err := agentRuntime.Run(ctx); err != nil && err != context.Canceled {
				return fmt.Errorf("agent: %w", err)
			}

			return nil
		},
	}

	startCmd.Flags().StringVar(&controllerURL, "controller", "", "Controller URL (e.g. http://controller:8080)")
	startCmd.Flags().StringVar(&apiKey, "api-key", "", "API key for authentication")
	startCmd.Flags().StringVar(&name, "name", "", "Agent name (default: hostname)")

	cmd.AddCommand(startCmd)
	return cmd
}
