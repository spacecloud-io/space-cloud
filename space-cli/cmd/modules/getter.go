package modules

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/spaceuptech/space-cloud/space-cli/cmd/modules/auth"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/modules/database"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/modules/eventing"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/modules/filestore"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/modules/ingress"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/modules/letsencrypt"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/modules/project"
	remoteservices "github.com/spaceuptech/space-cloud/space-cli/cmd/modules/remote-services"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/modules/services"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils"
)

// FetchGetSubCommands fetches all the generatesubcommands from different modules
func FetchGetSubCommands() *cobra.Command {
	var getCmd = &cobra.Command{
		Use:              "get",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {},
		Short:            "",
		SilenceErrors:    true,
	}
	getCmd.PersistentFlags().StringSliceP("filter", "", []string{}, "Filter ingress routes based on services, target-host, request-host & url")
	err := viper.BindPFlag("filter", getCmd.PersistentFlags().Lookup("filter"))
	if err != nil {
		_ = utils.LogError("Unable to bind the flag ('filter')", nil)
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
