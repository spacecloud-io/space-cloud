package main

import (
	"context"
	"os"

	"github.com/spaceuptech/helpers"

	"github.com/urfave/cli"

	"github.com/spaceuptech/space-cloud/runner-proxy/server"
)

func actionRunner(c *cli.Context) error {
	// Get proxy-runner config flags
	port := c.String("port")
	loglevel := c.String("log-level")
	logFormat := c.String("log-format")
	secret := c.String("admin-secret")

	// Set the log level
	if err := helpers.InitLogger(loglevel, logFormat, false); err != nil {
		return helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to initialize loggers", err, nil)
	}

	// Create a new runner object
	r, err := server.New(secret)
	if err != nil {
		_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Failed to start runner", err, nil)
		os.Exit(-1)
	}

	return r.Start(port)
}
