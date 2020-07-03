package database

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/spaceuptech/space-cli/cmd/utils"
)

// GenerateSubCommands is the list of commands the database module exposes
func GenerateSubCommands() []*cobra.Command {

	var generaterule = &cobra.Command{
		Use:  "db-rules",
		RunE: actionGenerateDBRule,
	}

	var generateconfig = &cobra.Command{
		Use:  "db-config",
		RunE: actionGenerateDBConfig,
	}

	var generateschema = &cobra.Command{
		Use:  "db-schema",
		RunE: actionGenerateDBSchema,
	}
	return []*cobra.Command{generaterule, generateconfig, generateschema}
}

// GetSubCommands is the list of commands the database module exposes
func GetSubCommands() []*cobra.Command {

	var getrules = &cobra.Command{
		Use:     "db-rules",
		Aliases: []string{"db-rule"},
		RunE:    actionGetDbRules,
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			switch len(args) {
			case 0:
				fmt.Printf("case 0")
				project, check := utils.GetProjectID()
				if !check {
					_ = utils.LogError("Project not specified in flag", nil)
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
				fmt.Printf("case 1")
				project, check := utils.GetProjectID()
				if !check {
					_ = utils.LogError("Project not specified in flag", nil)
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
		},
	}

	var getconfigs = &cobra.Command{
		Use:     "db-configs",
		Aliases: []string{"db-config"},
		RunE:    actionGetDbConfig,
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			project, check := utils.GetProjectID()
			if !check {
				_ = utils.LogError("Project not specified in flag", nil)
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
		},
	}

	var getschemas = &cobra.Command{
		Use:     "db-schemas",
		Aliases: []string{"db-schema"},
		RunE:    actionGetDbSchema,
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {

			switch len(args) {
			case 0:
				project, check := utils.GetProjectID()
				if !check {
					_ = utils.LogError("Project not specified in flag", nil)
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
					_ = utils.LogError("Project not specified in flag", nil)
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
		},
	}

	return []*cobra.Command{getrules, getconfigs, getschemas}
}

func actionGetDbRules(cmd *cobra.Command, args []string) error {
	// Get the project and url parameters
	project, check := utils.GetProjectID()
	if !check {
		return utils.LogError("Project not specified in flag", nil)
	}
	commandName := "db-rule"

	params := map[string]string{}
	switch len(args) {
	case 1:
		params["dbAlias"] = args[0]
	case 2:
		params["dbAlias"] = args[0]
		params["col"] = args[1]
	}
	objs, err := GetDbRule(project, commandName, params)
	if err != nil {
		return err
	}

	if err := utils.PrintYaml(objs); err != nil {
		return err
	}
	return nil
}

func actionGetDbConfig(cmd *cobra.Command, args []string) error {
	// Get the project and url parameters
	project, check := utils.GetProjectID()
	if !check {
		return utils.LogError("Project not specified in flag", nil)
	}
	commandName := "db-config"

	params := map[string]string{}
	if len(args) != 0 {
		params["dbAlias"] = args[0]
	}
	objs, err := GetDbConfig(project, commandName, params)
	if err != nil {
		return err
	}

	if err := utils.PrintYaml(objs); err != nil {
		return err
	}
	return nil
}

func actionGetDbSchema(cmd *cobra.Command, args []string) error {
	// Get the project and url parameters
	project, check := utils.GetProjectID()
	if !check {
		return utils.LogError("Project not specified in flag", nil)
	}
	commandName := "db-schema"

	params := map[string]string{}
	switch len(args) {
	case 1:
		params["dbAlias"] = args[0]
	case 2:
		params["dbAlias"] = args[0]
		params["col"] = args[1]
	}

	objs, err := GetDbSchema(project, commandName, params)
	if err != nil {
		return err
	}
	if err := utils.PrintYaml(objs); err != nil {
		return err
	}
	return nil
}

func actionGenerateDBRule(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return utils.LogError("incorrect number of arguments", nil)
	}
	dbruleConfigFile := args[0]
	dbrule, err := generateDBRule()
	if err != nil {
		return err
	}

	return utils.AppendConfigToDisk(dbrule, dbruleConfigFile)
}

func actionGenerateDBConfig(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return utils.LogError("incorrect number of arguments", nil)
	}
	dbruleConfigFile := args[0]
	dbrule, err := generateDBConfig()
	if err != nil {
		return err
	}

	return utils.AppendConfigToDisk(dbrule, dbruleConfigFile)
}

func actionGenerateDBSchema(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return utils.LogError("incorrect number of arguments", nil)
	}
	dbruleConfigFile := args[0]
	dbrule, err := generateDBSchema()
	if err != nil {
		return err
	}

	return utils.AppendConfigToDisk(dbrule, dbruleConfigFile)
}
