package accounts

import (
	"strings"

	"github.com/spaceuptech/space-cli/cmd/model"
	"github.com/spaceuptech/space-cli/cmd/utils"
)

func deleteAccount(prefix string) error {
	credential, err := utils.GetCredentials()
	if err != nil {
		return err
	}

	prefix = strings.ToLower(prefix)
	prefix, err = filterAccounts(credential.Accounts, prefix)
	if err != nil {
		return err
	}

	if prefix == credential.SelectedAccount {
		return utils.LogError("Chosen account cannot be deleted. Use space-cli accounts set to change the selected account", nil)
	}

	for i, v := range credential.Accounts {
		if v.ID == prefix {
			credential.Accounts = removeAccount(credential.Accounts, i)
		}
	}

	if err := utils.GenerateAccountsFile(credential); err != nil {
		return utils.LogError("Couldn't update accounts.yaml file", err)
	}

	return nil
}

func removeAccount(accounts []*model.Account, index int) []*model.Account {
	return append(accounts[:index], accounts[index+1:]...)
}
