package main

import (
	"fmt"
	"os"

	"github.com/spaceuptech/space-cli/modules/addons"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "space-cli",
	Version: "0.16.0",
}

func execute() error {
	return rootCmd.Execute()
}

func init() {
	//cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringP("log-level", "", "info", "Sets the log level of the command")
	addons.AddCmd.PersistentFlags().StringP("project", "", "", "The project to add the add-on to")
	addons.RemoveCmd.PersistentFlags().StringP("project", "", "", "The project to remove the add-on from")

	rootCmd.AddCommand(addons.AddCmd)
	rootCmd.AddCommand(addons.RemoveCmd)

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
