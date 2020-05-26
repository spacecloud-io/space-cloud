package deploy

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/spaceuptech/space-cli/cmd/utils"
)

// Commands deploys a service
func Commands() []*cobra.Command {
	var commandDeploy = &cobra.Command{
		Use: "deploy",
		PreRun: func(cmd *cobra.Command, args []string) {
			err := viper.BindPFlag("project", cmd.Flags().Lookup("project"))
			if err != nil {
				_ = utils.LogError("Unable to bind the flag ('project')", err)
			}
			err = viper.BindPFlag("docker-file", cmd.Flags().Lookup("docker-file"))
			if err != nil {
				_ = utils.LogError("Unable to bind the flag ('docker-file')", err)
			}
			err = viper.BindPFlag("service-file", cmd.Flags().Lookup("service-file"))
			if err != nil {
				_ = utils.LogError("Unable to bind the flag ('service-file')", err)
			}
			err = viper.BindPFlag("prepare", cmd.Flags().Lookup("prepare"))
			if err != nil {
				_ = utils.LogError("Unable to bind the flag ('prepare')", err)
			}
		},
		RunE:          actionDeploy,
		SilenceErrors: true,
	}

	commandDeploy.Flags().StringP("project", "", "", "The project to deploy the service to.")
	commandDeploy.Flags().StringP("docker-file", "", "Dockerfile", "The path of the docker file")
	commandDeploy.Flags().StringP("service-file", "", "service.yaml", "The path of the service config file")
	commandDeploy.Flags().BoolP("prepare", "", false, "Prepare the configuration used for deploying service")

	return []*cobra.Command{commandDeploy}
}

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
