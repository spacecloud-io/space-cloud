package config

// Types for project config
type (
	// AdminProject is the config which the admin module uses
	AdminProject struct {
		Config  AdminProjectConfig `json:"config"`
		AesKey  AdminProjectAesKey `json:"aes"`
		Secrets []*Secret          `json:"secrets"`
	}

	// AdminProjectConfig stores adminman config
	AdminProjectConfig struct {
		ID                 string `json:"id" jsonschema:"required"`
		Name               string `json:"name"`
		DockerRegistry     string `json:"dockerRegistry"`
		ContextTimeGraphQL int    `json:"contextTimeGraphQL"`
	}

	// AdminProjectAesKey stores adminman project aes key
	AdminProjectAesKey struct {
		Key string `json:"key" jsonschema:"required"`
	}
)
