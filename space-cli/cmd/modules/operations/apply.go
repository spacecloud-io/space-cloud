package operations

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"

	"github.com/spaceuptech/space-cloud/space-cli/cmd/model"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils"
)

// Apply reads the config file(s) from the provided file / directory and applies it to the server
func Apply(applyName string, isForceApply bool, delay time.Duration, retry int) error {
	if !strings.HasSuffix(applyName, ".yaml") {
		dirName := applyName
		if err := os.Chdir(dirName); err != nil {
			return utils.LogError(fmt.Sprintf("Unable to switch to directory %s", dirName), err)
		}
		// list file of current directory. Note : we have changed the directory with the help of above function
		files, err := ioutil.ReadDir(".")
		if err != nil {
			return utils.LogError(fmt.Sprintf("Unable to fetch config files from %s", dirName), err)
		}

		account, token, err := utils.LoginWithSelectedAccount()
		if err != nil {
			return utils.LogError("Couldn't get account details or login token", err)
		}

		var fileNames []string
		// filter directories
		for _, fileInfo := range files {
			if !fileInfo.IsDir() {
				fileNames = append(fileNames, fileInfo.Name())
			}
		}
		// sort file names alphanumerically

		for _, fileName := range fileNames {
			if strings.HasSuffix(fileName, ".yaml") {
				specs, err := utils.ReadSpecObjectsFromFile(fileName)
				if err != nil {
					return utils.LogError("Unable to read spec objects from file", err)
				}

				// Apply all spec
				for _, spec := range specs {
					if spec.Type == "db-schema" && !isForceApply {
						prompt := &survey.Confirm{
							Message: "Changing the schema can cause data loss (this option will be applied to all resources of type db-schema).\n Are you sure you want to continue?",
						}
						_ = survey.AskOne(prompt, &isForceApply)
						if !isForceApply {
							utils.LogInfo(fmt.Sprintf("Skipping the resource (db-schema) having meta %v", spec.Meta))
							continue
						}
					}

					currentRetryCount := 0
					for {
						err = ApplySpec(token, account, spec)
						if err == nil {
							break
						}
						err = utils.LogError(fmt.Sprintf("Unable to apply file (%s) spec object with id (%v) type (%v)", fileName, spec.Meta["id"], spec.Type), err)
						if currentRetryCount == retry {
							return err
						}
						currentRetryCount++
						utils.LogInfo(fmt.Sprintf("Retrying spec object with id (%v) and type (%v)", spec.Meta["id"], spec.Type))
					}
					time.Sleep(delay)
				}
			}
		}
		return nil
	}

	account, token, err := utils.LoginWithSelectedAccount()
	if err != nil {
		return utils.LogError("Couldn't get account details or login token", err)
	}

	specs, err := utils.ReadSpecObjectsFromFile(applyName)
	if err != nil {
		return utils.LogError("Unable to read spec objects from file", err)
	}

	// Apply all spec
	for _, spec := range specs {
		if spec.Type == "db-schema" && !isForceApply {
			prompt := &survey.Confirm{
				Message: "Changing the schema can cause data loss (this option will be applied to all resources of type db-schema).\n Are you sure you want to continue?",
			}
			_ = survey.AskOne(prompt, &isForceApply)
			if !isForceApply {
				utils.LogInfo(fmt.Sprintf("Skipping the resource (db-schema) having meta %v", spec.Meta))
				continue
			}
		}
		currentRetryCount := 0
		for {
			err = ApplySpec(token, account, spec)
			if err == nil {
				break
			}
			err = utils.LogError(fmt.Sprintf("Unable to apply spec object with id (%v) type (%v)", spec.Meta["id"], spec.Type), err)
			if currentRetryCount == retry {
				return err
			}
			currentRetryCount++
			utils.LogInfo(fmt.Sprintf("Retrying spec object with id (%v) and type (%v)", spec.Meta["id"], spec.Type))
		}
		time.Sleep(delay)
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

	if resp.StatusCode == http.StatusAccepted {
		// Make checker send this status
		utils.LogInfo(fmt.Sprintf("Successfully queued %s", specObj.Type))
	} else if resp.StatusCode == http.StatusOK {
		utils.LogInfo(fmt.Sprintf("Successfully applied %s", specObj.Type))
	} else {
		_ = utils.LogError(fmt.Sprintf("error while applying service got http status code %s - %s", resp.Status, v["error"]), nil)
		return fmt.Errorf("%v", v["error"])
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
