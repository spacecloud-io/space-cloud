package syncman

import "github.com/spaceuptech/space-cloud/gateway/config"

// AdminSyncmanInterface is an interface consisting of functions of admin module used by eventing module
type AdminSyncmanInterface interface {
	GetInternalAccessToken() (string, error)
	IsTokenValid(token string) error
	ValidateSyncOperation(c *config.Config, project *config.Project) bool
	SetConfig(admin *config.Admin)
}

type preparedQueryResponse struct {
	ID        string   `json:"id"`
	SQL       string   `json:"sql"`
	Arguments []string `json:"arguments" yaml:"arguments"`
}

type dbRulesResponse struct {
	IsRealTimeEnabled bool                    `json:"isRealtimeEnabled"`
	Rules             map[string]*config.Rule `json:"rules"`
}

type dbSchemaResponse struct {
	Schema string `json:"schema"`
}
