package utils

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
)

// HTTPMeta holds the http meta parameters
type HTTPMeta struct {
	ProjectID, Token string
}

// GetMetaInfo retrieves the http meta data from the request
func GetMetaInfo(r *http.Request) *HTTPMeta {
	// Get path parameters
	vars := mux.Vars(r)
	projectID := vars["projectID"]

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

// SendErrorResponse sends an error http response
func SendErrorResponse(w http.ResponseWriter, r *http.Request, status int, err error) {
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{"error": err.Error()}); err != nil {
		logrus.Errorf("Error while sending error response for %s %s - %s", r.Method, r.URL.String(), err.Error())
	}
}

// SendEmptySuccessResponse sends an empty http ok response
func SendEmptySuccessResponse(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{}); err != nil {
		logrus.Errorf("Error while sending error response for %s %s - %s", r.Method, r.URL.String(), err.Error())
	}
}

// CloseReaderCloser closes an io read closer while explicitly ignoring the error
func CloseReaderCloser(r io.Closer) {
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
