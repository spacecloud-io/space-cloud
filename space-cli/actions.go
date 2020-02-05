package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

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
	args := os.Args
	if len(args) != 3 {
		return fmt.Errorf("incorrect number of arguments provided")
	}
	data, err := ioutil.ReadFile(args[2])
	if err != nil {
		return err
	}
	fileContent := new(model.GitOp)
	if err := json.Unmarshal(data, fileContent); err != nil {
		return err
	}
	if fileContent.Type != "service" {
		return fmt.Errorf("invalid type found in serivce file type should be (service)")
	}
	service := fileContent.Spec.(*model.Service)
	service.ID = fileContent.Meta["id"]
	service.ProjectID = fileContent.Meta["projectId"]
	service.Version = fileContent.Meta["version"]

	requestBody, err := json.Marshal(&fileContent.Spec)
	if err != nil {
		logrus.Error("error in apply service unable to marshal data - %v", err)
		return err
	}
	urlPath := strings.Replace(fileContent.Api, "{projectId}", service.ProjectID, 1)
	resp, err := http.Post(fmt.Sprintf("http://localhost:4122%s", urlPath), "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		logrus.Error("error in apply service unable to send http request - %v", err)
		return err
	}
	defer resp.Body.Close()
	v := map[string]interface{}{}
	_ = json.NewDecoder(resp.Body).Decode(&v)
	if resp.StatusCode != 200 {
		return fmt.Errorf("%v", v["error"])
	}
	return nil
}

func actionGenerateService(c *cli.Context) error {
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
