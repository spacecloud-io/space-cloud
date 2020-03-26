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

	setup.Flags().StringP("id", "", "", "The unique id for the cluster")
	err := viper.BindPFlag("id", setup.Flags().Lookup("id"))
	if err != nil {
		_ = utils.LogError(fmt.Sprintf("Unable to bind the flag ('id')"), nil)
	}
	err = viper.BindEnv("id", "CLUSTER_ID")
	if err != nil {
		_ = utils.LogError(fmt.Sprintf("Unable to bind lag ('id') to environment variables"), nil)
	}

	setup.Flags().StringP("username", "", "", "The username used for login")
	err = viper.BindPFlag("username", setup.Flags().Lookup("username"))
	if err != nil {
		_ = utils.LogError(fmt.Sprintf("Unable to bind the flag ('username')"), nil)
	}
	err = viper.BindEnv("username", "USER_NAME")
	if err != nil {
		_ = utils.LogError(fmt.Sprintf("Unable to bind flag ('username') to environment variables"), nil)
	}

	setup.Flags().StringP("key", "", "", "The access key used for login")
	err = viper.BindPFlag("key", setup.Flags().Lookup("key"))
	if err != nil {
		_ = utils.LogError(fmt.Sprintf("Unable to bind the flag ('key')"), nil)
	}
	err = viper.BindEnv("key", "KEY")
	if err != nil {
		_ = utils.LogError(fmt.Sprintf("Unable to bind flag ('key' to environment variables"), nil)
	}

	setup.Flags().StringP("config", "", "", "The config used to bind config file")
	err = viper.BindPFlag("config", setup.Flags().Lookup("config"))
	if err != nil {
		_ = utils.LogError(fmt.Sprintf("Unable to bind the flag ('config')"), nil)
	}
	err = viper.BindEnv("config", "CONFIG")
	if err != nil {
		_ = utils.LogError(fmt.Sprintf("Unable to bind flag ('config') to environment variables"), nil)
	}

	setup.Flags().StringP("version", "", "", "The version is used to set SC version")
	err = viper.BindPFlag("version", setup.Flags().Lookup("version"))
	if err != nil {
		_ = utils.LogError(fmt.Sprintf("Unable to bind the flag ('version')"), nil)
	}
	err = viper.BindEnv("version", "VERSION")
	if err != nil {
		_ = utils.LogError(fmt.Sprintf("Unable to bind flag ('version') to environment variables"), nil)
	}

	setup.Flags().StringP("secret", "", "", "The jwt secret to start space-cloud with")
	err = viper.BindPFlag("secret", setup.Flags().Lookup("secret"))
	if err != nil {
		_ = utils.LogError(fmt.Sprintf("Unable to bind the flag ('secret')"), nil)
	}
	err = viper.BindEnv("secret", "JWT_SECRET")
	if err != nil {
		_ = utils.LogError(fmt.Sprintf("Unable to bind flag ('secret') to environment variables"), nil)
	}

	setup.Flags().BoolP("dev", "", false, "Run space cloud in development mode")
	err = viper.BindPFlag("dev", setup.Flags().Lookup("dev"))
	if err != nil {
		_ = utils.LogError(fmt.Sprintf("Unable to bind the flag ('dev')"), nil)
	}

	setup.Flags().Int64P("port-http", "", 4122, "The port to use for HTTP")
	err = viper.BindPFlag("port-http", setup.Flags().Lookup("port-http"))
	if err != nil {
		_ = utils.LogError(fmt.Sprintf("Unable to bind the flag ('port-http')"), nil)
	}
	err = viper.BindEnv("port-http", "PORT_HTTP")
	if err != nil {
		_ = utils.LogError(fmt.Sprintf("Unable to bind flag ('port-http') to environment variables"), nil)
	}

	setup.Flags().Int64P("port-https", "", 4126, "The port to use for HTTPS")
	err = viper.BindPFlag("port-https", setup.Flags().Lookup("port-https"))
	if err != nil {
		_ = utils.LogError(fmt.Sprintf("Unable to bind the flag ('port-https')"), nil)
	}
	err = viper.BindEnv("port-https", "PORT_HTTPS")
	if err != nil {
		_ = utils.LogError(fmt.Sprintf("Unable to bind flag ('port-https') to environment variables"), nil)
	}

	setup.Flags().StringSliceP("v", "", []string{}, "Volumes to be attached to gateway")
	err = viper.BindPFlag("v", setup.Flags().Lookup("v"))
	if err != nil {
		_ = utils.LogError(fmt.Sprintf("Unable to bind the flag ('v')"), nil)
	}

	setup.Flags().StringSliceP("e", "", []string{}, "Environment variables to be provided to gateway")
	err = viper.BindPFlag("e", setup.Flags().Lookup("e"))
	if err != nil {
		_ = utils.LogError(fmt.Sprintf("Unable to bind the flag ('e')"), nil)
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

	return []*cobra.Command{setup, upgrade, destroy, apply, start}

}

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

	return CodeSetup(id, userName, key, config, version, secret, local, portHTTP, portHTTPS, volumes, environmentVariables)
}

func actionUpgrade(cmd *cobra.Command, args []string) error {
	return Upgrade()
}

func actionDestroy(cmd *cobra.Command, args []string) error {
	return Destroy()
}

func actionApply(cmd *cobra.Command, args []string) error {
	if len(args) != 3 {
		_ = utils.LogError("error while applying service incorrect number of arguments provided", nil)
		return fmt.Errorf("incorrect number of arguments provided")
	}

	dirName := args[2]
	return Apply(dirName)
}

func actionStart(cmd *cobra.Command, args []string) error {
	return DockerStart()
}

// SetLogLevel sets a single verbosity level for log messages.
func SetLogLevel(loglevel string) {
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
