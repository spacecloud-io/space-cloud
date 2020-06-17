package filestore

import (
	"github.com/spf13/cobra"

	"github.com/spaceuptech/space-cli/cmd/utils"
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
		Use:  "filestore-rule",
		RunE: actionGetFileStoreRule,
	}

	var getFileStoreConfig = &cobra.Command{
		Use:  "filestore-config",
		RunE: actionGetFileStoreConfig,
	}

	var getFileStoreRules = &cobra.Command{
		Use:  "filestore-rules",
		RunE: actionGetFileStoreRule,
	}

	var getFileStoreConfigs = &cobra.Command{
		Use:  "filestore-configs",
		RunE: actionGetFileStoreConfig,
	}

	return []*cobra.Command{getFileStoreRule, getFileStoreConfig, getFileStoreRules, getFileStoreConfigs}
}

func actionGetFileStoreConfig(cmd *cobra.Command, args []string) error {
	// Get the project and url parameters
	project, check := utils.GetProjectID()
	if !check {
		return utils.LogError("Project not specified in flag", nil)
	}
	commandName := "filestore-config"

	params := map[string]string{}
	obj, err := GetFileStoreConfig(project, commandName, params)
	if err != nil {
		return err
	}

	if err := utils.PrintYaml(obj); err != nil {
		return err
	}
	return nil
}

func actionGetFileStoreRule(cmd *cobra.Command, args []string) error {
	// Get the project and url parameters
	project, check := utils.GetProjectID()
	if !check {
		return utils.LogError("Project not specified in flag", nil)
	}
	commandName := "filestore-rule"

	params := map[string]string{}
	if len(args) != 0 {
		params["id"] = args[0]
	}

	objs, err := GetFileStoreRule(project, commandName, params)
	if err != nil {
		return err
	}

	if err := utils.PrintYaml(objs); err != nil {
		return err
	}
	return nil
}

func actionGenerateFilestoreRule(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return utils.LogError("incorrect number of arguments", nil)
	}
	dbruleConfigFile := args[0]
	dbrule, err := generateFilestoreRule()
	if err != nil {
		return err
	}

	return utils.AppendConfigToDisk(dbrule, dbruleConfigFile)
}

func actionGenerateFilestoreConfig(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return utils.LogError("incorrect number of arguments", nil)
	}
	dbruleConfigFile := args[0]
	dbrule, err := generateFilestoreConfig()
	if err != nil {
		return err
	}

	return utils.AppendConfigToDisk(dbrule, dbruleConfigFile)
}
