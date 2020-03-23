package addons

import (
	"github.com/urfave/cli"

	"github.com/spaceuptech/space-cli/utils"
)

// AddSubCommands is the list of commands the addon module exposes
var AddSubCommands = []cli.Command{
	{
		Name:   "registry",
		Usage:  "Add a docker registry",
		Action: ActionAddRegistry,
	},
	{
		Name:  "database",
		Usage: "Add a database",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "username, U", Usage: "provide the username"},
			cli.StringFlag{Name: "password, P", Usage: "provide the password"},
			cli.StringFlag{Name: "alias", Usage: "provide the alias for the database"},
			cli.StringFlag{Name: "version", Usage: "provide the version of the database", Value: "latest"},
		},
		Action: ActionAddDatabase,
	},
}

// RemoveSubCommand is the list of commands the addon module exposes
var RemoveSubCommand = []cli.Command{
	{
		Name:   "registry",
		Usage:  "Remove a docker registry",
		Action: ActionRemoveRegistry,
	},
	{
		Name:   "database",
		Usage:  "Remove a database",
		Action: ActionRemoveDatabase,
	},
}

// ActionAddRegistry adds a registry add on
func ActionAddRegistry(c *cli.Context) error {
	project := c.GlobalString("project")
	return addRegistry(project)
}

// ActionRemoveRegistry removes a registry add on
func ActionRemoveRegistry(c *cli.Context) error {
	project := c.GlobalString("project")
	return removeRegistry(project)
}

// ActionAddDatabase adds a database add on
func ActionAddDatabase(c *cli.Context) error {
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
func ActionRemoveDatabase(c *cli.Context) error {
	alias := c.Args().Get(0)
	if len(alias) == 0 {
		return utils.LogError("Database Alias not provided as an argument", nil)
	}
	return removeDatabase(alias)
}
