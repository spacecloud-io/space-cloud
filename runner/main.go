package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/spaceuptech/space-cloud/runner/model"
	"github.com/spaceuptech/space-cloud/runner/modules/routing"
	"github.com/spaceuptech/space-cloud/runner/modules/secrets"
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
	app.Name = "runner"
	app.Version = model.Version

	app.Commands = []cli.Command{
		{
			Name:  "generate",
			Usage: "generates service config",
			Subcommands: []cli.Command{
				{
					Name:   "service-routing",
					Action: routing.ActionGenerateServiceRouting,
				},
				{
					Name:   "apply-secrets",
					Action: secrets.ActionGenerateSecret,
				},
			},
		},
		{
			Name:  "start",
			Usage: "Starts a runner instance",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:   "disable-metrics",
					EnvVar: "DISABLE_METRICS",
					Usage:  "Disable anonymous metric collection",
				},
				cli.BoolFlag{
					Name:   "dev",
					EnvVar: "DEV",
					Usage:  "start runner without authentication",
				},
				cli.StringFlag{
					Name:   "port",
					EnvVar: "PORT",
					Usage:  "The port the runner will bind too",
					Value:  "4050",
				},
				cli.StringFlag{
					Name:   "proxy-port",
					EnvVar: "PROXY_PORT",
					Usage:  "The port the proxy will bind too",
					Value:  "4055",
				},
				cli.StringFlag{
					Name:   "log-level",
					EnvVar: "LOG_LEVEL",
					Usage:  "Set the log level [debug | info | error]",
					Value:  loglevelInfo,
				},

				// JWT config
				cli.StringFlag{
					Name:   "jwt-secret",
					EnvVar: "JWT_SECRET",
					Usage:  "The jwt secret to use when the algorithm is set to HS256",
					Value:  "some-secret",
				},
				cli.StringFlag{
					Name:   "jwt-proxy-secret",
					EnvVar: "JWT_PROXY_SECRET",
					Usage:  "The jwt secret to use for authenticating the proxy",
					Value:  "some-proxy-secret",
				},

				// Driver config
				cli.StringFlag{
					Name:   "driver",
					EnvVar: "DRIVER",
					Usage:  "The driver to use for deployment",
					Value:  "istio",
				},
				cli.StringFlag{
					Name:   "driver-config",
					EnvVar: "DRIVER_CONFIG",
					Usage:  "Driver config file path",
				},
				cli.StringFlag{
					Name:   "artifact-addr",
					EnvVar: "ARTIFACT_ADDR",
					Usage:  "The address used to reach the artifact serverl",
					Value:  "http://store.space-cloud.svc.cluster.local",
				},
				cli.BoolFlag{
					Name:   "outside-cluster",
					EnvVar: "OUTSIDE_CLUSTER",
					Usage:  "Indicates whether runner in running inside the cluster",
				},
			},
			Action: actionRunner,
		},
	}

	// Start the app
	if err := app.Run(os.Args); err != nil {
		logrus.Fatalln("Failed to start runner:", err)
	}
}
