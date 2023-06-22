package login

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
			url := viper.GetString("url")
			username := viper.GetString("username")
			password := viper.GetString("password")

			if !cmd.Flags().Changed("url") {
				prompt := &survey.Input{
					Message: "URL of SpaceCloud?",
					Default: "http://localhost:4122",
				}
				survey.AskOne(prompt, &url)
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

			client := &http.Client{}
			path := fmt.Sprintf("%s/sc/v1/login", url)
			payload := []byte(fmt.Sprintf(`{"username": %q, "password": %q}`, username, password))

			req, err := http.NewRequest(http.MethodPost, path, bytes.NewBuffer(payload))
			if err != nil {
				log.Fatal("Failed to create HTTP request: ", err)
			}

			req.Header.Set("Content-Type", "application/json")
			resp, err := client.Do(req)
			if err != nil {
				log.Fatal("Failed to send HTTP request: ", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				body, _ := ioutil.ReadAll(resp.Body)
				m := make(map[string]string)
				json.Unmarshal(body, &m)
				log.Fatal("Failed to authenticate with SpaceCloud: ", m["error"])
			}

			location, err := UpdateSpaceCloudCredsFile(username, password, url)
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

func UpdateSpaceCloudCredsFile(username, password, url string) (string, error) {
	homeDir, _ := os.UserHomeDir()
	dirPath := filepath.Join(homeDir, ".space-cloud")
	_ = os.Mkdir(dirPath, 0777)

	output := map[string]string{
		"username": base64.StdEncoding.EncodeToString([]byte(username)),
		"password": base64.StdEncoding.EncodeToString([]byte(password)),
		"url":      url,
	}

	b, err := json.Marshal(output)
	if err != nil {
		return "", err
	}

	location := filepath.Join(dirPath, "creds.json")
	_ = os.WriteFile(location, b, 0777)
	return location, nil
}
