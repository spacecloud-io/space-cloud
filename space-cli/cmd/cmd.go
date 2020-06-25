package cmd

import (
	"fmt"
	"strings"

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
	}

	var completionCmd = &cobra.Command{
		Use:   "completion [bash|zsh] [--no-descriptions]",
		Short: "Generate completion script",
		PreRun: func(cmd *cobra.Command, args []string) {
			err := viper.BindPFlag("file", cmd.Flags().Lookup("file"))
			if err != nil {
				_ = utils.LogError("Unable to bind the flag ('file')", err)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			file := viper.GetString("file")
			switch args[0] {
			case "bash":
				if file == "" {
					_ = utils.LogError("Creating file ('space-cli.sh') in current Directory", nil)
					rootCmd.GenBashCompletionFile("space-cli.sh")
					break
				}
				if !strings.HasSuffix(file, ".sh") {
					_ = utils.LogError("file path should end with .sh file", nil)
					break
				}
				rootCmd.GenBashCompletionFile(fmt.Sprintf("%s", file))
				break
			case "zsh":
				if file == "" {
					_ = utils.LogError("Creating file ('_space-cli') in current Directory", nil)
					rootCmd.GenBashCompletionFile("_space-cli")
					break
				}
				if !strings.HasSuffix(file, "_space-cli") {
					_ = utils.LogError("file path should end with _space-cli", nil)
					break
				}
				rootCmd.GenZshCompletionFile(fmt.Sprintf("%s", file))
				break
			}
		},
	}
	completionCmd.Flags().StringP("file", "", "", "")

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
	rootCmd.AddCommand(completionCmd)
	return rootCmd
}

// func main() {
// 	_ = GetRootCommand()
// }
