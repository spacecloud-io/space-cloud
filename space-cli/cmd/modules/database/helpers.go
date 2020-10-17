package database

import (
	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils"
	"github.com/spf13/cobra"
)

func validDBRulesArgsFunc(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	switch len(args) {
	case 0:
		project, check := utils.GetProjectID()
		if !check {
			utils.LogDebug("Project not specified in flag", nil)
			return nil, cobra.ShellCompDirectiveDefault
		}
		objs, err := GetDbRule(project, "db-rule", map[string]string{})
		if err != nil {
			return nil, cobra.ShellCompDirectiveDefault
		}
		var dbAlias []string
		for _, v := range objs {
			dbAlias = append(dbAlias, v.Meta["dbAlias"])
		}
		return dbAlias, cobra.ShellCompDirectiveDefault
	case 1:
		project, check := utils.GetProjectID()
		if !check {
			utils.LogDebug("Project not specified in flag", nil)
			return nil, cobra.ShellCompDirectiveDefault
		}
		objs, err := GetDbRule(project, "db-rule", map[string]string{})
		if err != nil {
			return nil, cobra.ShellCompDirectiveDefault
		}
		var col []string
		for _, v := range objs {
			col = append(col, v.Meta["col"])
		}
		return col, cobra.ShellCompDirectiveDefault
	}
	return nil, cobra.ShellCompDirectiveDefault
}

func validDBConfigArgsFunc(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	project, check := utils.GetProjectID()
	if !check {
		utils.LogDebug("Project not specified in flag", nil)
		return nil, cobra.ShellCompDirectiveDefault
	}
	objs, err := GetDbConfig(project, "db-config", map[string]string{})
	if err != nil {
		return nil, cobra.ShellCompDirectiveDefault
	}
	var dbAlias []string
	for _, v := range objs {
		dbAlias = append(dbAlias, v.Meta["dbAlias"])
	}
	return dbAlias, cobra.ShellCompDirectiveDefault
}

func validDBSchemasArgsFunc(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {

	switch len(args) {
	case 0:
		project, check := utils.GetProjectID()
		if !check {
			utils.LogDebug("Project not specified in flag", nil)
			return nil, cobra.ShellCompDirectiveDefault
		}
		objs, err := GetDbSchema(project, "db-schema", map[string]string{})
		if err != nil {
			return nil, cobra.ShellCompDirectiveDefault
		}
		var dbAlias []string
		for _, v := range objs {
			dbAlias = append(dbAlias, v.Meta["dbAlias"])
		}
		return dbAlias, cobra.ShellCompDirectiveDefault
	case 1:
		project, check := utils.GetProjectID()
		if !check {
			utils.LogDebug("Project not specified in flag", nil)
			return nil, cobra.ShellCompDirectiveDefault
		}
		objs, err := GetDbSchema(project, "db-schema", map[string]string{})
		if err != nil {
			return nil, cobra.ShellCompDirectiveDefault
		}
		var col []string
		for _, v := range objs {
			col = append(col, v.Meta["col"])
		}
		return col, cobra.ShellCompDirectiveDefault
	}
	return nil, cobra.ShellCompDirectiveDefault
}

func validDBPreparedQueriesArgsFunc(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	switch len(args) {
	case 0:
		project, check := utils.GetProjectID()
		if !check {
			utils.LogDebug("Project not specified in flag", nil)
			return nil, cobra.ShellCompDirectiveDefault
		}
		objs, err := GetDbSchema(project, "db-schema", map[string]string{})
		if err != nil {
			return nil, cobra.ShellCompDirectiveDefault
		}
		var dbAlias []string
		for _, v := range objs {
			dbAlias = append(dbAlias, v.Meta["dbAlias"])
		}
		return dbAlias, cobra.ShellCompDirectiveDefault
	case 1:
		project, check := utils.GetProjectID()
		if !check {
			utils.LogDebug("Project not specified in flag", nil)
			return nil, cobra.ShellCompDirectiveDefault
		}
		objs, err := GetDbPreparedQuery(project, "db-prepared-query", map[string]string{})
		if err != nil {
			return nil, cobra.ShellCompDirectiveDefault
		}
		var col []string
		for _, v := range objs {
			col = append(col, v.Meta["id"])
		}
		return col, cobra.ShellCompDirectiveDefault
	}
	return nil, cobra.ShellCompDirectiveDefault
}
