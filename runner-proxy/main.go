package main

import (
	"context"
	"fmt"
	"os"

	"github.com/spaceuptech/helpers"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "runner-proxy"

	app.Commands = []cli.Command{
		{
			Name:  "start",
			Usage: "Starts a proxy-runner instance",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "port",
					Usage: "The port the runner will bind too",
					Value: "4055",
				},
				cli.StringFlag{
					Name:   "log-level",
					EnvVar: "LOG_LEVEL",
					Usage:  "Set the log level [debug | info | error]",
					Value:  helpers.LogLevelInfo,
				},
				cli.StringFlag{
					Name:   "log-format",
					EnvVar: "LOG_FORMAT",
					Usage:  "Set the log format [json | console]",
					Value:  helpers.LogFormatJSON,
				},
				cli.StringFlag{
					Name:   "admin-secret",
					Usage:  "Set the admin secret",
					EnvVar: "ADMIN_SECRET",
					Value:  "some-secret",
				},
			},
			Action: actionRunner,
		},
	}

	// Start the app
	if err := app.Run(os.Args); err != nil {
		helpers.Logger.LogFatal(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Failed to start runner-proxy: %v", err), nil)
	}
}
