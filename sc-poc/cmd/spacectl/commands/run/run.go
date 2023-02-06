package run

import (
	"context"
	"log"
	"strings"

	"github.com/caddyserver/caddy/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/spacecloud-io/space-cloud/managers/configman"
)

// NewCommand get spacectl run command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "run",
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.AutomaticEnv()
			viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

			_ = viper.BindPFlag("caddy.log-level", cmd.Flags().Lookup("log-level"))
			_ = viper.BindPFlag("caddy.port", cmd.Flags().Lookup("port"))

			_ = viper.BindPFlag("config.adapter", cmd.Flags().Lookup("config-adapter"))
			_ = viper.BindPFlag("config.path", cmd.Flags().Lookup("config-path"))
			_ = viper.BindPFlag("config.debounce-interval", cmd.Flags().Lookup("debounce-interval"))

		},
		RunE: func(cmd *cobra.Command, args []string) error {
			configloader, err := configman.InitializeConfigLoader()
			if err != nil {
				log.Fatal("Unable to initialize config loader: ", err)
			}

			c, err := configloader.GetCaddyConfig()
			if err != nil {
				log.Fatal("Unable to load caddy config: ", err)
			}

			if err := caddy.Run(c); err != nil {
				log.Fatal("Unable to start caddy: ", err)
			}

			ctx := context.Background()
			go configloader.WatchChanges(ctx)

			select {}
		},
	}

	// Caddy config
	cmd.Flags().StringP("log-level", "", "DEBUG", "Set the log level [DEBUG | INFO | WARN | ERROR | PANIC | FATAL]")
	cmd.Flags().IntP("port", "", 4122, "The port to start SpaceCloud on")

	// Config loader
	cmd.Flags().StringP("config-adapter", "", "file", "Set the configuration loader to be used [file | k8s]")
	cmd.Flags().StringP("config-path", "", "./sc-config", "Directory to use to manage SpaceCloud configuration")
	cmd.Flags().StringP("debounce-interval", "", "500ms", "Debounce interval in milliseconds")

	return cmd
}
