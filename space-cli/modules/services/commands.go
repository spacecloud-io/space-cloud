package services

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/spaceuptech/space-cli/utils"
)

// Commands is the list of commands the services module exposes
func Commands() []*cobra.Command {
	var generateSubCommands = &cobra.Command{}

	var generateService = &cobra.Command{
		Use:  "services",
		RunE: actionGenerateService,
	}

	var getSubCommands = &cobra.Command{}

	var getServicesRoutes = &cobra.Command{
		Use:  "services-routes",
		RunE: actionGetServicesRoutes,
	}

	var getServicesSecrets = &cobra.Command{
		Use:  "services-secrets",
		RunE: actionGetServicesSecrets,
	}

	var getService = &cobra.Command{
		Use:  "services",
		RunE: actionGetServices,
	}

	generateSubCommands.AddCommand(generateService)
	getSubCommands.AddCommand(getServicesRoutes)
	getSubCommands.AddCommand(getServicesSecrets)
	getSubCommands.AddCommand(getService)

	command := make([]*cobra.Command, 0)
	command = append(command, generateSubCommands)
	command = append(command, getSubCommands)
	return command
}

// // GenerateSubCommands is the list of commands the services module exposes
// var GenerateSubCommands = []cli.Command{
// 	{
// 		Name:   "service",
// 		Action: actionGenerateService,
// 	},
// }

// // GetSubCommands is the list of commands the services module exposes
// var GetSubCommands = []cli.Command{
// 	{
// 		Name:   "services-routes",
// 		Action: actionGetServicesRoutes,
// 	},
// 	{
// 		Name:   "services-secrets",
// 		Action: actionGetServicesSecrets,
// 	},
// 	{
// 		Name:   "services",
// 		Action: actionGetServices,
// 	},
// }

func actionGetServicesRoutes(cmd *cobra.Command, args []string) error {
	// Get the project and url parameters
	project := viper.GetString("project")
	commandName := cmd.CalledAs()

	params := map[string]string{}
	if len(args) != 0 {
		params["id"] = args[0]
	}

	objs, err := GetServicesRoutes(project, commandName, params)
	if err != nil {
		return err
	}
	if err := utils.PrintYaml(objs); err != nil {
		return err
	}
	return nil
}

func actionGetServicesSecrets(cmd *cobra.Command, args []string) error {
	// Get the project and url parameters
	project := viper.GetString("project")
	commandName := cmd.CalledAs()

	params := map[string]string{}
	if len(args) != 0 {
		params["id"] = args[0]
	}

	objs, err := GetServicesSecrets(project, commandName, params)
	if err != nil {
		return err
	}
	if err := utils.PrintYaml(objs); err != nil {
		return err
	}
	return nil
}

func actionGetServices(cmd *cobra.Command, args []string) error {
	// Get the project and url parameters
	project := viper.GetString("project")
	commandName := cmd.CalledAs()

	params := map[string]string{}
	switch len(args) {
	case 1:
		params["serviceId"] = args[0]
	case 2:
		params["serviceId"] = args[0]
		params["version"] = args[1]
	}
	objs, err := GetServices(project, commandName, params)
	if err != nil {
		return err
	}
	if err := utils.PrintYaml(objs); err != nil {
		return err
	}
	return nil
}

func actionGenerateService(cmd *cobra.Command, args []string) error {
	// get filename from args in which service config will be stored
	argsArr := os.Args
	if len(argsArr) != 4 {
		return fmt.Errorf("incorrect number of arguments")
	}
	serviceConfigFile := argsArr[3]

	service, err := GenerateService("", "")
	if err != nil {
		return err
	}

	return utils.AppendConfigToDisk(service, serviceConfigFile)
}
