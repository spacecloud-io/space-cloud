package accounts

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/spaceuptech/space-cloud/space-cli/cmd/modules/project"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils"
)

// Commands are the set of account commands for space-cli
func Commands() []*cobra.Command {

	var accountsCmd = &cobra.Command{
		Use:   "accounts",
		Short: "Operations for space-cloud accounts",
	}
	autoCompleteFunc := func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		switch len(args) {
		case 0:
			credential, err := utils.GetCredentials()
			if err != nil {
				utils.LogDebug("Unable to get all the stored credentials", nil)
				return nil, cobra.ShellCompDirectiveDefault
			}
			accountIDs := []string{}
			for _, v := range credential.Accounts {
				accountIDs = append(accountIDs, v.ID)
			}
			return accountIDs, cobra.ShellCompDirectiveDefault
		}
		return nil, cobra.ShellCompDirectiveDefault
	}

	var setDefaultProjectCommand = &cobra.Command{
		Use:   "set-default-project",
		Short: "Sets the default project to be used if --project flag is not provided",
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			obj, err := project.GetProjectConfig("", "project", map[string]string{})
			if err != nil {
				return nil, cobra.ShellCompDirectiveDefault
			}
			var projects []string
			for _, v := range obj {
				projects = append(projects, v.Meta["project"])
			}
			return projects, cobra.ShellCompDirectiveDefault
		},
		SilenceErrors: true,
		RunE:          actionSetDefaultProject,
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
		ValidArgsFunction: autoCompleteFunc,
		SilenceErrors:     true,
		RunE:              actionViewAccount,
	}

	var setAccountCommand = &cobra.Command{
		Use:               "set",
		Short:             "set the given account as the selected account",
		SilenceErrors:     true,
		RunE:              actionSetAccount,
		ValidArgsFunction: autoCompleteFunc,
	}

	var deleteAccountCommand = &cobra.Command{
		Use:               "delete",
		Short:             "deletes the given account",
		SilenceErrors:     true,
		RunE:              actionDeleteAccount,
		ValidArgsFunction: autoCompleteFunc,
	}

	viewAccountsCommand.Flags().BoolP("show-keys", "", false, "shows the keys of the accounts")

	accountsCmd.AddCommand(viewAccountsCommand)
	accountsCmd.AddCommand(setAccountCommand)
	accountsCmd.AddCommand(deleteAccountCommand)
	accountsCmd.AddCommand(setDefaultProjectCommand)

	return []*cobra.Command{accountsCmd}
}

func actionSetDefaultProject(cmd *cobra.Command, args []string) error {
	project := ""
	if len(args) > 0 {
		project = args[0]
	}
	return utils.SetDefaultProject(project)
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
