package project

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func projectAutoCompletionFun(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	project := viper.GetString("project")
	objs, err := GetProjectConfig(project, "project", map[string]string{})
	if err != nil {
		return nil, cobra.ShellCompDirectiveDefault
	}
	var ids []string
	for _, v := range objs {
		ids = append(ids, v.Meta["id"])
	}
	return ids, cobra.ShellCompDirectiveDefault
}
