package letsencrypt

import (
	"fmt"

	"github.com/spaceuptech/space-cli/utils"
	"github.com/urfave/cli"
)

// ActionGenerateLetsEncryptDomain generates lets encrypt domain spec object
func ActionGenerateLetsEncryptDomain(c *cli.Context) error {
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
