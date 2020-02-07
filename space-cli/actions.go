package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"

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

func actionApply(c *cli.Context) error {
	return cmd.Apply()
}

func actionGenerateService(c *cli.Context) error {
	// get filename from args in which service config will be stored
	argsArr := os.Args
	if len(argsArr) != 4 {
		return fmt.Errorf("incorrect number of arguments")
	}
	serviceConfigFile := argsArr[3]

	service, err := cmd.GenerateService()
	if err != nil {
		return err
	}
	v := model.GitOp{
		Api:  "/v1/runner/{projectId}/services",
		Type: "service",
		Meta: map[string]string{
			"id":        service.ID,
			"projectId": service.ProjectID,
			"version":   "v1",
		},
	}
	service.ID = ""
	service.ProjectID = ""
	service.Version = ""
	v.Spec = service

	data, err := yaml.Marshal(v)
	if err != nil {
		logrus.Errorf("error pretty printing service struct - %s", err.Error())
		return err
	}

	if err := ioutil.WriteFile(serviceConfigFile, data, 0755); err != nil {
		return err
	}
	fmt.Printf(string(data))
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
