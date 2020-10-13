package modules

import (
	"github.com/spaceuptech/space-cloud/space-cli/cmd/modules/auth"
	"github.com/spf13/cobra"
)

// FetchDeleteSubCommands fetches all the delete subcommands from different modules
func FetchDeleteSubCommands() *cobra.Command {
	var generateCmd = &cobra.Command{
		Use:           "delete",
		Short:         "",
		SilenceErrors: true,
	}
	generateCmd.AddCommand(auth.DeleteSubCommands()...)

	return generateCmd
}
