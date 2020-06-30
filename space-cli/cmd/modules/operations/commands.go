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
			err := viper.BindPFlag("username", cmd.Flags().Lookup("username"))
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

			err = viper.BindPFlag("cluster-name", cmd.Flags().Lookup("cluster-name"))
			if err != nil {
				_ = utils.LogError("Unable to bind the flag ('cluster-name')", nil)
			}
		},
		RunE:          actionSetup,
		SilenceErrors: true,
	}

	setup.Flags().StringP("username", "", "", "The username used for login")
	err := viper.BindEnv("username", "USER_NAME")
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

	setup.Flags().StringP("cluster-name", "", "default", "The name of space-cloud cluster")
	err = viper.BindEnv("cluster-name", "CLUSTER_NAME")
	if err != nil {
		_ = utils.LogError("Unable to bind lag ('cluster-name') to environment variables", nil)
	}

	var upgrade = &cobra.Command{
		Use:   "upgrade",
		Short: "Upgrade development environment",
		PreRun: func(cmd *cobra.Command, args []string) {
			if err := viper.BindPFlag("cluster-name", cmd.Flags().Lookup("cluster-name")); err != nil {
				_ = utils.LogError("Unable to bind the flag ('cluster-name')", nil)
			}
		},
		SilenceErrors: true,
		RunE:          actionUpgrade,
	}
	upgrade.Flags().StringP("cluster-name", "", "default", "The name of space-cloud cluster")

	if err = viper.BindEnv("cluster-name", "CLUSTER_NAME"); err != nil {
		_ = utils.LogError("Unable to bind lag ('cluster-name') to environment variables", nil)
	}

	var destroy = &cobra.Command{
		Use:   "destroy",
		Short: "clean development environment & remove secrets",
		PreRun: func(cmd *cobra.Command, args []string) {
			if err := viper.BindPFlag("cluster-name", cmd.Flags().Lookup("cluster-name")); err != nil {
				_ = utils.LogError("Unable to bind the flag ('cluster-name')", nil)
			}
		},
		SilenceErrors: true,
		RunE:          actionDestroy,
	}
	destroy.Flags().StringP("cluster-name", "", "default", "The name of  space-cloud cluster")
	if err = viper.BindEnv("cluster-name", "CLUSTER_NAME"); err != nil {
		_ = utils.LogError("Unable to bind lag ('cluster-name') to environment variables", nil)
	}

	var apply = &cobra.Command{
		Use:           "apply",
		Short:         "deploys service",
		RunE:          actionApply,
		SilenceErrors: true,
	}

	var start = &cobra.Command{
		Use:   "start",
		Short: "Resumes the space-cloud docker environment",
		PreRun: func(cmd *cobra.Command, args []string) {
			if err := viper.BindPFlag("cluster-name", cmd.Flags().Lookup("cluster-name")); err != nil {
				_ = utils.LogError("Unable to bind the flag ('cluster-name')", nil)
			}
		},
		SilenceErrors: true,
		RunE:          actionStart,
	}
	start.Flags().StringP("cluster-name", "", "default", "The name of space-cloud cluster")
	if err = viper.BindEnv("cluster-name", "CLUSTER_NAME"); err != nil {
		_ = utils.LogError("Unable to bind lag ('cluster-name') to environment variables", nil)
	}

	var stop = &cobra.Command{
		Use:   "stop",
		Short: "Stops the space-cloud docker environment",
		PreRun: func(cmd *cobra.Command, args []string) {
			if err := viper.BindPFlag("cluster-name", cmd.Flags().Lookup("cluster-name")); err != nil {
				_ = utils.LogError("Unable to bind the flag ('cluster-name')", nil)
			}
		},
		SilenceErrors: true,
		RunE:          actionStop,
	}
	stop.Flags().StringP("cluster-name", "", "default", "The name of space-cloud cluster")
	if err = viper.BindEnv("cluster-name", "CLUSTER_NAME"); err != nil {
		_ = utils.LogError("Unable to bind lag ('cluster-name') to environment variables", nil)
	}
	return []*cobra.Command{setup, upgrade, destroy, apply, start, stop}

}

func actionSetup(cmd *cobra.Command, args []string) error {
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
	clusterName := viper.GetString("cluster-name")

	return Setup(userName, key, config, version, secret, clusterName, local, portHTTP, portHTTPS, volumes, environmentVariables)
}

func actionUpgrade(cmd *cobra.Command, args []string) error {
	clusterName := viper.GetString("cluster-name")
	return Upgrade(clusterName)
}

func actionDestroy(cmd *cobra.Command, args []string) error {
	clusterName := viper.GetString("cluster-name")
	return Destroy(clusterName)
}

func actionApply(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return utils.LogError("error while applying service incorrect number of arguments provided", nil)
	}

	dirName := args[0]
	return Apply(dirName)
}

func actionStart(cmd *cobra.Command, args []string) error {
	clusterName := viper.GetString("cluster-name")
	return DockerStart(clusterName)
}

func actionStop(cmd *cobra.Command, args []string) error {
	clusterName := viper.GetString("cluster-name")
	return DockerStop(clusterName)
}
