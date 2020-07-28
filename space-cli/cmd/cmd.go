package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/spaceuptech/space-cloud/space-cli/cmd/modules"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/modules/accounts"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/modules/addons"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/modules/deploy"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/modules/login"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/modules/logs"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/modules/operations"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/modules/project"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils"
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
			switch args[0] {
			case "bash":
				err := rootCmd.GenBashCompletion(os.Stdout)
				if err != nil {
					_ = utils.LogError(fmt.Sprintf("Error in generating Bash completion script-%s", err), nil)
				}

			case "zsh":
				err := rootCmd.GenZshCompletion(os.Stdout)
				if err != nil {
					_ = utils.LogError(fmt.Sprintf("Error in generating Zsh completion script- %s", err), nil)
				}
			}
		},
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
	err = rootCmd.RegisterFlagCompletionFunc("log-level", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"info", "debug", "error"}, cobra.ShellCompDirectiveDefault
	})
	if err != nil {
		utils.LogDebug("Unable to provide suggetion for flag ('log-level')", nil)
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
	err = rootCmd.RegisterFlagCompletionFunc("project", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		obj, err := project.GetProjectConfig("", "project", map[string]string{})
		if err != nil {
			return nil, cobra.ShellCompDirectiveDefault
		}
		var projects []string
		for _, v := range obj {
			projects = append(projects, v.Meta["project"])
		}
		return projects, cobra.ShellCompDirectiveDefault
	})
	if err != nil {
		utils.LogDebug("Unable to provide suggetion for flag ('project')", nil)
	}

	rootCmd.AddCommand(modules.FetchGenerateSubCommands())
	rootCmd.AddCommand(modules.FetchGetSubCommands())
	rootCmd.AddCommand(addons.Commands()...)
	rootCmd.AddCommand(deploy.Commands()...)
	rootCmd.AddCommand(operations.Commands()...)
	rootCmd.AddCommand(login.Commands()...)
	rootCmd.AddCommand(accounts.Commands()...)
	rootCmd.AddCommand(logs.GetSubCommands()...)
	rootCmd.AddCommand(completionCmd)
	return rootCmd
}

// func main() {
// 	_ = GetRootCommand()
// }
