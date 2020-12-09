package model

import (
	"net/http"
)

// RequestParams describes the params passed down in every request
type RequestParams struct {
	RequestID  string                 `json:"requestId"`
	Resource   string                 `json:"resource"`
	Op         string                 `json:"op"`
	Attributes map[string]string      `json:"attributes"`
	Headers    http.Header            `json:"headers"`
	Claims     map[string]interface{} `json:"claims"`
	Method     string                 `json:"method"`
	Path       string                 `json:"path"`
	Payload    interface{}            `json:"payload"`
}

// SpecObject describes the basic structure of config specifications
type SpecObject struct {
	API  string            `json:"api" yaml:"api"`
	Type string            `json:"type" yaml:"type"`
	Meta map[string]string `json:"meta" yaml:"meta"`
	Spec interface{}       `json:"spec" yaml:"spec,omitempty"`
}

type BatchSpecApplyRequest struct {
	Specs []*SpecObject `json:"specs" yaml:"specs"`
}
