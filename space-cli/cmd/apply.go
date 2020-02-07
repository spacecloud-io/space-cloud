package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/spaceuptech/space-cli/model"
)

func Apply() error {
	args := os.Args
	if len(args) != 3 {
		logrus.Errorf("error while applying service incorrect number of arguments provided")
		return fmt.Errorf("incorrect number of arguments provided")
	}
	fileName := args[2]
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		logrus.Errorf("error while applying service unable to read file (%s) - %s", fileName, err.Error())
		return err
	}
	fileContent := new(model.GitOp)
	if err := yaml.Unmarshal(data, &fileContent); err != nil {
		logrus.Errorf("error while applying service unable to unmarshal file (%s) - %s", fileName, err.Error())
		return err
	}
	projectId := fileContent.Meta["projectId"]
	spec := fileContent.Spec.(map[string]interface{})
	spec["id"] = fileContent.Meta["id"]
	spec["projectId"] = projectId
	spec["version"] = fileContent.Meta["version"]
	requestBody, err := json.Marshal(fileContent.Spec)
	if err != nil {
		logrus.Errorf("error while applying service unable to marshal spec - %s", err.Error())
		return err
	}

	account, err := getSelectedAccount()
	if err != nil {
		logrus.Errorf("error while applying service unable to get selected account - %s", err.Error())
		return err
	}
	login, err := login(account)
	if err != nil {
		logrus.Errorf("error while applying service unable to login - %s", err.Error())
		return err
	}
	//urlPathArr := strings.Split(fileContent.Api, "{")
	//path := strings.Split(urlPathArr[1], "}")
	//logrus.Print("path:",fmt.Sprintf("%s%s%s%s", account.ServerUrl, fileContent.Api))
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s%s", account.ServerUrl, fileContent.Api), bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", login.Token))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logrus.Errorf("error while applying service unable to send http request - %s", err.Error())
		return err
	}
	defer resp.Body.Close()
	v := map[string]interface{}{}
	_ = json.NewDecoder(resp.Body).Decode(&v)
	if resp.StatusCode != 200 {
		logrus.Errorf("error while applying service got http status code %s - %s", resp.Status, v["error"])
		return fmt.Errorf("%v", v["error"])
	}
	logrus.Infof("Successfully applied %s", fileContent.Type) // Why say service
	return nil
}
