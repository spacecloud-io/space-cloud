package addons

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/spaceuptech/space-cli/utils"
)

// Commands is the list of commands the addon module exposes
func Commands() []*cobra.Command {
	var addCmd = &cobra.Command{
		Use:   "add",
		Short: "Add a add-on to the environment",
	}

	var addRegistryCmd = &cobra.Command{
		Use:   "registry",
		Short: "Add a docker registry",
		RunE:  ActionAddRegistry,
	}

	var addDatabaseCmd = &cobra.Command{
		Use:   "database",
		Short: "Add a database",
		RunE:  ActionAddDatabase,
	}

	var removeCmd = &cobra.Command{
		Use:   "remove",
		Short: "Remove a add-on from the environment",
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

	addRegistryCmd.Flags().StringP("username", "U", "", "provide the username")
	err := viper.BindPFlag("username", addRegistryCmd.Flags().Lookup("username"))
	if err != nil {
		utils.LogError("", err)
	}

	addRegistryCmd.Flags().StringP("password", "P", "", "provide the password")
	err = viper.BindPFlag("password", addRegistryCmd.Flags().Lookup("password"))
	if err != nil {
		utils.LogError("", err)
	}

	addRegistryCmd.Flags().StringP("alias", "", "", "provide the alias for the database")
	err = viper.BindPFlag("alias", addRegistryCmd.Flags().Lookup("alias"))
	if err != nil {
		utils.LogError("", err)
	}

	addRegistryCmd.Flags().StringP("version", "", "latest", "provide the version of the database")
	err = viper.BindPFlag("version", addRegistryCmd.Flags().Lookup("version"))
	if err != nil {
		utils.LogError("", err)
	}

	command := make([]*cobra.Command, 0)
	command = append(command, addCmd)
	command = append(command, removeCmd)
	return command
}

// ActionAddRegistry adds a registry add on
func ActionAddRegistry(cmd *cobra.Command, args []string) error {
	project := viper.GetString("project")
	return addRegistry(project)
}

// ActionRemoveRegistry removes a registry add on
func ActionRemoveRegistry(cmd *cobra.Command, args []string) error {
	project := viper.GetString("project")
	return removeRegistry(project)
}

// ActionAddDatabase adds a database add on
func ActionAddDatabase(cmd *cobra.Command, args []string) error {
	dbtype := args
	if len(dbtype) == 0 {
		return utils.LogError("Database type not provided as an arguement", nil)
	}
	username := viper.GetString("username")
	if username == "" {
		switch dbtype[0] {
		case "postgres":
			username = "postgres"
		case "mysql":
			username = "root"
		}
	}
	password := viper.GetString("password")
	if password == "" {
		switch dbtype[0] {
		case "postgres":
			password = "mysecretpassword"
		case "mysql":
			password = "my-secret-pw"
		}
	}
	alias := viper.GetString("alias")

	version := viper.GetString("versio")

	return addDatabase(dbtype[0], username, password, alias, version)
}

// ActionRemoveDatabase removes a database add on
func ActionRemoveDatabase(cmd *cobra.Command, args []string) error {
	alias := args
	if len(alias) == 0 {
		return utils.LogError("Database Alias not provided as an argument", nil)
	}
	return removeDatabase(alias[0])
}
