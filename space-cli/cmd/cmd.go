package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/spaceuptech/space-cli/cmd/modules"
	"github.com/spaceuptech/space-cli/cmd/modules/accounts"
	"github.com/spaceuptech/space-cli/cmd/modules/addons"
	"github.com/spaceuptech/space-cli/cmd/modules/deploy"
	"github.com/spaceuptech/space-cli/cmd/modules/login"
	"github.com/spaceuptech/space-cli/cmd/modules/operations"
	"github.com/spaceuptech/space-cli/cmd/utils"
)

// GetRootCommand return the rootcmd
func GetRootCommand() *cobra.Command {

	var rootCmd = &cobra.Command{
		Use:     "space-cli",
		Version: version,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			utils.SetLogLevel(viper.GetString("log-level"))
		},
		SilenceUsage: true,
	}

	rootCmd.PersistentFlags().StringP("log-level", "", "info", "Sets the log level of the command")
	err := viper.BindPFlag("log-level", rootCmd.PersistentFlags().Lookup("log-level"))
	if err != nil {
		_ = utils.LogError("Unable to bind the flag ('log-level')", nil)
	}
	err = viper.BindEnv("log-level", "LOG_LEVEL")
	if err != nil {
		_ = utils.LogError("Unable to bind flag ('log-level') to environment variables", nil)
	}

	rootCmd.PersistentFlags().StringP("project", "", "", "The project id to perform the options in")
	err = viper.BindPFlag("project", rootCmd.PersistentFlags().Lookup("project"))
	if err != nil {
		_ = utils.LogError("Unable to bind the flag ('project')", nil)
	}
	err = viper.BindEnv("project", "PROJECT")
	if err != nil {
		_ = utils.LogError("Unable to bind flag ('project') to environment variables", nil)
	}

	rootCmd.AddCommand(modules.FetchGenerateSubCommands())
	rootCmd.AddCommand(modules.FetchGetSubCommands())
	rootCmd.AddCommand(addons.Commands()...)
	rootCmd.AddCommand(deploy.Commands()...)
	rootCmd.AddCommand(operations.Commands()...)
	rootCmd.AddCommand(login.Commands()...)
	rootCmd.AddCommand(accounts.Commands()...)
	return rootCmd
}

// func main() {
// 	_ = GetRootCommand()
// }
