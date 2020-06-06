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
		Short: "Operations for space-cloud accounts",
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
		SilenceErrors: true,
		RunE:          actionViewAccount,
	}

	var setAccountCommand = &cobra.Command{
		Use:           "set",
		Short:         "set the given account as the selected account",
		SilenceErrors: true,
		RunE:          actionSetAccount,
	}

	var deleteAccountCommand = &cobra.Command{
		Use:           "delete",
		Short:         "deletes the given account",
		SilenceErrors: true,
		RunE:          actionDeleteAccount,
	}

	viewAccountsCommand.Flags().BoolP("show-keys", "", false, "shows the keys of the accounts")

	accountsCmd.AddCommand(viewAccountsCommand)
	accountsCmd.AddCommand(setAccountCommand)
	accountsCmd.AddCommand(deleteAccountCommand)

	return []*cobra.Command{accountsCmd}
}

func actionViewAccount(cmd *cobra.Command, args []string) error {

	showKeys := viper.GetBool("show-keys")

	prefix := ""
	if len(args) > 0 {
		prefix = args[0]
	}

	return listAccounts(prefix, showKeys)
}

func actionSetAccount(cmd *cobra.Command, args []string) error {

	prefix := ""
	if len(args) > 0 {
		prefix = args[0]
	}

	return setAccount(prefix)
}

func actionDeleteAccount(cmd *cobra.Command, args []string) error {

	prefix := ""
	if len(args) > 0 {
		prefix = args[0]
	}

	return deleteAccount(prefix)
}
