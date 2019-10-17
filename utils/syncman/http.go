package syncman

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"golang.org/x/net/context"

	"github.com/spaceuptech/space-cloud/utils"
)

// MakeHTTPRequest fires an http request and returns a response
func (s *Manager) MakeHTTPRequest(ctx context.Context, kind utils.RequestKind, method, url, token string, params, vPtr interface{}) error {
	// Marshal json into byte array
	data, _ := json.Marshal(params)

	// Make a request object
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	// Add token header
	req.Header.Add("Authorization", "Bearer "+token)

	// Create a http client and fire the request
	client := &http.Client{}

	if kind == utils.RequestKindConsulConnect {
		client = s.consulService.HTTPClient()
	}

	req = req.WithContext(ctx)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(vPtr); err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("service responded with status code " + strconv.Itoa(resp.StatusCode))
	}

	return nil
}
