package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/rs/cors"
)

// HTTPRequest describes the request object
type HTTPRequest struct {
	Headers        headers
	Method, URL    string
	Token, SCToken string
	Params         interface{}
}

type headers interface {
	UpdateHeader(http.Header)
}

// MakeHTTPRequest fires an http request and returns a response
func MakeHTTPRequest(ctx context.Context, request *HTTPRequest, vPtr interface{}) (int, error) {
	// Marshal json into byte array
	data, _ := json.Marshal(request.Params)

	// Make a request object
	req, err := http.NewRequestWithContext(ctx, request.Method, request.URL, bytes.NewBuffer(data))
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Add the headers
	req.Header.Add("Content-Type", "application/json")

	// Add the token only if its provided
	if request.Token != "" {
		req.Header.Add("Authorization", "Bearer "+request.Token)
	}

	// Add the sc-token only if its provided
	if request.SCToken != "" {
		req.Header.Add("x-sc-token", "Bearer "+request.SCToken)
	}

	// Add the remaining headers
	request.Headers.UpdateHeader(req.Header)

	// Create a http client and fire the request
	client := &http.Client{}

	req = req.WithContext(ctx)
	resp, err := client.Do(req)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	defer CloseTheCloser(resp.Body)

	if resp.StatusCode != 204 {
		if err := json.NewDecoder(resp.Body).Decode(vPtr); err != nil {
			return resp.StatusCode, err
		}
	}

	return resp.StatusCode, nil
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
