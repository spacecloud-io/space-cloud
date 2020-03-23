package project

import (
	"github.com/urfave/cli"

	"github.com/spaceuptech/space-cli/model"
	"github.com/spaceuptech/space-cli/utils"
)

// Commands is the list of commands the project module exposes
var Commands = []cli.Command{
	{
		Name:  "get",
		Usage: "gets different services",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "project",
				Usage:  "The id of the project",
				EnvVar: "PROJECT_ID",
			},
		},
		Subcommands: []cli.Command{
			{
				Name:   "project",
				Action: actionGetProjectConfig,
			},
		},
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
	if err := utils.PrintYaml([]*model.SpecObject{obj}); err != nil {
		return err
	}
	return nil
}
