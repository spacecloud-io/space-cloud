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
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "log-level", Value: "info", Usage: "Sets the log level of the command", EnvVar: "LOG_LEVEL"},
	}
	app.Commands = append(app.Commands, addons.Commands...)
	app.Commands = append(app.Commands, auth.Commands...)
	app.Commands = append(app.Commands, database.Commands...)
	app.Commands = append(app.Commands, deploy.CommandDeploy)
	app.Commands = append(app.Commands, eventing.Commands...)
	app.Commands = append(app.Commands, filestore.Commands...)
	app.Commands = append(app.Commands, ingress.Commands...)
	app.Commands = append(app.Commands, letsencrypt.Commands...)
	app.Commands = append(app.Commands, operations.Commands...)
	app.Commands = append(app.Commands, project.Commands...)
	app.Commands = append(app.Commands, remoteservices.Commands...)
	app.Commands = append(app.Commands, services.Commands...)
	app.Commands = append(app.Commands, userman.Commands...)
	app.Commands = append(app.Commands, modules.Commands...)
	app.Commands = append(app.Commands, utils.LoginCommands...)

	// Start the app
	if err := app.Run(os.Args); err != nil {
		logrus.Fatalln("Failed to run execute command:", err)
	}
}
