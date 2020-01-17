package syncman

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"golang.org/x/net/context"

	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// MakeHTTPRequest fires an http request and returns a response
func (s *Manager) MakeHTTPRequest(ctx context.Context, method, url, token, scToken string, params, vPtr interface{}) error {
	// Marshal json into byte array
	data, _ := json.Marshal(params)

	// Make a request object
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	// Add the headers
	if token != "" {
		// Add the token only if its provided
		req.Header.Add("Authorization", "Bearer "+token)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-sc-token", "Bearer "+scToken)

	// Create a http client and fire the request
	client := &http.Client{}

	// if s.storeType && s.isConsulConnectEnabled && strings.Contains(url, "https") && strings.Contains(url, ".consul") {
	// 	 client = s.consulService.HTTPClient()
	// }

	req = req.WithContext(ctx)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer utils.CloseTheCloser(resp.Body)

	if err := json.NewDecoder(resp.Body).Decode(vPtr); err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New("service responded with status code " + strconv.Itoa(resp.StatusCode))
	}

	return nil
}
