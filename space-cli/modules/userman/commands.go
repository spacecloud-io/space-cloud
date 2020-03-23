package userman

import (
	"fmt"

	"github.com/spaceuptech/space-cli/utils"
	"github.com/urfave/cli"
)

// Commands is the list of commands the userman module exposes
var GenerateSubCommands = []cli.Command{
	{
		Name:   "auth-providers",
		Action: actionGenerateUserManagement,
	},
}

func actionGenerateUserManagement(c *cli.Context) error {
	argsArr := c.Args()
	if len(argsArr) != 1 {
		return fmt.Errorf("incorrect number of arguments")
	}
	dbruleConfigFile := argsArr[0]
	dbrule, err := generateUserManagement()
	if err != nil {
		return err
	}

	return utils.AppendConfigToDisk(dbrule, dbruleConfigFile)
}
