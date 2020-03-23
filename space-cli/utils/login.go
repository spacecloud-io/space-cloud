package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AlecAivazis/survey/v2"
	"github.com/urfave/cli"

	"github.com/spaceuptech/space-cli/model"
)

// LoginCommands is the list of commands the utils module exposes
var LoginCommands = []cli.Command{
	{
		Name:  "login",
		Usage: "Logs into space cloud",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "username",
				Usage:  "Accepts the username for login",
				EnvVar: "USER_NAME", // don't set environment variable as USERNAME -> defaults to username of host machine in linux
				Value:  "None",
			},
			cli.StringFlag{
				Name:   "key",
				Usage:  "Accepts the access key to be verified during login",
				EnvVar: "KEY",
				Value:  "None",
			},
			cli.StringFlag{
				Name:   "url",
				Usage:  "Accepts the URL of server",
				EnvVar: "URL",
				Value:  "http://localhost:4122",
			},
		},
		Action: actionLogin,
	},
}

func actionLogin(c *cli.Context) error {
	userName := c.String("username")
	key := c.String("key")
	url := c.String("url")

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
	err = json.NewDecoder(resp.Body).Decode(loginResp)

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
