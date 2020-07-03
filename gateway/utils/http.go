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
)

// HTTPRequest describes the request object
type HTTPRequest struct {
	Headers        map[string]string
	Method, URL    string
	Token, SCToken string
	Params         interface{}
}

// MakeHTTPRequest fires an http request and returns a response
func MakeHTTPRequest(ctx context.Context, request *HTTPRequest, vPtr interface{}) error {
	// Marshal json into byte array
	data, _ := json.Marshal(request.Params)

	// Make a request object
	req, err := http.NewRequestWithContext(ctx, request.Method, request.URL, bytes.NewBuffer(data))
	if err != nil {
		return err
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
	for k, v := range request.Headers {
		req.Header.Set(k, v)
	}

	// Create a http client and fire the request
	client := &http.Client{}

	req = req.WithContext(ctx)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer CloseTheCloser(resp.Body)

	if resp.StatusCode != 204 {
		if err := json.NewDecoder(resp.Body).Decode(vPtr); err != nil {
			return err
		}
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New("service responded with status code " + strconv.Itoa(resp.StatusCode))
	}

	return nil
}

// GetTokenFromHeader returns the token from the request header
func GetTokenFromHeader(r *http.Request) string {
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

// SendOkayResponse sends an Okay http response
func SendOkayResponse(w http.ResponseWriter) error {
	return SendResponse(w, 200, map[string]string{})
}

// SendErrorResponse sends an Error http response
func SendErrorResponse(w http.ResponseWriter, status int, message string) error {
	return SendResponse(w, status, map[string]string{"error": message})
}

// SendResponse sends an http response
func SendResponse(w http.ResponseWriter, status int, body interface{}) error {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(body)
}
