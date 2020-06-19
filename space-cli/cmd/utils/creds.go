package utils

import (
	"io/ioutil"

	"github.com/ghodss/yaml"
	"github.com/sirupsen/logrus"

	"github.com/spaceuptech/space-cli/cmd/model"
)

// getSelectedAccount gets the account information of the selected account
func getSelectedAccount() (*model.Account, error) {
	filePath := getAccountConfigPath()
	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		logrus.Errorf("error getting credential unable to read accounts config file - %s", err.Error())
		return nil, err
	}

	credential := new(model.Credential)
	if err := yaml.Unmarshal(yamlFile, credential); err != nil {
		logrus.Errorf("error getting credential unable to unmarshal accounts config file - %s", err.Error())
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

// StoreCredentials stores the credential in the accounts config file
func StoreCredentials(account *model.Account) error {
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
			val.ID, val.UserName, val.Key, val.ServerURL = account.ID, account.UserName, account.Key, account.ServerURL
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
