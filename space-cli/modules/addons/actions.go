package addons

import (
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
					cli.StringFlag{Name: "alias, A", Usage: "provide the alias for the database"},
					cli.StringFlag{Name: "version, V", Usage: "provide the version of the database"},
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
	username := c.GlobalString("username")
	password := c.GlobalString("password")
	alias := c.GlobalString("alias")
	version := c.GlobalString("version")
	return addDatabase(dbtype, username, password, alias, version)
}
