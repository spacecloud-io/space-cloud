package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/stretchr/testify/mock"
)

type transport interface {
	Get(method, url string, params map[string]string, vPtr interface{}) error
}

type def struct{}

//Client todo
var Client transport

func init() {
	Client = &def{}
}

func (d *def) Get(method, url string, params map[string]string, vPtr interface{}) error {
	account, err := GetSelectedAccount()
	if err != nil {
		return err
	}
	login, err := Login(account)
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

	defer CloseTheCloser(resp.Body)

	data, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		respBody := map[string]interface{}{}
		if err := json.Unmarshal(data, &respBody); err != nil {
			return err
		}
		_ = LogError(fmt.Sprintf("error while getting service got http status code %s - %s", resp.Status, respBody["error"]), nil)
		return fmt.Errorf("received invalid status code (%d)", resp.StatusCode)
	}

	if err := json.Unmarshal(data, vPtr); err != nil {
		return err
	}

	return nil
}

type MocketAuthProviders struct {
	mock.Mock
}

func (m *MocketAuthProviders) Get(method, url string, params map[string]string, vPtr interface{}) error {
	c := m.Called(method, url, params, vPtr)
	a, _ := json.Marshal(c[1])
	_ = json.Unmarshal(a, vPtr)
	return c.Error(0)
}
