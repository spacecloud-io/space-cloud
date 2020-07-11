package login

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils"
)

// Commands is the list of commands the utils module exposes
func Commands() []*cobra.Command {
	var loginCommands = &cobra.Command{
		Use:   "login",
		Short: "Logs into space cloud",
		PreRun: func(cmd *cobra.Command, args []string) {
			err := viper.BindPFlag("username", cmd.Flags().Lookup("username"))
			if err != nil {
				_ = utils.LogError("Unable to bind the flag ('username')", nil)
			}
			err = viper.BindPFlag("key", cmd.Flags().Lookup("key"))
			if err != nil {
				_ = utils.LogError("Unable to bind the flag ('key')", nil)
			}
			err = viper.BindPFlag("url", cmd.Flags().Lookup("url"))
			if err != nil {
				_ = utils.LogError("Unable to bind the flag ('url')", nil)
			}

		},
		RunE:          actionLogin,
		SilenceErrors: true,
	}
	loginCommands.Flags().StringP("username", "", "None", "Accepts the username for login")
	err := viper.BindEnv("username", "USER_NAME")
	if err != nil {
		_ = utils.LogError("Unable to bind flag ('username') to environment variables", nil)
	}

	loginCommands.Flags().StringP("id", "", "None", "Accepts the id for login")
	err = viper.BindEnv("id", "ID")
	if err != nil {
		_ = utils.LogError("Unable to bind flag ('id') to environment variables", nil)
	}

	loginCommands.Flags().StringP("key", "", "None", "Accepts the access key to be verified during login")
	err = viper.BindEnv("key", "KEY")
	if err != nil {
		_ = utils.LogError("Unable to bind flag ('key') to environment variables", nil)
	}

	loginCommands.Flags().StringP("url", "", "http://localhost:4122", "Accepts the URL of server")
	err = viper.BindEnv("url", "URL")
	if err != nil {
		_ = utils.LogError("Unable to bind flag ('url') to environment variables", nil)
	}

	return []*cobra.Command{loginCommands}
}

func actionLogin(cmd *cobra.Command, args []string) error {
	userName := viper.GetString("username")
	ID := viper.GetString("id")
	key := viper.GetString("key")
	url := viper.GetString("url")

	return utils.LoginStart(userName, ID, key, url)
}
