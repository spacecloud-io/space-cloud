package auth

import (
	"github.com/spf13/cobra"

	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils"
)

// GetSubCommands is the list of commands the auth module exposes
func GetSubCommands() []*cobra.Command {
	var getAuthProviders = &cobra.Command{
		Use:     "auth-providers",
		Aliases: []string{"auth-provider"},
		RunE:    actionGetAuthProviders,
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			project, check := utils.GetProjectID()
			if !check {
				utils.LogDebug("Project not specified in flag", nil)
				return nil, cobra.ShellCompDirectiveDefault
			}
			objs, err := GetAuthProviders(project, "auth-providers", map[string]string{})
			if err != nil {
				return nil, cobra.ShellCompDirectiveDefault
			}
			var ids []string
			for _, v := range objs {
				ids = append(ids, v.Meta["id"])
			}
			return ids, cobra.ShellCompDirectiveDefault
		},
	}
	return []*cobra.Command{getAuthProviders}
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
		Use:     "auth-provider [path to config file]",
		RunE:    actionGenerateUserManagement,
		Aliases: []string{"auth-providers"},
		Example: "space-cli generate auth-provider config.yaml --project myproject --log-level info",
	}

	return []*cobra.Command{generateUserManagement}
}

func actionGenerateUserManagement(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return utils.LogError("incorrect number of arguments. Use -h to check usage instructions", nil)
	}
	dbruleConfigFile := args[0]
	dbrule, err := generateUserManagement()
	if err != nil {
		return err
	}

	return utils.AppendConfigToDisk(dbrule, dbruleConfigFile)
}

// DeleteSubCommands is the list of commands the auth module exposes
func DeleteSubCommands() []*cobra.Command {

	var deleteAuthProvider = &cobra.Command{
		Use:     "auth-provider",
		Aliases: []string{"auth-providers"},
		RunE:    actionDeleteAuthProvider,
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			project, check := utils.GetProjectID()
			if !check {
				utils.LogDebug("Project not specified in flag", nil)
				return nil, cobra.ShellCompDirectiveDefault
			}
			objs, err := GetAuthProviders(project, "auth-providers", map[string]string{})
			if err != nil {
				return nil, cobra.ShellCompDirectiveDefault
			}
			var ids []string
			for _, v := range objs {
				ids = append(ids, v.Meta["id"])
			}
			return ids, cobra.ShellCompDirectiveDefault
		},
	}

	return []*cobra.Command{deleteAuthProvider}
}

func actionDeleteAuthProvider(cmd *cobra.Command, args []string) error {
	// Get the project and url parameters
	project, check := utils.GetProjectID()
	if !check {
		return utils.LogError("Project not specified in flag", nil)
	}

	prefix := ""
	if len(args) != 0 {
		prefix = args[0]
	}

	return deleteAuthProvider(project, prefix)
}
