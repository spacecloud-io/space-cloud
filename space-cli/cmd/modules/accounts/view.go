package accounts

import (
	"os"
	"strings"
	"unicode/utf8"

	"github.com/olekukonko/tablewriter"

	"github.com/spaceuptech/space-cli/cmd/model"
	"github.com/spaceuptech/space-cli/cmd/utils"
)

func listAccounts(prefix string, showKeys bool) error {
	credential, err := utils.GetCredentials()
	if err != nil {
		return err
	}
	if len(credential.Accounts) == 0 {
		utils.LogInfo("No accounts found. Try adding an account using `space-cli login`")
		return nil
	}

	accounts := []*model.Account{}
	for _, v := range credential.Accounts {
		if strings.HasPrefix(strings.ToLower(v.ID), strings.ToLower(prefix)) {
			accounts = append(accounts, v)
		}
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

	table.Render()

	return nil
}
