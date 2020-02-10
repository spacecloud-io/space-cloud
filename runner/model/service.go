package model

// Service describes a service's configurations
type Service struct {
	ID        string            `json:"id" yaml:"id"`
	Name      string            `json:"name" yaml:"name"`
	ProjectID string            `json:"projectId" yaml:"projectId"`
	Version   string            `json:"version" yaml:"version"`
	Scale     ScaleConfig       `json:"scale" yaml:"scale"`
	Labels    map[string]string `json:"labels" yaml:"labels"`
	Tasks     []Task            `json:"tasks" yaml:"tasks"`
	Affinity  []Affinity        `json:"affinity" yaml:"affinity"`
	Whitelist []Whitelist       `json:"whitelists" yaml:"whitelists"`
	Upstreams []Upstream        `json:"upstreams" yaml:"upstreams"`
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
	Secrets   []string          `json:"secrets" yaml:"secrets"`
	Runtime   Runtime           `json:"runtime" yaml:"runtime"`
}

// Port describes the port used by a task
type Port struct {
	Name     string   `json:"name" yaml:"name"`
	Protocol Protocol `json:"protocol" yaml:"protocol"`
	Port     int32    `json:"port" yaml:"port"`
}

// Protocol describes the protocol to be used for a port
type Protocol string

const (
	// TCP is used for tcp based workloads
	TCP Protocol = "tcp"

	// HTTP is used for http based workloads
	HTTP Protocol = "http"
)

// Resources describes the resources to be used by a task
type Resources struct {
	CPU    int64 `json:"cpu" yaml:"cpu"`
	Memory int64 `json:"memory" yaml:"memory"`
}

// Docker describes the docker configurations
type Docker struct {
	Image  string   `json:"image" yaml:"image"`
	Cmd    []string `json:"cmd" yaml:"cmd"`
	Secret string   `json:"secret" yaml:"secret"`
}

// Affinity describes the affinity rules of a service
type Affinity struct {
	Attribute string           `json:"attribute" yaml:"attribute"`
	Value     string           `json:"value" yaml:"value"`
	Weight    float32          `json:"weight" yaml:"weight"`
	Operator  AffinityOperator `json:"operator" yaml:"operator"`
}

// AffinityOperator describes the type of operator
type AffinityOperator string

const (
	// AffinityOperatorEQ is used for equal to operation
	AffinityOperatorEQ AffinityOperator = "=="

	// AffinityOperatorNEQ is used for not equal to operation
	AffinityOperatorNEQ AffinityOperator = "!="
)

// Upstream is the upstream dependencies of this service
type Upstream struct {
	ProjectID string `json:"projectId" yaml:"projectId"`
	Service   string `json:"service" yaml:"service"`
}

// Whitelist is the allowed downstream services
type Whitelist struct {
	ProjectID string `json:"projectId" yaml:"projectId"`
	Service   string `json:"service" yaml:"service"`
}

// Runtime is the container runtime
type Runtime string

const (
	// Image indicates that the user has provided a docker image
	Image Runtime = "image"

	// Code indicates that the user's code isn't containerized. We need to use a custom runtime for this
	Code Runtime = "code"
)

// Expose describes how an http service needs to be exposed
type Expose struct {
	Hosts []string     `json:"hosts" yaml:"hosts"`
	Rules []ExposeRule `json:"rules" yaml:"rules"`
}

// ExposeRule describes the rules for exposing an http service
type ExposeRule struct {
	URI  ExposeRuleURI `json:"uri" yaml:"uri"`
	Port int32         `json:"port" yaml:"port"`
}

// ExposeRuleURI describes the the http routes which need to be exposed
type ExposeRuleURI struct {
	Prefix  *string `json:"prefix" yaml:"prefix"`
	Exact   *string `json:"exact" yaml:"exact"`
	Rewrite *string `json:"rewrite" yaml:"rewrite"`
}
