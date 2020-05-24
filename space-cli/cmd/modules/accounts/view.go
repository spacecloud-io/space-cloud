package accounts

import (
	"os"
	"strings"
	"unicode/utf8"

	"github.com/olekukonko/tablewriter"

	"github.com/spaceuptech/space-cli/cmd/model"
	"github.com/spaceuptech/space-cli/cmd/utils"
)

func listAccounts(accountID string, showKeys bool) error {
	credential, err := utils.GetCredentials()
	if err != nil {
		return err
	}

	var accounts []*model.Account
	if accountID != "" {
		for _, v := range credential.Accounts {
			if v.ID == accountID {
				accounts = append(accounts, v)
			}
		}
	} else {
		accounts = append(accounts, credential.Accounts...)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Username", "Key"})

	table.SetBorder(false)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")

	for _, account := range accounts {
		if showKeys {
			table.Append([]string{account.ID, account.UserName, account.Key})
		} else {
			table.Append([]string{account.ID, account.UserName, strings.Repeat("*", utf8.RuneCountInString(account.Key))})
		}
	}

	return nil
}
