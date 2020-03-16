package services

import (
	"fmt"

	"github.com/spaceuptech/space-cli/utils"
	"github.com/urfave/cli"
)

// ActionGenerateService generates remote service spec object
func ActionGenerateService(c *cli.Context) error {
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
