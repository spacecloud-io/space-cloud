package model

// Deploy holds the config for the deployment
type Deploy struct {
	Name        string            `json:"name" yaml:"name"`
	Project     string            `json:"project" yaml:"project"`
	Runtime     *Runtime          `json:"runtime" yaml:"runtime"`
	Constraints *Constraints      `json:"constraints" yaml:"constraints"`
	Ports       []*Port           `json:"ports,omitempty" yaml:"ports,omitempty"`
	Env         map[string]string `json:"env" yaml:"env"`
}

// Runtime holds the runtime information
type Runtime struct {
	Name    string `json:"name" yaml:"name"`
	Install string `json:"install" yaml:"install"`
	Run     string `json:"run" yaml:"run"`
}

// Constraints holds the constraints information
type Constraints struct {
	Replicas int32    `json:"replicas" yaml:"replicas"`
	CPU      *float32 `json:"cpu,omitempty" yaml:"cpu,omitempty"`
	Memory   *int64   `json:"memory,omitempty" yaml:"memory,omitempty"`
}

// Port holds the container ports information
type Port struct {
	Port     int     `json:"port" yaml:"port"`
	Name     *string `json:"name,omitempty" yaml:"name,omitempty"`
	Protocol *string `json:"protocol,omitempty" yaml:"protocol,omitempty"`
}

// UploadResponse is the response received from the registry
type UploadResponse struct {
	Error  string  `json:"error"`
	Config *Deploy `json:"config,omitempty"`
}
