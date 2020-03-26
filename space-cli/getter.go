package main

import (
	"github.com/spaceuptech/space-cli/modules"
	"github.com/spaceuptech/space-cli/modules/auth"
	"github.com/spaceuptech/space-cli/modules/database"
	"github.com/spaceuptech/space-cli/modules/eventing"
	"github.com/spaceuptech/space-cli/modules/filestore"
	"github.com/spaceuptech/space-cli/modules/ingress"
	"github.com/spaceuptech/space-cli/modules/letsencrypt"
	"github.com/spaceuptech/space-cli/modules/project"
	remoteservices "github.com/spaceuptech/space-cli/modules/remote-services"
	"github.com/spaceuptech/space-cli/modules/services"
	"github.com/spf13/cobra"
)

func FetchGetSubCommands() *cobra.Command {
	var getCmd = &cobra.Command{
		Use:   "get",
		Short: "",
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
	getCmd.AddCommand(modules.GetSubCommands()...)

	return getCmd
}
