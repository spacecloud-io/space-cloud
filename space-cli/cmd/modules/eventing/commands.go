package eventing

import (
	"github.com/spf13/cobra"

	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils"
)

// GenerateSubCommands is the list of commands the eventing module exposes
func GenerateSubCommands() []*cobra.Command {

	var generatetrigger = &cobra.Command{
		Use:     "eventing-trigger [path to config file]",
		RunE:    actionGenerateEventingTrigger,
		Aliases: []string{"eventing-triggers"},
		Example: "space-cli generate eventing-trigger config.yaml --project myproject --log-level info",
	}

	var generateconfig = &cobra.Command{
		Use:     "eventing-config [path to config file]",
		RunE:    actionGenerateEventingConfig,
		Aliases: []string{"eventing-configs"},
		Example: "space-cli generate eventing-config config.yaml --project myproject --log-level info",
	}

	var generateschema = &cobra.Command{
		Use:     "eventing-schema [path to config file]",
		RunE:    actionGenerateEventingSchema,
		Aliases: []string{"eventing-schemas"},
		Example: "space-cli generate eventing-schema config.yaml --project myproject --log-level info",
	}

	var generaterule = &cobra.Command{
		Use:     "eventing-rule [path to config file]",
		RunE:    actionGenerateEventingRule,
		Aliases: []string{"eventing-rules"},
		Example: "space-cli generate eventing-rule config.yaml --project myproject --log-level info",
	}

	return []*cobra.Command{generatetrigger, generateconfig, generateschema, generaterule}
}

// GetSubCommands is the list of commands the eventing module exposes
func GetSubCommands() []*cobra.Command {

	var gettrigger = &cobra.Command{
		Use:  "eventing-trigger",
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

	var gettriggers = &cobra.Command{
		Use:  "eventing-triggers",
		RunE: actionGetEventingTrigger,
	}

	var getconfigs = &cobra.Command{
		Use:  "eventing-configs",
		RunE: actionGetEventingConfig,
	}

	var getschemas = &cobra.Command{
		Use:  "eventing-schemas",
		RunE: actionGetEventingSchema,
	}

	var getrules = &cobra.Command{
		Use:  "eventing-rules",
		RunE: actionGetEventingSecurityRule,
	}

	return []*cobra.Command{gettrigger, getconfig, getschema, getrule, gettriggers, getconfigs, getschemas, getrules}
}

func actionGetEventingTrigger(cmd *cobra.Command, args []string) error {
	// Get the project and url parameters
	project, check := utils.GetProjectID()
	if !check {
		return utils.LogError("Project not specified in flag", nil)
	}
	commandName := "eventing-trigger"

	params := map[string]string{}
	if len(args) != 0 {
		params["id"] = args[0]
	}
	objs, err := GetEventingTrigger(project, commandName, params)
	if err != nil {
		return err
	}
	if err := utils.PrintYaml(objs); err != nil {
		return err
	}
	return nil
}

func actionGetEventingConfig(cmd *cobra.Command, args []string) error {
	// Get the project and url parameters
	project, check := utils.GetProjectID()
	if !check {
		return utils.LogError("Project not specified in flag", nil)
	}
	commandName := "eventing-config"

	params := map[string]string{}
	obj, err := GetEventingConfig(project, commandName, params)
	if err != nil {
		return err
	}

	if err := utils.PrintYaml(obj); err != nil {
		return err
	}
	return nil
}

func actionGetEventingSchema(cmd *cobra.Command, args []string) error {
	// Get the project and url parameters
	project, check := utils.GetProjectID()
	if !check {
		return utils.LogError("Project not specified in flag", nil)
	}
	commandName := "eventing-schema"

	params := map[string]string{}
	if len(args) != 0 {
		params["id"] = args[0]
	}
	objs, err := GetEventingSchema(project, commandName, params)
	if err != nil {
		return err
	}

	if err := utils.PrintYaml(objs); err != nil {
		return err
	}
	return nil
}

func actionGetEventingSecurityRule(cmd *cobra.Command, args []string) error {
	// Get the project and url parameters
	project, check := utils.GetProjectID()
	if !check {
		return utils.LogError("Project not specified in flag", nil)
	}
	commandName := "eventing-rule"

	params := map[string]string{}
	if len(args) != 0 {
		params["id"] = args[0]
	}
	objs, err := GetEventingSecurityRule(project, commandName, params)
	if err != nil {
		return err
	}

	if err := utils.PrintYaml(objs); err != nil {
		return err
	}
	return nil
}

func actionGenerateEventingRule(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return utils.LogError("incorrect number of arguments. Use -h to check usage instructions", nil)
	}
	dbruleConfigFile := args[0]
	dbrule, err := generateEventingRule()
	if err != nil {
		return err
	}

	return utils.AppendConfigToDisk(dbrule, dbruleConfigFile)
}

func actionGenerateEventingSchema(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return utils.LogError("incorrect number of arguments. Use -h to check usage instructions", nil)
	}
	dbruleConfigFile := args[0]
	dbrule, err := generateEventingSchema()
	if err != nil {
		return err
	}

	return utils.AppendConfigToDisk(dbrule, dbruleConfigFile)
}

func actionGenerateEventingConfig(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return utils.LogError("incorrect number of arguments. Use -h to check usage instructions", nil)
	}
	dbruleConfigFile := args[0]
	dbrule, err := generateEventingConfig()
	if err != nil {
		return err
	}

	return utils.AppendConfigToDisk(dbrule, dbruleConfigFile)
}

func actionGenerateEventingTrigger(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return utils.LogError("incorrect number of arguments. Use -h to check usage instructions", nil)
	}
	dbruleConfigFile := args[0]
	dbrule, err := generateEventingTrigger()
	if err != nil {
		return err
	}

	return utils.AppendConfigToDisk(dbrule, dbruleConfigFile)
}
