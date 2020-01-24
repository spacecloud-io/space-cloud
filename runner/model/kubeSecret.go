package model

// Secret is an object for kubernetes! :|
type Secret struct {
	Name     string            `json:"name" yaml:"name"`
	Type     string            `json:"type" yaml:"type"`
	RootPath string            `json:"rootPath" yaml:"rootPath"`
	Data     map[string]string `json:"data" yaml:"data"`
}

// SecretValue is the non-encoded secret value in the request body
type SecretValue struct {
	Value string `json:"value" yaml:"value"`
}
