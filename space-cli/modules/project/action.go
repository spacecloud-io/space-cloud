package project

import (
	"github.com/urfave/cli"

	"github.com/spaceuptech/space-cli/model"
	"github.com/spaceuptech/space-cli/utils"
)

//ActionGetProjectConfig gets global config
func ActionGetProjectConfig(c *cli.Context) error {
	// Get the project and cmd parameters
	project := c.GlobalString("project")
	commandName := c.Command.Name

	params := map[string]string{}
	obj, err := GetProjectConfig(project, commandName, params)
	if err != nil {
		return err
	}
	if err := utils.PrintYaml([]*model.SpecObject{obj}); err != nil {
		return err
	}
	return nil
}
