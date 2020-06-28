package modules

import (
	"github.com/spf13/cobra"

	"github.com/spaceuptech/space-cli/cmd/modules/database"
	"github.com/spaceuptech/space-cli/cmd/modules/eventing"
	"github.com/spaceuptech/space-cli/cmd/modules/filestore"
	"github.com/spaceuptech/space-cli/cmd/modules/ingress"
	"github.com/spaceuptech/space-cli/cmd/modules/letsencrypt"
	"github.com/spaceuptech/space-cli/cmd/modules/project"
	remoteservices "github.com/spaceuptech/space-cli/cmd/modules/remote-services"
	"github.com/spaceuptech/space-cli/cmd/modules/services"
	"github.com/spaceuptech/space-cli/cmd/modules/userman"
)

// FetchGenerateSubCommands fetches all the generatesubcommands from different modules
func FetchGenerateSubCommands() *cobra.Command {
	var generateCmd = &cobra.Command{
		Use:           "generate",
		Short:         "",
		SilenceErrors: true,
	}
	generateCmd.AddCommand(database.GenerateSubCommands()...)
	generateCmd.AddCommand(eventing.GenerateSubCommands()...)
	generateCmd.AddCommand(filestore.GenerateSubCommands()...)
	generateCmd.AddCommand(ingress.GenerateSubCommands()...)
	generateCmd.AddCommand(letsencrypt.GenerateSubCommands()...)
	generateCmd.AddCommand(remoteservices.GenerateSubCommands()...)
	generateCmd.AddCommand(services.GenerateSubCommands()...)
	generateCmd.AddCommand(userman.GenerateSubCommands()...)
	generateCmd.AddCommand(project.GenerateSubCommands()...)

	return generateCmd
}
