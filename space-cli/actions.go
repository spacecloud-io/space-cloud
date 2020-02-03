package main

import (
	"encoding/json"
	"fmt"

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
	service, loginResp, err := cmd.CodeStart()
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
	service, loginResp, err := cmd.CodeStart()
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

func actionGenerateService(c *cli.Context) error {
	service, err := cmd.GenerateServiceConfigWithoutLogin()
	if err != nil {
		return err
	}
	data, _ := json.MarshalIndent(service, "", " ")
	fmt.Println(string(data))
	return nil
}

func actionLogin(c *cli.Context) error {
	userName := c.String("username")
	key := c.String("key")
	url := c.String("url")
	return cmd.LoginStart(userName, key, url)
}

func actionSetup(c *cli.Context) error {
	id := c.String("id")
	userName := c.String("username")
	key := c.String("key")
	secret := c.String("secret")
	local := c.Bool("dev")
	return cmd.CodeSetup(id, userName, key, secret, local)
}
