package model

// UsageQuotas describes the usage quotas
type UsageQuotas struct {
	MaxProjects      int     `json:"maxProjects"`
	MaxDatabases     int     `json:"maxDatabases"`
	IntegrationLevel float64 `json:"integrationLevel"`
}

// UsageQuotasResult describes the response received for fetch quota operation
type UsageQuotasResult struct {
	Result UsageQuotas `json:"result"`
	Error  string      `json:"error"`
}

// UpgradeResponse is the response of upgrade request
type UpgradeResponse struct {
	Error   string `json:"error"`
	License string `json:"license"`
}

// GraphqlFetchLicenseResponse is the response fetch license request
type GraphqlFetchLicenseResponse struct {
	Result  *UpgradeResponse `json:"result"`
	Error   string           `json:"error"`
	Message string           `json:"message"`
	Status  int              `json:"status"`
}

// GraphqlFetchPublicKeyResponse is the response fetch public key request
type GraphqlFetchPublicKeyResponse struct {
	Result  string `json:"result"`
	Error   string `json:"error"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

// License stores license info
type License struct {
	LicenseKey     string            `json:"licenseKey" mapstructure:"licenseKey" structs:"licenseKey"`
	LicenseRenewal string            `json:"licenseRenewal" mapstructure:"licenseRenewal" structs:"licenseRenewal"`
	SessionID      string            `json:"sessionId" mapstructure:"sessionId" structs:"sessionId"`
	Plan           string            `json:"plan" mapstructure:"plan" structs:"plan"`
	Meta           *LicenseTokenMeta `json:"meta" mapstructure:"meta" structs:"meta"`
}

// LicenseTokenMeta stores license token meta info
type LicenseTokenMeta struct {
	ProductMeta    *UsageQuotas    `json:"productMeta" mapstructure:"productMeta" structs:"productMeta"`
	LicenseKeyMeta *LicenseKeyMeta `json:"licenseKeyMeta" mapstructure:"licenseKeyMeta" structs:"licenseKeyMeta"`
}

// LicenseKeyMeta stores license key meta
type LicenseKeyMeta struct {
	ClusterName string `json:"clusterName" mapstructure:"clusterName" structs:"clusterName"`
}

// RenewLicense is the body for renew license request
type RenewLicense struct {
	ClusterName      string `json:"clusterName"`
	LicenseValue     string `json:"licenseValue"`
	LicenseKey       string `json:"licenseKey"`
	License          string `json:"license"`
	CurrentSessionID string `json:"sessionId"`
}
