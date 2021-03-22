package project

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils"
)

// GenerateSubCommands is the list of commands the project module exposes
func GenerateSubCommands() []*cobra.Command {

	var generateProject = &cobra.Command{
		Use:     "project [path to config file]",
		RunE:    actionGenerateProject,
		Aliases: []string{"projects"},
		Example: "space-cli generate project config.yaml --project myproject --log-level info",
	}
	return []*cobra.Command{generateProject}
}

// GetSubCommands is the list of commands the project module exposes
func GetSubCommands() []*cobra.Command {

	var getprojects = &cobra.Command{
		Use:               "projects",
		Aliases:           []string{"project"},
		RunE:              actionGetProjectConfig,
		ValidArgsFunction: projectAutoCompletionFun,
	}

	return []*cobra.Command{getprojects}
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
		return utils.LogError("incorrect number of arguments. Use -h to check usage instructions", nil)
	}
	projectFilePath := args[0]
	project, err := generateProject()
	if err != nil {
		return err
	}

	return utils.AppendConfigToDisk(project, projectFilePath)
}

// DeleteSubCommands is the list of commands the project module exposes
func DeleteSubCommands() []*cobra.Command {

	var getprojects = &cobra.Command{
		Use:               "projects",
		Aliases:           []string{"project"},
		RunE:              actionDeleteProjectConfig,
		ValidArgsFunction: projectAutoCompletionFun,
		Example:           "space-cli delete projects --project myproject",
	}

	return []*cobra.Command{getprojects}
}

func actionDeleteProjectConfig(cmd *cobra.Command, args []string) error {
	// Get the project
	project, check := utils.GetProjectID()
	if !check {
		return utils.LogError("Project not specified in flag", nil)
	}

	return DeleteProject(project)
}
