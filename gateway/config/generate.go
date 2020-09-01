package config

// GenerateEmptyConfig creates an empty config file
func GenerateEmptyConfig() *Config {

	return &Config{
		SSL: &SSL{Enabled: false},
		Admin: &Admin{
			ClusterConfig: &ClusterConfig{
				LetsEncryptEmail: "",
				EnableTelemetry:  true,
			},
		},
		Projects: []*Project{},
	}
}
