package model

import (
	"github.com/spaceuptech/space-cloud/config"
)

// StoreProject is a function used to store the project config
type StoreProject func(project *config.Project) error

// DeleteProject is a function used to delete a project
type DeleteProject func(projectID string)

// GetProjectIDs returns the ids of the projects
type GetProjectIDs func() []string

// ProjectCallbacks is used to set or delete a projects config
type ProjectCallbacks struct {
	Store      StoreProject
	Delete     DeleteProject
	ProjectIDs GetProjectIDs
}
