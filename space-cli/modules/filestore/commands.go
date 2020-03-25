package filestore

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/spaceuptech/space-cli/utils"
)

// Commands is the list of commands the filestore module exposes
func Commands() []*cobra.Command {
	var generateSubCommands = &cobra.Command{}

	var generaterule = &cobra.Command{
		Use:  "filestore-rules",
		RunE: actionGenerateFilestoreRule,
	}

	var generateconfig = &cobra.Command{
		Use:  "filestore-config",
		RunE: actionGenerateFilestoreConfig,
	}

	var getSubCommands = &cobra.Command{}

	var getFileStoreRule = &cobra.Command{
		Use:  "filestore-rules",
		RunE: actionGetFileStoreRule,
	}

	var getFileStoreConfig = &cobra.Command{
		Use:  "filestore-config",
		RunE: actionGetFileStoreConfig,
	}
	generateSubCommands.AddCommand(generaterule)
	generateSubCommands.AddCommand(generateconfig)
	getSubCommands.AddCommand(getFileStoreRule)
	getSubCommands.AddCommand(getFileStoreConfig)

	command := make([]*cobra.Command, 0)
	command = append(command, generateSubCommands)
	command = append(command, getSubCommands)
	return command
}

// // GetSubCommands is the list of commands the filestore module exposes
// var GetSubCommands = []cli.Command{
// 	{
// 		Name:   "filestore-config",
// 		Action: actionGetFileStoreConfig,
// 	},
// 	{
// 		Name:   "filestore-rules",
// 		Action: actionGetFileStoreRule,
// 	},
// }

// // GenerateSubCommands is the list of commands the filestore module exposes
// var GenerateSubCommands = []cli.Command{
// 	{
// 		Name:   "filestore-rules",
// 		Action: actionGenerateFilestoreRule,
// 	},
// 	{
// 		Name:   "filestore-config",
// 		Action: actionGenerateFilestoreConfig,
// 	},
// }

func actionGetFileStoreConfig(cmd *cobra.Command, args []string) error {
	// Get the project and url parameters
	project := viper.GetString("project")
	commandName := cmd.CalledAs()

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
	project := viper.GetString("project")
	commandName := cmd.CalledAs()

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
	argsArr := args
	if len(argsArr) != 1 {
		return fmt.Errorf("incorrect number of arguments")
	}
	dbruleConfigFile := argsArr[0]
	dbrule, err := generateFilestoreRule()
	if err != nil {
		return err
	}

	return utils.AppendConfigToDisk(dbrule, dbruleConfigFile)
}

func actionGenerateFilestoreConfig(cmd *cobra.Command, args []string) error {
	argsArr := args
	if len(argsArr) != 1 {
		return fmt.Errorf("incorrect number of arguments")
	}
	dbruleConfigFile := argsArr[0]
	dbrule, err := generateFilestoreConfig()
	if err != nil {
		return err
	}

	return utils.AppendConfigToDisk(dbrule, dbruleConfigFile)
}
