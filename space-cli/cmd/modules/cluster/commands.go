package cluster

import (
	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils"
	"github.com/spf13/cobra"
)

// GetSubCommands is the list of commands the cluster module exposes
func GetSubCommands() []*cobra.Command {
	var getClusterConfig = &cobra.Command{
		Use:     "cluster-config",
		Aliases: []string{"cluster-configs"},
		RunE:    actionGetClusterConfig,
	}

	var getIntegration = &cobra.Command{
		Use:     "integration",
		Aliases: []string{"integrations"},
		RunE:    actionGetIntegration,
	}

	return []*cobra.Command{getClusterConfig, getIntegration}
}

func actionGetClusterConfig(cmd *cobra.Command, args []string) error {

	objs, err := GetClusterConfig()
	if err != nil {
		return err
	}
	if err := utils.PrintYaml(objs); err != nil {
		return err
	}
	return nil
}

func actionGetIntegration(cmd *cobra.Command, args []string) error {

	objs, err := GetIntegration()
	if err != nil {
		return err
	}
	if err := utils.PrintYaml(objs); err != nil {
		return err
	}
	return nil
}
