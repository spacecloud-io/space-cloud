package ingress

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils"
)

// GenerateSubCommands is the list of commands the ingress module exposes
func GenerateSubCommands() []*cobra.Command {

	var generateroutes = &cobra.Command{
		Use:     "ingress-route [path to config file]",
		RunE:    actionGenerateIngressRouting,
		Aliases: []string{"ingress-routes"},
		Example: "space-cli generate ingress-route config.yaml --project myproject --log-level info",
	}

	var generateIngressGlobal = &cobra.Command{
		Use:     "ingress-global [path to config file]",
		RunE:    actionGenerateIngressGlobal,
		Example: "space-cli generate ingress-global config.yaml --project myproject --log-level info",
	}

	return []*cobra.Command{generateroutes, generateIngressGlobal}
}

// GetSubCommands is the list of commands the ingress module exposes
func GetSubCommands() []*cobra.Command {

	var getroutes = &cobra.Command{
		Use:               "ingress-routes",
		Aliases:           []string{"ingress-route"},
		RunE:              actionGetIngressRoutes,
		ValidArgsFunction: ingressRoutesAutoCompleteFun,
	}

	var getIngressGlobal = &cobra.Command{
		Use:  "ingress-global",
		RunE: actionGetIngressGlobal,
	}

	return []*cobra.Command{getroutes, getIngressGlobal}
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
	filters := viper.GetStringSlice("filter")
	objs, err := GetIngressRoutes(project, commandName, params, filters)
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
		return utils.LogError("Incorrect number of arguments. Use -h to check usage instructions", nil)
	}
	dbruleConfigFile := args[0]
	dbrule, err := generateIngressRouting()
	if err != nil {
		return err
	}

	return utils.AppendConfigToDisk(dbrule, dbruleConfigFile)
}

func actionGenerateIngressGlobal(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return utils.LogError("Incorrect number of arguments. Use -h to check usage instructions", nil)
	}
	dbruleConfigFile := args[0]
	dbrule, err := generateIngressGlobal()
	if err != nil {
		return err
	}

	return utils.AppendConfigToDisk(dbrule, dbruleConfigFile)
}

func actionGetIngressGlobal(cmd *cobra.Command, args []string) error {
	// Get the project and url parameters
	project, check := utils.GetProjectID()
	if !check {
		return utils.LogError("Project not specified in flag", nil)
	}

	objs, err := GetIngressGlobal(project, "ingress-global")
	if err != nil {
		return err
	}

	if err := utils.PrintYaml(objs); err != nil {
		return err
	}
	return nil
}

// DeleteSubCommands is the list of commands the ingress module exposes
func DeleteSubCommands() []*cobra.Command {

	var deleteRoutes = &cobra.Command{
		Use:               "ingress-routes",
		Aliases:           []string{"ingress-route"},
		RunE:              actionDeleteIngressRoutes,
		ValidArgsFunction: ingressRoutesAutoCompleteFun,
		Example:           "space-cli delete ingress-routes routeID --project myproject --log-level info",
	}

	var deleteIngressGlobal = &cobra.Command{
		Use:     "ingress-global",
		RunE:    actionDeleteIngressGlobal,
		Example: "space-cli delete ingress-global --project myproject --log-level info",
	}

	return []*cobra.Command{deleteRoutes, deleteIngressGlobal}
}

func actionDeleteIngressRoutes(cmd *cobra.Command, args []string) error {
	// Get the project and url parameters
	project, check := utils.GetProjectID()
	if !check {
		return utils.LogError("Project not specified in flag", nil)
	}

	prefix := ""
	if len(args) != 0 {
		prefix = args[0]
	}

	return deleteIngressRoute(project, prefix)
}

func actionDeleteIngressGlobal(cmd *cobra.Command, args []string) error {
	// Get the project
	project, check := utils.GetProjectID()
	if !check {
		return utils.LogError("Project not specified in flag", nil)
	}

	return deleteIngressGlobalConfig(project)
}
