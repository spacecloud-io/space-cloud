package utils

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// HTTPMeta holds the http meta parameters
type HTTPMeta struct {
	ProjectID, Token string
}

// GetMetaInfo retrieves the http meta data from the request
func GetMetaInfo(r *http.Request) *HTTPMeta {
	// Get path parameters
	vars := mux.Vars(r)
	projectID := vars["project"]

	// Get the JWT token from header
	token := GetToken(r)

	// Return the http meta object
	return &HTTPMeta{ProjectID: projectID, Token: token}
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
