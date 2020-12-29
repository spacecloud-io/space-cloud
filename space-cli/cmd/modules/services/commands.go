package services

import (
	"github.com/spf13/cobra"

	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils"
)

// GenerateSubCommands is the list of commands the services module exposes
func GenerateSubCommands() []*cobra.Command {

	var generateService = &cobra.Command{
		Use:     "service [path to config file]",
		RunE:    actionGenerateService,
		Aliases: []string{"services"},
		Example: "space-cli generate service config.yaml --project myproject --log-level info",
	}

	var generateServiceRoute = &cobra.Command{
		Use:     "service-route [path to config file]",
		RunE:    actionGenerateServiceRoute,
		Aliases: []string{"service-route"},
		Example: "space-cli generate service-route config.yaml --project myproject --log-level info",
	}

	var generateServiceRole = &cobra.Command{
		Use:     "service-role [path to config file]",
		RunE:    actionGenerateServiceRole,
		Aliases: []string{"service-roles"},
		Example: "space-cli generate service-role config.yaml --project myproject",
	}

	return []*cobra.Command{generateService, generateServiceRoute, generateServiceRole}

}

// GetSubCommands is the list of commands the services module exposes
func GetSubCommands() []*cobra.Command {

	var getServicesRoutes = &cobra.Command{
		Use:               "service-routes",
		RunE:              actionGetServicesRoutes,
		ValidArgsFunction: serviceRoutesAutoCompleteFun,
	}

	var getServicesRole = &cobra.Command{
		Use:               "service-role",
		RunE:              actionGetServicesRole,
		ValidArgsFunction: serviceRoleAutoCompleteFun,
	}

	var getServicesSecrets = &cobra.Command{
		Use:               "secrets",
		RunE:              actionGetServicesSecrets,
		ValidArgsFunction: secretsAutoCompleteFun,
	}

	var getServices = &cobra.Command{
		Use:               "services",
		RunE:              actionGetServices,
		ValidArgsFunction: servicesAutoCompleteFun,
	}

	return []*cobra.Command{getServicesRoutes, getServicesSecrets, getServices, getServicesRole}
}

func actionGetServicesRoutes(cmd *cobra.Command, args []string) error {
	// Get the project and url parameters
	project, check := utils.GetProjectID()
	if !check {
		return utils.LogError("Project not specified in flag", nil)
	}
	commandName := "service-route"

	params := map[string]string{}
	if len(args) != 0 {
		params["id"] = args[0]
	}

	objs, err := GetServicesRoutes(project, commandName, params)
	if err != nil {
		return err
	}
	if err := utils.PrintYaml(objs); err != nil {
		return err
	}
	return nil
}

func actionGetServicesRole(cmd *cobra.Command, args []string) error {
	// Get the project and url parameters
	project, check := utils.GetProjectID()
	if !check {
		return utils.LogError("Project not specified in flag", nil)
	}
	commandName := "service-role"

	params := map[string]string{}
	switch len(args) {
	case 1:
		params["serviceID"] = args[0]
	case 2:
		params["serviceID"] = args[0]
		params["roleID"] = args[1]
	}

	objs, err := GetServicesRole(project, commandName, params)
	if err != nil {
		return err
	}
	if err := utils.PrintYaml(objs); err != nil {
		return err
	}
	return nil
}

func actionGetServicesSecrets(cmd *cobra.Command, args []string) error {
	// Get the project and url parameters
	project, check := utils.GetProjectID()
	if !check {
		return utils.LogError("Project not specified in flag", nil)
	}
	commandName := "secret"

	params := map[string]string{}
	if len(args) != 0 {
		params["id"] = args[0]
	}

	objs, err := GetServicesSecrets(project, commandName, params)
	if err != nil {
		return err
	}
	if err := utils.PrintYaml(objs); err != nil {
		return err
	}
	return nil
}

func actionGetServices(cmd *cobra.Command, args []string) error {
	// Get the project and url parameters
	project, check := utils.GetProjectID()
	if !check {
		return utils.LogError("Project not specified in flag", nil)
	}
	commandName := "service"

	params := map[string]string{}
	switch len(args) {
	case 1:
		params["serviceId"] = args[0]
	case 2:
		params["serviceId"] = args[0]
		params["version"] = args[1]
	}
	objs, err := GetServices(project, commandName, params)
	if err != nil {
		return err
	}

	if err := utils.PrintYaml(objs); err != nil {
		return err
	}
	return nil
}

func actionGenerateService(cmd *cobra.Command, args []string) error {
	// get filename from args in which service config will be stored
	if len(args) != 1 {
		return utils.LogError("incorrect number of arguments. Use -h to check usage instructions", nil)
	}
	serviceConfigFile := args[0]

	service, err := GenerateService("", "")
	if err != nil {
		return err
	}

	return utils.AppendConfigToDisk(service, serviceConfigFile)
}

func actionGenerateServiceRoute(cmd *cobra.Command, args []string) error {

	// get filename from args in which service config will be stored
	if len(args) != 1 {
		return utils.LogError("incorrect number of arguments. Use -h to check usage instructions", nil)
	}
	serviceConfigFile := args[0]

	project, _ := utils.GetProjectID()

	serviceRoute, err := GenerateServiceRoute(project)
	if err != nil {
		return err
	}

	return utils.AppendConfigToDisk(serviceRoute, serviceConfigFile)
}

func actionGenerateServiceRole(cmd *cobra.Command, args []string) error {

	// get filename from args in which service config will be stored
	if len(args) != 1 {
		return utils.LogError("incorrect number of arguments. Use -h to check usage instructions", nil)
	}
	serviceConfigFile := args[0]

	project, _ := utils.GetProjectID()

	serviceRole, err := GenerateServiceRole(project)
	if err != nil {
		return err
	}

	return utils.AppendConfigToDisk(serviceRole, serviceConfigFile)
}

// DeleteSubCommands is the list of commands the services module exposes
func DeleteSubCommands() []*cobra.Command {
	var deleteServicesSecrets = &cobra.Command{
		Use:               "secrets",
		Aliases:           []string{"secret"},
		RunE:              actionDeleteServicesSecrets,
		ValidArgsFunction: secretsAutoCompleteFun,
		Example:           "space-cli delete secrets secretID --project myproject",
	}

	var deleteService = &cobra.Command{
		Use:               "service",
		Aliases:           []string{"services"},
		RunE:              actionDeleteService,
		ValidArgsFunction: servicesAutoCompleteFun,
		Example:           "space-cli delete service serviceId version --project myproject",
	}

	var deleteServiceRole = &cobra.Command{
		Use:               "service-role",
		Aliases:           []string{"service-roles"},
		RunE:              actionDeleteServiceRole,
		ValidArgsFunction: serviceRoleAutoCompleteFun,
		Example:           "space-cli delete service-role serviceId roleId --project myproject",
	}

	return []*cobra.Command{deleteServicesSecrets, deleteService, deleteServiceRole}
}

func actionDeleteServicesSecrets(cmd *cobra.Command, args []string) error {
	// Get the project
	project, check := utils.GetProjectID()
	if !check {
		return utils.LogError("Project not specified in flag", nil)
	}

	prefix := ""
	if len(args) != 0 {
		prefix = args[0]
	}

	return deleteSecret(project, prefix)
}

func actionDeleteService(cmd *cobra.Command, args []string) error {
	// Get the project
	project, check := utils.GetProjectID()
	if !check {
		return utils.LogError("Project not specified in flag", nil)
	}

	prefix := map[string]string{}
	switch len(args) {
	case 1:
		prefix["serviceId"] = args[0]
	case 2:
		prefix["serviceId"] = args[0]
		prefix["version"] = args[1]
	}

	return deleteService(project, prefix)
}

func actionDeleteServiceRole(cmd *cobra.Command, args []string) error {
	// Get the project
	project, check := utils.GetProjectID()
	if !check {
		return utils.LogError("Project not specified in flag", nil)
	}

	prefix := map[string]string{}
	switch len(args) {
	case 1:
		prefix["serviceId"] = args[0]
	case 2:
		prefix["serviceId"] = args[0]
		prefix["roleId"] = args[1]
	}

	return deleteServiceRole(project, prefix)
}
