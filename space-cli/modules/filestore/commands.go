package filestore

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/spaceuptech/space-cli/utils"
)

// GenerateSubCommands is the list of commands the filestore module exposes
func GenerateSubCommands() []*cobra.Command {

	var generaterule = &cobra.Command{
		Use:  "filestore-rules",
		RunE: actionGenerateFilestoreRule,
	}

	var generateconfig = &cobra.Command{
		Use:  "filestore-config",
		RunE: actionGenerateFilestoreConfig,
	}

	return []*cobra.Command{generaterule, generateconfig}
}

// GetSubCommands is the list of commands the filestore module exposes
func GetSubCommands() []*cobra.Command {

	var getFileStoreRule = &cobra.Command{
		Use:  "filestore-rules",
		RunE: actionGetFileStoreRule,
	}

	var getFileStoreConfig = &cobra.Command{
		Use:  "filestore-config",
		RunE: actionGetFileStoreConfig,
	}

	return []*cobra.Command{getFileStoreRule, getFileStoreConfig}
}

func actionGetFileStoreConfig(cmd *cobra.Command, args []string) error {
	// Get the project and url parameters
	project := viper.GetString("project")
	commandName := cmd.Use

	params := map[string]string{}
	obj, err := GetFileStoreConfig(project, commandName, params)
	if err != nil {
		return nil
	}

	if err := utils.PrintYaml(obj); err != nil {
		return nil
	}
	return nil
}

func actionGetFileStoreRule(cmd *cobra.Command, args []string) error {
	// Get the project and url parameters
	project := viper.GetString("project")
	commandName := cmd.Use

	params := map[string]string{}
	if len(args) != 0 {
		params["id"] = args[0]
	}

	objs, err := GetFileStoreRule(project, commandName, params)
	if err != nil {
		return nil
	}

	if err := utils.PrintYaml(objs); err != nil {
		return nil
	}
	return nil
}

func actionGenerateFilestoreRule(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("incorrect number of arguments")
	}
	dbruleConfigFile := args[0]
	dbrule, err := generateFilestoreRule()
	if err != nil {
		return nil
	}

	_ = utils.AppendConfigToDisk(dbrule, dbruleConfigFile)
	return nil
}

func actionGenerateFilestoreConfig(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("incorrect number of arguments")
	}
	dbruleConfigFile := args[0]
	dbrule, err := generateFilestoreConfig()
	if err != nil {
		return nil
	}

	_ = utils.AppendConfigToDisk(dbrule, dbruleConfigFile)
	return nil
}
