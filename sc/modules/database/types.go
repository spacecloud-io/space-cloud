package database

import "github.com/spacecloud-io/space-cloud/config"

// Config describes the configuration required by a single database
type Config struct {
	Connector       *config.DatabaseConfig         `json:"connector"`
	Schemas         config.DatabaseSchemas         `json:"schemas"`
	PreparedQueries config.DatabasePreparedQueries `json:"preparedQueries"`
}
