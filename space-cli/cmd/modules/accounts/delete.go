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

	doesAccountExist := false
	for i, v := range credential.Accounts {
		if v.ID == accountID {
			credential.Accounts = removeAccount(credential.Accounts, i)
			doesAccountExist = true
		}
	}

	if !doesAccountExist {
		return utils.LogError("Account ID not found in accounts.yaml", nil)
	}

	if err := utils.GenerateAccountsFile(credential); err != nil {
		return utils.LogError("Couldn't update accounts.yaml file", err)
	}

	return nil
}

func removeAccount(accounts []*model.Account, index int) []*model.Account {
	return append(accounts[:index], accounts[index+1:]...)
}
