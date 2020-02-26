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

func actionApply(c *cli.Context) error {
	return cmd.Apply()
}

func actionDestroy(c *cli.Context) error {
	return cmd.Destroy()
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
	fmt.Printf("%s", string(data))
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
	portHTTP := c.Int64("port-http")
	portHTTPS := c.Int64("port-https")
	volumes := c.StringSlice("v")
	environmentVariables := c.StringSlice("e")

	return cmd.CodeSetup(id, userName, key, secret, local, portHTTP, portHTTPS, volumes, environmentVariables)
}
