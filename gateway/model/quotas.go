package model

// UsageQuotas describes the usage quotas
type UsageQuotas struct {
	MaxProjects  int `json:"maxProjects"`
	MaxDatabases int `json:"maxDatabases"`
	Version      int `json:"version,omitempty"`
}
