package model

// Service describes a service's configurations
type Service struct {
	ID                     string            `json:"id,omitempty" yaml:"id,omitempty"`
	Name                   string            `json:"name,omitempty" yaml:"name,omitempty"`
	ProjectID              string            `json:"projectId,omitempty" yaml:"projectId,omitempty"`
	Version                string            `json:"version,omitempty" yaml:"version,omitempty"`
	Scale                  *ScaleConfig      `json:"scale,omitempty" yaml:"scale,omitempty"`
	AutoScale              *AutoScaleConfig  `json:"autoScale,omitempty" yaml:"autoScale,omitempty"`
	Labels                 map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	Tasks                  []Task            `json:"tasks" yaml:"tasks"`
	Affinity               []Affinity        `json:"affinity,omitempty" yaml:"affinity,omitempty"`
	Whitelist              []Whitelist       `json:"whitelists" yaml:"whitelists"`
	Upstreams              []Upstream        `json:"upstreams" yaml:"upstreams"`
	StatsInclusionPrefixes string            `json:"statsInclusionPrefixes" yaml:"statsInclusionPrefixes"`
}

// ScaleConfig describes the config used to scale a service
type ScaleConfig struct {
	Replicas    int32  `json:"replicas" yaml:"replicas"`
	MinReplicas int32  `json:"minReplicas" yaml:"minReplicas"`
	MaxReplicas int32  `json:"maxReplicas" yaml:"maxReplicas"`
	Concurrency int32  `json:"concurrency" yaml:"concurrency"`
	Mode        string `json:"mode" yaml:"mode"`
}

// AutoScaleConfig describes the config used to scale a service
type AutoScaleConfig struct {
	PollingInterval  int32 `json:"pollingInterval" yaml:"pollingInterval"`   // Default 15 (in seconds)
	CoolDownInterval int32 `json:"coolDownInterval" yaml:"coolDownInterval"` // Default 120 (in seconds)

	MinReplicas int32 `json:"minReplicas" yaml:"minReplicas"` // Default 1
	MaxReplicas int32 `json:"maxReplicas" yaml:"maxReplicas"` // Default 100

	Triggers []AutoScaleTrigger `json:"triggers,omitempty" yaml:"triggers,omitempty"`
}

// AutoScaleTrigger describes the config of a scaler
type AutoScaleTrigger struct {
	Name             string            `json:"name" yaml:"name"`
	Type             string            `json:"type" yaml:"type"`
	MetaData         map[string]string `json:"metadata" yaml:"metadata"`
	AuthenticatedRef *AutoScaleAuthRef `json:"authRef" yaml:"authRef"` // Authentication ref is optional
}

// AutoScaleAuthRef describes the authentication reference for a scaler
type AutoScaleAuthRef struct {
	SecretMapping []AutoScaleAuthRefMapping `json:"secretMapping" yaml:"secretMapping"`
}

// AutoScaleAuthRefMapping describes the mapping between the keys of the secret and parameters of the scaler
type AutoScaleAuthRefMapping struct {
	Parameter string `json:"param" yaml:"param"`
	Key       string `json:"key" yaml:"key"`
}

// Task describes the configuration of a task
type Task struct {
	ID        string            `json:"id" yaml:"id"`
	Name      string            `json:"name,omitempty" yaml:"name,omitempty"`
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
	CPU    int64        `json:"cpu" yaml:"cpu"`
	Memory int64        `json:"memory" yaml:"memory"`
	GPU    *GPUResource `json:"gpu,omitempty" yaml:"gpu,omitempty"`
}

// GPUResource describes the GPUs required by a task
type GPUResource struct {
	Type  string `json:"type" yaml:"type"`
	Value int64  `json:"value" yaml:"value"`
}

// Docker describes the docker configurations
type Docker struct {
	Image           string          `json:"image" yaml:"image"`
	Cmd             []string        `json:"cmd" yaml:"cmd"`
	Secret          string          `json:"secret" yaml:"secret"`
	ImagePullPolicy ImagePullPolicy `json:"imagePullPolicy" yaml:"imagePullPolicy"`
}

// ImagePullPolicy describes the image pull policy for docker config
type ImagePullPolicy string

const (
	// PullAlways is used for always pull policy
	PullAlways ImagePullPolicy = "always"

	// PullIfNotExists is use for pull if not exist locally pull policy
	PullIfNotExists ImagePullPolicy = "pull-if-not-exists"
)

// Affinity describes the affinity rules of a service
type Affinity struct {
	ID               string             `json:"id" yaml:"id"`
	Type             AffinityType       `json:"type" yaml:"type"`         // node or service
	Weight           int32              `json:"weight" yaml:"weight"`     // -100 to 100
	Operator         AffinityOperator   `json:"operator" yaml:"operator"` // preferred or required
	TopologyKey      string             `json:"topologyKey" yaml:"topologyKey"`
	Projects         []string           `json:"projects" yaml:"projects"`
	MatchExpressions []MatchExpressions `json:"matchExpressions" yaml:"matchExpressions"`
}

// MatchExpressions are a set of rules defined for scheduling pods on specifc nodes
type MatchExpressions struct {
	Key       string   `json:"key" yaml:"key"`
	Values    []string `json:"values" yaml:"values"`
	Attribute string   `json:"attribute" yaml:"attribute"`
	Operator  string   `json:"operator" yaml:"operator"`
}

// AffinityOperator describes the type of operator
type AffinityOperator string

const (
	// AffinityOperatorPreferred scheduling on particular node is preferred
	AffinityOperatorPreferred AffinityOperator = "preferred"

	// AffinityOperatorRequired scheduling on partiuclar node is required
	AffinityOperatorRequired AffinityOperator = "required"
)

// AffinityType describe the type of affinity it can be either node (corresponds to kubernetes nodes) or service (corresponds to kubernetes pod)
type AffinityType string

const (
	// AffinityTypeNode is used for kubernetes node resource
	AffinityTypeNode AffinityType = "node"

	// AffinityTypeService is used for kubernetes service resource
	AffinityTypeService AffinityType = "service"
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

// SpecObject describes the basic structure of config specifications
type SpecObject struct {
	API  string            `yaml:"api"`
	Type string            `yaml:"type"`
	Meta map[string]string `yaml:"meta"`
	Spec interface{}       `yaml:"spec,omitempty"`
}

// ServiceStatus describes structure of service status
type ServiceStatus struct {
	ServiceID       string         `json:"serviceId" yaml:"serviceId"`
	Version         string         `json:"version" yaml:"version"`
	DesiredReplicas interface{}    `json:"desiredReplicas" yaml:"desiredReplicas"`
	Replicas        []*ReplicaInfo `json:"replicas" yaml:"replicas"`
}

// ReplicaInfo describes structure of replica info
type ReplicaInfo struct {
	ID     string `json:"id" yaml:"id"`
	Status string `json:"status" yaml:"status"`
}
