package cmd

import (
	"fmt"

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

	// var completionCmd = &cobra.Command{
	// 	Use:   "completion",
	// 	Short: "Shell completion code",
	// 	// RunE: func(cmd *cobra.Command, args []string) error {
	// 	// 	return rootCmd.GenBashCompletionFile("/home/prithvi/project/space-cloud/out.sh")
	// 	// },
	// }
	// var bashCmd = &cobra.Command{
	// 	Use:   "bash",
	// 	Short: "Bash shell completion code",
	// 	RunE: func(cmd *cobra.Command, args []string) error {
	// 		return rootCmd.GenBashCompletionFile("out.sh")
	// 	},
	// }
	// var zshCmd = &cobra.Command{
	// 	Use:   "zsh",
	// 	Short: "Zsh shell completion code",
	// 	RunE: func(cmd *cobra.Command, args []string) error {
	// 		return rootCmd.GenZshCompletionFile("out1.sh")
	// 	},
	// }
	// var powershellCmd = &cobra.Command{
	// 	Use:   "powershell",
	// 	Short: "Powershell shell completion code",
	// 	RunE: func(cmd *cobra.Command, args []string) error {
	// 		return rootCmd.GenPowerShellCompletionFile("out2.sh")
	// 	},
	// }

	var completionCmd = &cobra.Command{
		Use:   "completion [bash|zsh] [--no-descriptions]",
		Short: "Generate completion script",
		// DisableFlagsInUseLine: true,
		// ValidArgs:             []string{"bash", "zsh1\tOriginal zsh completion", "zsh2\tV2 zsh completion", "fish"},
		// Args:                  cobra.ExactValidArgs(1),
		PreRun: func(cmd *cobra.Command, args []string) {
			err := viper.BindPFlag("file", cmd.Flags().Lookup("file"))
			if err != nil {
				_ = utils.LogError("Unable to bind the flag ('file')", err)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			switch args[0] {
			case "bash":
				file := viper.GetString("file")
				rootCmd.GenBashCompletionFile(fmt.Sprintf("%s/out.sh", file))
				break
			case "zsh":
				file := viper.GetString("file")
				fmt.Printf(fmt.Sprintf("%s/out.sh", file))
				rootCmd.GenZshCompletionFile(fmt.Sprintf("%s/_space-cli", file))
				break
			}
		},
	}
	completionCmd.Flags().StringP("file", "", "/etc/bash_completion.d", "")

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
