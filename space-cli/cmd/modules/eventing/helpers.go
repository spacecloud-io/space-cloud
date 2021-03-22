package eventing

import (
	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils"
	"github.com/spf13/cobra"
)

func eventingTriggersAutoCompleteFun(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	project, check := utils.GetProjectID()
	if !check {
		utils.LogDebug("Project not specified in flag", nil)
		return nil, cobra.ShellCompDirectiveDefault
	}
	objs, err := GetEventingTrigger(project, "eventing-trigger", map[string]string{})
	if err != nil {
		return nil, cobra.ShellCompDirectiveDefault
	}
	var ids []string
	for _, v := range objs {
		ids = append(ids, v.Meta["id"])
	}
	return ids, cobra.ShellCompDirectiveDefault
}

func eventingSchemasAutoCompleteFun(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	project, check := utils.GetProjectID()
	if !check {
		utils.LogDebug("Project not specified in flag", nil)
		return nil, cobra.ShellCompDirectiveDefault
	}
	objs, err := GetEventingSchema(project, "eventing-schema", map[string]string{})
	if err != nil {
		return nil, cobra.ShellCompDirectiveDefault
	}
	var ids []string
	for _, v := range objs {
		ids = append(ids, v.Meta["id"])
	}
	return ids, cobra.ShellCompDirectiveDefault
}

func eventingRulesAutoCompleteFun(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	project, check := utils.GetProjectID()
	if !check {
		utils.LogDebug("Project not specified in flag", nil)
		return nil, cobra.ShellCompDirectiveDefault
	}
	objs, err := GetEventingSecurityRule(project, "eventing-rule", map[string]string{})
	if err != nil {
		return nil, cobra.ShellCompDirectiveDefault
	}
	var ids []string
	for _, v := range objs {
		ids = append(ids, v.Meta["id"])
	}
	return ids, cobra.ShellCompDirectiveDefault
}
