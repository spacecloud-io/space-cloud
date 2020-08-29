package secrets

import (
	"github.com/spaceuptech/helpers"
	"github.com/urfave/cli"

	"github.com/spaceuptech/space-cloud/runner/utils"
)

// ActionGenerateSecret creates spec object for service routing
func ActionGenerateSecret(c *cli.Context) error {
	argsArr := c.Args()
	if len(argsArr) != 1 {
		return helpers.Logger.LogError(helpers.GetInternalRequestID(), "Incorrect number of arguments", nil, nil)
	}
	dbruleConfigFile := argsArr[0]
	dbrule, err := generateSecrets()
	if err != nil {
		return err
	}

	return utils.AppendConfigToDisk(dbrule, dbruleConfigFile)
}
