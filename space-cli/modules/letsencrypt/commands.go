package letsencrypt

import (
	"fmt"

	"github.com/urfave/cli"

	"github.com/spaceuptech/space-cli/utils"
)

// GenerateSubCommands is the list of commands the letsencrypt module exposes
var GenerateSubCommands = []cli.Command{
	{
		Name:   "letsencrypt",
		Action: actionGenerateLetsEncryptDomain,
	},
}

// GetSubCommands is the list of commands the letsencrypt module exposes
var GetSubCommands = []cli.Command{{
	Name:   "letsencrypt",
	Action: actionGetLetsEncrypt,
}}

func actionGetLetsEncrypt(c *cli.Context) error {
	// Get the project and url parameters
	project := c.GlobalString("project")
	commandName := c.Command.Name

	params := map[string]string{}
	obj, err := GetLetsEncryptDomain(project, commandName, params)
	if err != nil {
		return err
	}
	if err := utils.PrintYaml(obj); err != nil {
		return err
	}
	return nil
}

func actionGenerateLetsEncryptDomain(c *cli.Context) error {
	argsArr := c.Args()
	if len(argsArr) != 1 {
		return fmt.Errorf("incorrect number of arguments")
	}
	dbruleConfigFile := argsArr[0]
	dbrule, err := generateLetsEncryptDomain()
	if err != nil {
		return err
	}

	return utils.AppendConfigToDisk(dbrule, dbruleConfigFile)
}
