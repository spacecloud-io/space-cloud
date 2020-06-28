package addons

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/spaceuptech/space-cli/cmd/utils"
)

// Commands is the list of commands the addon module exposes
func Commands() []*cobra.Command {
	var addCmd = &cobra.Command{
		Use:           "add",
		Short:         "Add a add-on to the environment",
		SilenceErrors: true,
	}

	var addRegistryCmd = &cobra.Command{
		Use:   "registry",
		Short: "Add a docker registry",
		RunE:  ActionAddRegistry,
	}

	var addDatabaseCmd = &cobra.Command{
		Use:   "database",
		Short: "Add a database",
		PreRun: func(cmd *cobra.Command, args []string) {
			err := viper.BindPFlag("username", cmd.Flags().Lookup("username"))
			if err != nil {
				_ = utils.LogError("Unable to bind the flag ('username')", nil)
			}
			err = viper.BindPFlag("password", cmd.Flags().Lookup("password"))
			if err != nil {
				_ = utils.LogError("Unable to bind the flag ('password')", nil)
			}
			err = viper.BindPFlag("alias", cmd.Flags().Lookup("alias"))
			if err != nil {
				_ = utils.LogError("Unable to bind the flag ('alias')", nil)
			}
			err = viper.BindPFlag("version", cmd.Flags().Lookup("version"))
			if err != nil {
				_ = utils.LogError("Unable to bind the flag ('version')", nil)
			}
			err = viper.BindPFlag("auto-apply", cmd.Flags().Lookup("auto-apply"))
			if err != nil {
				_ = utils.LogError("Unable to bind the flag ('auto-apply')", nil)
			}
		},
		RunE: ActionAddDatabase,
	}

	addDatabaseCmd.Flags().StringP("username", "U", "", "provide the username")
	addDatabaseCmd.Flags().StringP("password", "P", "", "provide the password")
	addDatabaseCmd.Flags().StringP("alias", "", "", "provide the alias for the database")
	addDatabaseCmd.Flags().StringP("version", "", "latest", "provide the version of the database")
	addDatabaseCmd.Flags().BoolP("auto-apply", "", false, "add database in space cloud config")

	var removeCmd = &cobra.Command{
		Use:           "remove",
		Short:         "Remove a add-on from the environment",
		SilenceErrors: true,
	}

	var removeRegistryCmd = &cobra.Command{
		Use:   "registry",
		Short: "Remove a docker registry",
		RunE:  ActionRemoveRegistry,
	}

	var removeDatabaseCmd = &cobra.Command{
		Use:   "database",
		Short: "Remove a database",
		RunE:  ActionRemoveDatabase,
	}
	addCmd.AddCommand(addRegistryCmd)
	addCmd.AddCommand(addDatabaseCmd)
	removeCmd.AddCommand(removeRegistryCmd)
	removeCmd.AddCommand(removeDatabaseCmd)

	return []*cobra.Command{addCmd, removeCmd}
}

// ActionAddRegistry adds a registry add on
func ActionAddRegistry(cmd *cobra.Command, args []string) error {
	project, check := utils.GetProjectID()
	if !check {
		return utils.LogError("Project not specified in flag", nil)
	}
	return addRegistry(project)
}

// ActionRemoveRegistry removes a registry add on
func ActionRemoveRegistry(cmd *cobra.Command, args []string) error {
	project, check := utils.GetProjectID()
	if !check {
		return utils.LogError("Project not specified in flag", nil)
	}
	return removeRegistry(project)
}

// ActionAddDatabase adds a database add on
func ActionAddDatabase(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return utils.LogError("Database type not provided as an arguement", nil)
	}
	dbtype := args[0]
	username := viper.GetString("username")
	if username == "" {
		switch dbtype {
		case "postgres":
			username = "postgres"
		case "mysql":
			username = "root"
		case "sqlserver":
			username = "sa"
		}
	}
	password := viper.GetString("password")
	if password == "" {
		switch dbtype {
		case "postgres":
			password = "mysecretpassword"
		case "mysql":
			password = "my-secret-pw"
		case "sqlserver":
			password = "yourStrong(!)Password"
		}
	}
	alias := viper.GetString("alias")
	version := viper.GetString("version")
	return addDatabase(dbtype, username, password, alias, version)
}

// ActionRemoveDatabase removes a database add on
func ActionRemoveDatabase(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return utils.LogError("Database Alias not provided as an argument", nil)
	}
	return removeDatabase(args[0])
}
