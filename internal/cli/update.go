package cli

import (
	"fmt"
	"os"

	"github.com/ghoststack/ghoststack/internal/update"
	"github.com/spf13/cobra"
)

func newUpdateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Manage GhostStack updates",
	}

	cmd.AddCommand(newUpdateCheckCommand())
	cmd.AddCommand(newUpdateRollbackCommand())

	return cmd
}

func newUpdateCheckCommand() *cobra.Command {
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "check",
		Short: "Check for updates",
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := update.NewUpdateManager(nil)
			result, err := manager.Check(cmd.Context())
			if err != nil {
				return fmt.Errorf("check update: %w", err)
			}

			if !result.Available {
				fmt.Fprintln(os.Stdout, "No updates available")
				return nil
			}

			fmt.Fprintf(os.Stdout, "Update available: %s\n", result.Manifest.Version)

			if dryRun {
				fmt.Fprintln(os.Stdout, "Dry run - no changes made")
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "perform a dry run")
	return cmd
}

func newUpdateRollbackCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rollback",
		Short: "Rollback to previous version",
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := update.NewUpdateManager(nil)
			if err := manager.Rollback(cmd.Context()); err != nil {
				return fmt.Errorf("rollback: %w", err)
			}
			fmt.Fprintln(os.Stdout, "Rollback completed")
			return nil
		},
	}

	return cmd
}
