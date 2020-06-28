package main

import (
	"errors"
	"os"
	"strings"

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
	app.Name = "metric-proxy"
	app.Version = "0.2.0"

	app.Commands = []cli.Command{
		{
			Name:  "start",
			Usage: "Starts the proxy to collect metrics directly from envoy",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "addr",
					Usage:  "Address of the runner instance",
					EnvVar: "ADDR",
					Value:  "runner.space-cloud.svc.cluster.local:4050",
				},
				cli.StringFlag{
					Name:   "token",
					Usage:  "The token to be used for authentication",
					EnvVar: "TOKEN",
				},
				cli.StringFlag{
					Name:   "log-level",
					EnvVar: "LOG_LEVEL",
					Usage:  "Set the log level [debug | info | error]",
					Value:  loglevelInfo,
				},
				cli.StringFlag{
					Name:   "mode",
					EnvVar: "MODE",
					Usage:  "The collection mode [parallel | per-second]",
					Value:  "per-second",
				},
			},
			Action: func(c *cli.Context) error {
				// Get all flags
				addr := c.String("addr")
				token := c.String("token")
				loglevel := c.String("log-level")
				mode := c.String("mode")

				// Set the log level
				setLogLevel(loglevel)

				// Throw an error if invalid token provided
				if len(strings.Split(token, ".")) != 3 {
					return errors.New("invalid token provided")
				}

				// Start the proxy
				p := New(addr, token, mode)
				return p.Start()
			},
		},
	}

	// Start the app
	if err := app.Run(os.Args); err != nil {
		logrus.Fatalln("Failed to start runner:", err)
	}
}
