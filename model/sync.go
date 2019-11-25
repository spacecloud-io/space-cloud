package model

import (
	"github.com/spaceuptech/space-cloud/config"
)

// StoreProject is a function used to store the project config
type StoreProject func(project *config.Project) error
type StoreProjectIgnoreError func(project *config.Project) error

type SetGlobalConfig func(projectID, secret string) error
type SetCrudConfig func(projectID string, crud config.Crud) error
type SetServicesConfig func(projectID string, services *config.ServicesModule) error
type SetFileStorageConfig func(projectID string, fileStore *config.FileStore) error
type SetEventingConfig func(projectID string, eventing *config.Eventing) error
type SetUserManConfig func(projectID string, userMan config.Auth) error

// DeleteProject is a function used to delete a project
type DeleteProject func(projectID string)

// GetProjectIDs returns the ids of the projects
type GetProjectIDs func() []string

// ProjectCallbacks is used to set or delete a projects config
type ProjectCallbacks struct {
	Store            StoreProject
	StoreIgnoreError StoreProjectIgnoreError

	SetGlobalConfig      SetGlobalConfig
	SetCrudConfig        SetCrudConfig
	SetServicesConfig    SetServicesConfig
	SetFileStorageConfig SetFileStorageConfig
	SetEventingConfig    SetEventingConfig
	SetUserManConfig     SetUserManConfig

	Delete     DeleteProject
	ProjectIDs GetProjectIDs
}
