package deploy

import "github.com/urfave/cli"

// CommandDeploy deploys a service
var CommandDeploy = cli.Command{
	Name: "deploy",
	Flags: []cli.Flag{
		cli.StringFlag{Name: "project", Usage: "The project to deploy the service to."},
		cli.StringFlag{Name: "dockerfile", Usage: "The path of the docker file", Value: "Dockerfile"},
		cli.StringFlag{Name: "service-file", Usage: "The path of the service config file", Value: "service.yaml"},
		cli.BoolFlag{Name: "prepare", Usage: "Prepare the configuration used for deploying service"},
	},
	Action: actionDeploy,
}

func actionDeploy(c *cli.Context) error {
	projectID := c.String("project")
	dockerFilePath := c.String("dockerfile")
	serviceFilePath := c.String("service-file")
	prepare := c.Bool("prepare")

	// Prepare configuration files
	if prepare {
		return prepareService(projectID, dockerFilePath, serviceFilePath)
	}

	return deployService(dockerFilePath, serviceFilePath)
}
