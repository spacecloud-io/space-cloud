package main

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spaceuptech/space-cli/modules"
	"github.com/spaceuptech/space-cli/modules/addons"
	"github.com/spaceuptech/space-cli/modules/auth"
	"github.com/spaceuptech/space-cli/modules/database"
	"github.com/spaceuptech/space-cli/modules/deploy"
	"github.com/spaceuptech/space-cli/modules/eventing"
	"github.com/spaceuptech/space-cli/modules/filestore"
	"github.com/spaceuptech/space-cli/modules/ingress"
	"github.com/spaceuptech/space-cli/modules/letsencrypt"
	"github.com/spaceuptech/space-cli/modules/operations"
	"github.com/spaceuptech/space-cli/modules/project"
	"github.com/spaceuptech/space-cli/modules/services"
	"github.com/spaceuptech/space-cli/modules/userman"
	"github.com/spaceuptech/space-cli/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Used for flags.
	cfgFile     string
	userLicense string
)

var rootCmd = &cobra.Command{
	Use:     "space-cli",
	Version: "0.16.0",
}

func execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringP("log-level", "", "info", "Sets the log level of the command")
	viper.BindPFlag("log-level", rootCmd.PersistentFlags().Lookup("log-level"))

	rootCmd.PersistentFlags().StringP("project", "", "", "The project to add the add-on to")
	viper.BindPFlag("project", rootCmd.PersistentFlags().Lookup("project"))

	rootCmd.PersistentFlags().StringP("project", "", "", "The project to remove the add-on from")
	viper.BindPFlag("project", rootCmd.PersistentFlags().Lookup("project"))

	rootCmd.AddCommand(addons.Commands()...)
	rootCmd.AddCommand(auth.Commands()...)
	rootCmd.AddCommand(database.Commands()...)
	rootCmd.AddCommand(deploy.Commands()...)
	rootCmd.AddCommand(eventing.Commands()...)
	rootCmd.AddCommand(filestore.Commands()...)
	rootCmd.AddCommand(ingress.Commands()...)
	rootCmd.AddCommand(letsencrypt.Commands()...)
	rootCmd.AddCommand(operations.Commands()...)
	rootCmd.AddCommand(project.Commands()...)
	rootCmd.AddCommand(services.Commands()...)
	rootCmd.AddCommand(userman.Commands()...)
	rootCmd.AddCommand(modules.Commands()...)
	rootCmd.AddCommand(utils.Commands()...)
}

func er(msg interface{}) {
	fmt.Println("Error:", msg)
	os.Exit(1)
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			er(err)
		}

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".cobra")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func main() {

	err := execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// // Setup logrus
	// logrus.SetFormatter(&logrus.TextFormatter{})
	// logrus.SetOutput(os.Stdout)

	// app := cli.NewApp()
	// app.EnableBashCompletion = true
	// app.Name = "space-cli"
	// app.Version = "0.16.0"
	// app.Flags = []cli.Flag{
	// 	cli.StringFlag{Name: "log-level", Value: "info", Usage: "Sets the log level of the command", EnvVar: "LOG_LEVEL"},
	// }

	// app.Commands = append(app.Commands, addons.Commands...)
	// app.Commands = append(app.Commands, auth.Commands...)
	// app.Commands = append(app.Commands, database.Commands...)
	// app.Commands = append(app.Commands, deploy.CommandDeploy)
	// app.Commands = append(app.Commands, eventing.Commands...)
	// app.Commands = append(app.Commands, filestore.Commands...)
	// app.Commands = append(app.Commands, ingress.Commands...)
	// app.Commands = append(app.Commands, letsencrypt.Commands...)
	// app.Commands = append(app.Commands, operations.Commands...)
	// app.Commands = append(app.Commands, project.Commands...)
	// app.Commands = append(app.Commands, remoteservices.Commands...)
	// app.Commands = append(app.Commands, services.Commands...)
	// app.Commands = append(app.Commands, userman.Commands...)
	// app.Commands = append(app.Commands, modules.Commands...)
	// app.Commands = append(app.Commands, utils.LoginCommands...)

	// // Start the app
	// if err := app.Run(os.Args); err != nil {
	// 	logrus.Fatalln("Failed to run execute command:", err)
	// }
}

// func fetchAddSubCommands() []cli.Command {
// 	v := []cli.Command{}
// 	v = append(v, addons.AddSubCommands...)
// 	return v
// }

// func fetchRemoveSubCommands() []cli.Command {
// 	v := []cli.Command{}
// 	v = append(v, addons.RemoveSubCommand...)
// 	return v
// }

// func fetchGetSubCommands() []cli.Command {
// 	v := []cli.Command{}
// 	v = append(v, auth.GetSubCommands...)
// 	v = append(v, database.GetSubCommands...)
// 	v = append(v, eventing.GetSubCommands...)
// 	v = append(v, filestore.GetSubCommands...)
// 	v = append(v, letsencrypt.GetSubCommands...)
// 	v = append(v, project.GetSubCommands...)
// 	v = append(v, remoteservices.GetSubCommands...)
// 	v = append(v, services.GetSubCommands...)
// 	v = append(v, modules.GetSubCommands...)

// 	return v
// }

// func fetchGenerateSubCommands() []cli.Command {
// 	v := []cli.Command{}
// 	v = append(v, database.GenerateSubCommands...)
// 	v = append(v, eventing.GenerateSubCommands...)
// 	v = append(v, filestore.GenerateSubCommands...)
// 	v = append(v, ingress.GenerateSubCommands...)
// 	v = append(v, letsencrypt.GenerateSubCommands...)
// 	v = append(v, remoteservices.GenerateSubCommands...)
// 	v = append(v, services.GenerateSubCommands...)
// 	v = append(v, userman.GenerateSubCommands...)

// 	return v
// }
