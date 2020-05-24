package accounts

import (
	"github.com/spaceuptech/space-cli/cmd/utils"
)

func setAccount(accountID string) error {
	credential, err := utils.GetCredentials()
	if err != nil {
		return err
	}
	exists := false
	for _, v := range credential.Accounts {
		if v.ID == accountID {
			exists = true
		}
	}
	if exists == false {
		_ = utils.LogError("No account exists with the given account ID", nil)
		return nil
	}
	credential.SelectedAccount = accountID

	return nil
}
