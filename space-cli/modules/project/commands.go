package project

import (
	"github.com/urfave/cli"

	"github.com/spaceuptech/space-cli/utils"
)

// GetSubCommands is the list of commands the project module exposes
var GetSubCommands = []cli.Command{
	{
		Name:   "project",
		Action: actionGetProjectConfig,
	},
}

func actionGetProjectConfig(c *cli.Context) error {
	// Get the project and cmd parameters
	project := c.GlobalString("project")
	commandName := c.Command.Name

	params := map[string]string{}
	obj, err := GetProjectConfig(project, commandName, params)
	if err != nil {
		return err
	}
	if err := utils.PrintYaml(obj); err != nil {
		return err
	}
	return nil
}
