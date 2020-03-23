package eventing

import (
	"fmt"

	"github.com/urfave/cli"

	"github.com/spaceuptech/space-cli/utils"
)

// GetSubCommands is the list of commands the eventing module exposes
var GetSubCommands = []cli.Command{
	{
		Name:   "eventing-triggers",
		Action: actionGetEventingTrigger,
	},
	{
		Name:   "eventing-config",
		Action: actionGetEventingConfig,
	},
	{
		Name:   "eventing-schema",
		Action: actionGetEventingSchema,
	},
	{
		Name:   "eventing-rule",
		Action: actionGetEventingSecurityRule,
	},
}

// GenerateSubCommands is the list of commands the eventing module exposes
var GenerateSubCommands = []cli.Command{
	{
		Name:   "eventing-rule",
		Action: actionGenerateEventingRule,
	},
	{
		Name:   "eventing-schema",
		Action: actionGenerateEventingSchema,
	},
	{
		Name:   "eventing-config",
		Action: actionGenerateEventingConfig,
	},
	{
		Name:   "eventing-triggers",
		Action: actionGenerateEventingTrigger,
	},
}

func actionGetEventingTrigger(c *cli.Context) error {
	// Get the project and url parameters
	project := c.GlobalString("project")
	commandName := c.Command.Name

	params := map[string]string{}
	if len(c.Args()) != 0 {
		params["id"] = c.Args()[0]
	}
	objs, err := GetEventingTrigger(project, commandName, params)
	if err != nil {
		return err
	}
	if err := utils.PrintYaml(objs); err != nil {
		return err
	}
	return nil
}

func actionGetEventingConfig(c *cli.Context) error {
	// Get the project and url parameters
	project := c.GlobalString("project")
	commandName := c.Command.Name

	params := map[string]string{}
	obj, err := GetEventingConfig(project, commandName, params)
	if err != nil {
		return err
	}
	if err := utils.PrintYaml(obj); err != nil {
		return err
	}
	return nil
}

func actionGetEventingSchema(c *cli.Context) error {
	// Get the project and url parameters
	project := c.GlobalString("project")
	commandName := c.Command.Name

	params := map[string]string{}
	if len(c.Args()) != 0 {
		params["id"] = c.Args()[0]
	}
	objs, err := GetEventingSchema(project, commandName, params)
	if err != nil {
		return err
	}
	if err := utils.PrintYaml(objs); err != nil {
		return err
	}
	return nil
}

func actionGetEventingSecurityRule(c *cli.Context) error {
	// Get the project and url parameters
	project := c.GlobalString("project")
	commandName := c.Command.Name

	params := map[string]string{}
	if len(c.Args()) != 0 {
		params["id"] = c.Args()[0]
	}
	objs, err := GetEventingSecurityRule(project, commandName, params)
	if err != nil {
		return err
	}
	if err := utils.PrintYaml(objs); err != nil {
		return err
	}
	return nil
}

func actionGenerateEventingRule(c *cli.Context) error {
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

func actionGenerateEventingSchema(c *cli.Context) error {
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

func actionGenerateEventingConfig(c *cli.Context) error {
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

func actionGenerateEventingTrigger(c *cli.Context) error {
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
