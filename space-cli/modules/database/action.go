package database

import (
	"github.com/urfave/cli"

	"github.com/spaceuptech/space-cli/model"
	"github.com/spaceuptech/space-cli/utils"
)

//ActionGetDbRule gets database rule
func ActionGetDbRule(c *cli.Context) error {
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
	objs, err := getDbRule(project, commandName, params)
	if err != nil {
		return err
	}
	if err := utils.PrintYaml(objs); err != nil {
		return err
	}
	return nil
}

//ActionGetDbConfig gets database config
func ActionGetDbConfig(c *cli.Context) error {
	// Get the project and url parameters
	project := c.GlobalString("project")
	commandName := c.Command.Name

	params := map[string]string{}
	if len(c.Args()) != 0 {
		params["dbAlias"] = c.Args()[0]
		//params["col"] = c.Args()[1]
	}

	obj, err := getDbConfig(project, commandName, params)
	if err != nil {
		return err
	}
	objs := []*model.SpecObject{obj}
	if err := utils.PrintYaml(objs); err != nil {
		return err
	}
	return nil
}

//ActionGetDbSchema gets database schema
func ActionGetDbSchema(c *cli.Context) error {
	// Get the project and url parameters
	project := c.GlobalString("project")
	commandName := c.Command.Name

	params := map[string]string{}
	switch len(c.Args()) {
	case 1:
		params["dbType"] = c.Args()[0]
	case 2:
		params["dbType"] = c.Args()[0]
		params["col"] = c.Args()[1]
	}

	objs, err := getDbSchema(project, commandName, params)
	if err != nil {
		return err
	}
	if err := utils.PrintYaml(objs); err != nil {
		return err
	}
	return nil
}
