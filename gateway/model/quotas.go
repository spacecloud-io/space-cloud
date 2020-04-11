package model

// UsageQuotas describes the usage quotas
type UsageQuotas struct {
	MaxProjects  int `json:"maxProjects"`
	MaxDatabases int `json:"maxDatabases"`
	Version      int `json:"version,omitempty"`
}

// UsageQuotasResult describes the response received for fetch quota operation
type UsageQuotasResult struct {
	Result UsageQuotas `json:"result"`
	Error  string      `json:"error"`
}
