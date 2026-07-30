package cli

import (
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ghost",
		Short: "GhostStack privacy orchestration CLI",
	}

	cmd.AddCommand(newStartCommand())
	cmd.AddCommand(newStopCommand())
	cmd.AddCommand(newStatusCommand())
	cmd.AddCommand(newDiagnoseCommand())
	cmd.AddCommand(newSecurityCommand())
	cmd.AddCommand(newPluginCommand())
	cmd.AddCommand(newConfigCommand())
	cmd.AddCommand(newProviderCommand())
	cmd.AddCommand(newUpdateCommand())

	return cmd
}

func execute(cmd *cobra.Command, args []string) error {
	return cmd.Help()
}
