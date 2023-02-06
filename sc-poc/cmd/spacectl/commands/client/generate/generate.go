package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewCommand get spacectl client generate command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "generate",
		Aliases: []string{"gen", "g"},
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.AutomaticEnv()
			viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

			_ = viper.BindPFlag("config", cmd.Flags().Lookup("config"))
			_ = viper.BindPFlag("output", cmd.Flags().Lookup("output"))
			_ = viper.BindPFlag("name", cmd.Flags().Lookup("name"))
			_ = viper.BindPFlag("base-url", cmd.Flags().Lookup("base-url"))
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			config := viper.GetString("config")
			output := viper.GetString("output")
			name := viper.GetString("name")
			baseURL := viper.GetString("base-url")

			doc, err := openapi3.NewLoader().LoadFromFile(config)
			if err != nil {
				fmt.Println("Unable to openapi doc:", err)
				os.Exit(1)
			}

			_ = os.MkdirAll(output, 0777)
			_ = os.WriteFile(filepath.Join(output, "types.ts"), []byte(generateTypes(doc)), 0777)
			_ = os.WriteFile(filepath.Join(output, "helpers.ts"), []byte(helperTS), 0777)
			_ = os.WriteFile(filepath.Join(output, "api.ts"), []byte(generateAPI(name, baseURL, doc)), 0777)
			_ = os.WriteFile(filepath.Join(output, "index.ts"), []byte(indexTS), 0777)
			return nil
		},
	}

	cmd.Flags().StringP("config", "c", "openapi.yaml", "The openapi yaml file to generate the client from.")
	cmd.Flags().StringP("output", "o", "client", "The directory to output the client to.")
	cmd.Flags().StringP("name", "n", "my-api", "Name for the API.")
	cmd.Flags().StringP("base-url", "b", "/", "Base path for the url.")

	return cmd
}
