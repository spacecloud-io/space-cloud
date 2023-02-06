package utils

import (
	"net/http"
	"net/url"
	"strings"
)

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
