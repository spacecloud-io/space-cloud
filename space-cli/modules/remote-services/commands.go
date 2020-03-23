package remoteservices

import (
	"fmt"

	"github.com/urfave/cli"

	"github.com/spaceuptech/space-cli/utils"
)

// Commands is the list of commands the remoteservices module exposes
var Commands = []cli.Command{
	{
		Name:  "generate",
		Usage: "generates service config",
		Subcommands: []cli.Command{
			{
				Name:   "remote-services",
				Action: actionGenerateService,
			},
		},
	},
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
				Name:   "remote-services",
				Action: actionGetRemoteServices,
			},
		},
	},
}

func actionGetRemoteServices(c *cli.Context) error {
	// Get the project and url parameters
	project := c.GlobalString("project")
	commandName := c.Command.Name

	params := map[string]string{}
	if len(c.Args()) != 0 {
		params["id"] = c.Args()[0]
	}

	objs, err := GetRemoteServices(project, commandName, params)
	if err != nil {
		return err
	}
	if err := utils.PrintYaml(objs); err != nil {
		return err
	}
	return nil
}

func actionGenerateService(c *cli.Context) error {
	argsArr := c.Args()
	if len(argsArr) != 1 {
		return fmt.Errorf("incorrect number of arguments")
	}
	dbruleConfigFile := argsArr[0]
	dbrule, err := generateService()
	if err != nil {
		return err
	}

	return utils.AppendConfigToDisk(dbrule, dbruleConfigFile)
}
