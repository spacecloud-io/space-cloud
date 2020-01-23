package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

const (
	loglevelDebug = "debug"
	loglevelInfo  = "info"
	logLevelError = "error"
)

func main() {

	// Setup logrus
	logrus.SetFormatter(&logrus.TextFormatter{})
	logrus.SetOutput(os.Stdout)

	app := cli.NewApp()
	app.Name = "galaxy"
	app.Version = "0.1.0"
	app.Commands = []cli.Command{
		{
			Name:  "code",
			Usage: "Commands to work with non dockerized code",
			Subcommands: []cli.Command{
				{
					Name: "start",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:   "env",
							Usage:  "Builds and deploys a codebase",
							EnvVar: "ENV",
							Value:  "none",
						},
					},
					Action: actionStartCode,
				},
				{
					Name: "build",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:   "env",
							Usage:  "Builds a codebase",
							EnvVar: "ENV",
							Value:  "none",
						},
					},
					Action: actionBuildCode,
				},
			},
		},
		{
			Name:  "login",
			Usage: "Commands to log in",
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
					Value:  "localhost:4122",
				},
				cli.BoolFlag{
					Name:   "local",
					Usage:  "Determines whether local URL is to be used as server URL",
					EnvVar: "LOCAL",
				},
			},
			Action: actionLogin,
		},
		{
			Name:  "setup",
			Usage: "setup development environment",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "username",
					Usage:  "Accepts the username for login",
					EnvVar: "USER_NAME", // don't set environment variable as USERNAME -> defaults to username of host machine in linux
					Value:  "",
				},
				cli.StringFlag{
					Name:   "key",
					Usage:  "Accepts the access key to be verified during login",
					EnvVar: "KEY",
					Value:  "",
				},
				cli.StringFlag{
					Name:   "url",
					Usage:  "Accepts the URL of server",
					EnvVar: "URL",
					Value:  "localhost:4122",
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
		logrus.Fatalln("Failed to start galaxy:", err)
	}
}
