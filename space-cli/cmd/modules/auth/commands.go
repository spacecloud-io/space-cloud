package auth

import (
	"github.com/spf13/cobra"

	"github.com/spaceuptech/space-cli/cmd/utils"
)

// GetSubCommands is the list of commands the auth module exposes
func GetSubCommands() []*cobra.Command {
	var getAuthProvider = &cobra.Command{
		Use:  "auth-provider",
		RunE: actionGetAuthProviders,
	}
	var getAuthProviders = &cobra.Command{
		Use:  "auth-providers",
		RunE: actionGetAuthProviders,
	}
	return []*cobra.Command{getAuthProvider, getAuthProviders}
}

func actionGetAuthProviders(cmd *cobra.Command, args []string) error {
	// Get the project and url parameters
	project, check := utils.GetProjectID()
	if !check {
		return utils.LogError("Project not specified in flag", nil)
	}
	commandName := "auth-provider"

	params := map[string]string{}
	if len(args) != 0 {
		params["id"] = args[0]
	}

	objs, err := GetAuthProviders(project, commandName, params)
	if err != nil {
		return err
	}

	if err := utils.PrintYaml(objs); err != nil {
		return err
	}
	return nil
}

// GenerateSubCommands dis the list of commands the project module exposes
func GenerateSubCommands() []*cobra.Command {

	var generateUserManagement = &cobra.Command{
		Use:  "auth-provider",
		RunE: actionGenerateUserManagement,
	}

	return []*cobra.Command{generateUserManagement}
}

func actionGenerateUserManagement(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return utils.LogError("incorrect number of arguments", nil)
	}
	dbruleConfigFile := args[0]
	dbrule, err := generateUserManagement()
	if err != nil {
		return err
	}

	return utils.AppendConfigToDisk(dbrule, dbruleConfigFile)
}
