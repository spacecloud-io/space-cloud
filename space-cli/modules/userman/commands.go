package userman

import (
	"fmt"

	"github.com/spaceuptech/space-cli/utils"
	"github.com/spf13/cobra"
)

// Commands dis the list of commands the project module exposes
func Commands() []*cobra.Command {

	var GenerateSubCommands = &cobra.Command{}

	var generateUserManagement = &cobra.Command{
		Use:  "auth-providers",
		RunE: actionGenerateUserManagement,
	}

	GenerateSubCommands.AddCommand(generateUserManagement)

	command := make([]*cobra.Command, 0)
	command = append(command, GenerateSubCommands)
	return command
}

// GenerateSubCommands is the list of commands the userman module exposes
// var GenerateSubCommands = []cli.Command{
// 	{
// 		Name:   "auth-providers",
// 		Action: actionGenerateUserManagement,
// 	},
// }

func actionGenerateUserManagement(cmd *cobra.Command, args []string) error {
	argsArr := args
	if len(argsArr) != 1 {
		return fmt.Errorf("incorrect number of arguments")
	}
	dbruleConfigFile := argsArr[0]
	dbrule, err := generateUserManagement()
	if err != nil {
		return err
	}

	return utils.AppendConfigToDisk(dbrule, dbruleConfigFile)
}
