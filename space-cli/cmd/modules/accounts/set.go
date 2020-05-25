package accounts

import (
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spaceuptech/space-cli/cmd/utils"
)

func setAccount(prefix string) error {
	credential, err := utils.GetCredentials()
	if err != nil {
		return err
	}
	exists := false
	accountIDOptions := []string{}
	if prefix != "" {
		for _, v := range credential.Accounts {
			accountIDOptions = append(accountIDOptions, v.ID)
			if strings.HasPrefix(v.ID, prefix) {
				prefix = v.ID
				exists = true
			}
		}
	}
	if !exists {
		if prefix != "" {
			utils.LogInfo("Warning! No account found for prefix provided, showing all")
		}
		if err := survey.AskOne(&survey.Select{Message: "Choose the account ID to be set: ", Options: accountIDOptions}, &prefix); err != nil {
			return err
		}
	}
	credential.SelectedAccount = prefix
	if err := utils.GenerateYamlFile(credential); err != nil {
		_ = utils.LogError("Could not update yaml file while setting selected account", nil)
		return nil
	}

	return nil
}
