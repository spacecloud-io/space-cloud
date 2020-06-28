package remoteservices

import (
	"github.com/spf13/cobra"

	"github.com/spaceuptech/space-cli/cmd/utils"
)

// GenerateSubCommands is the list of commands the remote-services module exposes
func GenerateSubCommands() []*cobra.Command {

	var generateService = &cobra.Command{
		Use:  "remote-services",
		RunE: actionGenerateService,
	}
	return []*cobra.Command{generateService}
}

// GetSubCommands is the list of commands the remote-services module exposes
func GetSubCommands() []*cobra.Command {

	var getService = &cobra.Command{
		Use:  "remote-service",
		RunE: actionGetRemoteServices,
	}

	var getServices = &cobra.Command{
		Use:  "remote-services",
		RunE: actionGetRemoteServices,
	}
	return []*cobra.Command{getService, getServices}
}

func actionGetRemoteServices(cmd *cobra.Command, args []string) error {
	// Get the project and url parameters
	project, check := utils.GetProjectID()
	if !check {
		return utils.LogError("Project not specified in flag", nil)
	}
	commandName := "remote-service"

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
	if len(args) != 1 {
		return utils.LogError("incorrect number of arguments", nil)
	}
	dbruleConfigFile := args[0]
	dbrule, err := generateService()
	if err != nil {
		return err
	}

	return utils.AppendConfigToDisk(dbrule, dbruleConfigFile)
}
