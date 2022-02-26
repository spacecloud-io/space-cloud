package commands

import (
	"github.com/caddyserver/caddy/v2"
	"github.com/spf13/cobra"

	"github.com/spacecloud-io/space-cloud/utils"
)

func GetRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:          "space-cli",
		Version:      "v0.22.0",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			c := utils.LoadAdminConfig(true)
			if err := caddy.Run(c); err != nil {
				return err
			}

			select {}
		},
	}

	return rootCmd
}
