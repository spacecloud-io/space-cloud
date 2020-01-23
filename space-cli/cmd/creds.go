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

func getSpaceCliDirectory() string {
	return fmt.Sprintf("%s/space-cli", getSpaceCloudDirectory())
}

func getAccountConfigPath() string {
	return fmt.Sprintf("%s/accounts.yaml", getSpaceCliDirectory())
}

func getSelectedAccount(credential *model.Credential) *model.Account {
	var selectedaccount model.Account
	for _, v := range credential.Accounts {
		if credential.SelectedAccount == v.ID {
			selectedaccount = v
		}
	}
	return &selectedaccount
}

func getCreds() (*model.Credential, error) {
	fileName := getAccountConfigPath()
	yamlFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		logrus.Error("error getting credential unable to read accounts config file - %v", err)
		return nil, err
	}

	credential := new(model.Credential)
	if err := yaml.Unmarshal(yamlFile, credential); err != nil {
		logrus.Error("error getting credential unable to unmarshal accounts config file - %v", err)
		return nil, err
	}
	return credential, nil
}

func checkCred(selectedAccount *model.Account) error {
	fileName := getAccountConfigPath()
	yamlFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		credential := model.Credential{
			Accounts:        []model.Account{*selectedAccount},
			SelectedAccount: selectedAccount.ID,
		}
		if err := generateYamlFile(&credential); err != nil {
			logrus.Errorf("error in checking credentials unable to create accounts yaml file - %v", err)
			return err
		}
	}
	credential := new(model.Credential)
	if err := yaml.Unmarshal(yamlFile, credential); err != nil {
		return err
	}
	for _, val := range credential.Accounts {
		if val.ID == selectedAccount.ID {
			val.ID, val.UserName, val.Key, val.ServerUrl = selectedAccount.ID, selectedAccount.UserName, selectedAccount.Key, selectedAccount.ServerUrl
			if err := generateYamlFile(credential); err != nil {
				logrus.Errorf("error in checking credentials unable to update accounts yaml file - %v", err)
				return err
			}
			return nil
		}
	}
	credential.Accounts = append(credential.Accounts, *selectedAccount)
	credential.SelectedAccount = selectedAccount.ID
	if err := generateYamlFile(credential); err != nil {
		logrus.Errorf("error in checking credentials unable to update accounts yaml file - %v", err)
		return err
	}
	return nil
}
