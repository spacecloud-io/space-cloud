package model

// ID describes the configuration of a project
type Project struct {
	ID          string `json:"id" yaml:"id"`
	Environment string `json:"env" yaml:"env"`
}
