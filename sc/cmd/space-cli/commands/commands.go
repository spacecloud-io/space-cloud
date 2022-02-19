package commands

import (
	"github.com/spf13/cobra"

	"github.com/spacecloud-io/space-cloud/cmd/space-cli/commands/run"
)

// NewRootCommand returns space-cli command
func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "space-cli",
		Version:      "v0.22.0",
		SilenceUsage: true,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, args)
		},
	}

	// Add all sub commands
	cmd.AddCommand(run.NewCommand())

	return cmd
}
