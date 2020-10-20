package model

type Role struct {
	ID      string `json:"id" yaml:"id"`
	Project string `json:"project" yaml:"project"`
	Service string `json:"service" yaml:"service"`
	Type    string `json:"type" yaml:"type"`
	Rules   []Rule `json:"rules" yaml:"rules"`
}

type Rule struct {
	APIGroups []string `json:"apiGroups" yaml:"apiGroups"`
	Verbs     []string `json:"verbs" yaml:"verbs"`
	Resources []string `json:"resources" yaml:"resources"`
}
