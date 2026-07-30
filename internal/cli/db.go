package cli

import (
	"fmt"
	"os"

	"github.com/ghoststack/ghoststack/internal/storage"
	"github.com/spf13/cobra"
)

func newDBCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "db",
		Short: "Database inspection commands",
	}

	cmd.AddCommand(newDBStatusCommand())
	cmd.AddCommand(newDBAuditCommand())

	return cmd
}

func newDBStatusCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show persisted daemon status",
		RunE: func(cmd *cobra.Command, args []string) error {
			store := storage.NewSQLiteProvider()
			if err := store.Open(cmd.Context(), ""); err != nil {
				return fmt.Errorf("open storage: %w", err)
			}
			defer store.Close(cmd.Context())

			state, err := store.LoadRuntimeState(cmd.Context())
			if err != nil {
				return fmt.Errorf("load: %w", err)
			}
			if state == nil {
				fmt.Fprintln(os.Stdout, "No persisted state found")
				return nil
			}

			fmt.Fprintf(os.Stdout, "Status:     %s\n", state.Status)
			fmt.Fprintf(os.Stdout, "Mode:       %s\n", state.Mode)
			fmt.Fprintf(os.Stdout, "Started:    %d\n", state.StartedAt)
			fmt.Fprintf(os.Stdout, "Updated:    %d\n", state.UpdatedAt)

			providers, err := store.LoadProviderStates(cmd.Context())
			if err != nil {
				return fmt.Errorf("load providers: %w", err)
			}
			if len(providers) > 0 {
				fmt.Fprintln(os.Stdout, "\nProviders:")
				for _, p := range providers {
					fmt.Fprintf(os.Stdout, "  %s: %s\n", p.Name, p.State)
				}
			}
			return nil
		},
	}
}

func newDBAuditCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "audit",
		Short: "Show audit log",
		RunE: func(cmd *cobra.Command, args []string) error {
			store := storage.NewSQLiteProvider()
			if err := store.Open(cmd.Context(), ""); err != nil {
				return fmt.Errorf("open storage: %w", err)
			}
			defer store.Close(cmd.Context())

			entries, err := store.QueryAuditLog(cmd.Context(), 50)
			if err != nil {
				return fmt.Errorf("query: %w", err)
			}
			if len(entries) == 0 {
				fmt.Fprintln(os.Stdout, "No audit entries found")
				return nil
			}

			for _, e := range entries {
				line := fmt.Sprintf("[%d] %s", e.Timestamp, e.Action)
				if e.Source != "" {
					line += fmt.Sprintf(" (source: %s)", e.Source)
				}
				if e.Detail != "" {
					line += fmt.Sprintf(" — %s", e.Detail)
				}
				fmt.Fprintln(os.Stdout, line)
			}
			return nil
		},
	}
}
