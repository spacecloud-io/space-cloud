package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/spacecloud-io/space-cloud/cmd/spacectl/commands/client/generate/driver"
	"github.com/spacecloud-io/space-cloud/cmd/spacectl/commands/client/generate/driver/golang"
	"github.com/spacecloud-io/space-cloud/cmd/spacectl/commands/client/generate/driver/rtk"
	"github.com/spacecloud-io/space-cloud/cmd/spacectl/commands/client/generate/driver/typescript"
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
			_ = viper.BindPFlag("lang", cmd.Flags().Lookup("lang"))
			_ = viper.BindPFlag("package", cmd.Flags().Lookup("package"))
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			config := viper.GetString("config")
			output := viper.GetString("output")
			name := viper.GetString("name")
			lang := viper.GetString("lang")
			pkgName := viper.GetString("package")

			doc, err := openapi3.NewLoader().LoadFromFile(config)
			if err != nil {
				fmt.Println("Unable to openapi doc:", err)
				os.Exit(1)
			}

			var driver driver.Driver

			switch lang {
			case "rtk":
				driver = rtk.MakeRTKDriver(name)

				_ = os.MkdirAll(output, 0777)
				_ = os.WriteFile(filepath.Join(output, "helpers.ts"), []byte(rtk.HelperTS), 0777)
				_ = os.WriteFile(filepath.Join(output, "index.ts"), []byte(rtk.IndexTS), 0777)
				if _, err = os.Stat(filepath.Join(output, "http.config.ts")); os.IsNotExist(err) {
					_ = os.WriteFile(filepath.Join(output, "http.config.ts"), []byte(rtk.ConfigTS), 0777)
				}
			case "go":
				driver = golang.MakeGoDriver(pkgName)
			case "typescript":
				driver = typescript.MakeTSDriver()
			default:
				fmt.Printf("Invalid language name or language %s not supported.\n", lang)
				return nil
			}

			_ = os.MkdirAll(output, 0777)
			api, fileName, err := driver.GenerateAPIs(doc)
			if err != nil {
				fmt.Printf("error generating api: %s\n", err)
				os.Exit(1)
			}
			_ = os.WriteFile(filepath.Join(output, fileName), []byte(api), 0777)

			types, fileName, err := driver.GenerateTypes(doc)
			if err != nil {
				fmt.Printf("error generating types: %s\n", err)
				os.Exit(1)
			}
			_ = os.WriteFile(filepath.Join(output, fileName), []byte(types), 0777)

			return nil
		},
	}

	cmd.Flags().StringP("config", "c", "openapi.yaml", "The openapi yaml file to generate the client from.")
	cmd.Flags().StringP("output", "o", "client", "The directory to output the client to.")
	cmd.Flags().StringP("name", "n", "my-api", "Name for the API.")
	cmd.Flags().StringP("lang", "l", "go", "Language in which to generate code. Supported languages are 'rtk', 'Go'")
	cmd.Flags().StringP("package", "p", "openapi", "The name of the package to generate.")

	return cmd
}
