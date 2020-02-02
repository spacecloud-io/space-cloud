package model

// ProxyMessage is the payload send by the proxy
type ProxyMessage struct {
	ActiveRequests int32  `json:"active,omitempty"`
	Project        string `json:"project,omitempty"`
	Service        string `json:"service,omitempty"`
	Environment    string `json:"env,omitempty"`
	NodeID         string `json:"id,omitempty"`
	Version        string `json:"version,omitempty"`
}

// EnvoyMetrics is the metrics collected from envoy
type EnvoyMetrics struct {
	Stats []EnvoyStat `json:"stats"`
}

// EnvoyStat describes the stats received from envoy
type EnvoyStat struct {
	Name  string `json:"name"`
	Value uint64 `json:"value"`
}
