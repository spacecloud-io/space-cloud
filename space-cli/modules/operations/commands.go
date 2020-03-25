package operations

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spaceuptech/space-cli/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Commands is the list of commands the operations module exposes
func Commands() []*cobra.Command {
	var setup = &cobra.Command{
		Use:   "setup",
		Short: "setup development environment",
		RunE:  actionSetup,
	}
	var upgrade = &cobra.Command{
		Use:   "upgrade",
		Short: "Upgrade development environment",
		RunE:  actionUpgrade,
	}
	var destroy = &cobra.Command{
		Use:   "destroy",
		Short: "clean development environment & remove secrets",
		RunE:  actionDestroy,
	}
	var apply = &cobra.Command{
		Use:   "apply",
		Short: "deploys service",
		RunE:  actionApply,
	}
	var start = &cobra.Command{
		Use:   "start",
		Short: "Resumes the space-cloud docker environment",
		RunE:  actionStart,
	}

	var b bool
	var s []string
	var t []string
	setup.Flags().StringP("id", "", "", "The unique id for the cluster")
	viper.BindPFlag("id", setup.Flags().Lookup("id"))

	setup.Flags().StringP("username", "", "", "The username used for login")
	viper.BindPFlag("username", setup.Flags().Lookup("username"))

	setup.Flags().StringP("key", "", "", "The access key used for login")
	viper.BindPFlag("key", setup.Flags().Lookup("key"))

	setup.Flags().StringP("config", "", "", "The config used to bind config file")
	viper.BindPFlag("config", setup.Flags().Lookup("config"))

	setup.Flags().StringP("version", "", "", "The version is used to set SC version")
	viper.BindPFlag("version", setup.Flags().Lookup("version"))

	setup.Flags().StringP("secret", "", "", "The jwt secret to start space-cloud with")
	viper.BindPFlag("secret", setup.Flags().Lookup("secret"))

	setup.Flags().BoolP("dev", "", b, "Run space cloud in development mode")
	viper.BindPFlag("dev", setup.Flags().Lookup("dev"))

	setup.Flags().Int64P("port-http", "", 4122, "The port to use for HTTP")
	viper.BindPFlag("port-http", setup.Flags().Lookup("port-http"))

	setup.Flags().Int64P("port-https", "", 4126, "The port to use for HTTPS")
	viper.BindPFlag("port-https", setup.Flags().Lookup("port-https"))

	setup.Flags().StringSliceP("v", "", s, "Volumes to be attached to gateway")
	viper.BindPFlag("v", setup.Flags().Lookup("v"))

	setup.Flags().StringSliceP("e", "", t, "Environment variables to be provided to gateway")
	viper.BindPFlag("e", setup.Flags().Lookup("e"))

	command := make([]*cobra.Command, 0)
	command = append(command, setup)
	command = append(command, upgrade)
	command = append(command, destroy)
	command = append(command, apply)
	command = append(command, start)
	return command

}

// // Commands is the list of commands the operations module exposes
// //var Commands = []cli.Command{
// 	{
// 		Name:  "setup",
// 		Usage: "setup development environment",
// 		Flags: []cli.Flag{
// 			cli.StringFlag{
// 				Name:   "id",
// 				Usage:  "The unique id for the cluster",
// 				EnvVar: "CLUSTER_ID",
// 				Value:  "",
// 			},
// 			cli.StringFlag{
// 				Name:   "username",
// 				Usage:  "The username used for login",
// 				EnvVar: "USER_NAME", // don't set environment variable as USERNAME -> defaults to username of host machine in linux
// 				Value:  "",
// 			},
// 			cli.StringFlag{
// 				Name:   "key",
// 				Usage:  "The access key used for login",
// 				EnvVar: "KEY",
// 				Value:  "",
// 			},
// 			cli.StringFlag{
// 				Name:   "config",
// 				Usage:  "The config used to bind config file",
// 				EnvVar: "CONFIG",
// 				Value:  "",
// 			},
// 			cli.StringFlag{
// 				Name:   "version",
// 				Usage:  "The version is used to set SC version",
// 				EnvVar: "VERSION",
// 				Value:  "",
// 			},
// 			cli.StringFlag{
// 				Name:   "secret",
// 				Usage:  "The jwt secret to start space-cloud with",
// 				EnvVar: "JWT_SECRET",
// 				Value:  "",
// 			},
// 			cli.BoolFlag{
// 				Name:  "dev",
// 				Usage: "Run space cloud in development mode",
// 			},
// 			cli.Int64Flag{
// 				Name:   "port-http",
// 				Usage:  "The port to use for HTTP",
// 				EnvVar: "PORT_HTTP",
// 				Value:  4122,
// 			},
// 			cli.Int64Flag{
// 				Name:   "port-https",
// 				Usage:  "The port to use for HTTPS",
// 				EnvVar: "PORT_HTTPS",
// 				Value:  4126,
// 			},
// 			cli.StringSliceFlag{
// 				Name:  "v",
// 				Usage: "Volumes to be attached to gateway",
// 			},
// 			cli.StringSliceFlag{
// 				Name:  "e",
// 				Usage: "Environment variables to be provided to gateway",
// 			},
// 		},
// 		Action: actionSetup,
// 	},
// 	{
// 		Name:   "upgrade",
// 		Usage:  "Upgrade development environment",
// 		Action: actionUpgrade,
// 	},
// 	{
// 		Name:   "destroy",
// 		Usage:  "clean development environment & remove secrets",
// 		Action: actionDestroy,
// 	},
// 	{
// 		Name:   "apply",
// 		Usage:  "deploys service",
// 		Action: actionApply,
// 	},
// 	{
// 		Name:   "start",
// 		Usage:  "Resumes the space-cloud docker environment",
// 		Action: actionStart,
// 	},
// }

func actionSetup(cmd *cobra.Command, args []string) error {
	id := viper.GetString("id")
	userName := viper.GetString("username")
	key := viper.GetString("key")
	config := viper.GetString("config")
	version := viper.GetString("version")
	secret := viper.GetString("secret")
	local := viper.GetBool("dev")
	portHTTP := viper.GetInt64("port-http")
	portHTTPS := viper.GetInt64("port-https")
	volumes := viper.GetStringSlice("v")
	environmentVariables := viper.GetStringSlice("e")

	setLogLevel(viper.GetString("log-level"))

	return CodeSetup(id, userName, key, config, version, secret, local, portHTTP, portHTTPS, volumes, environmentVariables)
}

func actionUpgrade(cmd *cobra.Command, args []string) error {
	return Upgrade()
}

func actionDestroy(cmd *cobra.Command, args []string) error {
	return Destroy()
}

func actionApply(cmd *cobra.Command, args []string) error {
	//args := os.Args
	if len(args) != 3 {
		_ = utils.LogError("error while applying service incorrect number of arguments provided", nil)
		return fmt.Errorf("incorrect number of arguments provided")
	}

	fileName := args[2]

	return Apply(fileName)
}

func actionStart(cmd *cobra.Command, args []string) error {
	return DockerStart()
}

func setLogLevel(loglevel string) {
	switch loglevel {
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	default:
		_ = utils.LogError(fmt.Sprintf("Invalid log level (%s) provided", loglevel), nil)
		utils.LogInfo("Defaulting to `info` level")
		logrus.SetLevel(logrus.InfoLevel)
	}
}
