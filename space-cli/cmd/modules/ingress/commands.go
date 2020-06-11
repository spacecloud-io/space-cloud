package ingress

import (
	"github.com/spf13/cobra"

	"github.com/spaceuptech/space-cli/cmd/utils"
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
	project, check := utils.GetProjectID()
	if !check {
		_ = utils.LogError("Project not specified in flag", nil)
		return nil
	}
	commandName := "ingress-route"

	params := map[string]string{}
	if len(args) != 0 {
		params["id"] = args[0]
	}

	objs, err := GetIngressRoutes(project, commandName, params)
	if err != nil {
		return nil
	}

	if err := utils.PrintYaml(objs); err != nil {
		return nil
	}
	return nil
}

func actionGenerateIngressRouting(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		_ = utils.LogError("incorrect number of arguments", nil)
		return nil
	}
	dbruleConfigFile := args[0]
	dbrule, err := generateIngressRouting()
	if err != nil {
		return nil
	}

	_ = utils.AppendConfigToDisk(dbrule, dbruleConfigFile)
	return nil
}
