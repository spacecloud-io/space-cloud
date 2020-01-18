package model

// Secret is an object for kubernetes! :|
type Secret struct {
	SecretName string            `json:"secret-name" yaml:"secret-name"`
	SecretType string            `json:"secret-type" yaml:"secret-type"`
	RoutePath  string            `json:"route-path" yaml:"route-path"`
	Data       map[string]string `json:"data" yaml:"data"`
}

// SecretValue is the non-encoded secret value in the request body
type SecretValue struct {
	Value string
}
