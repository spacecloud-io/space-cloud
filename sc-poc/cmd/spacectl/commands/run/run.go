package run

import (
	"context"
	"log"
	"strings"

	"github.com/AlecAivazis/survey/v2"
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

			_ = viper.BindPFlag("id", cmd.Flags().Lookup("id"))

			_ = viper.BindPFlag("caddy.log-level", cmd.Flags().Lookup("log-level"))
			_ = viper.BindPFlag("caddy.port", cmd.Flags().Lookup("port"))

			_ = viper.BindPFlag("config.adapter", cmd.Flags().Lookup("config-adapter"))
			_ = viper.BindPFlag("config.path", cmd.Flags().Lookup("config-path"))
			_ = viper.BindPFlag("config.debounce-interval", cmd.Flags().Lookup("debounce-interval"))

			_ = viper.BindPFlag("admin.secret", cmd.Flags().Lookup("admin.secret"))
			_ = viper.BindPFlag("admin.username", cmd.Flags().Lookup("admin.username"))
			_ = viper.BindPFlag("admin.password", cmd.Flags().Lookup("admin.password"))

		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("admin.secret") {
				secret := ""
				prompt := &survey.Input{
					Message: "HSA256 Secret?",
					Default: "your-256-bit-secret",
				}
				survey.AskOne(prompt, &secret)
				viper.Set("admin.secret", secret)
			}

			if !cmd.Flags().Changed("admin.username") {
				username := ""
				prompt := &survey.Input{
					Message: "Username?",
					Default: "admin",
				}
				survey.AskOne(prompt, &username)
				viper.Set("admin.username", username)
			}

			if !cmd.Flags().Changed("admin.password") {
				password := ""
				prompt := &survey.Input{
					Message: "Password?",
					Default: "admin",
				}
				survey.AskOne(prompt, &password)
				viper.Set("admin.password", password)
			}

			if err := configman.InitializeConfigLoader(); err != nil {
				log.Fatal("Unable to initialize config loader: ", err)
			}

			c, err := configman.GetCaddyConfig()
			if err != nil {
				log.Fatal("Unable to load caddy config: ", err)
			}

			if err := caddy.Run(c); err != nil {
				log.Fatal("Unable to start caddy: ", err)
			}

			ctx := context.Background()
			go configman.WatchChanges(ctx)

			select {}
		},
	}

	cmd.Flags().String("id", "sc-id", "Set a unique id for this SpaceCloud instance")

	// Caddy config
	cmd.Flags().StringP("log-level", "", "DEBUG", "Set the log level [DEBUG | INFO | WARN | ERROR | PANIC | FATAL]")
	cmd.Flags().IntP("port", "", 4122, "The port to start SpaceCloud on")

	// Config loader
	cmd.Flags().StringP("config-adapter", "", "file", "Set the configuration loader to be used [file | k8s]")
	cmd.Flags().StringP("config-path", "", "./sc-config", "Directory to use to manage SpaceCloud configuration")
	cmd.Flags().StringP("debounce-interval", "", "500ms", "Debounce interval in milliseconds")

	// Auth
	cmd.Flags().StringP("admin.secret", "", "", "Set admin secret")
	cmd.Flags().StringP("admin.username", "", "", "Set admin username")
	cmd.Flags().StringP("admin.password", "", "", "Set admin password")

	return cmd
}
