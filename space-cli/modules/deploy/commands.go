package deploy

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Commands deploys a service
func Commands() []*cobra.Command {
	var commandDeploy = &cobra.Command{
		Use:  "deploy",
		RunE: actionDeploy,
	}
	commandDeploy.Flags().StringP("project", "", "", "The project to deploy the service to.")
	viper.BindPFlag("project", commandDeploy.Flags().Lookup("project"))

	commandDeploy.Flags().StringP("docker-file", "", "Dockerfile", "The path of the docker file")
	viper.BindPFlag("docker-file", commandDeploy.Flags().Lookup("docker-file"))

	commandDeploy.Flags().StringP("service-file", "", "service.yaml", "The path of the service config file")
	viper.BindPFlag("service-file", commandDeploy.Flags().Lookup("service-file"))

	commandDeploy.Flags().StringP("prepare", "", "", "Prepare the configuration used for deploying service")
	viper.BindPFlag("prepare", commandDeploy.Flags().Lookup("prepare"))

	command := make([]*cobra.Command, 0)
	command = append(command, commandDeploy)
	return command
}

// // CommandDeploy deploys a service
// var CommandDeploy = cli.Command{
// 	Name: "deploy",
// 	Flags: []cli.Flag{
// 		cli.StringFlag{Name: "project", Usage: "The project to deploy the service to."},
// 		cli.StringFlag{Name: "docker-file", Usage: "The path of the docker file", Value: "Dockerfile"},
// 		cli.StringFlag{Name: "service-file", Usage: "The path of the service config file", Value: "service.yaml"},
// 		cli.BoolFlag{Name: "prepare", Usage: "Prepare the configuration used for deploying service"},
// 	},
// 	Action: actionDeploy,
// }

func actionDeploy(cmd *cobra.Command, args []string) error {
	projectID := viper.GetString("project")
	dockerFilePath := viper.GetString("docker-file")
	serviceFilePath := viper.GetString("service-file")
	prepare := viper.GetBool("prepare")

	// Prepare configuration files
	if prepare {
		return prepareService(projectID, dockerFilePath, serviceFilePath)
	}

	return deployService(dockerFilePath, serviceFilePath)
}
