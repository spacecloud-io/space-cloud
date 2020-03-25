package project

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/spaceuptech/space-cli/utils"
)

// Commands dis the list of commands the project module exposes
func Commands() []*cobra.Command {

	var getSubCommands = &cobra.Command{}

	var getproject = &cobra.Command{
		Use:  "project",
		RunE: actionGetProjectConfig,
	}

	getSubCommands.AddCommand(getproject)

	command := make([]*cobra.Command, 0)
	command = append(command, getSubCommands)
	return command
}

// // GetSubCommands is the list of commands the project module exposes
// var GetSubCommands = []cli.Command{
// 	{
// 		Name:   "project",
// 		Action: actionGetProjectConfig,
// 	},
// }

func actionGetProjectConfig(cmd *cobra.Command, args []string) error {
	// Get the project and cmd parameters
	project := viper.GetString("project")
	commandName := cmd.CalledAs()

	params := map[string]string{}
	obj, err := GetProjectConfig(project, commandName, params)
	if err != nil {
		return err
	}
	if err := utils.PrintYaml(obj); err != nil {
		return err
	}
	return nil
}
