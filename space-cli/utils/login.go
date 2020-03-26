package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/spaceuptech/space-cli/model"
)

// Commands is the list of commands the utils module exposes
func Commands() []*cobra.Command {
	var loginCommands = &cobra.Command{
		Use:   "login",
		Short: "Logs into space cloud",
		RunE:  actionLogin,
	}
	loginCommands.Flags().StringP("username", "", "None", "Accepts the username for login")
	err := viper.BindPFlag("username", loginCommands.Flags().Lookup("username"))
	if err != nil {
		_ = LogError(fmt.Sprintf("Unable to bind the flag ('username')"), nil)
	}
	err = viper.BindEnv("username", "USER_NAME")
	if err != nil {
		_ = LogError(fmt.Sprintf("Unable to bind flag ('username') to environment variables"), nil)
	}

	loginCommands.Flags().StringP("key", "", "None", "Accepts the access key to be verified during login")
	err = viper.BindPFlag("key", loginCommands.Flags().Lookup("key"))
	if err != nil {
		_ = LogError(fmt.Sprintf("Unable to bind the flag ('key')"), nil)
	}
	err = viper.BindEnv("key", "KEY")
	if err != nil {
		_ = LogError(fmt.Sprintf("Unable to bind flag ('key') to environment variables"), nil)
	}

	loginCommands.Flags().StringP("url", "", "http://localhost:4122", "Accepts the URL of server")
	err = viper.BindPFlag("url", loginCommands.Flags().Lookup("url"))
	if err != nil {
		_ = LogError(fmt.Sprintf("Unable to bind the flag ('url')"), nil)
	}
	err = viper.BindEnv("url", "URL")
	if err != nil {
		_ = LogError(fmt.Sprintf("Unable to bind flag ('url') to environment variables"), nil)
	}

	return []*cobra.Command{loginCommands}
}

func actionLogin(cmd *cobra.Command, args []string) error {
	userName := viper.GetString("username")
	key := viper.GetString("key")
	url := viper.GetString("url")

	return loginStart(userName, key, url)
}

// Login logs the user in
func Login(selectedAccount *model.Account) (*model.LoginResponse, error) {
	requestBody, err := json.Marshal(map[string]string{
		"user": selectedAccount.UserName,
		"key":  selectedAccount.Key,
	})
	if err != nil {
		_ = LogError(fmt.Sprintf("error in login unable to marshal data - %s", err.Error()), nil)
		return nil, err
	}

	resp, err := http.Post(fmt.Sprintf("%s/v1/config/login?cli=true", selectedAccount.ServerURL), "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		_ = LogError(fmt.Sprintf("error in login unable to send http request - %s", err.Error()), nil)
		return nil, err
	}
	defer CloseTheCloser(resp.Body)

	loginResp := new(model.LoginResponse)
	_ = json.NewDecoder(resp.Body).Decode(loginResp)

	if resp.StatusCode != 200 {
		_ = LogError(fmt.Sprintf("error in login got http status code %v with error message - %v", resp.StatusCode, loginResp.Error), nil)
		return nil, fmt.Errorf("error in login got http status code %v with error message - %v", resp.StatusCode, loginResp.Error)
	}
	return loginResp, err
}

func loginStart(userName, key, url string) error {
	if userName == "None" {
		if err := survey.AskOne(&survey.Input{Message: "Enter username:"}, &userName); err != nil {
			_ = LogError(fmt.Sprintf("error in login start unable to get username - %v", err), nil)
			return err
		}
	}
	if key == "None" {
		if err := survey.AskOne(&survey.Password{Message: "Enter key:"}, &key); err != nil {
			_ = LogError(fmt.Sprintf("error in login start unable to get key - %v", err), nil)
			return err
		}
	}
	account := model.Account{
		UserName:  userName,
		Key:       key,
		ServerURL: url,
	}
	_, err := Login(&account)
	if err != nil {
		_ = LogError(fmt.Sprintf("error in login start unable to login - %v", err), nil)
		return err
	}
	account = model.Account{
		ID:        userName,
		UserName:  userName,
		Key:       key,
		ServerURL: url,
	}
	// write credentials into accounts.yaml file
	if err := StoreCredentials(&account); err != nil {
		_ = LogError(fmt.Sprintf("error in login start unable to check credentials - %v", err), nil)
		return err
	}
	fmt.Printf("Login Successful\n")
	return nil
}
