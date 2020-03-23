package services

import (
	"fmt"
	"os"

	"github.com/urfave/cli"

	"github.com/spaceuptech/space-cli/utils"
)

// Commands is the list of commands the services module exposes
var GenerateSubCommands = []cli.Command{
	{
		Name:   "service",
		Action: actionGenerateService,
	},
}

var GetSubCommands = []cli.Command{
	{
		Name:   "services-routes",
		Action: actionGetServicesRoutes,
	},
	{
		Name:   "services-secrets",
		Action: actionGetServicesSecrets,
	},
	{
		Name:   "services",
		Action: actionGetServices,
	},
}

func actionGetServicesRoutes(c *cli.Context) error {
	// Get the project and url parameters
	project := c.GlobalString("project")
	commandName := c.Command.Name

	params := map[string]string{}
	if len(c.Args()) != 0 {
		params["id"] = c.Args()[0]
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

func actionGetServicesSecrets(c *cli.Context) error {
	// Get the project and url parameters
	project := c.GlobalString("project")
	commandName := c.Command.Name

	params := map[string]string{}
	if len(c.Args()) != 0 {
		params["id"] = c.Args()[0]
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

func actionGetServices(c *cli.Context) error {
	// Get the project and url parameters
	project := c.GlobalString("project")
	commandName := c.Command.Name

	params := map[string]string{}
	switch len(c.Args()) {
	case 1:
		params["serviceId"] = c.Args()[0]
	case 2:
		params["serviceId"] = c.Args()[0]
		params["version"] = c.Args()[1]
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

func actionGenerateService(_ *cli.Context) error {
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
