package auth

import (
	"github.com/urfave/cli"

	"github.com/spaceuptech/space-cli/utils"
)

// Commands is the list of commands the operations module exposes
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
				Name:   "auth-providers",
				Action: actionGetAuthProviders,
			},
		},
	},
}

func actionGetAuthProviders(c *cli.Context) error {
	// Get the project and url parameters
	project := c.GlobalString("project")
	commandName := c.Command.Name

	params := map[string]string{}
	if len(c.Args()) != 0 {
		params["provider"] = c.Args()[0]
	}

	objs, err := GetAuthProviders(project, commandName, params)
	if err != nil {
		return err
	}
	if err := utils.PrintYaml(objs); err != nil {
		return err
	}
	return nil
}
