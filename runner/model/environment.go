package model

// Environment describes the configuration of a project
type Environment struct {
	ID      string `json:"id" yaml:"id"`
	Project string `json:"projectId" yaml:"projectId"`
}

// Environments describes a collection on environments
type Environments []Environment
