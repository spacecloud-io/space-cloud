package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/spaceuptech/helpers"
)

// HTTPRequest describes the request object
type HTTPRequest struct {
	Headers        headers
	Method, URL    string
	Token, SCToken string
	Params         io.Reader
}

type headers interface {
	UpdateHeader(http.Header)
}

// MakeHTTPRequest fires an http request and returns a response
func MakeHTTPRequest(ctx context.Context, request *HTTPRequest, vPtr interface{}) (int, error) {
	// Make a request object
	if request.Method == http.MethodGet || request.Method == http.MethodDelete {
		request.Params = nil
	}

	req, err := http.NewRequestWithContext(ctx, request.Method, request.URL, request.Params)
	if err != nil {
		return http.StatusInternalServerError, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to create http request for url (%s)", request.URL), err, nil)
	}

	// Add the token only if its provided
	if request.Token != "" {
		req.Header.Add("Authorization", "Bearer "+request.Token)
	}

	// Add the sc-token only if its provided
	if request.SCToken != "" {
		req.Header.Add("x-sc-token", "Bearer "+request.SCToken)
	}

	// Add the remaining headers
	if request.Headers != nil {
		request.Headers.UpdateHeader(req.Header)
	}

	// Create a http client and fire the request
	client := &http.Client{}

	req = req.WithContext(ctx)
	resp, err := client.Do(req)
	if err != nil {
		return http.StatusInternalServerError, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to make http request for url (%s)", request.URL), err, nil)
	}
	defer CloseTheCloser(resp.Body)

	if resp.StatusCode != 204 {
		if err := json.NewDecoder(resp.Body).Decode(vPtr); err != nil {
			return resp.StatusCode, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to json unmarshal response of http request url (%s)", request.URL), err, nil)
		}
	}

	return resp.StatusCode, nil
}

// GetQueryParams generates query params map
func GetQueryParams(queryParams url.Values) map[string]string {
	params := make(map[string]string)

	for key, val := range queryParams {
		params[key] = val[0]
	}

	return params
}

// GetTokenFromHeader returns the token from the request header
func GetTokenFromHeader(r *http.Request) string {
	// Get the JWT token from header
	tokens, ok := r.Header["Authorization"]
	if !ok {
		tokens = []string{""}
	}

	arr := strings.Split(tokens[0], " ")
	if strings.ToLower(arr[0]) == "bearer" {
		return arr[1]
	}

	return ""
}

// CloseTheCloser closes the closer
func CloseTheCloser(c io.Closer) {
	_ = c.Close()
}
