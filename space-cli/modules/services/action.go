package services

import (
	"github.com/urfave/cli"

	"github.com/spaceuptech/space-cli/utils"
)

//ActionGetServicesRoutes gets services routes
func ActionGetServicesRoutes(c *cli.Context) error {
	// Get the project and url parameters
	project := c.GlobalString("project")
	commandName := c.Command.Name

	params := map[string]string{}
	if len(c.Args()) != 0 {
		params["serviceId"] = c.Args()[0]
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

//ActionGetServicesSecrets gets services routes
func ActionGetServicesSecrets(c *cli.Context) error {
	// Get the project and url parameters
	project := c.GlobalString("project")
	commandName := c.Command.Name

	params := map[string]string{}
	if len(c.Args()) != 0 {
		params["name"] = c.Args()[0]
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

//ActionGetServices gets runner services
func ActionGetServices(c *cli.Context) error {
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
