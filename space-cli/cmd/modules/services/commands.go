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

	return []*cobra.Command{generateService}

}

// GetSubCommands is the list of commands the services module exposes
func GetSubCommands() []*cobra.Command {

	var getServicesRoutes = &cobra.Command{
		Use:  "service-routes",
		RunE: actionGetServicesRoutes,
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			project, check := utils.GetProjectID()
			if !check {
				utils.LogDebug("Project not specified in flag", nil)
				return nil, cobra.ShellCompDirectiveDefault
			}
			obj, err := GetServicesRoutes(project, "service-route", map[string]string{})
			if err != nil {
				return nil, cobra.ShellCompDirectiveDefault
			}
			var ids []string
			for _, v := range obj {
				ids = append(ids, v.Meta["id"])
			}
			return ids, cobra.ShellCompDirectiveDefault
		},
	}

	var getServicesSecrets = &cobra.Command{
		Use:  "secrets",
		RunE: actionGetServicesSecrets,
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			project, check := utils.GetProjectID()
			if !check {
				utils.LogDebug("Project not specified in flag", nil)
				return nil, cobra.ShellCompDirectiveDefault
			}
			obj, err := GetServicesSecrets(project, "secret", map[string]string{})
			if err != nil {
				return nil, cobra.ShellCompDirectiveDefault
			}
			var ids []string
			for _, v := range obj {
				ids = append(ids, v.Meta["id"])
			}
			return ids, cobra.ShellCompDirectiveDefault
		},
	}

	var getServices = &cobra.Command{
		Use:  "services",
		RunE: actionGetServices,
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			switch len(args) {
			case 0:
				project, check := utils.GetProjectID()
				if !check {
					utils.LogDebug("Project not specified in flag", nil)
					return nil, cobra.ShellCompDirectiveDefault
				}
				objs, err := GetServices(project, "service", map[string]string{})
				if err != nil {
					return nil, cobra.ShellCompDirectiveDefault
				}
				var serviceIds []string
				for _, v := range objs {
					serviceIds = append(serviceIds, v.Meta["serviceId"])
				}
				return serviceIds, cobra.ShellCompDirectiveDefault
			case 1:
				project, check := utils.GetProjectID()
				if !check {
					utils.LogDebug("Project not specified in flag", nil)
					return nil, cobra.ShellCompDirectiveDefault
				}
				objs, err := GetServices(project, "service", map[string]string{})
				if err != nil {
					return nil, cobra.ShellCompDirectiveDefault
				}
				var versions []string
				for _, v := range objs {
					versions = append(versions, v.Meta["version"])
				}
				return versions, cobra.ShellCompDirectiveDefault
			}
			return nil, cobra.ShellCompDirectiveDefault
		},
	}

	return []*cobra.Command{getServicesRoutes, getServicesSecrets, getServices}
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
