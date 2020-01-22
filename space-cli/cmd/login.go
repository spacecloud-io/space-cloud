package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/AlecAivazis/survey/v2"

	"github.com/spaceuptech/space-cli/model"
)

func login(selectedAccount *model.Account) (*model.LoginResponse, error) {
	requestBody, err := json.Marshal(map[string]string{
		"username": selectedAccount.UserName,
		"key":      selectedAccount.Key,
	})
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(fmt.Sprintf("%s/v1/galaxy/login", selectedAccount.ServerUrl), "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	loginResp := new(model.LoginResponse)
	if err := json.Unmarshal(body, loginResp); err != nil {
		return nil, err
	}
	return loginResp, nil

}

// LoginStart logs the user in galaxy
func LoginStart(userName, key, url string, local bool) error {
	if userName == "None" {
		if err := survey.AskOne(&survey.Input{Message: "Enter username"}, &userName); err != nil {
			return err
		}
	}
	if key == "None" {
		if err := survey.AskOne(&survey.Input{Message: "Enter key"}, &key); err != nil {
			return err
		}
	}
	selectedAccount := model.Account{
		UserName:  userName,
		Key:       key,
		ServerUrl: url,
	}
	loginRes, err := login(&selectedAccount)
	if err != nil {
		return err
	}
	selectedAccount = model.Account{
		ID:        loginRes.Token,
		UserName:  userName,
		Key:       key,
		ServerUrl: url,
	}
	if err := checkCred(&selectedAccount); err != nil {
		return err
	}
	return nil
}
