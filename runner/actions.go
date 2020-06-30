package main

import (
	"os"
	"strings"

	"github.com/spaceuptech/space-cloud/runner/model"
	"github.com/spaceuptech/space-cloud/runner/utils/auth"
	"github.com/spaceuptech/space-cloud/runner/utils/driver"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/spaceuptech/space-cloud/runner/server"
)

func actionRunner(c *cli.Context) error {
	// Get runner config flags
	port := c.String("port")
	proxyPort := c.String("proxy-port")
	loglevel := c.String("log-level")

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
	clusterID := os.Getenv("CLUSTER_ID")
	if clusterID == "" {
		logrus.Error("Failed to setup runner: CLUSTER_ID environment variable not provided")
		return nil
	}
	clusterName := strings.Split(clusterID, "--")[0]
	// Set the log level
	setLogLevel(loglevel)

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
		logrus.Errorf("Failed to start runner - %s", err.Error())
		os.Exit(-1)
	}

	return r.Start()
}

func setLogLevel(loglevel string) {
	switch loglevel {
	case loglevelDebug:
		logrus.SetLevel(logrus.DebugLevel)
	case loglevelInfo:
		logrus.SetLevel(logrus.InfoLevel)
	case logLevelError:
		logrus.SetLevel(logrus.ErrorLevel)
	default:
		logrus.Errorf("Invalid log level (%s) provided", loglevel)
		logrus.Infoln("Defaulting to `info` level")
		logrus.SetLevel(logrus.InfoLevel)
	}
}
