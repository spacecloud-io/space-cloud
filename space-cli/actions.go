package main

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/spaceuptech/space-cli/cmd"
	"github.com/spaceuptech/space-cli/utils"
)

func actionDestroy(_ *cli.Context) error {
	return cmd.Destroy()
}

func actionLogin(c *cli.Context) error {
	userName := c.String("username")
	key := c.String("key")
	url := c.String("url")

	return utils.LoginStart(userName, key, url)
}

func actionSetup(c *cli.Context) error {
	id := c.String("id")
	userName := c.String("username")
	key := c.String("key")
	config := c.String("config")
	version := c.String("version")
	secret := c.String("secret")
	local := c.Bool("dev")
	portHTTP := c.Int64("port-http")
	portHTTPS := c.Int64("port-https")
	volumes := c.StringSlice("v")
	environmentVariables := c.StringSlice("e")

	setLogLevel(c.GlobalString("log-level"))

	return cmd.CodeSetup(id, userName, key, config, version, secret, local, portHTTP, portHTTPS, volumes, environmentVariables)
}

func actionUpgrade(_ *cli.Context) error {
	return cmd.Upgrade()
}

func setLogLevel(loglevel string) {
	switch loglevel {
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	default:
		logrus.Errorf("Invalid log level (%s) provided", loglevel)
		logrus.Infoln("Defaulting to `info` level")
		logrus.SetLevel(logrus.InfoLevel)
	}
}
