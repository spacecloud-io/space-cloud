package database

import (
	"fmt"
	"github.com/spaceuptech/space-cli/utils"
	"github.com/urfave/cli"
)

// Commands is the list of commands the database module exposes
var GenerateSubCommands = []cli.Command{
	{
		Name:   "db-rules",
		Action: actionGenerateDBRule,
	},
	{
		Name:   "db-config",
		Action: actionGenerateDBConfig,
	},
	{
		Name:   "db-schema",
		Action: actionGenerateDBSchema,
	},
}

var GetSubCommands = []cli.Command{
	{
		Name:   "db-rules",
		Action: actionGetDbRules,
	},
	{
		Name:   "db-config",
		Action: actionGetDbConfig,
	},
	{
		Name:   "db-schema",
		Action: actionGetDbSchema,
	},
}

func actionGetDbRules(c *cli.Context) error {
	// Get the project and url parameters
	project := c.GlobalString("project")
	commandName := c.Command.Name

	params := map[string]string{}
	switch len(c.Args()) {
	case 1:
		params["dbAlias"] = c.Args()[0]
	case 2:
		params["dbAlias"] = c.Args()[0]
		params["col"] = c.Args()[1]
	}
	objs, err := GetDbRule(project, commandName, params)
	if err != nil {
		return err
	}
	if err := utils.PrintYaml(objs); err != nil {
		return err
	}
	return nil
}

func actionGetDbConfig(c *cli.Context) error {
	// Get the project and url parameters
	project := c.GlobalString("project")
	commandName := c.Command.Name

	params := map[string]string{}
	if len(c.Args()) != 0 {
		params["dbAlias"] = c.Args()[0]
	}
	objs, err := GetDbConfig(project, commandName, params)
	if err != nil {
		return err
	}

	if err := utils.PrintYaml(objs); err != nil {
		return err
	}
	return nil
}

func actionGetDbSchema(c *cli.Context) error {
	// Get the project and url parameters
	project := c.GlobalString("project")
	commandName := c.Command.Name

	params := map[string]string{}
	switch len(c.Args()) {
	case 1:
		params["dbAlias"] = c.Args()[0]
	case 2:
		params["dbAlias"] = c.Args()[0]
		params["col"] = c.Args()[1]
	}

	objs, err := GetDbSchema(project, commandName, params)
	if err != nil {
		return err
	}
	if err := utils.PrintYaml(objs); err != nil {
		return err
	}
	return nil
}

func actionGenerateDBRule(c *cli.Context) error {
	argsArr := c.Args()
	if len(argsArr) != 1 {
		return fmt.Errorf("incorrect number of arguments")
	}
	dbruleConfigFile := argsArr[0]
	dbrule, err := generateDBRule()
	if err != nil {
		return err
	}

	return utils.AppendConfigToDisk(dbrule, dbruleConfigFile)
}

func actionGenerateDBConfig(c *cli.Context) error {
	argsArr := c.Args()
	if len(argsArr) != 1 {
		return fmt.Errorf("incorrect number of arguments")
	}
	dbruleConfigFile := argsArr[0]
	dbrule, err := generateDBConfig()
	if err != nil {
		return err
	}

	return utils.AppendConfigToDisk(dbrule, dbruleConfigFile)
}

func actionGenerateDBSchema(c *cli.Context) error {
	argsArr := c.Args()
	if len(argsArr) != 1 {
		return fmt.Errorf("incorrect number of arguments")
	}
	dbruleConfigFile := argsArr[0]
	dbrule, err := generateDBSchema()
	if err != nil {
		return err
	}

	return utils.AppendConfigToDisk(dbrule, dbruleConfigFile)
}
