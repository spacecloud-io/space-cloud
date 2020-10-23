package filestore

import (
	"github.com/spf13/cobra"

	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils"
)

// GenerateSubCommands is the list of commands the filestore module exposes
func GenerateSubCommands() []*cobra.Command {

	var generaterule = &cobra.Command{
		Use:     "filestore-rule [path to config file]",
		RunE:    actionGenerateFilestoreRule,
		Aliases: []string{"filestore-rules"},
		Example: "space-cli generate filestore-rule config.yaml --project myproject --log-level info",
	}

	var generateconfig = &cobra.Command{
		Use:     "filestore-config [path to config file]",
		RunE:    actionGenerateFilestoreConfig,
		Aliases: []string{"filestore-configs"},
		Example: "space-cli generate filestore-config config.yaml --project myproject --log-level info",
	}

	return []*cobra.Command{generaterule, generateconfig}
}

// GetSubCommands is the list of commands the filestore module exposes
func GetSubCommands() []*cobra.Command {

	var getFileStoreRules = &cobra.Command{
		Use:               "filestore-rules",
		Aliases:           []string{"filestore-rule"},
		RunE:              actionGetFileStoreRule,
		ValidArgsFunction: fileStoreRulesAutoCompleteFun,
	}

	var getFileStoreConfigs = &cobra.Command{
		Use:     "filestore-configs",
		Aliases: []string{"filestore-config"},
		RunE:    actionGetFileStoreConfig,
	}

	return []*cobra.Command{getFileStoreRules, getFileStoreConfigs}
}

// DeleteSubCommands is the list of commands the filestore module exposes
func DeleteSubCommands() []*cobra.Command {

	var deleteFileStoreRules = &cobra.Command{
		Use:               "filestore-rules",
		Aliases:           []string{"filestore-rule"},
		RunE:              actionDeleteFileStoreRule,
		ValidArgsFunction: fileStoreRulesAutoCompleteFun,
		Example:           "space-cli delete filestore-rules ruleID --project myproject --log-level info",
	}

	var deleteFileStoreConfigs = &cobra.Command{
		Use:     "filestore-configs",
		Aliases: []string{"filestore-config"},
		RunE:    actionDeleteFileStoreConfig,
		Example: "space-cli delete filestore-configs --project myproject --log-level info",
	}

	return []*cobra.Command{deleteFileStoreRules, deleteFileStoreConfigs}
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
		return utils.LogError("incorrect number of arguments. Use -h to check usage instructions", nil)
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
		return utils.LogError("incorrect number of arguments. Use -h to check usage instructions", nil)
	}
	dbruleConfigFile := args[0]
	dbrule, err := generateFilestoreConfig()
	if err != nil {
		return err
	}

	return utils.AppendConfigToDisk(dbrule, dbruleConfigFile)
}

func actionDeleteFileStoreConfig(cmd *cobra.Command, args []string) error {
	// Get the project
	project, check := utils.GetProjectID()
	if !check {
		return utils.LogError("Project not specified in flag", nil)
	}

	return deleteFileStoreConfig(project)
}

func actionDeleteFileStoreRule(cmd *cobra.Command, args []string) error {
	// Get the project and url parameters
	project, check := utils.GetProjectID()
	if !check {
		return utils.LogError("Project not specified in flag", nil)
	}

	prefix := ""
	if len(args) != 0 {
		prefix = args[0]
	}

	return deleteFileStoreRule(project, prefix)
}
