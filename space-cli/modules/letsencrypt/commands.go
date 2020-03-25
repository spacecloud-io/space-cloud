package letsencrypt

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/spaceuptech/space-cli/utils"
)

// Commands is the list of commands the ingress module exposes
func Commands() []*cobra.Command {
	var generateSubCommands = &cobra.Command{}

	var generateletsencrypt = &cobra.Command{
		Use:  "letsencrypt",
		RunE: actionGenerateLetsEncryptDomain,
	}

	var getSubCommands = &cobra.Command{}

	var getletsencrypt = &cobra.Command{
		Use:  "letsencrypt",
		RunE: actionGetLetsEncrypt,
	}

	generateSubCommands.AddCommand(generateletsencrypt)
	getSubCommands.AddCommand(getletsencrypt)

	command := make([]*cobra.Command, 0)
	command = append(command, generateSubCommands)
	command = append(command, getSubCommands)
	return command
}

// GenerateSubCommands is the list of commands the letsencrypt module exposes
// var GenerateSubCommands = []cli.Command{
// 	{
// 		Name:   "letsencrypt",
// 		Action: actionGenerateLetsEncryptDomain,
// 	},
// }

// // GetSubCommands is the list of commands the letsencrypt module exposes
// var GetSubCommands = []cli.Command{{
// 	Name:   "letsencrypt",
// 	Action: actionGetLetsEncrypt,
// }}

func actionGetLetsEncrypt(cmd *cobra.Command, args []string) error {
	// Get the project and url parameters
	project := viper.GetString("project")
	commandName := cmd.CalledAs()

	params := map[string]string{}
	obj, err := GetLetsEncryptDomain(project, commandName, params)
	if err != nil {
		return err
	}
	if err := utils.PrintYaml(obj); err != nil {
		return err
	}
	return nil
}

func actionGenerateLetsEncryptDomain(cmd *cobra.Command, args []string) error {
	argsArr := args
	if len(argsArr) != 1 {
		return fmt.Errorf("incorrect number of arguments")
	}
	dbruleConfigFile := argsArr[0]
	dbrule, err := generateLetsEncryptDomain()
	if err != nil {
		return err
	}

	return utils.AppendConfigToDisk(dbrule, dbruleConfigFile)
}
