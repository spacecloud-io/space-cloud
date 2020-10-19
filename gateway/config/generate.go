package config

// GenerateEmptyConfig creates an empty config file
func GenerateEmptyConfig() *Config {
	return &Config{
		Projects:         make(Projects),
		SSL:              &SSL{Enabled: false},
		ClusterConfig:    new(ClusterConfig),
		Integrations:     make(Integrations),
		IntegrationHooks: make(IntegrationHooks),
		License:          new(License),
	}
}

// GenerateEmptyProject creates a empty project
func GenerateEmptyProject(project *ProjectConfig) *Project {
	return &Project{
		ProjectConfig:           project,
		DatabaseConfigs:         make(map[string]*DatabaseConfig),
		DatabaseSchemas:         make(map[string]*DatabaseSchema),
		DatabaseRules:           make(map[string]*DatabaseRule),
		DatabasePreparedQueries: make(map[string]*DatbasePreparedQuery),
		EventingConfig:          new(EventingConfig),
		EventingSchemas:         make(map[string]*EventingSchema),
		EventingRules:           make(map[string]*Rule),
		EventingTriggers:        make(map[string]*EventingTrigger),
		FileStoreConfig:         new(FileStoreConfig),
		FileStoreRules:          FileStoreRules{},
		Auths:                   make(Auths),
		LetsEncrypt:             new(LetsEncrypt),
		IngressRoutes:           make(IngressRoutes),
		IngressGlobal:           new(GlobalRoutesConfig),
		RemoteService:           make(Services),
	}
}
