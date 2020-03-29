package userman

import (
	"github.com/spaceuptech/space-cli/utils"
	"github.com/spf13/cobra"
)

// GenerateSubCommands dis the list of commands the project module exposes
func GenerateSubCommands() []*cobra.Command {

	var generateUserManagement = &cobra.Command{
		Use:  "auth-providers",
		RunE: actionGenerateUserManagement,
	}

	return []*cobra.Command{generateUserManagement}
}

func actionGenerateUserManagement(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		_ = utils.LogError("incorrect number of arguments", nil)
		return nil
	}
	dbruleConfigFile := args[0]
	dbrule, err := generateUserManagement()
	if err != nil {
		return nil
	}

	_ = utils.AppendConfigToDisk(dbrule, dbruleConfigFile)
	return nil
}
