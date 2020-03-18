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
	if len(data) == 0 {
		logrus.Infoln("empty file provided")
		return nil
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
