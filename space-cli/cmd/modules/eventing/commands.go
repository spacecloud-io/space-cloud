package eventing

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/spaceuptech/space-cli/cmd/utils"
)

// GenerateSubCommands is the list of commands the eventing module exposes
func GenerateSubCommands() []*cobra.Command {

	var generatetrigger = &cobra.Command{
		Use:  "eventing-triggers",
		RunE: actionGenerateEventingTrigger,
	}

	var generateconfig = &cobra.Command{
		Use:  "eventing-config",
		RunE: actionGenerateEventingConfig,
	}

	var generateschema = &cobra.Command{
		Use:  "eventing-schema",
		RunE: actionGenerateEventingSchema,
	}

	var generaterule = &cobra.Command{
		Use:  "eventing-rule",
		RunE: actionGenerateEventingRule,
	}

	return []*cobra.Command{generatetrigger, generateconfig, generateschema, generaterule}
}

// GetSubCommands is the list of commands the eventing module exposes
func GetSubCommands() []*cobra.Command {

	var gettrigger = &cobra.Command{
		Use:  "eventing-triggers",
		RunE: actionGetEventingTrigger,
	}

	var getconfig = &cobra.Command{
		Use:  "eventing-config",
		RunE: actionGetEventingConfig,
	}

	var getschema = &cobra.Command{
		Use:  "eventing-schema",
		RunE: actionGetEventingSchema,
	}

	var getrule = &cobra.Command{
		Use:  "eventing-rule",
		RunE: actionGetEventingSecurityRule,
	}

	return []*cobra.Command{gettrigger, getconfig, getschema, getrule}
}

func actionGetEventingTrigger(cmd *cobra.Command, args []string) error {
	// Get the project and url parameters
	project := viper.GetString("project")
	commandName := cmd.Use

	params := map[string]string{}
	if len(args) != 0 {
		params["id"] = args[0]
	}
	objs, err := GetEventingTrigger(project, commandName, params)
	if err != nil {
		return nil
	}
	if err := utils.PrintYaml(objs); err != nil {
		return nil
	}
	return nil
}

func actionGetEventingConfig(cmd *cobra.Command, args []string) error {
	// Get the project and url parameters
	project := viper.GetString("project")
	commandName := cmd.Use

	params := map[string]string{}
	obj, err := GetEventingConfig(project, commandName, params)
	if err != nil {
		return nil
	}

	if err := utils.PrintYaml(obj); err != nil {
		return nil
	}
	return nil
}

func actionGetEventingSchema(cmd *cobra.Command, args []string) error {
	// Get the project and url parameters
	project := viper.GetString("project")
	commandName := cmd.Use

	params := map[string]string{}
	if len(args) != 0 {
		params["id"] = args[0]
	}
	objs, err := GetEventingSchema(project, commandName, params)
	if err != nil {
		return nil
	}

	if err := utils.PrintYaml(objs); err != nil {
		return nil
	}
	return nil
}

func actionGetEventingSecurityRule(cmd *cobra.Command, args []string) error {
	// Get the project and url parameters
	project := viper.GetString("project")
	commandName := cmd.Use

	params := map[string]string{}
	if len(args) != 0 {
		params["id"] = args[0]
	}
	objs, err := GetEventingSecurityRule(project, commandName, params)
	if err != nil {
		return nil
	}

	if err := utils.PrintYaml(objs); err != nil {
		return nil
	}
	return nil
}

func actionGenerateEventingRule(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		_ = utils.LogError("incorrect number of arguments", nil)
		return nil
	}
	dbruleConfigFile := args[0]
	dbrule, err := generateEventingRule()
	if err != nil {
		return nil
	}

	_ = utils.AppendConfigToDisk(dbrule, dbruleConfigFile)
	return nil
}

func actionGenerateEventingSchema(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		_ = utils.LogError("incorrect number of arguments", nil)
		return nil
	}
	dbruleConfigFile := args[0]
	dbrule, err := generateEventingSchema()
	if err != nil {
		return nil
	}

	_ = utils.AppendConfigToDisk(dbrule, dbruleConfigFile)
	return nil
}

func actionGenerateEventingConfig(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		_ = utils.LogError("incorrect number of arguments", nil)
		return nil
	}
	dbruleConfigFile := args[0]
	dbrule, err := generateEventingConfig()
	if err != nil {
		return nil
	}

	_ = utils.AppendConfigToDisk(dbrule, dbruleConfigFile)
	return nil
}

func actionGenerateEventingTrigger(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		_ = utils.LogError("incorrect number of arguments", nil)
		return nil
	}
	dbruleConfigFile := args[0]
	dbrule, err := generateEventingTrigger()
	if err != nil {
		return nil
	}

	_ = utils.AppendConfigToDisk(dbrule, dbruleConfigFile)
	return nil
}
