package model

// Role describes the configuration for the service role
type Role struct {
	ID      string `json:"id" yaml:"id"`
	Project string `json:"project" yaml:"project"`
	Service string `json:"service" yaml:"service"`
	Type    string `json:"type" yaml:"type"` // Can be `project` or `cluster`
	Rules   []Rule `json:"rules" yaml:"rules"`
}

// Rule describe rule for service role
type Rule struct {
	APIGroups []string `json:"apiGroups" yaml:"apiGroups"`
	Verbs     []string `json:"verbs" yaml:"verbs"`
	Resources []string `json:"resources" yaml:"resources"`
}

// ServiceRoleProject is used to provide Type Project to role
const ServiceRoleProject string = "project"

// ServiceRoleCluster is used to provide Type Cluster to role
const ServiceRoleCluster string = "cluster"
