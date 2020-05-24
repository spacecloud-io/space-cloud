package utils

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/AlecAivazis/survey/v2"

	"github.com/spaceuptech/space-cli/cmd/model"
	"github.com/spaceuptech/space-cli/cmd/utils/file"
)

// Login logs the user in
func login(selectedAccount *model.Account) (*model.LoginResponse, error) {
	requestBody, err := json.Marshal(map[string]string{
		"user": selectedAccount.UserName,
		"key":  selectedAccount.Key,
	})
	if err != nil {
		_ = LogError(fmt.Sprintf("error in login unable to marshal data - %s", err.Error()), nil)
		return nil, err
	}

	resp, err := file.File.Post(fmt.Sprintf("%s/v1/config/login?cli=true", selectedAccount.ServerURL), "application/json", bytes.NewBuffer(requestBody))
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

// LoginStart take info of the user
func LoginStart(userName, ID, key, url string) error {
	if userName == "None" {
		if err := survey.AskOne(&survey.Input{Message: "Enter username:"}, &userName); err != nil {
			_ = LogError(fmt.Sprintf("error in login start unable to get username - %v", err), nil)
			return err
		}
	}

	if ID == "None" {
		ID = userName
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
	_, err := login(&account)
	if err != nil {
		_ = LogError(fmt.Sprintf("error in login start unable to login - %v", err), nil)
		return err
	}
	account = model.Account{
		ID:        ID,
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

// LoginWithSelectedAccount returns selected account & login token
func LoginWithSelectedAccount() (*model.Account, string, error) {
	account, err := getSelectedAccount()
	if err != nil {
		return nil, "", err
	}
	login, err := login(account)
	if err != nil {
		return nil, "", err
	}
	return account, login.Token, nil
}
