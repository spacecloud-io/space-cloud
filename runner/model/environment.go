package model

// Project describes the configuration of a project
type Project struct {
	ID          string `json:"id" yaml:"id"`
	Kind        string `json:"kind" yaml:"kind"`
	Environment string `json:"env" yaml:"env"`
}
