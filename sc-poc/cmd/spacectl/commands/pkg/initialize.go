package pkg

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newCommandInitialize() *cobra.Command {
	wd, _ := os.Getwd()
	dirName := filepath.Base(wd)

	cmd := &cobra.Command{
		Use:     "init",
		Aliases: []string{"i"},
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.AutomaticEnv()
			viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

			_ = viper.BindPFlag("name", cmd.Flags().Lookup("name"))
			_ = viper.BindPFlag("lang", cmd.Flags().Lookup("lang"))
			_ = viper.BindPFlag("output-dir", cmd.Flags().Lookup("output-dir"))
			_ = viper.BindPFlag("resource-dir", cmd.Flags().Lookup("resource-dir"))
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			name := viper.GetString("name")
			language := viper.GetString("lang")
			outputDir := viper.GetString("output-dir")
			resourceDir := viper.GetString("resource-dir")

			if !cmd.Flags().Changed("name") {
				prompt := &survey.Input{
					Message: "Name of the package?",
					Default: dirName,
				}
				survey.AskOne(prompt, &name)
			}

			if !cmd.Flags().Changed("lang") {
				prompt := &survey.Input{
					Message: "Client gen language?",
					Default: "go",
				}
				survey.AskOne(prompt, &language)
			}

			if !cmd.Flags().Changed("output-dir") {
				prompt := &survey.Input{
					Message: "Directory of generated files?",
					Default: "sc/output",
				}
				survey.AskOne(prompt, &outputDir)
			}

			if !cmd.Flags().Changed("resource-dir") {
				prompt := &survey.Input{
					Message: "Directory of source files?",
					Default: "sc/resources",
				}
				survey.AskOne(prompt, &resourceDir)
			}

			_ = os.MkdirAll(outputDir, 0777)
			_ = os.MkdirAll(resourceDir, 0777)

			cfg := Config{
				Name: name,
				Output: Output{
					Language:  language,
					OutputDir: outputDir,
				},
				ResourceDir: resourceDir,
			}
			CreateConfig(cfg)
			return nil
		},
	}

	cmd.Flags().StringP("name", "n", dirName, "Name of the package")
	cmd.Flags().StringP("lang", "l", "go", "Client gen language")
	cmd.Flags().StringP("output-dir", "o", "sc/output", "Directory of generated files")
	cmd.Flags().StringP("resource-dir", "r", "sc/resources", "Directory of source files")

	return cmd
}
