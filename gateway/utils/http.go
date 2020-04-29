package utils

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
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

// SendErrorResponse sends an error http response
func SendErrorResponse(w http.ResponseWriter, r *http.Request, status int, err error) {
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{"error": err.Error()}); err != nil {
		logrus.Errorf("Error while sending error response for %s %s - %s", r.Method, r.URL.String(), err.Error())
	}
}
