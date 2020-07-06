package model

import (
	"net/http"
)

// RequestParams describes the params passed down in every request
type RequestParams struct {
	Resource, Op string
	Attributes   map[string]string
	Headers      http.Header
	Claims       map[string]interface{}
}
