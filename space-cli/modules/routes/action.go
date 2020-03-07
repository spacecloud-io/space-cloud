package routes

import (
	"github.com/urfave/cli"

	"github.com/spaceuptech/space-cli/utils"
)

//ActionGetRoutes gets routes
func ActionGetRoutes(c *cli.Context) error {
	// Get the project and url parameters
	project := c.GlobalString("project")
	commandName := c.Command.Name

	params := map[string]string{}
	if len(c.Args()) != 0 {
		params["routesId"] = c.Args()[0]
	}

	objs, err := getRoutes(project, commandName, params)
	if err != nil {
		return err
	}
	if err := utils.PrintYaml(objs); err != nil {
		return err
	}
	return nil
}
