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

type LicenseUpgradeRequest struct {
	LicenseKey   string `json:"licenseKey" mapstructure:"licenseKey"`
	LicenseValue string `json:"licenseValue" mapstructure:"licenseValue"`
	ClusterName  string `json:"clusterName" mapstructure:"clusterName"`
}
