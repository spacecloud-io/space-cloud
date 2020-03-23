package addons

import (
	"github.com/spf13/cobra"

	"github.com/spaceuptech/space-cli/utils"
)

// AddCmd is the list of commands the addon module exposes
var AddCmd = &cobra.Command{
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

// RemoveCmd is the list of commands the addon module exposes
var RemoveCmd = &cobra.Command{
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

func init() {
	AddCmd.AddCommand(addRegistryCmd)
	AddCmd.AddCommand(addDatabaseCmd)
	RemoveCmd.AddCommand(removeRegistryCmd)
	RemoveCmd.AddCommand(removeDatabaseCmd)

	addRegistryCmd.Flags().StringP("username", "U", "", "provide the username")
	addRegistryCmd.Flags().StringP("password", "P", "", "provide the password")
	addRegistryCmd.Flags().StringP("alias", "", "", "provide the alias for the database")
	addRegistryCmd.Flags().StringP("version", "", "latest", "provide the version of the database")
}

// Commands is the list of commands the addon module exposes
// var Commands = []cli.Command{
// 	{
// 		Name:  "add",
// 		Usage: "Add a add-on to the environment",
// 		Flags: []cli.Flag{cli.StringFlag{Name: "project", Usage: "The project to add the add-on to"}},
// 		Subcommands: []cli.Command{
// 			{
// 				Name:   "registry",
// 				Usage:  "Add a docker registry",
// 				Action: ActionAddRegistry,
// 			},
// 			{
// 				Name:  "database",
// 				Usage: "Add a database",
// 				Flags: []cli.Flag{
// 					cli.StringFlag{Name: "username, U", Usage: "provide the username"},
// 					cli.StringFlag{Name: "password, P", Usage: "provide the password"},
// 					cli.StringFlag{Name: "alias", Usage: "provide the alias for the database"},
// 					cli.StringFlag{Name: "version", Usage: "provide the version of the database", Value: "latest"},
// 				},
// 				Action: ActionAddDatabase,
// 			},
// 		},
// 	},
// 	{
// 		Name:  "remove",
// 		Usage: "Remove a add-on from the environment",
// 		Flags: []cli.Flag{cli.StringFlag{Name: "project", Usage: "The project to remove the add-on from"}},
// 		Subcommands: []cli.Command{
// 			{
// 				Name:   "registry",
// 				Usage:  "Remove a docker registry",
// 				Action: ActionRemoveRegistry,
// 			},
// 			{
// 				Name:   "database",
// 				Usage:  "Remove a database",
// 				Action: ActionRemoveDatabase,
// 			},
// 		},
// 	},
// }

// ActionAddRegistry adds a registry add on
func ActionAddRegistry(cmd *cobra.Command, args []string) error {
	project := c.GlobalString("project")
	return addRegistry(project)
}

// ActionRemoveRegistry removes a registry add on
func ActionRemoveRegistry(cmd *cobra.Command, args []string) error {
	project := c.GlobalString("project")
	return removeRegistry(project)
}

// ActionAddDatabase adds a database add on
func ActionAddDatabase(cmd *cobra.Command, args []string) error {
	dbtype := c.Args().Get(0)
	if len(dbtype) == 0 {
		return utils.LogError("Database type not provided as an arguement", nil)
	}
	username := c.String("username")
	if username == "" {
		switch dbtype {
		case "postgres":
			username = "postgres"
		case "mysql":
			username = "root"
		}
	}
	password := c.String("password")
	if password == "" {
		switch dbtype {
		case "postgres":
			password = "mysecretpassword"
		case "mysql":
			password = "my-secret-pw"
		}
	}
	alias := c.String("alias")
	version := c.String("version")
	return addDatabase(dbtype, username, password, alias, version)
}

// ActionRemoveDatabase removes a database add on
func ActionRemoveDatabase(cmd *cobra.Command, args []string) error {
	alias := c.Args().Get(0)
	if len(alias) == 0 {
		return utils.LogError("Database Alias not provided as an argument", nil)
	}
	return removeDatabase(alias)
}
