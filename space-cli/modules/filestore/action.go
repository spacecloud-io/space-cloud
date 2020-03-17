package filestore

import (
	"github.com/urfave/cli"

	"github.com/spaceuptech/space-cli/utils"
)

// ActionGetFileStoreConfig gets filestore config
func ActionGetFileStoreConfig(c *cli.Context) error {
	// Get the project and url parameters
	project := c.GlobalString("project")
	commandName := c.Command.Name

	params := map[string]string{}
	objs, err := GetFileStoreConfig(project, commandName, params)
	if err != nil {
		return err
	}
	if err := utils.PrintYaml(objs); err != nil {
		return err
	}
	return nil
}

// ActionGetFileStoreRule gets filestore rule
func ActionGetFileStoreRule(c *cli.Context) error {
	// Get the project and url parameters
	project := c.GlobalString("project")
	commandName := c.Command.Name

	params := map[string]string{}
	if len(c.Args()) != 0 {
		params["id"] = c.Args()[0]
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
