package ingress

import (
	"fmt"

	"github.com/spaceuptech/space-cli/utils"
	"github.com/urfave/cli"
)

// ActionGetIngressRoutes gets routes
func ActionGetIngressRoutes(c *cli.Context) error {
	// Get the project and url parameters
	project := c.GlobalString("project")
	commandName := c.Command.Name

	params := map[string]string{}
	if len(c.Args()) != 0 {
		params["routesID"] = c.Args()[0]
	}

	objs, err := GetIngressRoutes(project, commandName, params)
	if err != nil {
		return err
	}
	if err := utils.PrintYaml(objs); err != nil {
		return err
	}
	return nil
}

// ActionGenerateIngressRouting creates spec object for service routing
func ActionGenerateIngressRouting(c *cli.Context) error {
	argsArr := c.Args()
	if len(argsArr) != 1 {
		return fmt.Errorf("incorrect number of arguments")
	}
	dbruleConfigFile := argsArr[0]
	dbrule, err := generateIngressRouting()
	if err != nil {
		return err
	}

	return utils.AppendConfigToDisk(dbrule, dbruleConfigFile)
}
