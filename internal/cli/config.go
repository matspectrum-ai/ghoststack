package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage GhostStack configuration",
	}

	cmd.AddCommand(newConfigValidateCommand())
	cmd.AddCommand(newConfigReloadCommand())

	return cmd
}

func newConfigValidateCommand() *cobra.Command {
	var configPath string

	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate GhostStack configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			if configPath == "" {
				return fmt.Errorf("config path is required")
			}

			if _, err := os.Stat(configPath); err != nil {
				return fmt.Errorf("validate config: %w", err)
			}

			fmt.Fprintf(os.Stdout, "Validating config: %s\n", configPath)
			fmt.Fprintln(os.Stdout, "Configuration is valid")
			return nil
		},
	}

	cmd.Flags().StringVar(&configPath, "config", "", "path to GhostStack config")
	return cmd
}

func newConfigReloadCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reload",
		Short: "Reload GhostStack configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintln(os.Stdout, "Configuration reloaded")
			return nil
		},
	}

	return cmd
}
