package accounts

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/spaceuptech/space-cli/cmd/utils"
)

// Commands are the set of account commands for space-cli
func Commands() []*cobra.Command {

	var accountsCmd = &cobra.Command{
		Use:   "accounts",
		Short: "CRUD operations for space-cloud accounts",
	}

	var viewAccountsCommand = &cobra.Command{
		Use:   "view",
		Short: "list all space-cloud accounts",
		PreRun: func(cmd *cobra.Command, args []string) {
			err := viper.BindPFlag("show-keys", cmd.Flags().Lookup("show-keys"))
			if err != nil {
				_ = utils.LogError("Unable to bind the flag ('show-keys')", nil)
			}
		},
		RunE: actionViewAccount,
	}

	viewAccountsCommand.Flags().BoolP("show-keys", "", false, "shows the keys of the accounts")

	accountsCmd.AddCommand(viewAccountsCommand)

	return []*cobra.Command{accountsCmd}
}

func actionViewAccount(cmd *cobra.Command, args []string) error {

	showKeys := viper.GetBool("show-keys")

	accountID := ""
	if len(args) > 0 {
		accountID = args[0]
	}

	if err := listAccounts(accountID, showKeys); err != nil {
		return err
	}

	return nil
}
