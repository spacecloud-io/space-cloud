package openapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// Client which conforms to the OpenAPI3 specification for SpaceCloud.
type Client struct {
	// The endpoint of the server. All the paths in
	// the swagger spec will be appended to the server.
	Server string

	// Client for performing requests.
	Client *http.Client
}

// Creates a new SpaceCloud Client
func NewClient(server string) (*Client, error) {
	// create a client with sane default values
	client := Client{
		Server: server,
	}

	// create httpClient, if not already present
	if client.Client == nil {
		client.Client = &http.Client{}
	}
	return &client, nil
}

// GetAllTodos
func (c *Client) GetAllTodos(ctx context.Context, params GetAllTodosRequest) (*GetAllTodosResponse, error) {
	path := c.Server + "/v1/todos"

	url, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	queryValues := url.Query()

	b, err := json.Marshal(params.Eq)
	if err != nil {
		return nil, err
	}

	queryValues.Add("_eq", fmt.Sprint(string(b)))
	url.RawQuery = queryValues.Encode()

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	var obj GetAllTodosResponse
	err = json.NewDecoder(resp.Body).Decode(&obj)
	if err != nil {
		return nil, err
	}

	return &obj, nil
}
