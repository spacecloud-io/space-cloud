package modules

import (
	"github.com/spaceuptech/space-cloud/space-cli/cmd/modules/auth"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/modules/database"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/modules/eventing"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/modules/filestore"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/modules/ingress"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/modules/letsencrypt"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/modules/project"
	remoteservices "github.com/spaceuptech/space-cloud/space-cli/cmd/modules/remote-services"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/modules/services"
	"github.com/spf13/cobra"
)

// FetchGetSubCommands fetches all the generatesubcommands from different modules
func FetchGetSubCommands() *cobra.Command {
	var getCmd = &cobra.Command{
		Use:           "get",
		Short:         "",
		SilenceErrors: true,
	}
	getCmd.AddCommand(auth.GetSubCommands()...)
	getCmd.AddCommand(database.GetSubCommands()...)
	getCmd.AddCommand(eventing.GetSubCommands()...)
	getCmd.AddCommand(filestore.GetSubCommands()...)
	getCmd.AddCommand(ingress.GetSubCommands()...)
	getCmd.AddCommand(letsencrypt.GetSubCommands()...)
	getCmd.AddCommand(project.GetSubCommands()...)
	getCmd.AddCommand(remoteservices.GetSubCommands()...)
	getCmd.AddCommand(services.GetSubCommands()...)
	getCmd.AddCommand(getSubCommands()...)

	return getCmd
}
