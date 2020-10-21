package modules

import (
	"github.com/spaceuptech/space-cloud/space-cli/cmd/modules/auth"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/modules/database"
	"github.com/spf13/cobra"
)

// FetchDeleteSubCommands fetches all the delete subcommands from different modules
func FetchDeleteSubCommands() *cobra.Command {
	var deleteCmd = &cobra.Command{
		Use:           "delete",
		Short:         "",
		SilenceErrors: true,
	}
	deleteCmd.AddCommand(auth.DeleteSubCommands()...)
	deleteCmd.AddCommand(database.DeleteSubCommands()...)

	return deleteCmd
}
