package accounts

import (
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spaceuptech/space-cli/cmd/model"
	"github.com/spaceuptech/space-cli/cmd/utils"
	"github.com/spaceuptech/space-cli/cmd/utils/input"
)

func setAccount(prefix string) error {
	credential, err := utils.GetCredentials()
	if err != nil {
		return err
	}

	prefix = strings.ToLower(prefix)
	prefix, err = filterAccounts(credential.Accounts, prefix)
	if err != nil {
		return err
	}

	credential.SelectedAccount = prefix

	if err := utils.GenerateAccountsFile(credential); err != nil {
		return utils.LogError("Couldn't update accounts.yaml file", err)
	}

	return nil
}

func filterAccounts(accounts []*model.Account, prefix string) (string, error) {
	filteredAccountIDs := []string{}
	doesAccountExists := false

	accountIDs := []string{}
	for _, v := range accounts {
		accountIDs = append(accountIDs, v.ID)
	}

	for _, account := range accounts {
		if prefix != "" && strings.HasPrefix(strings.ToLower(account.ID), prefix) {
			filteredAccountIDs = append(filteredAccountIDs, account.ID)
			doesAccountExists = true
		}
	}
	if doesAccountExists {
		if err := input.Survey.AskOne(&survey.Select{Message: "Choose the account ID: ", Options: filteredAccountIDs, Default: filteredAccountIDs[0]}, &prefix); err != nil {
			return "", err
		}
	} else {
		if prefix != "" {
			utils.LogInfo("Warning! No account found for prefix provided, showing all")
		}
		if err := input.Survey.AskOne(&survey.Select{Message: "Choose the account ID: ", Options: accountIDs, Default: accountIDs[0]}, &prefix); err != nil {
			return "", err
		}
	}

	return prefix, nil
}
