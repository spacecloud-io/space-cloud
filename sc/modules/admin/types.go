package admin

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
