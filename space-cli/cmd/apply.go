package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/sirupsen/logrus"

	"github.com/spaceuptech/space-cli/model"
)

// Apply reads the config file(s) from the provided file / directory and applies it to the server
func Apply() error {
	args := os.Args
	if len(args) != 3 {
		logrus.Errorf("error while applying service incorrect number of arguments provided")
		return fmt.Errorf("incorrect number of arguments provided")
	}

	fileName := args[2]
	if strings.HasSuffix(fileName, ".yaml") {
		err := ApplyWithFileName(fileName)
		if err != nil {
			return err
		}
	} else {
		err := actionApplyAll(fileName)
		if err != nil {
			return err
		}
	}

	return nil
}

func actionApplyAll(fileName string) error {
	files, err := ioutil.ReadDir("./" + fileName)
	if err != nil {
		return err
	}
	allYamlFiles := []string{}
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".yaml") {
			allYamlFiles = append(allYamlFiles, f.Name())

		}
	}

	orderList := []string{"1-project.yaml", "2-db-config.yaml", "3-db-rules.yaml", "4-db-schema.yaml", "5-filestore-config.yaml", "6-filestore-rule.yaml", "7-eventing-config.yaml", "8-eventing-triggers.yaml", "9-eventing-rule.yaml", "10-eventing-schema.yaml", "11-remote-services.yaml", "12-services.yaml", "13-services-routes.yaml", "14-services-secrets", "15--ingress-routes.yaml", "16-auth-providers.yaml", "17-letsencrypt.yaml"}
	for _, file := range orderList {
		for _, presentfile := range allYamlFiles {
			if strings.HasSuffix(presentfile, file) {
				log.Println("filename", presentfile, "file", file)
				err := ApplyWithFileName(fileName + "/" + presentfile)
				if err != nil {
					return err
				}
				break
			}
		}
	}
	return nil
}

// ApplyWithFileName does apply function by taking input as the name of file
func ApplyWithFileName(fileName string) error {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		logrus.Errorf("error while applying service unable to read file (%s) - %s", fileName, err.Error())
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

	dataStrings := strings.Split(string(data[:len(data)-4]), "---")
	for _, dataString := range dataStrings {
		fileContent := new(model.SpecObject)
		if err := yaml.Unmarshal([]byte(dataString), &fileContent); err != nil {
			logrus.Errorf("error while applying service unable to unmarshal file (%s) - %s", fileName, err.Error())
			return err
		}
		requestBody, err := json.Marshal(fileContent.Spec)
		if err != nil {
			logrus.Errorf("error while applying service unable to marshal spec - %s", err.Error())
			return err
		}
		url, err := adjustPath(fmt.Sprintf("%s%s", account.ServerURL, fileContent.API), fileContent.Meta)
		if err != nil {
			return err
		}

		log.Println("URL:", url)
		req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestBody))
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", login.Token))
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			logrus.Errorf("error while applying service unable to send http request - %s", err.Error())
			return err
		}
		defer CloseTheCloser(req.Body)

		v := map[string]interface{}{}
		_ = json.NewDecoder(resp.Body).Decode(&v)
		if resp.StatusCode != 200 {
			logrus.Errorf("error while applying service got http status code %s - %s", resp.Status, v["error"])
			return fmt.Errorf("%v", v["error"])
		}
		logrus.Infof("Successfully applied %s", fileContent.Type) // Why say service
	}

	return nil
}

func adjustPath(path string, meta map[string]string) (string, error) {
	newPath := path
	for {
		pre := strings.IndexRune(newPath, '{')
		if pre < 0 {
			return newPath, nil
		}
		post := strings.IndexRune(newPath, '}')

		key := strings.TrimSuffix(strings.TrimPrefix(newPath[pre:post], "{"), "}")
		value, p := meta[key]
		if !p {
			return "", fmt.Errorf("provided key (%s) does not exist in metadata", key)
		}

		newPath = newPath[:pre] + value + newPath[post+1:]
	}
}
