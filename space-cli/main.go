package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
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
			},
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
			Subcommands: []cli.Command{
				{
					Name:   "global-config",
					Action: actionGetGlobalConfig,
				},
				{
					Name:   "remote-services",
					Action: actionGetRemoteServices,
				},
				{
					Name:   "auth-providers",
					Action: actionGetAuthProviders,
				},
				{
					Name:   "eventing-triggers",
					Action: actionGetEventingTrigger,
				},
				{
					Name:   "eventing-config",
					Action: actionGetEventingConfig,
				},
				{
					Name:   "eventing-schema",
					Action: actionGetEventingSchema,
				},
				{
					Name:   "eventing-rule",
					Action: actionGetEventingSecurityRule,
				},
				{
					Name:   "filestore-config",
					Action: actionGetFileStoreConfig,
				},
				{
					Name:   "filestore-rules",
					Action: actionGetFileStoreRule,
				},
				{
					Name:   "db-rule",
					Action: actionGetDbRule,
				},
				{
					Name:   "db-config",
					Action: actionGetDbConfig,
				},
				{
					Name:   "db-schema",
					Action: actionGetDbSchema,
				},
				{
					Name:   "letsencrypt-domain",
					Action: actionGetLetsEncryptDomain,
				},
				{
					Name:   "routes",
					Action: actionGetRoutes,
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
			},
			Action: actionSetup,
		},
	}

	// Start the app
	if err := app.Run(os.Args); err != nil {
		logrus.Fatalln("Failed to run execute command:", err)
	}
}
