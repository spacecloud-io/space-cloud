package admin

import (
	"github.com/spacecloud-io/space-cloud/config"
	"github.com/spacecloud-io/space-cloud/utils/auth"
)

// Types for the global operations
type (
	loadEnvResponse struct {
		IsProd      bool   `json:"isProd" jsonschema:"required"`
		Version     string `json:"version" jsonschema:"required"`
		ClusterType string `json:"clusterType" jsonschema:"required"`
		LoginURL    string `json:"loginUrl" jsonschema:"required"`
	}

	tokenPayload struct {
		Token string `json:"token" jsonschema:"required"`
	}

	loginRequest struct {
		User string `json:"user" jsonschema:"required"`
		Key  string `json:"key" jsonschema:"required"`
	}
)

// Types for project config
type (
	// Project is the config which the admin module uses
	Project struct {
		Config  projectConfig    `json:"config"`
		AesKey  projectAesKey    `json:"aes"`
		Secrets []*config.Secret `json:"secrets"`

		// Auth module for the project
		auth *auth.Module `json:"-"`
	}

	projectConfig struct {
		ID                 string `json:"id" jsonschema:"required"`
		Name               string `json:"name"`
		DockerRegistry     string `json:"dockerRegistry"`
		ContextTimeGraphQL int    `json:"contextTimeGraphQL"`
	}

	projectAesKey struct {
		Key string `json:"key" jsonschema:"required"`
	}
)
