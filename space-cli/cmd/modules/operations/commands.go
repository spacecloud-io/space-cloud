package operations

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/spaceuptech/space-cli/cmd/utils"
)

// Commands is the list of commands the operations module exposes
func Commands() []*cobra.Command {
	var setup = &cobra.Command{
		Use:   "setup",
		Short: "setup development environment",
		PreRun: func(cmd *cobra.Command, args []string) {
			err := viper.BindPFlag("id", cmd.Flags().Lookup("id"))
			if err != nil {
				_ = utils.LogError("Unable to bind the flag ('id')", nil)
			}
			err = viper.BindPFlag("username", cmd.Flags().Lookup("username"))
			if err != nil {
				_ = utils.LogError("Unable to bind the flag ('username')", nil)
			}
			err = viper.BindPFlag("key", cmd.Flags().Lookup("key"))
			if err != nil {
				_ = utils.LogError("Unable to bind the flag ('key')", nil)
			}
			err = viper.BindPFlag("config", cmd.Flags().Lookup("config"))
			if err != nil {
				_ = utils.LogError("Unable to bind the flag ('config')", nil)
			}
			err = viper.BindPFlag("version", cmd.Flags().Lookup("version"))
			if err != nil {
				_ = utils.LogError("Unable to bind the flag ('version')", nil)
			}
			err = viper.BindPFlag("secret", cmd.Flags().Lookup("secret"))
			if err != nil {
				_ = utils.LogError("Unable to bind the flag ('secret')", nil)
			}
			err = viper.BindPFlag("dev", cmd.Flags().Lookup("dev"))
			if err != nil {
				_ = utils.LogError("Unable to bind the flag ('dev')", nil)
			}
			err = viper.BindPFlag("port-http", cmd.Flags().Lookup("port-http"))
			if err != nil {
				_ = utils.LogError("Unable to bind the flag ('port-http", nil)
			}
			err = viper.BindPFlag("port-https", cmd.Flags().Lookup("port-https"))
			if err != nil {
				_ = utils.LogError("Unable to bind the flag ('port-https", nil)
			}
			err = viper.BindPFlag("volume", cmd.Flags().Lookup("volume"))
			if err != nil {
				_ = utils.LogError("Unable to bind the flag ('v')", nil)
			}
			err = viper.BindPFlag("env", cmd.Flags().Lookup("env"))
			if err != nil {
				_ = utils.LogError("Unable to bind the flag ('e')", nil)
			}
		},
		RunE:          actionSetup,
		SilenceErrors: true,
	}

	setup.Flags().StringP("id", "", "", "The unique id for the cluster")
	err := viper.BindEnv("id", "CLUSTER_ID")
	if err != nil {
		_ = utils.LogError("Unable to bind lag ('id') to environment variables", nil)
	}

	setup.Flags().StringP("username", "", "", "The username used for login")
	err = viper.BindEnv("username", "USER_NAME")
	if err != nil {
		_ = utils.LogError("Unable to bind flag ('username') to environment variables", nil)
	}

	setup.Flags().StringP("key", "", "", "The access key used for login")
	err = viper.BindEnv("key", "KEY")
	if err != nil {
		_ = utils.LogError("Unable to bind flag ('key' to environment variables", nil)
	}

	setup.Flags().StringP("config", "", "", "The config used to bind config file")
	err = viper.BindEnv("config", "CONFIG")
	if err != nil {
		_ = utils.LogError("Unable to bind flag ('config') to environment variables", nil)
	}

	setup.Flags().StringP("version", "", "", "The version is used to set SC version")
	err = viper.BindEnv("version", "VERSION")
	if err != nil {
		_ = utils.LogError("Unable to bind flag ('version') to environment variables", nil)
	}

	setup.Flags().StringP("secret", "", "", "The jwt secret to start space-cloud with")
	err = viper.BindEnv("secret", "JWT_SECRET")
	if err != nil {
		_ = utils.LogError("Unable to bind flag ('secret') to environment variables", nil)
	}

	setup.Flags().BoolP("dev", "", false, "Run space cloud in development mode")

	setup.Flags().Int64P("port-http", "", 4122, "The port to use for HTTP")
	err = viper.BindEnv("port-http", "PORT_HTTP")
	if err != nil {
		_ = utils.LogError("Unable to bind flag ('port-http') to environment variables", nil)
	}

	setup.Flags().Int64P("port-https", "", 4126, "The port to use for HTTPS")
	err = viper.BindEnv("port-https", "PORT_HTTPS")
	if err != nil {
		_ = utils.LogError("Unable to bind flag ('port-https') to environment variables", nil)
	}

	setup.Flags().StringSliceP("volume", "v", []string{}, "Volumes to be attached to gateway")

	setup.Flags().StringSliceP("env", "e", []string{}, "Environment variables to be provided to gateway")

	var upgrade = &cobra.Command{
		Use:           "upgrade",
		Short:         "Upgrade development environment",
		RunE:          actionUpgrade,
		SilenceErrors: true,
	}
	var destroy = &cobra.Command{
		Use:           "destroy",
		Short:         "clean development environment & remove secrets",
		RunE:          actionDestroy,
		SilenceErrors: true,
	}
	var apply = &cobra.Command{
		Use:           "apply",
		Short:         "deploys service",
		RunE:          actionApply,
		SilenceErrors: true,
	}
	var start = &cobra.Command{
		Use:           "start",
		Short:         "Resumes the space-cloud docker environment",
		RunE:          actionStart,
		SilenceErrors: true,
	}
	var stop = &cobra.Command{
		Use:           "stop",
		Short:         "Stops the space-cloud docker environment",
		RunE:          actionStop,
		SilenceErrors: true,
	}

	return []*cobra.Command{setup, upgrade, destroy, apply, start, stop}

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
	volumes := viper.GetStringSlice("volume")
	environmentVariables := viper.GetStringSlice("env")

	return Setup(id, userName, key, config, version, secret, local, portHTTP, portHTTPS, volumes, environmentVariables)
}

func actionUpgrade(cmd *cobra.Command, args []string) error {
	return Upgrade()
}

func actionDestroy(cmd *cobra.Command, args []string) error {
	return Destroy()
}

func actionApply(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return utils.LogError("error while applying service incorrect number of arguments provided", nil)
	}

	dirName := args[0]
	return Apply(dirName)
}

func actionStart(cmd *cobra.Command, args []string) error {
	return DockerStart()
}

func actionStop(cmd *cobra.Command, args []string) error {
	return DockerStop()
}
