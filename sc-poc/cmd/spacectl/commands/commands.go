package commands

import (
	"github.com/spf13/cobra"

	"github.com/spacecloud-io/space-cloud/cmd/spacectl/commands/client"
	"github.com/spacecloud-io/space-cloud/cmd/spacectl/commands/login"
	"github.com/spacecloud-io/space-cloud/cmd/spacectl/commands/pkg"
	"github.com/spacecloud-io/space-cloud/cmd/spacectl/commands/run"
)

// NewRootCommand returns spacectl command
func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "spacectl",
		Version: "v0.22.0",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, args)
		},
	}

	// Add all sub commands
	cmd.AddCommand(run.NewCommand())
	cmd.AddCommand(client.NewCommand())
	cmd.AddCommand(pkg.NewCommand())
	cmd.AddCommand(login.NewCommand())

	return cmd
}
