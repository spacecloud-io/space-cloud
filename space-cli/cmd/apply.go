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

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/spaceuptech/space-cli/model"
	"github.com/spaceuptech/space-cli/utils"
)

// ActionApply applies the file / folder
func ActionApply(cli *cli.Context) error {
	args := os.Args
	if len(args) != 3 {
		logrus.Errorf("error while applying service incorrect number of arguments provided")
		return fmt.Errorf("incorrect number of arguments provided")
	}

	fileName := args[2]

	return applyAll(fileName)
}

func applyAll(fileName string) error {
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
	account, err := utils.GetSelectedAccount()
	if err != nil {
		return utils.LogError("Unable to fetch account information", "apply", "", err)
	}
	login, err := utils.Login(account)
	if err != nil {
		return utils.LogError("Unable to login", "apply", "", err)
	}

	specs, err := utils.ReadSpecObjectsFromFile(fileName)
	if err != nil {
		return utils.LogError("Unable to read spec objects from file", "apply", "", err)
	}

	// Apply all spec
	for _, spec := range specs {
		if err := ApplySpec(login.Token, account, spec); err != nil {
			return err
		}
	}

	return nil
}

// Apply reads the config file(s) from the provided file / directory and applies it to the server
func Apply(fileName string) error {
	account, err := utils.GetSelectedAccount()
	if err != nil {
		return utils.LogError("Unable to fetch account information", err)
	}
	login, err := utils.Login(account)
	if err != nil {
		return utils.LogError("Unable to login", err)
	}

	specs, err := utils.ReadSpecObjectsFromFile(fileName)
	if err != nil {
		return utils.LogError("Unable to read spec objects from file", err)
	}

	// Apply all spec
	for _, spec := range specs {
		if err := ApplySpec(login.Token, account, spec); err != nil {
			return err
		}
	}

	return nil
}

// ApplySpec takes a spec object and applies it
func ApplySpec(token string, account *model.Account, specObj *model.SpecObject) error {
	requestBody, err := json.Marshal(specObj.Spec)
	if err != nil {
		logrus.Errorf("error while applying service unable to marshal spec - %s", err.Error())
		return err
	}
	url, err := adjustPath(fmt.Sprintf("%s%s", account.ServerURL, specObj.API), specObj.Meta)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logrus.Errorf("error while applying service unable to send http request - %s", err.Error())
		return err
	}

	v := map[string]interface{}{}
	_ = json.NewDecoder(resp.Body).Decode(&v)
	utils.CloseTheCloser(req.Body)
	if resp.StatusCode != 200 {
		logrus.Errorf("error while applying service got http status code %s - %s", resp.Status, v["error"])
		return fmt.Errorf("%v", v["error"])
	}
	logrus.Infof("Successfully applied %s", specObj.Type)
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
