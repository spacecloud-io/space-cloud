package model

// UsageQuotas describes the usage quotas
type UsageQuotas struct {
	MaxProjects  int `json:"maxProjects"`
	MaxDatabases int `json:"maxDatabases"`
}

// UsageQuotasResult describes the response received for fetch quota operation
type UsageQuotasResult struct {
	Result UsageQuotas `json:"result"`
	Error  string      `json:"error"`
}

type UpgradeResponse struct {
	Error      string `json:"error"`
	ClusterID  string `json:"clusterId"`
	ClusterKey string `json:"clusterKey"`
	License    string `json:"license"`
}

type GraphqlFetchLicenseResponse struct {
	Result struct {
		Result  *UpgradeResponse `json:"result"`
		Error   string           `json:"error"`
		Message string           `json:"message"`
		Status  int              `json:"status"`
	} `json:"result"`
	Error string `json:"error"`
}

type GraphqlFetchPublicKeyResponse struct {
	Result struct {
		Result  string `json:"result"`
		Error   string `json:"error"`
		Message string `json:"message"`
		Status  int    `json:"status"`
	} `json:"result"`
}

type License struct {
	ClusterID string       `json:"clusterId" mapstructure:"clusterId" structs:"clusterId"`
	SessionID string       `json:"sessionId" mapstructure:"sessionId" structs:"sessionId"`
	Plan      string       `json:"plan" mapstructure:"plan" structs:"plan"`
	Quotas    *UsageQuotas `json:"quotas" mapstructure:"quotas" structs:"quotas"`
}

type RenewLicense struct {
	ClusterID        string `json:"clusterId"`
	ClusterKey       string `json:"clusterKey"`
	License          string `json:"license"`
	CurrentSessionID string `json:"sessionId"`
}
