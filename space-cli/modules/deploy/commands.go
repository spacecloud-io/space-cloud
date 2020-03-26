package deploy

import (
	"fmt"

	"github.com/spaceuptech/space-cli/utils"
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
	err := viper.BindPFlag("project", commandDeploy.Flags().Lookup("project"))
	if err != nil {
		_ = utils.LogError(fmt.Sprintf("Unable to bind the flag ('project')"), err)
	}

	commandDeploy.Flags().StringP("docker-file", "", "Dockerfile", "The path of the docker file")
	err = viper.BindPFlag("docker-file", commandDeploy.Flags().Lookup("docker-file"))
	if err != nil {
		_ = utils.LogError(fmt.Sprintf("Unable to bind the flag ('docker-file')"), err)
	}

	commandDeploy.Flags().StringP("service-file", "", "service.yaml", "The path of the service config file")
	err = viper.BindPFlag("service-file", commandDeploy.Flags().Lookup("service-file"))
	if err != nil {
		_ = utils.LogError(fmt.Sprintf("Unable to bind the flag ('service-file')"), err)
	}

	commandDeploy.Flags().StringP("prepare", "", "", "Prepare the configuration used for deploying service")
	err = viper.BindPFlag("prepare", commandDeploy.Flags().Lookup("prepare"))
	if err != nil {
		_ = utils.LogError(fmt.Sprintf("Unable to bind the flag ('prepare')"), err)
	}

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
