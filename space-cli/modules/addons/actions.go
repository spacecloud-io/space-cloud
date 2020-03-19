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
