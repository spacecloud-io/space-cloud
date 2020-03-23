package operations

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spaceuptech/space-cli/utils"
	"github.com/urfave/cli"
)

// Commands is the list of commands the operations module exposes
var Commands = []cli.Command{
	{
		Name:  "setup",
		Usage: "setup development environment",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "id",
				Usage:  "The unique id for the cluster",
				EnvVar: "CLUSTER_ID",
				Value:  "",
			},
			cli.StringFlag{
				Name:   "username",
				Usage:  "The username used for login",
				EnvVar: "USER_NAME", // don't set environment variable as USERNAME -> defaults to username of host machine in linux
				Value:  "",
			},
			cli.StringFlag{
				Name:   "key",
				Usage:  "The access key used for login",
				EnvVar: "KEY",
				Value:  "",
			},
			cli.StringFlag{
				Name:   "config",
				Usage:  "The config used to bind config file",
				EnvVar: "CONFIG",
				Value:  "",
			},
			cli.StringFlag{
				Name:   "version",
				Usage:  "The version is used to set SC version",
				EnvVar: "VERSION",
				Value:  "",
			},
			cli.StringFlag{
				Name:   "secret",
				Usage:  "The jwt secret to start space-cloud with",
				EnvVar: "JWT_SECRET",
				Value:  "",
			},
			cli.BoolFlag{
				Name:  "dev",
				Usage: "Run space cloud in development mode",
			},
			cli.Int64Flag{
				Name:   "port-http",
				Usage:  "The port to use for HTTP",
				EnvVar: "PORT_HTTP",
				Value:  4122,
			},
			cli.Int64Flag{
				Name:   "port-https",
				Usage:  "The port to use for HTTPS",
				EnvVar: "PORT_HTTPS",
				Value:  4126,
			},
			cli.StringSliceFlag{
				Name:  "v",
				Usage: "Volumes to be attached to gateway",
			},
			cli.StringSliceFlag{
				Name:  "e",
				Usage: "Environment variables to be provided to gateway",
			},
		},
		Action: actionSetup,
	},
	{
		Name:   "upgrade",
		Usage:  "Upgrade development environment",
		Action: actionUpgrade,
	},
	{
		Name:   "destroy",
		Usage:  "clean development environment & remove secrets",
		Action: actionDestroy,
	},
	{
		Name:   "apply",
		Usage:  "deploys service",
		Action: actionApply,
	},
	{
		Name:   "start",
		Usage:  "Resumes the space-cloud docker environment",
		Action: actionStart,
	},
}

func actionSetup(c *cli.Context) error {
	id := c.String("id")
	userName := c.String("username")
	key := c.String("key")
	config := c.String("config")
	version := c.String("version")
	secret := c.String("secret")
	local := c.Bool("dev")
	portHTTP := c.Int64("port-http")
	portHTTPS := c.Int64("port-https")
	volumes := c.StringSlice("v")
	environmentVariables := c.StringSlice("e")

	setLogLevel(c.GlobalString("log-level"))

	return CodeSetup(id, userName, key, config, version, secret, local, portHTTP, portHTTPS, volumes, environmentVariables)
}

func actionUpgrade(_ *cli.Context) error {
	return Upgrade()
}

func actionDestroy(_ *cli.Context) error {
	return Destroy()
}

func actionApply(cli *cli.Context) error {
	args := os.Args
	if len(args) != 3 {
		_ = utils.LogError("error while applying service incorrect number of arguments provided", nil)
		return fmt.Errorf("incorrect number of arguments provided")
	}

	fileName := args[2]

	return Apply(fileName)
}

func actionStart(_ *cli.Context) error {
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
		utils.LogInfo(fmt.Sprintf("Defaulting to `info` level"))
		logrus.SetLevel(logrus.InfoLevel)
	}
}
