package filestore

import (
	"fmt"

	"github.com/urfave/cli"

	"github.com/spaceuptech/space-cli/model"
	"github.com/spaceuptech/space-cli/utils"
)

// Commands is the list of commands the filestore module exposes
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
				Name:   "filestore-config",
				Action: actionGetFileStoreConfig,
			},
			{
				Name:   "filestore-rules",
				Action: actionGetFileStoreRule,
			},
		},
	},
	{
		Name:  "generate",
		Usage: "generates service config",
		Subcommands: []cli.Command{
			{
				Name:   "filestore-rules",
				Action: actionGenerateFilestoreRule,
			},
			{
				Name:   "filestore-config",
				Action: actionGenerateFilestoreConfig,
			},
		},
	},
}

func actionGetFileStoreConfig(c *cli.Context) error {
	// Get the project and url parameters
	project := c.GlobalString("project")
	commandName := c.Command.Name

	params := map[string]string{}
	obj, err := GetFileStoreConfig(project, commandName, params)
	if err != nil {
		return err
	}
	objs := []*model.SpecObject{obj}
	if err := utils.PrintYaml(objs); err != nil {
		return err
	}
	return nil
}

func actionGetFileStoreRule(c *cli.Context) error {
	// Get the project and url parameters
	project := c.GlobalString("project")
	commandName := c.Command.Name

	params := map[string]string{}
	if len(c.Args()) != 0 {
		params["ruleName"] = c.Args()[0]
	}

	objs, err := GetFileStoreRule(project, commandName, params)
	if err != nil {
		return err
	}
	if err := utils.PrintYaml(objs); err != nil {
		return err
	}
	return nil
}

func actionGenerateFilestoreRule(c *cli.Context) error {
	argsArr := c.Args()
	if len(argsArr) != 1 {
		return fmt.Errorf("incorrect number of arguments")
	}
	dbruleConfigFile := argsArr[0]
	dbrule, err := generateFilestoreRule()
	if err != nil {
		return err
	}

	return utils.AppendConfigToDisk(dbrule, dbruleConfigFile)
}

func actionGenerateFilestoreConfig(c *cli.Context) error {
	argsArr := c.Args()
	if len(argsArr) != 1 {
		return fmt.Errorf("incorrect number of arguments")
	}
	dbruleConfigFile := argsArr[0]
	dbrule, err := generateFilestoreConfig()
	if err != nil {
		return err
	}

	return utils.AppendConfigToDisk(dbrule, dbruleConfigFile)
}
