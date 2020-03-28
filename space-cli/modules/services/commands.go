package services

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/spaceuptech/space-cli/utils"
)

// GenerateSubCommands is the list of commands the services module exposes
func GenerateSubCommands() []*cobra.Command {

	var generateService = &cobra.Command{
		Use:  "services",
		RunE: actionGenerateService,
	}

	return []*cobra.Command{generateService}

}

// GetSubCommands is the list of commands the services module exposes
func GetSubCommands() []*cobra.Command {

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

	return []*cobra.Command{getServicesRoutes, getServicesSecrets, getService}
}

func actionGetServicesRoutes(cmd *cobra.Command, args []string) error {
	// Get the project and url parameters
	project := viper.GetString("project")
	commandName := cmd.Use

	params := map[string]string{}
	if len(args) != 0 {
		params["id"] = args[0]
	}

	objs, _ := GetServicesRoutes(project, commandName, params)
	_ = utils.PrintYaml(objs)
	return nil
}

func actionGetServicesSecrets(cmd *cobra.Command, args []string) error {
	// Get the project and url parameters
	project := viper.GetString("project")
	commandName := cmd.Use

	params := map[string]string{}
	if len(args) != 0 {
		params["id"] = args[0]
	}

	objs, _ := GetServicesSecrets(project, commandName, params)
	_ = utils.PrintYaml(objs)
	return nil
}

func actionGetServices(cmd *cobra.Command, args []string) error {
	// Get the project and url parameters
	project := viper.GetString("project")
	commandName := cmd.Use

	params := map[string]string{}
	switch len(args) {
	case 1:
		params["serviceId"] = args[0]
	case 2:
		params["serviceId"] = args[0]
		params["version"] = args[1]
	}
	objs, _ := GetServices(project, commandName, params)
	_ = utils.PrintYaml(objs)
	return nil
}

func actionGenerateService(cmd *cobra.Command, args []string) error {
	// get filename from args in which service config will be stored
	if len(args) != 4 {
		return fmt.Errorf("incorrect number of arguments")
	}
	serviceConfigFile := args[3]

	service, _ := GenerateService("", "")

	_ = utils.AppendConfigToDisk(service, serviceConfigFile)
	return nil
}
