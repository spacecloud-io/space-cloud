package letsencrypt

import (
	"github.com/urfave/cli"

	"github.com/spaceuptech/space-cli/model"
	"github.com/spaceuptech/space-cli/utils"
)

//ActionGetLetsEncrypt gets encrypt domain
func ActionGetLetsEncrypt(c *cli.Context) error {
	// Get the project and url parameters
	project := c.GlobalString("project")
	commandName := c.Command.Name

	params := map[string]string{}
	obj, err := GetLetsEncryptDomain(project, commandName, params)
	if err != nil {
		return err
	}
	objs := []*model.SpecObject{obj}
	if err := utils.PrintYaml(objs); err != nil {
		return err
	}
	return nil
}
