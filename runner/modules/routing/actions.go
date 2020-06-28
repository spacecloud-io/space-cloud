package routing

import (
	"fmt"

	"github.com/spaceuptech/space-cloud/runner/utils"
	"github.com/urfave/cli"
)

// ActionGenerateServiceRouting creates spec object for service routing
func ActionGenerateServiceRouting(c *cli.Context) error {
	argsArr := c.Args()
	if len(argsArr) != 1 {
		return fmt.Errorf("incorrect number of arguments")
	}
	dbruleConfigFile := argsArr[0]
	dbrule, err := generateServiceRouting()
	if err != nil {
		return err
	}

	return utils.AppendConfigToDisk(dbrule, dbruleConfigFile)
}
