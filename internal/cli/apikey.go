package cli

import (
	"fmt"
	"os"
	"time"

	"github.com/ghoststack/ghoststack/internal/api"
	"github.com/ghoststack/ghoststack/internal/storage"
	"github.com/spf13/cobra"
)

func newAPIKeyCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "apikey",
		Short: "Manage API keys for remote access",
	}

	cmd.AddCommand(newAPIKeyCreateCommand())
	cmd.AddCommand(newAPIKeyListCommand())
	cmd.AddCommand(newAPIKeyRevokeCommand())

	return cmd
}

func newAPIKeyCreateCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "create <name>",
		Short: "Create a new API key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := openStorage(cmd)
			if err != nil {
				return err
			}
			defer store.Close(cmd.Context())

			keyStore := api.NewAPIKeyStore(store)
			raw, err := keyStore.Generate(cmd.Context(), args[0])
			if err != nil {
				return fmt.Errorf("generate key: %w", err)
			}

			fmt.Fprintf(os.Stdout, "Key: %s\n", raw)
			fmt.Fprintln(os.Stdout, "Store it safely — it will not be shown again.")
			return nil
		},
	}
}

func newAPIKeyListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all API keys",
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := openStorage(cmd)
			if err != nil {
				return err
			}
			defer store.Close(cmd.Context())

			keyStore := api.NewAPIKeyStore(store)
			keys, err := keyStore.List(cmd.Context())
			if err != nil {
				return fmt.Errorf("list keys: %w", err)
			}

			if len(keys) == 0 {
				fmt.Fprintln(os.Stdout, "No API keys found.")
				return nil
			}

			for _, k := range keys {
				revoked := ""
				if k.Revoked {
					revoked = " (revoked)"
				}
				lastUsed := "never"
				if k.LastUsedAt > 0 {
					lastUsed = time.Unix(k.LastUsedAt, 0).Format(time.RFC3339)
				}
				fmt.Fprintf(os.Stdout, "  %-20s  created: %s  last used: %s%s\n",
					k.Name,
					time.Unix(k.CreatedAt, 0).Format(time.RFC3339),
					lastUsed,
					revoked,
				)
			}
			return nil
		},
	}
}

func newAPIKeyRevokeCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "revoke <name>",
		Short: "Revoke an API key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := openStorage(cmd)
			if err != nil {
				return err
			}
			defer store.Close(cmd.Context())

			keyStore := api.NewAPIKeyStore(store)
			if err := keyStore.Revoke(cmd.Context(), args[0]); err != nil {
				return fmt.Errorf("revoke key: %w", err)
			}

			fmt.Fprintf(os.Stdout, "Key %q revoked.\n", args[0])
			return nil
		},
	}
}

func openStorage(cmd *cobra.Command) (*storage.SQLiteProvider, error) {
	store := storage.NewSQLiteProvider()
	if err := store.Open(cmd.Context(), ""); err != nil {
		return nil, fmt.Errorf("open storage: %w", err)
	}
	if err := store.Migrate(cmd.Context()); err != nil {
		store.Close(cmd.Context())
		return nil, fmt.Errorf("migrate: %w", err)
	}
	return store, nil
}
