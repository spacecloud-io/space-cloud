package migrate

import (
	"fmt"
	"os"
	"strings"

	"github.com/spacecloud-io/space-cloud/model"
	"github.com/spacecloud-io/space-cloud/utils"
	"github.com/spaceuptech/helpers"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewCommand get space-cli migrate command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "migrate",
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.AutomaticEnv()
			viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
			_ = viper.BindPFlag("input-config-path", cmd.Flags().Lookup("input-config-path"))
			_ = viper.BindPFlag("output-config-path", cmd.Flags().Lookup("output-config-path"))
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = helpers.InitLogger("debug", "text", true)
			configPath := viper.GetString("input-config-path")
			outConfigPath := viper.GetString("output-config-path")

			resources, err := migrate(configPath)
			if err != nil {
				fmt.Println("Unable to migrate config:", err)
				os.Exit(1)
			}

			if err := utils.StoreConfigToFile(resources, outConfigPath); err != nil {
				fmt.Println("Unable to store config:", err)
				os.Exit(1)
			}

			return nil
		},
	}

	cmd.Flags().StringP("input-config-path", "", "", "The path to configs dir")
	cmd.Flags().StringP("output-config-path", "", "config.yaml", "The path to output new config file")

	return cmd
}

func migrate(configPath string) (*model.SCConfig, error) {
	resource := new(model.SCConfig)
	resource.Config = make(map[string]model.ConfigModule)

	if err := getAdminProjectConfig(resource, configPath); err != nil {
		return nil, err
	}

	if err := getDBConfig(resource, configPath); err != nil {
		return nil, err
	}

	if err := getDBRule(resource, configPath); err != nil {
		return nil, err
	}

	if err := getDBSchema(resource, configPath); err != nil {
		return nil, err
	}

	if err := getDBPreparedQuery(resource, configPath); err != nil {
		return nil, err
	}

	if err := getRemoteServices(resource, configPath); err != nil {
		return nil, err
	}

	return resource, nil
}
