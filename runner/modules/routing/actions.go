package routing

import (
	"context"

	"github.com/spaceuptech/helpers"
	"github.com/urfave/cli"

	"github.com/spaceuptech/space-cloud/runner/utils"
)

// ActionGenerateServiceRouting creates spec object for service routing
func ActionGenerateServiceRouting(c *cli.Context) error {
	argsArr := c.Args()
	if len(argsArr) != 1 {
		return helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "incorrect number of arguments", nil, nil)
	}
	dbruleConfigFile := argsArr[0]
	dbrule, err := generateServiceRouting()
	if err != nil {
		return err
	}

	return utils.AppendConfigToDisk(dbrule, dbruleConfigFile)
}
