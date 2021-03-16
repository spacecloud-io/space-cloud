package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/rs/cors"
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
func MakeHTTPRequest(ctx context.Context, method, url, token, scToken string, params, vPtr interface{}) error {
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
	defer CloseTheCloser(resp.Body)

	if err := json.NewDecoder(resp.Body).Decode(vPtr); err != nil {
		return helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable decode response", err, nil)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New("service responded with status code " + strconv.Itoa(resp.StatusCode))
	}

	return nil
}

// GetToken retrieves the json web token present in the request
func GetToken(r *http.Request) (token string) {
	// Get the JWT token from header
	tokens, ok := r.Header["Authorization"]
	if !ok {
		tokens = []string{""}
	}

	return strings.TrimPrefix(tokens[0], "Bearer ")
}

// CreateCorsObject creates a cors object with the required config
func CreateCorsObject() *cors.Cors {
	return cors.New(cors.Options{
		AllowCredentials: true,
		AllowOriginFunc: func(s string) bool {
			return true
		},
		AllowedMethods: []string{"GET", "PUT", "POST", "DELETE"},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
		ExposedHeaders: []string{"Authorization", "Content-Type"},
	})
}

// CloseTheCloser closes the closer
func CloseTheCloser(c io.Closer) {
	_ = c.Close()
}
