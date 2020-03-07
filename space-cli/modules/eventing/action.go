package eventing

import (
	"github.com/urfave/cli"

	"github.com/spaceuptech/space-cli/model"
	"github.com/spaceuptech/space-cli/utils"
)

//ActionGetEventingTrigger gets eventing trigger
func ActionGetEventingTrigger(c *cli.Context) error {
	// Get the project and url parameters
	project := c.GlobalString("project")
	commandName := c.Command.Name

	params := map[string]string{}
	if len(c.Args()) != 0 {
		params["ruleName"] = c.Args()[0]
	}
	objs, err := getEventingTrigger(project, commandName, params)
	if err != nil {
		return err
	}
	if err := utils.PrintYaml(objs); err != nil {
		return err
	}
	return nil
}

//ActionGetEventingConfig gets eventing config
func ActionGetEventingConfig(c *cli.Context) error {
	// Get the project and url parameters
	project := c.GlobalString("project")
	commandName := c.Command.Name

	params := map[string]string{}
	obj, err := getEventingConfig(project, commandName, params)
	if err != nil {
		return err
	}
	objs := []*model.SpecObject{obj}
	if err := utils.PrintYaml(objs); err != nil {
		return err
	}
	return nil
}

//ActionGetEventingSchema gets eventing schema
func ActionGetEventingSchema(c *cli.Context) error {
	// Get the project and url parameters
	project := c.GlobalString("project")
	commandName := c.Command.Name

	params := map[string]string{}
	if len(c.Args()) != 0 {
		params["type"] = c.Args()[0]
	}
	objs, err := getEventingSchema(project, commandName, params)
	if err != nil {
		return err
	}
	if err := utils.PrintYaml(objs); err != nil {
		return err
	}
	return nil
}

//ActionGetEventingSecurityRule gets eventing security rule
func ActionGetEventingSecurityRule(c *cli.Context) error {
	// Get the project and url parameters
	project := c.GlobalString("project")
	commandName := c.Command.Name

	params := map[string]string{}
	if len(c.Args()) != 0 {
		params["type"] = c.Args()[0]
	}
	objs, err := getEventingSecurityRule(project, commandName, params)
	if err != nil {
		return err
	}
	if err := utils.PrintYaml(objs); err != nil {
		return err
	}
	return nil
}
