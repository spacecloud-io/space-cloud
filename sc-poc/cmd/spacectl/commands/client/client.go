package client

import (
	"github.com/spf13/cobra"

	"github.com/spacecloud-io/space-cloud/cmd/spacectl/commands/client/generate"
)

// NewCommand get spacectl client command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "client",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, args)
		},
	}

	// Add all sub commands
	cmd.AddCommand(generate.NewCommand())

	return cmd
}
