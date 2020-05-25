package accounts

import (
	"github.com/spaceuptech/space-cli/cmd/model"
	"github.com/spaceuptech/space-cli/cmd/utils"
)

func deleteAccount(accountID string) error {
	credential, err := utils.GetCredentials()
	if err != nil {
		return err
	}
	for i, v := range credential.Accounts {
		if v.ID == accountID {
			credential.Accounts = removeAccount(credential.Accounts, i)
		}
	}

	if err := utils.GenerateAccountsFile(credential); err != nil {
		_ = utils.LogError("Could not update yaml file while deleting selected account", nil)
		return nil
	}

	return nil
}

func removeAccount(accounts []*model.Account, index int) []*model.Account {
	return append(accounts[:index], accounts[index+1:]...)
}
