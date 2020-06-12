package project

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/spaceuptech/space-cli/cmd/utils"
)

// GenerateSubCommands is the list of commands the project module exposes
func GenerateSubCommands() []*cobra.Command {

	var generateService = &cobra.Command{
		Use:  "project",
		RunE: actionGenerateProject,
	}
	return []*cobra.Command{generateService}
}

// GetSubCommands dis the list of commands the project module exposes
func GetSubCommands() []*cobra.Command {

	var getproject = &cobra.Command{
		Use:  "project",
		RunE: actionGetProjectConfig,
	}

	var getprojects = &cobra.Command{
		Use:  "projects",
		RunE: actionGetProjectConfig,
	}

	return []*cobra.Command{getproject, getprojects}
}

func actionGetProjectConfig(cmd *cobra.Command, args []string) error {
	// Get the project and cmd parameters
	project := viper.GetString("project")
	commandName := "project"

	params := map[string]string{}
	if len(args) != 0 {
		params["id"] = args[0]
	}
	obj, err := GetProjectConfig(project, commandName, params)
	if err != nil {
		return err
	}

	if err := utils.PrintYaml(obj); err != nil {
		return err
	}
	return nil
}

func actionGenerateProject(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return utils.LogError("incorrect number of arguments", nil)
	}
	projectFilePath := args[0]
	project, err := generateProject()
	if err != nil {
		return err
	}

	return utils.AppendConfigToDisk(project, projectFilePath)
}
