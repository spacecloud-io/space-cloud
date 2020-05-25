package accounts

import (
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spaceuptech/space-cli/cmd/model"
	"github.com/spaceuptech/space-cli/cmd/utils"
)

func setAccount(prefix string) error {
	credential, err := utils.GetCredentials()
	if err != nil {
		return err
	}
	accountIDOptions := []string{}
	for _, v := range credential.Accounts {
		v.ID = strings.ToLower(v.ID)
		accountIDOptions = append(accountIDOptions, v.ID)
	}

	prefix = strings.ToLower(prefix)
	filteredAccountIDOptions, exists := filterAccounts(credential.Accounts, prefix)

	if exists {
		if err := survey.AskOne(&survey.Select{Message: "Choose the account ID to be set: ", Options: filteredAccountIDOptions}, &prefix); err != nil {
			return err
		}
	} else {
		if prefix != "" {
			utils.LogInfo("Warning! No account found for prefix provided, showing all")
		}
		if err := survey.AskOne(&survey.Select{Message: "Choose the account ID to be set: ", Options: accountIDOptions}, &prefix); err != nil {
			return err
		}
	}

	credential.SelectedAccount = prefix

	if err := utils.GenerateAccountsFile(credential); err != nil {
		_ = utils.LogError("Could not update yaml file while setting selected account", nil)
		return nil
	}

	return nil
}

func filterAccounts(accounts []*model.Account, prefix string) ([]string, bool) {
	filteredAccountOptions := []string{}
	exists := false
	for _, account := range accounts {
		if prefix != "" && strings.HasPrefix(account.ID, prefix) {
			filteredAccountOptions = append(filteredAccountOptions, account.ID)
			exists = true
		}
	}

	return filteredAccountOptions, exists
}
