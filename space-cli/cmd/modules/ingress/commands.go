package ingress

import (
	"github.com/spf13/cobra"

	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils"
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

	var getroute = &cobra.Command{
		Use:  "ingress-route",
		RunE: actionGetIngressRoutes,
	}

	var getroutes = &cobra.Command{
		Use:  "ingress-routes",
		RunE: actionGetIngressRoutes,
	}

	return []*cobra.Command{getroute, getroutes}
}

func actionGetIngressRoutes(cmd *cobra.Command, args []string) error {
	// Get the project and url parameters
	project, check := utils.GetProjectID()
	if !check {
		return utils.LogError("Project not specified in flag", nil)
	}
	commandName := "ingress-route"

	params := map[string]string{}
	if len(args) != 0 {
		params["id"] = args[0]
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

func actionGenerateIngressRouting(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return utils.LogError("incorrect number of arguments", nil)
	}
	dbruleConfigFile := args[0]
	dbrule, err := generateIngressRouting()
	if err != nil {
		return err
	}

	return utils.AppendConfigToDisk(dbrule, dbruleConfigFile)
}
