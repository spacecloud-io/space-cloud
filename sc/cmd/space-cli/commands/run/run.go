package run

import (
	"fmt"
	"os"
	"strings"

	"github.com/caddyserver/caddy/v2"
	"github.com/spacecloud-io/space-cloud/utils"
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
			// Blob store
			_ = viper.BindPFlag("loading-interval", cmd.Flags().Lookup("loading-interval"))
			_ = viper.BindPFlag("log-level", cmd.Flags().Lookup("log-level"))
			_ = viper.BindPFlag("store-type", cmd.Flags().Lookup("store-type"))
			_ = viper.BindPFlag("config-path", cmd.Flags().Lookup("config-path"))
			_ = viper.BindPFlag("cluster-id", cmd.Flags().Lookup("cluster-id"))
			_ = viper.BindPFlag("port", cmd.Flags().Lookup("port"))
			_ = viper.BindPFlag("ssl-cert", cmd.Flags().Lookup("ssl-cert"))
			_ = viper.BindPFlag("ssl-key", cmd.Flags().Lookup("ssl-key"))

		},
		RunE: func(cmd *cobra.Command, args []string) error {
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

	cmd.Flags().StringP("loading-interval", "", "60s", "The interval to pull config")
	cmd.Flags().StringP("log-level", "", "DEBUG", "Set the log level [DEBUG | INFO | WARN | ERROR | PANIC | FATAL]")
	cmd.Flags().StringP("store-type", "", "file", "The config store to use for storing project configs and other meta data eg. file, kube, db")
	cmd.Flags().StringP("config-path", "", "", "The path to config file")
	cmd.Flags().StringP("cluster-id", "", "", "The cluster id to start space-cloud with")
	cmd.Flags().IntP("port", "p", 4122, "run xlr8s server")
	cmd.Flags().StringP("ssl-cert", "", "none", "Load ssl certificate from `FILE`")
	cmd.Flags().StringP("ssl-key", "", "none", "Load ssl key from `FILE`")

	return cmd
}
