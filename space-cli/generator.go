package main

import (
	"github.com/spaceuptech/space-cli/modules/database"
	"github.com/spaceuptech/space-cli/modules/eventing"
	"github.com/spaceuptech/space-cli/modules/filestore"
	"github.com/spaceuptech/space-cli/modules/ingress"
	"github.com/spaceuptech/space-cli/modules/letsencrypt"
	remoteservices "github.com/spaceuptech/space-cli/modules/remote-services"
	"github.com/spaceuptech/space-cli/modules/services"
	"github.com/spaceuptech/space-cli/modules/userman"
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
