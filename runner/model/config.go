package model

import (
	"net/http"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Config is a map with key = projectId-serviceId and the value being the routes([]Route)
type Config map[string]Routes // key = projectId-serviceId

// Response is the object returned by every handler to client
type Response struct {
	Error  string      `json:"error,omitempty"`
	Result interface{} `json:"result,omitempty"`
}

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

// LogRequest represent log request structure
type LogRequest struct {
	TaskID    string       `json:"taskId"`
	ReplicaID string       `json:"replicaId"`
	Since     *int64       `json:"since"`
	SinceTime *metav1.Time `json:"sinceTime"`
	Tail      *int64       `json:"tail"`
	IsFollow  bool         `json:"isFollow"`
}
