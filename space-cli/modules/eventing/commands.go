package eventing

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/spaceuptech/space-cli/utils"
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
	objs, _ := GetEventingTrigger(project, commandName, params)
	_ = utils.PrintYaml(objs)
	return nil
}

func actionGetEventingConfig(cmd *cobra.Command, args []string) error {
	// Get the project and url parameters
	project := viper.GetString("project")
	commandName := cmd.Use

	params := map[string]string{}
	obj, _ := GetEventingConfig(project, commandName, params)

	_ = utils.PrintYaml(obj)
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
	objs, _ := GetEventingSchema(project, commandName, params)
	_ = utils.PrintYaml(objs)
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
	objs, _ := GetEventingSecurityRule(project, commandName, params)
	_ = utils.PrintYaml(objs)
	return nil
}

func actionGenerateEventingRule(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("incorrect number of arguments")
	}
	dbruleConfigFile := args[0]
	dbrule, _ := generateEventingRule()

	_ = utils.AppendConfigToDisk(dbrule, dbruleConfigFile)
	return nil
}

func actionGenerateEventingSchema(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("incorrect number of arguments")
	}
	dbruleConfigFile := args[0]
	dbrule, _ := generateEventingSchema()

	_ = utils.AppendConfigToDisk(dbrule, dbruleConfigFile)
	return nil
}

func actionGenerateEventingConfig(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("incorrect number of arguments")
	}
	dbruleConfigFile := args[0]
	dbrule, _ := generateEventingConfig()

	_ = utils.AppendConfigToDisk(dbrule, dbruleConfigFile)
	return nil
}

func actionGenerateEventingTrigger(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("incorrect number of arguments")
	}
	dbruleConfigFile := args[0]
	dbrule, _ := generateEventingTrigger()

	_ = utils.AppendConfigToDisk(dbrule, dbruleConfigFile)
	return nil
}
