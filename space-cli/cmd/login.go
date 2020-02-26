package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AlecAivazis/survey/v2"
	"github.com/sirupsen/logrus"

	"github.com/spaceuptech/space-cli/model"
)

func login(selectedAccount *model.Account) (*model.LoginResponse, error) {
	requestBody, err := json.Marshal(map[string]string{
		"user": selectedAccount.UserName,
		"key":  selectedAccount.Key,
	})
	if err != nil {
		logrus.Errorf("error in login unable to marshal data - %s", err.Error())
		return nil, err
	}

	resp, err := http.Post(fmt.Sprintf("%s/v1/config/login?cli=true", selectedAccount.ServerURL), "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		logrus.Errorf("error in login unable to send http request - %s", err.Error())
		return nil, err
	}
	defer CloseTheCloser(resp.Body)

	loginResp := new(model.LoginResponse)
	err = json.NewDecoder(resp.Body).Decode(loginResp)

	if resp.StatusCode != 200 {
		logrus.Errorf("error in login got http status code %v with error message - %v", resp.StatusCode, loginResp.Error)
		return nil, fmt.Errorf("error in login got http status code %v with error message - %v", resp.StatusCode, loginResp.Error)
	}
	return loginResp, err
}

// LoginStart logs the user in space cloud
func LoginStart(userName, key, url string) error {
	if userName == "None" {
		if err := survey.AskOne(&survey.Input{Message: "Enter username:"}, &userName); err != nil {
			logrus.Errorf("error in login start unable to get username - %v", err)
			return err
		}
	}
	if key == "None" {
		if err := survey.AskOne(&survey.Password{Message: "Enter key:"}, &key); err != nil {
			logrus.Errorf("error in login start unable to get key - %v", err)
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
		logrus.Errorf("error in login start unable to login - %v", err)
		return err
	}
	account = model.Account{
		ID:        userName,
		UserName:  userName,
		Key:       key,
		ServerURL: url,
	}
	// write credentials into accounts.yaml file
	if err := checkCred(&account); err != nil {
		logrus.Errorf("error in login start unable to check credentials - %v", err)
		return err
	}
	fmt.Printf("Login Successful\n")
	return nil
}
