package main

import (
	"context"
	"fmt"
	"os"

	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/runner/model"
	"github.com/spaceuptech/space-cloud/runner/utils/auth"
	"github.com/spaceuptech/space-cloud/runner/utils/driver"

	"github.com/urfave/cli"

	"github.com/spaceuptech/space-cloud/runner/server"
)

func actionRunner(c *cli.Context) error {
	// Get runner config flags
	port := c.String("port")
	proxyPort := c.String("proxy-port")
	loglevel := c.String("log-level")
	logFormat := c.String("log-format")

	// Get jwt config
	jwtSecret := c.String("jwt-secret")
	jwtProxySecret := c.String("jwt-proxy-secret")

	// Get driver config
	driverType := c.String("driver")
	driverConfig := c.String("driver-config")
	outsideCluster := c.Bool("outside-cluster")

	isDev := c.Bool("dev")
	isMetricDisabled := c.Bool("disable-metrics")

	artifactAddr := c.String("artifact-addr")
	clusterName := c.String("cluster-name")
	if driverType == model.DockerType {
		helpers.Logger.LogInfo(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Runner is starting in cluster (%s)", clusterName), nil)
	}
	// Set the log level
	if err := helpers.InitLogger(loglevel, logFormat, isDev); err != nil {
		return helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to initialize loggers", err, nil)
	}

	// Create a new runner object
	r, err := server.New(&server.Config{
		Port:             port,
		ProxyPort:        proxyPort,
		IsMetricDisabled: isMetricDisabled,
		Auth: &auth.Config{
			Secret:      jwtSecret,
			ProxySecret: jwtProxySecret,
			IsDev:       isDev,
		},
		Driver: &driver.Config{
			DriverType:     model.DriverType(driverType),
			ConfigFilePath: driverConfig,
			IsInCluster:    !outsideCluster,
			ArtifactAddr:   artifactAddr,
			ClusterName:    clusterName,
		},
	})
	if err != nil {
		_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Failed to start runner", err, nil)
		os.Exit(-1)
	}

	return r.Start()
}
