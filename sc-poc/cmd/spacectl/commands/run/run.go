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

			_ = viper.BindPFlag("caddy.log-level", cmd.Flags().Lookup("log-level"))
			_ = viper.BindPFlag("caddy.loading-interval", cmd.Flags().Lookup("loading-interval"))
			_ = viper.BindPFlag("caddy.port", cmd.Flags().Lookup("port"))

			_ = viper.BindPFlag("config.loader", cmd.Flags().Lookup("config-loader"))
			_ = viper.BindPFlag("config.path", cmd.Flags().Lookup("config-path"))

		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, _ := utils.LoadAdminConfig(true)
			if err := caddy.Run(c); err != nil {
				fmt.Println("Unable to start caddy:", err)
				os.Exit(1)
			}

			select {}
		},
	}

	// Caddy config
	cmd.Flags().StringP("log-level", "", "DEBUG", "Set the log level [DEBUG | INFO | WARN | ERROR | PANIC | FATAL]")
	cmd.Flags().StringP("loading-interval", "", "60s", "The interval to pull config")
	cmd.Flags().IntP("port", "", 4122, "The port to start SpaceCloud on")

	// Config loader
	cmd.Flags().StringP("config-loader", "", "file", "Set the configuration loader to be used [file | k8s]")
	cmd.Flags().StringP("config-path", "", "./sc-config", "Directory to use to manage SpaceCloud configuration")

	return cmd
}
