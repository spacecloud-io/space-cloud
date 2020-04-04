package modules

import (
	"github.com/spaceuptech/space-cli/cmd/modules/database"
	"github.com/spaceuptech/space-cli/cmd/modules/eventing"
	"github.com/spaceuptech/space-cli/cmd/modules/filestore"
	"github.com/spaceuptech/space-cli/cmd/modules/ingress"
	"github.com/spaceuptech/space-cli/cmd/modules/letsencrypt"
	remoteservices "github.com/spaceuptech/space-cli/cmd/modules/remote-services"
	"github.com/spaceuptech/space-cli/cmd/modules/services"
	"github.com/spaceuptech/space-cli/cmd/modules/userman"
	"github.com/spf13/cobra"
)

// FetchGenerateSubCommands fetches all the generatesubcommands from different modules
func FetchGenerateSubCommands() *cobra.Command {
	var generateCmd = &cobra.Command{
		Use:   "generate",
		Short: "",
	}
	generateCmd.AddCommand(database.GenerateSubCommands()...)
	generateCmd.AddCommand(eventing.GenerateSubCommands()...)
	generateCmd.AddCommand(filestore.GenerateSubCommands()...)
	generateCmd.AddCommand(ingress.GenerateSubCommands()...)
	generateCmd.AddCommand(letsencrypt.GenerateSubCommands()...)
	generateCmd.AddCommand(remoteservices.GenerateSubCommands()...)
	generateCmd.AddCommand(services.GenerateSubCommands()...)
	generateCmd.AddCommand(userman.GenerateSubCommands()...)

	return generateCmd
}
