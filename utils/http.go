package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

// MakeHTTPRequest fires an http request and returns a response
func MakeHTTPRequest(ctx context.Context, method, url, token string, params, vPtr interface{}) error {
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

// GetTokenFromHeader returns the token from the request header
func GetTokenFromHeader(r *http.Request) string {
	// Get the JWT token from header
	tokens, ok := r.Header["Authorization"]
	if !ok {
		tokens = []string{""}
	}

	return strings.TrimPrefix(tokens[0], "Bearer ")
}
