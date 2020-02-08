package cmd

import (
	"fmt"
	"io/ioutil"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"github.com/spaceuptech/space-cli/model"
)

func getSpaceCloudDirectory() string {
	return fmt.Sprintf("%s/.space-cloud", getHomeDirectory())
}

func getAccountConfigPath() string {
	return fmt.Sprintf("%s/accounts.yaml", getSpaceCloudDirectory())
}

func getSelectedAccount() (*model.Account, error) {
	filePath := getAccountConfigPath()
	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		logrus.Error("error getting credential unable to read accounts config file - %s", err.Error())
		return nil, err
	}

	credential := new(model.Credential)
	if err := yaml.Unmarshal(yamlFile, credential); err != nil {
		logrus.Error("error getting credential unable to unmarshal accounts config file - %s", err.Error())
		return nil, err
	}

	var account *model.Account
	for _, v := range credential.Accounts {
		if credential.SelectedAccount == v.ID {
			account = v
		}
	}
	return account, nil
}

func checkCred(account *model.Account) error {
	yamlFile, err := ioutil.ReadFile(getAccountConfigPath())
	if err != nil {
		// accounts.yaml file doesn't exist create new one
		credential := model.Credential{
			Accounts:        []*model.Account{account},
			SelectedAccount: account.ID,
		}
		if err := generateYamlFile(&credential); err != nil {
			logrus.Errorf("error in checking credentials unable to create accounts yaml file - %v", err)
			return err
		}
	}
	// file already exists, read data from accounts.yaml file
	credential := new(model.Credential)
	if err := yaml.Unmarshal(yamlFile, credential); err != nil {
		return err
	}
	for _, val := range credential.Accounts {
		// update account if already exists
		if val.ID == account.ID {
			val.ID, val.UserName, val.Key, val.ServerUrl = account.ID, account.UserName, account.Key, account.ServerUrl
			credential.SelectedAccount = account.ID
			if err := generateYamlFile(credential); err != nil {
				logrus.Errorf("error in checking credentials unable to update accounts yaml file - %v", err)
				return err
			}
			return nil
		}
	}
	// add new account to already existing accounts.yaml file
	credential.Accounts = append(credential.Accounts, account)
	credential.SelectedAccount = account.ID
	if err := generateYamlFile(credential); err != nil {
		logrus.Errorf("error in checking credentials unable to update accounts yaml file - %v", err)
		return err
	}
	return nil
}
