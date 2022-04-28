package run

import (
	"fmt"
	"os"
	"strings"

	"github.com/caddyserver/caddy/v2"
	"github.com/spacecloud-io/space-cloud/utils"
	"github.com/spaceuptech/helpers"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewCommand get space-cli run command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "run",
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.AutomaticEnv()
			viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
			_ = viper.BindPFlag("loading-interval", cmd.Flags().Lookup("loading-interval"))
			_ = viper.BindPFlag("store-type", cmd.Flags().Lookup("store-type"))
			_ = viper.BindPFlag("config-path", cmd.Flags().Lookup("config-path"))

			_ = viper.BindPFlag("dev", cmd.Flags().Lookup("dev"))
			_ = viper.BindPFlag("log-level", cmd.Flags().Lookup("log-level"))
			_ = viper.BindPFlag("cluster-id", cmd.Flags().Lookup("cluster-id"))
			_ = viper.BindPFlag("port", cmd.Flags().Lookup("port"))
			_ = viper.BindPFlag("ssl-cert", cmd.Flags().Lookup("ssl-cert"))
			_ = viper.BindPFlag("ssl-key", cmd.Flags().Lookup("ssl-key"))

			_ = viper.BindPFlag("admin-user", cmd.Flags().Lookup("admin-user"))
			_ = viper.BindPFlag("admin-pass", cmd.Flags().Lookup("admin-pass"))
			_ = viper.BindPFlag("admin-secret", cmd.Flags().Lookup("admin-secret"))

		},
		RunE: func(cmd *cobra.Command, args []string) error {
			helpers.InitLogger("debug", "text", true)
			c, err := utils.LoadAdminConfig(true)
			if err != nil {
				fmt.Println("Unable to load admin config:", err)
				os.Exit(1)
			}

			if err := caddy.Run(c); err != nil {
				fmt.Println("Unable to start caddy:", err)
				os.Exit(1)
			}

			select {}
		},
	}

	// Config store config
	cmd.Flags().StringP("loading-interval", "", "60s", "The interval to pull config")
	cmd.Flags().StringP("store-type", "", "file", "The config store to use for storing project configs and other meta data eg. file, kube, db")
	cmd.Flags().StringP("config-path", "", "./config.yaml", "The path to config file")

	// Server Config
	cmd.Flags().Bool("dev", false, "Run SpaceCloud in development mode")
	cmd.Flags().StringP("log-level", "", "DEBUG", "Set the log level [DEBUG | INFO | WARN | ERROR | PANIC | FATAL]")
	cmd.Flags().StringP("cluster-id", "", "", "The cluster id to start space-cloud with")
	cmd.Flags().IntP("port", "p", 4122, "Port to start space cloud server on")
	cmd.Flags().StringP("ssl-cert", "", "none", "Load ssl certificate from `FILE`")
	cmd.Flags().StringP("ssl-key", "", "none", "Load ssl key from `FILE`")

	// Config related to admin module
	cmd.Flags().StringP("admin-user", "", "admin", "Set the admin user name")
	cmd.Flags().StringP("admin-pass", "", "1234", "Set the admin password")
	cmd.Flags().StringP("admin-secret", "", "my-secretive-secret", "Set the admin jwt hmac secret")

	return cmd
}
