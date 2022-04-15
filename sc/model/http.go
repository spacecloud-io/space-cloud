package model

import "net/http"

// ErrorResponse describes the standard error response
type ErrorResponse struct {
	Error        string   `json:"error" jsonschema:"required"`
	SchemaErrors []string `json:"schemaErrors,omitempty"`
}

// HTTPParams describes the http params of the request
type HTTPParams struct {
	Headers http.Header `json:"headers"`
	Method  string      `json:"method"`
	Path    string      `json:"path"`
	Payload interface{} `json:"payload"`
}
