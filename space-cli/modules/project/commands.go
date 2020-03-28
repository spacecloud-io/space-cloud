package project

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/spaceuptech/space-cli/utils"
)

// GetSubCommands dis the list of commands the project module exposes
func GetSubCommands() []*cobra.Command {

	var getproject = &cobra.Command{
		Use:  "project",
		RunE: actionGetProjectConfig,
	}

	return []*cobra.Command{getproject}
}

func actionGetProjectConfig(cmd *cobra.Command, args []string) error {
	// Get the project and cmd parameters
	project := viper.GetString("project")
	commandName := cmd.Use

	params := map[string]string{}
	obj, _ := GetProjectConfig(project, commandName, params)
	_ = utils.PrintYaml(obj)
	return nil
}
