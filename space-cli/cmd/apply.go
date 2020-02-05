package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/sirupsen/logrus"

	"github.com/spaceuptech/space-cli/model"
)

func Apply() error {
	args := os.Args
	if len(args) != 3 {
		return fmt.Errorf("incorrect number of arguments provided")
	}
	data, err := ioutil.ReadFile(args[2])
	if err != nil {
		return err
	}
	fileContent := new(model.GitOp)
	if err := yaml.Unmarshal(data, &fileContent); err != nil {
		return err
	}
	projectId := fileContent.Meta["projectId"]
	spec := fileContent.Spec.(map[string]interface{})
	spec["id"] = fileContent.Meta["id"]
	spec["projectId"] = projectId
	spec["version"] = fileContent.Meta["version"]

	requestBody, err := json.Marshal(fileContent.Spec)
	if err != nil {
		logrus.Error("error in apply service unable to marshal data - %v", err)
		return err
	}

	account, err := getSelectedAccount()
	if err != nil {
		return err
	}
	urlPathArr := strings.Split(fileContent.Api, "{")
	path := strings.Split(urlPathArr[1], "}")
	resp, err := http.Post(fmt.Sprintf("%s%s%s%s", account.ServerUrl, urlPathArr[0], projectId, path[1]), "application/json", bytes.NewBuffer(requestBody))
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
