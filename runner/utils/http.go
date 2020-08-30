package utils

import (
	"io"
	"net/http"
	"strings"

	"github.com/rs/cors"
)

// GetToken retrieves the json web token present in the request
func GetToken(r *http.Request) (token string) {
	// Get the JWT token from header
	tokens, ok := r.Header["Authorization"]
	if !ok {
		tokens = []string{""}
	}

	return strings.TrimPrefix(tokens[0], "Bearer ")
}

// CloseTheCloser closes an io read closer while explicitly ignoring the error
func CloseTheCloser(r io.Closer) {
	_ = r.Close()
}

// CreateCorsObject returns a new cors object
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
