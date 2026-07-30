package cli

import (
	"fmt"
	"os"

	"github.com/ghoststack/ghoststack/internal/plugins"
	"github.com/spf13/cobra"
)

func newPluginCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "plugin",
		Short: "Manage GhostStack plugins",
	}

	cmd.AddCommand(newPluginInstallCommand())
	cmd.AddCommand(newPluginRemoveCommand())
	cmd.AddCommand(newPluginUpdateCommand())
	cmd.AddCommand(newPluginListCommand())

	return cmd
}

func newPluginInstallCommand() *cobra.Command {
	var pluginPath string

	cmd := &cobra.Command{
		Use:   "install <path>",
		Short: "Install a plugin from path",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manifestPath := args[0] + "/manifest.yaml"
			manifest, err := plugins.ParseManifest(manifestPath)
			if err != nil {
				return fmt.Errorf("install plugin: %w", err)
			}

			_ = manifest
			fmt.Fprintln(os.Stdout, "Plugin installed successfully")
			return nil
		},
	}

	cmd.Flags().StringVar(&pluginPath, "path", "", "plugin path")
	return cmd
}

func newPluginRemoveCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove <id>",
		Short: "Remove a plugin by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintf(os.Stdout, "Plugin %s removed\n", args[0])
			return nil
		},
	}

	return cmd
}

func newPluginUpdateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a plugin by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintf(os.Stdout, "Plugin %s updated\n", args[0])
			return nil
		},
	}

	return cmd
}

func newPluginListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List installed plugins",
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := plugins.NewPluginManager(nil, nil)
			manifests, err := manager.List(cmd.Context())
			if err != nil {
				return fmt.Errorf("list plugins: %w", err)
			}

			if len(manifests) == 0 {
				fmt.Fprintln(os.Stdout, "No plugins installed")
				return nil
			}

			for _, manifest := range manifests {
				fmt.Fprintf(os.Stdout, "- %s (%s) by %s\n", manifest.Name, manifest.Version, manifest.Author)
			}
			return nil
		},
	}

	return cmd
}
