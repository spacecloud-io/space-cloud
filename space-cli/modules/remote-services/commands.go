package remoteservices

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/spaceuptech/space-cli/utils"
)

// Commands is the list of commands the remote-services module exposes
func Commands() []*cobra.Command {
	var generateSubCommands = &cobra.Command{}

	var generateService = &cobra.Command{
		Use:  "remote-services",
		RunE: actionGenerateService,
	}

	var getSubCommands = &cobra.Command{}

	var getService = &cobra.Command{
		Use:  "remote-services",
		RunE: actionGetRemoteServices,
	}

	generateSubCommands.AddCommand(generateService)
	getSubCommands.AddCommand(getService)

	command := make([]*cobra.Command, 0)
	command = append(command, generateSubCommands)
	command = append(command, getSubCommands)
	return command
}

// // GenerateSubCommands is the list of commands the remoteservices module exposes
// var GenerateSubCommands = []cli.Command{
// 	{
// 		Name:   "remote-services",
// 		Action: actionGenerateService,
// 	},
// }

// // GetSubCommands is the list of commands the remoteservices module exposes
// var GetSubCommands = []cli.Command{
// 	{
// 		Name:   "remote-services",
// 		Action: actionGetRemoteServices,
// 	},
// }

func actionGetRemoteServices(cmd *cobra.Command, args []string) error {
	// Get the project and url parameters
	project := viper.GetString("project")
	commandName := cmd.CalledAs()

	params := map[string]string{}
	if len(args) != 0 {
		params["id"] = args[0]
	}

	objs, err := GetRemoteServices(project, commandName, params)
	if err != nil {
		return err
	}
	if err := utils.PrintYaml(objs); err != nil {
		return err
	}
	return nil
}

func actionGenerateService(cmd *cobra.Command, args []string) error {
	argsArr := args
	if len(argsArr) != 1 {
		return fmt.Errorf("incorrect number of arguments")
	}
	dbruleConfigFile := argsArr[0]
	dbrule, err := generateService()
	if err != nil {
		return err
	}

	return utils.AppendConfigToDisk(dbrule, dbruleConfigFile)
}
