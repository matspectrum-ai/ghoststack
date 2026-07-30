package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/ghoststack/ghoststack/internal/secrets"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

func newSecretsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "secrets",
		Short: "Manage encrypted secrets",
		Long: `Encrypted secret storage using AES-256-GCM with argon2id key derivation.

Secrets are stored in ~/.ghoststack/secrets.enc and never written to
disk in plaintext. Use 'ghost init --passphrase' to initialize.`,
	}

	cmd.AddCommand(newSecretsInitCmd())
	cmd.AddCommand(newSecretsSetCmd())
	cmd.AddCommand(newSecretsGetCmd())
	cmd.AddCommand(newSecretsListCmd())
	cmd.AddCommand(newSecretsDeleteCmd())

	return cmd
}

func readPassphrase() (string, error) {
	fmt.Fprint(os.Stderr, "Passphrase: ")
	bytepw, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Fprintln(os.Stderr)
	if err != nil {
		return "", fmt.Errorf("read passphrase: %w", err)
	}
	pass := strings.TrimSpace(string(bytepw))
	if pass == "" {
		return "", fmt.Errorf("passphrase must not be empty")
	}
	return pass, nil
}

func homeDir() string {
	if d := os.Getenv("GHOSTSTACK_HOME"); d != "" {
		return d
	}
	return filepath.Join(os.Getenv("HOME"), ".ghoststack")
}

func newSecretsInitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initialize encrypted secrets store",
		RunE: func(cmd *cobra.Command, args []string) error {
			passphrase, err := readPassphrase()
			if err != nil {
				return err
			}

			sm := secrets.NewSecretsManager(homeDir())
			if err := sm.Init(passphrase); err != nil {
				return fmt.Errorf("init secrets: %w", err)
			}

			fmt.Fprintln(os.Stderr, "Secrets store initialized at", filepath.Join(homeDir(), "secrets.enc"))
			return nil
		},
	}
}

func newSecretsSetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "set <name> <value>",
		Short: "Set a secret",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			passphrase, err := readPassphrase()
			if err != nil {
				return err
			}

			sm := secrets.NewSecretsManager(homeDir())
			if err := sm.Load(passphrase); err != nil {
				return fmt.Errorf("load secrets: %w", err)
			}

			if err := sm.Set(args[0], args[1]); err != nil {
				return fmt.Errorf("set: %w", err)
			}

			if err := sm.Save(passphrase); err != nil {
				return fmt.Errorf("save: %w", err)
			}

			fmt.Fprintf(os.Stderr, "Secret %q set\n", args[0])
			return nil
		},
	}
}

func newSecretsGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <name>",
		Short: "Get a secret value",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			passphrase, err := readPassphrase()
			if err != nil {
				return err
			}

			sm := secrets.NewSecretsManager(homeDir())
			if err := sm.Load(passphrase); err != nil {
				return fmt.Errorf("load secrets: %w", err)
			}

			val, err := sm.Get(args[0])
			if err != nil {
				return fmt.Errorf("get: %w", err)
			}

			fmt.Println(val)
			return nil
		},
	}
}

func newSecretsListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List secret names",
		RunE: func(cmd *cobra.Command, args []string) error {
			passphrase, err := readPassphrase()
			if err != nil {
				return err
			}

			sm := secrets.NewSecretsManager(homeDir())
			if err := sm.Load(passphrase); err != nil {
				return fmt.Errorf("load secrets: %w", err)
			}

			keys := sm.List()
			if len(keys) == 0 {
				fmt.Fprintln(os.Stderr, "No secrets stored")
				return nil
			}

			for _, k := range keys {
				fmt.Println(k)
			}
			return nil
		},
	}
}

func newSecretsDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <name>",
		Short: "Delete a secret",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			passphrase, err := readPassphrase()
			if err != nil {
				return err
			}

			sm := secrets.NewSecretsManager(homeDir())
			if err := sm.Load(passphrase); err != nil {
				return fmt.Errorf("load secrets: %w", err)
			}

			if err := sm.Delete(args[0]); err != nil {
				return fmt.Errorf("delete: %w", err)
			}

			if err := sm.Save(passphrase); err != nil {
				return fmt.Errorf("save: %w", err)
			}

			fmt.Fprintf(os.Stderr, "Secret %q deleted\n", args[0])
			return nil
		},
	}
}
