package transport

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils"
	"github.com/stretchr/testify/mock"
)

type transport interface {
	Get(method, url string, params map[string]string, vPtr interface{}) error
	GetLogs(method, url string) error
}

type def struct{}

//Client todo
var Client transport

func init() {
	Client = &def{}
}

// Get gets spec object
func (d *def) Get(method, url string, params map[string]string, vPtr interface{}) error {
	account, token, err := utils.LoginWithSelectedAccount()
	if err != nil {
		return utils.LogError("Couldn't get account details or login token", err)
	}
	url = fmt.Sprintf("%s%s", account.ServerURL, url)

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return err
	}
	if token != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
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
		_ = utils.LogError(fmt.Sprintf("error while getting service got http status code %s - %s", resp.Status, respBody["error"]), nil)
		return fmt.Errorf("received invalid status code (%d)", resp.StatusCode)
	}

	if err := json.Unmarshal(data, vPtr); err != nil {
		return err
	}

	return nil
}

func (d *def) GetLogs(method, url string) error {
	account, token, err := utils.LoginWithSelectedAccount()
	if err != nil {
		return utils.LogError("Couldn't get account details or login token", err)
	}
	url = fmt.Sprintf("%s%s", account.ServerURL, url)

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return err
	}
	if token != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer CloseTheCloser(resp.Body)

	if resp.StatusCode != 200 {
		data, _ := ioutil.ReadAll(resp.Body)
		respBody := map[string]interface{}{}
		if err := json.Unmarshal(data, &respBody); err != nil {
			return err
		}
		_ = utils.LogError(fmt.Sprintf("error while getting service logs got error %s", respBody["error"]), nil)
		return fmt.Errorf("received invalid status code (%d)", resp.StatusCode)
	}

	rd := bufio.NewReader(resp.Body)

	utils.LogInfo("Press ctrl + c to exit")
	for {
		str, _ := rd.ReadString('\n')
		if str != "\n" {
			fmt.Println(str)
			time.Sleep(500 * time.Millisecond)
		}
	}
}

// CloseTheCloser closes the closer
func CloseTheCloser(c io.Closer) {
	_ = c.Close()
}

//MocketAuthProviders used during test
type MocketAuthProviders struct {
	mock.Mock
}

// Get gets spec object during test
func (m *MocketAuthProviders) Get(method, url string, params map[string]string, vPtr interface{}) error {
	c := m.Called(method, url, params, vPtr)
	a, _ := json.Marshal(c[1])
	_ = json.Unmarshal(a, vPtr)
	return c.Error(0)
}

// GetLogs gets logs of service during test
func (m *MocketAuthProviders) GetLogs(method, url string) error {
	c := m.Called(method, url)
	return c.Error(0)
}
