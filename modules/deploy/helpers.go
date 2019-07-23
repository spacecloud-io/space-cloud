package deploy

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/spaceuptech/space-cloud/model"
)

type simpleAuth struct {
	ID  string `json:"id"`
	Key string `json:"key"`
}

type simpleAuthResponse struct {
	Token string `json:"token"`
	Error string `json:"error"`
}

func (m *Module) signIn() error {
	// Return if token exists
	if m.config.Registry.Token != nil {
		return nil
	}

	url := m.config.Registry.URL + "/v1/registry/auth/login"

	data, err := json.Marshal(&simpleAuth{ID: m.config.Registry.ID, Key: m.config.Registry.Key})
	if err != nil {
		return err
	}

	res, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// Parse the response
	obj := new(simpleAuthResponse)
	if err := json.NewDecoder(res.Body).Decode(&obj); err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return errors.New(obj.Error)
	}

	// Set the token
	m.config.Registry.Token = &obj.Token
	return nil
}

func (m *Module) upload(token string, r *http.Request) (*model.Deploy, error) {
	url := m.config.Registry.URL + "/v1/registry/deployment/upload"

	// Create a new http request
	req, err := http.NewRequest("POST", url, r.Body)
	if err != nil {
		return nil, err
	}

	// Set the http headers
	req.Header = make(http.Header)
	if contentType, p := r.Header["Content-Type"]; p {
		req.Header["Content-Type"] = contentType
	}

	// Add token header
	req.Header.Add("Authorization", "Bearer "+token)

	res, err := m.client.Do(req)
	if err != nil {
		return nil, err
	}

	// Load the body of the request
	uploadResponse := new(model.UploadResponse)
	if err := json.NewDecoder(res.Body).Decode(uploadResponse); err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.New(uploadResponse.Error)
	}

	return uploadResponse.Config, nil
}
