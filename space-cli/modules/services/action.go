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
