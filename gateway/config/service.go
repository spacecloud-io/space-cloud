package config

// Service describes a service's configurations
type RunnerService struct {
	ID          string            `json:"id" yaml:"id"`
	Name        string            `json:"name" yaml:"name"`
	ProjectID   string            `json:"projectId" yaml:"projectId"`
	Version     string            `json:"version" yaml:"version"`
	Environment string            `json:"env" yaml:"env"`
	Scale       ScaleConfig       `json:"scale" yaml:"scale"`
	Labels      map[string]string `json:"labels" yaml:"labels"`
	Tasks       []Task            `json:"tasks" yaml:"tasks"`
	Affinity    []Affinity        `json:"affinity" yaml:"affinity"`
	Whitelist   []string          `json:"whitelists" yaml:"whitelists"`
	Upstreams   []Upstream        `json:"upstreams" yaml:"upstreams"`
	Runtime     string            `json:"runtime" yaml:"runtime"`
}

// ScaleConfig describes the config used to scale a service
type ScaleConfig struct {
	Replicas    int32 `json:"replicas" yaml:"replicas"`
	MinReplicas int32 `json:"minReplicas" yaml:"minReplicas"`
	MaxReplicas int32 `json:"maxReplicas" yaml:"maxReplicas"`
	Concurrency int32 `json:"concurrency" yaml:"concurrency"`
}

// Task describes the configuration of a task
type Task struct {
	ID        string            `json:"id" yaml:"id"`
	Name      string            `json:"name" yaml:"name"`
	Ports     []Port            `json:"ports" yaml:"ports"`
	Resources Resources         `json:"resources" yaml:"resources"`
	Docker    Docker            `json:"docker" yaml:"docker"`
	Env       map[string]string `json:"env" yaml:"env"`
}

// Port describes the port used by a task
type Port struct {
	Name     string `json:"name" yaml:"name"`
	Protocol string `json:"protocol" yaml:"protocol"`
	Port     int32  `json:"port" yaml:"port"`
}

// Resources describes the resources to be used by a task
type Resources struct {
	CPU    int64 `json:"cpu" yaml:"cpu"`
	Memory int64 `json:"memory" yaml:"memory"`
}

// Docker describes the docker configurations
type Docker struct {
	Image string   `json:"image" yaml:"image"`
	Cmd   []string `json:"cmd" yaml:"cmd"`
}

// Affinity describes the affinity rules of a service
type Affinity struct {
	Attribute string  `json:"attribute" yaml:"attribute"`
	Value     string  `json:"value" yaml:"value"`
	Operator  string  `json:"operator" yaml:"operator"`
	Weight    float32 `json:"weight" yaml:"weight"`
}

// Upstream is the upstream dependencies of this service
type Upstream struct {
	ProjectID string `json:"projectId" yaml:"projectId"`
	Service   string `json:"service" yaml:"service"`
}
