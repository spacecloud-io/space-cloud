package eventing

import (
	"fmt"

	"github.com/spaceuptech/space-cli/utils"

	"github.com/urfave/cli"
)

// ActionGenerateEventingRule creates spec object for eventing rule
func ActionGenerateEventingRule(c *cli.Context) error {
	argsArr := c.Args()
	if len(argsArr) != 1 {
		return fmt.Errorf("incorrect number of arguments")
	}
	dbruleConfigFile := argsArr[0]
	dbrule, err := generateEventingRule()
	if err != nil {
		return err
	}

	return utils.AppendConfigToDisk(dbrule, dbruleConfigFile)
}

// ActionGenerateEventingSchema creates spec object for eventing Schema
func ActionGenerateEventingSchema(c *cli.Context) error {
	argsArr := c.Args()
	if len(argsArr) != 1 {
		return fmt.Errorf("incorrect number of arguments")
	}
	dbruleConfigFile := argsArr[0]
	dbrule, err := generateEventingSchema()
	if err != nil {
		return err
	}

	return utils.AppendConfigToDisk(dbrule, dbruleConfigFile)
}

// ActionGenerateEventingConfig creates spec object for eventing Config
func ActionGenerateEventingConfig(c *cli.Context) error {
	argsArr := c.Args()
	if len(argsArr) != 1 {
		return fmt.Errorf("incorrect number of arguments")
	}
	dbruleConfigFile := argsArr[0]
	dbrule, err := generateEventingConfig()
	if err != nil {
		return err
	}

	return utils.AppendConfigToDisk(dbrule, dbruleConfigFile)
}

// ActionGenerateEventingTrigger creates spec object for eventing Trigger
func ActionGenerateEventingTrigger(c *cli.Context) error {
	argsArr := c.Args()
	if len(argsArr) != 1 {
		return fmt.Errorf("incorrect number of arguments")
	}
	dbruleConfigFile := argsArr[0]
	dbrule, err := generateEventingTrigger()
	if err != nil {
		return err
	}

	return utils.AppendConfigToDisk(dbrule, dbruleConfigFile)
}
