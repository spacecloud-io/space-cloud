package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func Get(method, url string, params map[string]string, vPtr interface{}) error {
	account, err := getSelectedAccount()
	if err != nil {
		return err
	}
	login, err := login(account)
	if err != nil {
		return err
	}

	url = fmt.Sprintf("%s%s", account.ServerURL, url)
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return err
	}
	if login.Token != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", login.Token))
	}
	q := req.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	data, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(data, vPtr); err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("received invalid status code (%d)", resp.StatusCode)
	}

	return nil
}
