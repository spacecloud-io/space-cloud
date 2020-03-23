package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/spaceuptech/space-cli/modules"
	"github.com/spaceuptech/space-cli/modules/addons"
	"github.com/spaceuptech/space-cli/modules/auth"
	"github.com/spaceuptech/space-cli/modules/database"
	"github.com/spaceuptech/space-cli/modules/deploy"
	"github.com/spaceuptech/space-cli/modules/eventing"
	"github.com/spaceuptech/space-cli/modules/filestore"
	"github.com/spaceuptech/space-cli/modules/ingress"
	"github.com/spaceuptech/space-cli/modules/letsencrypt"
	"github.com/spaceuptech/space-cli/modules/operations"
	"github.com/spaceuptech/space-cli/modules/project"
	remoteservices "github.com/spaceuptech/space-cli/modules/remote-services"
	"github.com/spaceuptech/space-cli/modules/services"
	"github.com/spaceuptech/space-cli/modules/userman"
	"github.com/spaceuptech/space-cli/utils"
)

func main() {

	// Setup logrus
	logrus.SetFormatter(&logrus.TextFormatter{})
	logrus.SetOutput(os.Stdout)

	app := cli.NewApp()
	app.EnableBashCompletion = true
	app.Name = "space-cli"
	app.Version = "0.16.0"
	app.Flags = []cli.Flag{cli.StringFlag{Name: "log-level", Value: "info", Usage: "Sets the log level of the command", EnvVar: "LOG_LEVEL"}}
	app.Commands = []cli.Command{
		{
			Name:        "add",
			Usage:       "Add a add-on to the environment",
			Subcommands: fetchAddSubCommands(),
		},
		{
			Name:        "remove",
			Usage:       "Remove a add-on from the environment",
			Subcommands: fetchRemoveSubCommands(),
		},
		{
			Name:  "get",
			Usage: "gets different services",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "project",
					Usage:  "The id of the project",
					EnvVar: "PROJECT_ID",
				},
			},
			Subcommands: fetchGetSubCommands(),
		},
		{
			Name:        "generate",
			Usage:       "generates service config",
			Subcommands: fetchGenerateSubCommands(),
		},
	}
	app.Commands = append(app.Commands, deploy.CommandDeploy)
	app.Commands = append(app.Commands, operations.Commands...)
	app.Commands = append(app.Commands, utils.LoginCommands...)

	// Start the app
	if err := app.Run(os.Args); err != nil {
		logrus.Fatalln("Failed to run execute command:", err)
	}
}

func fetchAddSubCommands() []cli.Command {
	v := []cli.Command{}
	v = append(v, addons.AddSubCommands...)
	return v
}

func fetchRemoveSubCommands() []cli.Command {
	v := []cli.Command{}
	v = append(v, addons.RemoveSubCommand...)
	return v
}

func fetchGetSubCommands() []cli.Command {
	v := []cli.Command{}
	v = append(v, auth.GetSubCommands...)
	v = append(v, database.GetSubCommands...)
	v = append(v, eventing.GetSubCommands...)
	v = append(v, filestore.GetSubCommands...)
	v = append(v, letsencrypt.GetSubCommands...)
	v = append(v, project.GetSubCommands...)
	v = append(v, remoteservices.GetSubCommands...)
	v = append(v, services.GetSubCommands...)
	v = append(v, modules.GetSubCommands...)

	return v
}

func fetchGenerateSubCommands() []cli.Command {
	v := []cli.Command{}
	v = append(v, database.GenerateSubCommands...)
	v = append(v, eventing.GenerateSubCommands...)
	v = append(v, filestore.GenerateSubCommands...)
	v = append(v, ingress.GenerateSubCommands...)
	v = append(v, letsencrypt.GenerateSubCommands...)
	v = append(v, remoteservices.GenerateSubCommands...)
	v = append(v, services.GenerateSubCommands...)
	v = append(v, userman.GenerateSubCommands...)

	return v
}
