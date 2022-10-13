package utils

import (
	"context"
	"encoding/json"
	"net/http"
)

// SendOkayResponse sends an Okay http response
func SendOkayResponse(ctx context.Context, w http.ResponseWriter, statusCode int) error {
	return SendResponse(ctx, w, statusCode, map[string]string{})
}

// SendErrorResponse sends an Error http response
func SendErrorResponse(ctx context.Context, w http.ResponseWriter, statusCode int, err error) error {
	return SendResponse(ctx, w, statusCode, map[string]string{"error": err.Error()})
}

// SendResponse sends an http response
func SendResponse(ctx context.Context, w http.ResponseWriter, statusCode int, body interface{}) error {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(body)
}
