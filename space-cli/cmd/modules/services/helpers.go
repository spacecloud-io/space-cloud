package services

import (
	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils"
	"github.com/spf13/cobra"
)

func secretsAutoCompleteFun(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
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
}

func serviceRoutesAutoCompleteFun(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
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
}

func serviceRoleAutoCompleteFun(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	switch len(args) {
	case 0:
		project, check := utils.GetProjectID()
		if !check {
			utils.LogDebug("Project not specified in flag", nil)
			return nil, cobra.ShellCompDirectiveDefault
		}
		objs, err := GetServicesRole(project, "service-role", map[string]string{})
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
		objs, err := GetServicesRole(project, "service-role", map[string]string{})
		if err != nil {
			return nil, cobra.ShellCompDirectiveDefault
		}
		var roleID []string
		for _, v := range objs {
			roleID = append(roleID, v.Meta["roleId"])
		}
		return roleID, cobra.ShellCompDirectiveDefault
	}
	return nil, cobra.ShellCompDirectiveDefault
}

func servicesAutoCompleteFun(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
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
}
