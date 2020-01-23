package main

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/spaceuptech/space-cli/cmd"
	"github.com/spaceuptech/space-cli/model"
)

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

func actionStartCode(c *cli.Context) error {
	envID := c.String("env")
	service, loginResp, err := cmd.CodeStart(envID)
	if err != nil {
		return err
	}
	actionCodeStruct := &model.ActionCode{
		Service:  service,
		IsDeploy: true, //
	}
	if err := cmd.RunDockerFile(actionCodeStruct, loginResp); err != nil {
		return err
	}
	return nil
}

func actionBuildCode(c *cli.Context) error {
	envID := c.String("env")
	service, loginResp, err := cmd.CodeStart(envID)
	if err != nil {
		return err
	}
	actionCodeStruct := &model.ActionCode{
		Service:  service,
		IsDeploy: false,
	}
	if err := cmd.RunDockerFile(actionCodeStruct, loginResp); err != nil {
		return err
	}
	return nil
}

func actionLogin(c *cli.Context) error {
	userName := c.String("username")
	key := c.String("key")
	serverUrl := c.String("url")
	local := c.Bool("local")
	url := ""
	if local {
		url = "localhost:4122"
	}
	if serverUrl != "default url" { // todo get default url
		url = serverUrl
	}
	return cmd.LoginStart(userName, key, url, local)
}

func actionSetup(c *cli.Context) error {
	userName := c.String("username")
	key := c.String("key")
	serverUrl := c.String("url")
	local := c.Bool("dev")
	return cmd.CodeSetup(userName, key, serverUrl, local)
}
