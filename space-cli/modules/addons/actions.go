package addons

import (
	"fmt"

	"github.com/urfave/cli"
)

// Commands is the list of commands the addon module exposes
var Commands = []cli.Command{
	{
		Name:  "add",
		Usage: "Add a add-on to the environment",
		Flags: []cli.Flag{cli.StringFlag{Name: "project", Usage: "The project to add the add-on to"}},
		Subcommands: []cli.Command{
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
					cli.StringFlag{Name: "version", Usage: "provide the version of the database"},
				},
				Action: ActionAddDatabase,
			},
		},
	},
	{
		Name:  "remove",
		Usage: "Remove a add-on from the environment",
		Flags: []cli.Flag{cli.StringFlag{Name: "project", Usage: "The project to remove the add-on from"}},
		Subcommands: []cli.Command{
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
		},
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
		return fmt.Errorf("Database type not provided as an arguement")
	}
	username := c.GlobalString("username")
	if username == "" {
		switch dbtype {
		case "mongo":
			username = "mongodb"
		case "postgres":
			username = "postgres"
		case "mysql":
			username = "root"
		}
	}
	password := c.GlobalString("password")
	if password == "" {
		switch dbtype {
		case "mongo":
			password = ""
		case "postgres":
			password = "mysecretpassword"
		case "mysql":
			password = "my-secret-pw"
		}
	}
	alias := c.GlobalString("alias")
	version := c.GlobalString("version")
	return addDatabase(dbtype, username, password, alias, version)
}

// ActionRemoveDatabase removes a database add on
func ActionRemoveDatabase(c *cli.Context) error {
	alias := c.Args().Get(0)
	if len(alias) == 0 {
		return fmt.Errorf("Database Alias not provided as an arguement")
	}
	return removeDatabase(alias)
}
