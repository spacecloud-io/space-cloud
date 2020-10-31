package modules

import (
	"github.com/spaceuptech/space-cloud/space-cli/cmd/modules/auth"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/modules/database"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/modules/eventing"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/modules/filestore"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/modules/ingress"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/modules/letsencrypt"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/modules/project"
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
	deleteCmd.AddCommand(ingress.DeleteSubCommands()...)
	deleteCmd.AddCommand(filestore.DeleteSubCommands()...)
	deleteCmd.AddCommand(eventing.DeleteSubCommands()...)
	deleteCmd.AddCommand(letsencrypt.DeleteSubCommands()...)
	deleteCmd.AddCommand(project.DeleteSubCommands()...)

	return deleteCmd
}
