package model

// Secret is an object for kubernetes! :|
type Secret struct {
	ID       string            `json:"id" yaml:"id"`
	Type     string            `json:"type" yaml:"type"`
	RootPath string            `json:"rootPath" yaml:"rootPath"`
	Data     map[string]string `json:"data" yaml:"data"`
}

// SecretValue is the non-encoded secret value in the request body
type SecretValue struct {
	Value string `json:"value" yaml:"value"`
}

// FileType is the secretType:file
const FileType string = "file"

// EnvType is the secretType:file
const EnvType string = "env"

// DockerType is the secretType:file
const DockerType string = "docker"
