package login

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	clientutils "github.com/spacecloud-io/space-cloud/utils/client"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "login",
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.AutomaticEnv()
			viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

			_ = viper.BindPFlag("url", cmd.Flags().Lookup("url"))
			_ = viper.BindPFlag("username", cmd.Flags().Lookup("username"))
			_ = viper.BindPFlag("password", cmd.Flags().Lookup("password"))
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			baseUrl := viper.GetString("url")
			username := viper.GetString("username")
			password := viper.GetString("password")

			if !cmd.Flags().Changed("url") {
				prompt := &survey.Input{
					Message: "URL of SpaceCloud?",
					Default: "http://localhost:4122",
				}
				survey.AskOne(prompt, &baseUrl)
			}

			if !cmd.Flags().Changed("username") {
				prompt := &survey.Input{
					Message: "Username?",
					Default: "admin",
				}
				survey.AskOne(prompt, &username)
			}

			if !cmd.Flags().Changed("password") {
				prompt := &survey.Input{
					Message: "Password?",
					Default: "admin",
				}
				survey.AskOne(prompt, &password)
			}

			httpClient := &http.Client{}
			creds := clientutils.Credentials{
				Username: base64.StdEncoding.EncodeToString([]byte(username)),
				Password: base64.StdEncoding.EncodeToString([]byte(password)),
				BaseUrl:  baseUrl,
			}
			_, err := clientutils.Login(httpClient, creds)
			if err != nil {
				log.Fatal("Failed to authenticate with SpaceCloud", err)
			}

			location, err := clientutils.UpdateSpaceCloudCredsFile(creds)
			if err != nil {
				log.Fatal("Could not create creds.json file: ", err)
			}

			fmt.Printf("Successfully log into SpaceCloud. Credentials saved at %s\n", location)
			return nil
		},
	}

	cmd.Flags().StringP("url", "", "", "Base URL where SpaceCloud is running")
	cmd.Flags().StringP("username", "", "", "Username to log into SpaceCloud")
	cmd.Flags().StringP("password", "", "", "Password to log into SpaceCloud")

	return cmd
}
