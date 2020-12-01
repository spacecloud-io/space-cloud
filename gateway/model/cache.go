package model

import (
	"net/http"

	"github.com/spaceuptech/space-cloud/gateway/config"
)

// CachePurgeRequest describes the payload for cache purge request
type CachePurgeRequest struct {
	Resource  config.Resource `json:"resource,omitempty"`
	DbAlias   string          `json:"dbAlias,omitempty"`
	ServiceId string          `json:"serviceId,omitempty"`
	ID        string          `json:"id,omitempty"`
}

// CacheIngressRoute corresponds to a value of ingress route key
type CacheIngressRoute struct {
	Body    []byte      `json:"body"`
	Headers http.Header `json:"headers"`
}

type CacheDatabaseResult struct {
	Result      interface{} `json:"result"`
	MetricCount int64       `json:"metricCount"`
}
