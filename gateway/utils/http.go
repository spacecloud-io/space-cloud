package utils

import (
	"encoding/json"
	"io"
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

// CloseTheCloser closes the closer
func CloseTheCloser(c io.Closer) {
	_ = c.Close()
}

// SendOkayResponse sends an Okay http response
func SendOkayResponse(w http.ResponseWriter, r *http.Request) error {
	return SendResponse(w, r, 200, map[string]string{})
}

// SendErrorResponse sends an Error http response
func SendErrorResponse(w http.ResponseWriter, r *http.Request, status int, err error) error {
	return SendResponse(w, r, status, map[string]string{"error": err.Error()})
}

// SendResponse sends an http response
func SendResponse(w http.ResponseWriter, r *http.Request, status int, body interface{}) error {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(body)
	return err
}
