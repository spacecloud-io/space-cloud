package ingress

import (
	"fmt"

	"github.com/spaceuptech/space-cli/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Commands is the list of commands the ingress module exposes
func Commands() []*cobra.Command {
	var generateSubCommands = &cobra.Command{}

	var generateroutes = &cobra.Command{
		Use:  "ingress-routes",
		RunE: actionGenerateIngressRouting,
	}

	var getSubCommands = &cobra.Command{}

	var getroutes = &cobra.Command{
		Use:  "ingress-routes",
		RunE: actionGetIngressRoutes,
	}

	generateSubCommands.AddCommand(generateroutes)
	getSubCommands.AddCommand(getroutes)

	command := make([]*cobra.Command, 0)
	command = append(command, generateSubCommands)
	command = append(command, getSubCommands)
	return command
}

// // GenerateSubCommands is the list of commands the ingress module exposes
// var GenerateSubCommands = []cli.Command{
// 	{
// 		Name:   "ingress-routes",
// 		Action: actionGenerateIngressRouting,
// 	},
// }

// // GetSubCommands is the list of commands the ingress module exposes
// var GetSubCommands = []cli.Command{
// 	{
// 		Name:   "ingress-routes",
// 		Action: actionGetIngressRoutes,
// 	},
// }

func actionGetIngressRoutes(cmd *cobra.Command, args []string) error {
	// Get the project and url parameters
	project := viper.GetString("project")
	commandName := cmd.CalledAs()

	params := map[string]string{}
	if len(args) != 0 {
		params["id"] = args[0]
	}

	objs, err := GetIngressRoutes(project, commandName, params)
	if err != nil {
		return err
	}
	if err := utils.PrintYaml(objs); err != nil {
		return err
	}
	return nil
}

func actionGenerateIngressRouting(cmd *cobra.Command, args []string) error {
	argsArr := args
	if len(argsArr) != 1 {
		return fmt.Errorf("incorrect number of arguments")
	}
	dbruleConfigFile := argsArr[0]
	dbrule, err := generateIngressRouting()
	if err != nil {
		return err
	}

	return utils.AppendConfigToDisk(dbrule, dbruleConfigFile)
}
