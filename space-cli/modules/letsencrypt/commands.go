package letsencrypt

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/spaceuptech/space-cli/utils"
)

// GenerateSubCommands is the list of commands the ingress module exposes
func GenerateSubCommands() []*cobra.Command {

	var generateletsencrypt = &cobra.Command{
		Use:  "letsencrypt",
		RunE: actionGenerateLetsEncryptDomain,
	}
	return []*cobra.Command{generateletsencrypt}
}

// GetSubCommands is the list of commands the ingress module exposes
func GetSubCommands() []*cobra.Command {

	var getletsencrypt = &cobra.Command{
		Use:  "letsencrypt",
		RunE: actionGetLetsEncrypt,
	}

	return []*cobra.Command{getletsencrypt}
}

func actionGetLetsEncrypt(cmd *cobra.Command, args []string) error {
	// Get the project and url parameters
	project := viper.GetString("project")
	commandName := cmd.Use

	params := map[string]string{}
	obj, _ := GetLetsEncryptDomain(project, commandName, params)
	_ = utils.PrintYaml(obj)
	return nil
}

func actionGenerateLetsEncryptDomain(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("incorrect number of arguments")
	}
	dbruleConfigFile := args[0]
	dbrule, _ := generateLetsEncryptDomain()

	_ = utils.AppendConfigToDisk(dbrule, dbruleConfigFile)
	return nil
}
