package utils

import (
	"net/http"
	"strings"

	"github.com/rs/cors"
)

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
