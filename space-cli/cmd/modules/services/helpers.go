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
