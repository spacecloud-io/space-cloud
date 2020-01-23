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
		logrus.Error("error in login unable to marshal data - %v", err)
		return nil, err
	}

	resp, err := http.Post(fmt.Sprintf("http://%s/v1/config/login", selectedAccount.ServerUrl), "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		logrus.Error("error in login unable to send http request - %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	loginResp := new(model.LoginResponse)
	err = json.NewDecoder(resp.Body).Decode(loginResp)

	if resp.StatusCode != 200 {
		logrus.Errorf("error in login got http status code %v with error message - %v", resp.StatusCode, loginResp.Error)
		return nil, fmt.Errorf("error in login got http status code %v with error message - %v", resp.StatusCode, loginResp.Error)
	}
	return loginResp, err
}

// LoginStart logs the user in galaxy
func LoginStart(userName, key, url string, local bool) error {
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
	selectedAccount := model.Account{
		UserName:  userName,
		Key:       key,
		ServerUrl: url, // todo server url is like localhost:4122
	}
	_, err := login(&selectedAccount)
	if err != nil {
		logrus.Errorf("error in login start unable to login - %v", err)
		return err
	}
	selectedAccount = model.Account{
		ID:        userName,
		UserName:  userName,
		Key:       key,
		ServerUrl: url,
	}
	if err := checkCred(&selectedAccount); err != nil {
		logrus.Errorf("error in login start unable to check credentials - %v", err)
		return err
	}
	fmt.Printf("Login Successful\n")
	return nil
}
