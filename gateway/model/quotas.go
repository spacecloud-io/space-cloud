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
	Error   string `json:"error"`
	License string `json:"license"`
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
	LicenseKey     string            `json:"licenseKey" mapstructure:"licenseKey" structs:"licenseKey"`
	LicenseRenewal string            `json:"licenseRenewal" mapstructure:"licenseRenewal" structs:"licenseRenewal"`
	SessionID      string            `json:"sessionId" mapstructure:"sessionId" structs:"sessionId"`
	Plan           string            `json:"plan" mapstructure:"plan" structs:"plan"`
	Meta           *LicenseTokenMeta `json:"meta" mapstructure:"meta" structs:"meta"`
}
type LicenseTokenMeta struct {
	ProductMeta    *UsageQuotas    `json:"productMeta" mapstructure:"productMeta" structs:"productMeta"`
	LicenseKeyMeta *LicenseKeyMeta `json:"licenseKeyMeta" mapstructure:"licenseKeyMeta" structs:"licenseKeyMeta"`
}

type LicenseKeyMeta struct {
	ClusterName string `json:"clusterName" mapstructure:"clusterName" structs:"clusterName"`
}

type RenewLicense struct {
	ClusterName      string `json:"clusterName"`
	LicenseValue     string `json:"licenseValue"`
	LicenseKey       string `json:"licenseKey"`
	License          string `json:"license"`
	CurrentSessionID string `json:"sessionId"`
}
