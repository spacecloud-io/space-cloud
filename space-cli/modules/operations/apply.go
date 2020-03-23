package operations

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/urfave/cli"

	"github.com/spaceuptech/space-cli/model"
	"github.com/spaceuptech/space-cli/utils"
)

// ActionApply applies the file / folder
func ActionApply(cli *cli.Context) error {
	args := os.Args
	if len(args) != 3 {
		_ = utils.LogError("error while applying service incorrect number of arguments provided", nil)
		return fmt.Errorf("incorrect number of arguments provided")
	}

	fileName := args[2]

	return Apply(fileName)
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
		_ = utils.LogError(fmt.Sprintf("error while applying service unable to marshal spec - %s", err.Error()), nil)
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
		_ = utils.LogError(fmt.Sprintf("error while applying service unable to send http request - %s", err.Error()), nil)
		return err
	}

	v := map[string]interface{}{}
	_ = json.NewDecoder(resp.Body).Decode(&v)
	utils.CloseTheCloser(req.Body)
	if resp.StatusCode != 200 {
		_ = utils.LogError(fmt.Sprintf("error while applying service got http status code %s - %s", resp.Status, v["error"]), nil)
		return fmt.Errorf("%v", v["error"])
	}
	utils.LogInfo(fmt.Sprintf("Successfully applied %s", specObj.Type))
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
