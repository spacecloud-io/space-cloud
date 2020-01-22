package cmd

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"

	"github.com/spaceuptech/space-cli/model"
)

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
	fileName := fmt.Sprintf("/%s/galaxy/config.yaml", getHomeDirectory())
	yamlFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("error reading yaml file: %s", err)
	}

	credential := new(model.Credential)
	if err := yaml.Unmarshal(yamlFile, credential); err != nil {
		return nil, err
	}
	return credential, nil
}

func checkCred(selectedAccount *model.Account) error {
	fileName := fmt.Sprintf("/%s/galaxy/config.yaml", getHomeDirectory())
	yamlFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		credential := model.Credential{
			Accounts:        []model.Account{*selectedAccount},
			SelectedAccount: selectedAccount.ID,
		}
		if err := generateYamlFile(&credential); err != nil {
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
				return err
			}
			return nil
		}
	}
	credential.Accounts = append(credential.Accounts, *selectedAccount)
	credential.SelectedAccount = selectedAccount.ID
	if err := generateYamlFile(credential); err != nil {
		return err
	}
	return nil
}
