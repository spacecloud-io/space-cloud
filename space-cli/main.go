package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	eventing "github.com/spaceuptech/space-cli/modules/Eventing"
	"github.com/spaceuptech/space-cli/modules/database"
	"github.com/spaceuptech/space-cli/modules/filestore"
	"github.com/spaceuptech/space-cli/modules/userman"
)

func main() {

	// Setup logrus
	logrus.SetFormatter(&logrus.TextFormatter{})
	logrus.SetOutput(os.Stdout)

	app := cli.NewApp()
	app.EnableBashCompletion = true
	app.Name = "space-cli"
	app.Version = "0.16.0"
	app.Commands = []cli.Command{
		{
			Name:  "generate",
			Usage: "generates service config",
			Subcommands: []cli.Command{
				{
					Name:   "service",
					Action: actionGenerateService,
				},
				{
					Name:   "db-rule",
					Action: database.ActionGenerateDBRule,
				},
				{
					Name:   "db-config",
					Action: database.ActionGenerateDBConfig,
				},
				{
					Name:   "db-schema",
					Action: database.ActionGenerateDBSchema,
				},
				{
					Name:   "filestore-rules",
					Action: filestore.ActionGenerateFilestoreRule,
				},
				{
					Name:   "filestore-config",
					Action: filestore.ActionGenerateFilestoreConfig,
				},
				{
					Name:   "eventing-rule",
					Action: eventing.ActionGenerateEventingRule,
				},
				{
					Name:   "eventing-schema",
					Action: eventing.ActionGenerateEventingSchema,
				},
				{
					Name:   "eventing-config",
					Action: eventing.ActionGenerateEventingConfig,
				},
				{
					Name:   "eventing-trigger",
					Action: eventing.ActionGenerateEventingTrigger,
				},
				{
					Name:   "user-management",
					Action: userman.ActionGenerateUserManagement,
				},
			},
		},
		{
			Name:   "apply",
			Usage:  "deploys service",
			Action: actionApply,
		},
		{
			Name:   "destroy",
			Usage:  "clean development environment & remove secrets",
			Action: actionDestroy,
		},
		{
			Name:  "login",
			Usage: "Logs into space cloud",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "username",
					Usage:  "Accepts the username for login",
					EnvVar: "USER_NAME", // don't set environment variable as USERNAME -> defaults to username of host machine in linux
					Value:  "None",
				},
				cli.StringFlag{
					Name:   "key",
					Usage:  "Accepts the access key to be verified during login",
					EnvVar: "KEY",
					Value:  "None",
				},
				cli.StringFlag{
					Name:   "url",
					Usage:  "Accepts the URL of server",
					EnvVar: "URL",
					Value:  "http://localhost:4122",
				},
			},
			Action: actionLogin,
		},
		{
			Name:  "setup",
			Usage: "setup development environment",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "id",
					Usage:  "The unique id for the cluster",
					EnvVar: "CLUSTER_ID",
					Value:  "",
				},
				cli.StringFlag{
					Name:   "username",
					Usage:  "The username used for login",
					EnvVar: "USER_NAME", // don't set environment variable as USERNAME -> defaults to username of host machine in linux
					Value:  "",
				},
				cli.StringFlag{
					Name:   "key",
					Usage:  "The access key used for login",
					EnvVar: "KEY",
					Value:  "",
				},
				cli.StringFlag{
					Name:   "secret",
					Usage:  "The jwt secret to start space-cloud with",
					EnvVar: "JWT_SECRET",
					Value:  "",
				},
				cli.BoolFlag{
					Name:  "dev",
					Usage: "Run space cloud in development mode",
				},
				cli.Int64Flag{
					Name:   "port-http",
					Usage:  "The port to use for HTTP",
					EnvVar: "PORT_HTTP",
					Value:  4122,
				},
				cli.Int64Flag{
					Name:   "port-https",
					Usage:  "The port to use for HTTPS",
					EnvVar: "PORT_HTTPS",
					Value:  4126,
				},
				cli.StringSliceFlag{
					Name:  "v",
					Usage: "Volumes to be attached to gateway",
				},
				cli.StringSliceFlag{
					Name:  "e",
					Usage: "Environment variables to be provided to gateway",
				},
			},
			Action: actionSetup,
		},
	}

	// Start the app
	if err := app.Run(os.Args); err != nil {
		logrus.Fatalln("Failed to start space cli:", err)
	}
}
