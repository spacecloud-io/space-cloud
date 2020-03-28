package ingress

import (
	"fmt"

	"github.com/spaceuptech/space-cli/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// GenerateSubCommands is the list of commands the ingress module exposes
func GenerateSubCommands() []*cobra.Command {

	var generateroutes = &cobra.Command{
		Use:  "ingress-routes",
		RunE: actionGenerateIngressRouting,
	}

	return []*cobra.Command{generateroutes}
}

// GetSubCommands is the list of commands the ingress module exposes
func GetSubCommands() []*cobra.Command {

	var getroutes = &cobra.Command{
		Use:  "ingress-routes",
		RunE: actionGetIngressRoutes,
	}

	return []*cobra.Command{getroutes}
}

func actionGetIngressRoutes(cmd *cobra.Command, args []string) error {
	// Get the project and url parameters
	project := viper.GetString("project")
	commandName := cmd.Use

	params := map[string]string{}
	if len(args) != 0 {
		params["id"] = args[0]
	}

	objs, _ := GetIngressRoutes(project, commandName, params)
	_ = utils.PrintYaml(objs)
	return nil
}

func actionGenerateIngressRouting(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("incorrect number of arguments")
	}
	dbruleConfigFile := args[0]
	dbrule, _ := generateIngressRouting()

	_ = utils.AppendConfigToDisk(dbrule, dbruleConfigFile)
	return nil
}
