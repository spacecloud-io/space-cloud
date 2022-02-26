package model

import (
	"net/http"

	"github.com/spacecloud-io/space-cloud/config"
)

// CachePurgeRequest describes the payload for cache purge request
type CachePurgeRequest struct {
	Resource  config.Resource `json:"resource,omitempty"`
	DbAlias   string          `json:"dbAlias,omitempty"`
	ServiceID string          `json:"serviceId,omitempty"`
	ID        string          `json:"id,omitempty"`
}

// CacheIngressRoute corresponds to a value of ingress route key
type CacheIngressRoute struct {
	Body    []byte      `json:"body"`
	Headers http.Header `json:"headers"`
}

// CacheDatabaseResult is used to store cached database result
type CacheDatabaseResult struct {
	Result      interface{} `json:"result"`
	MetricCount int64       `json:"metricCount"`
}
