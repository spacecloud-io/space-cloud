package model

import (
	"context"
)

// MetricCrudHook is used to log a database operation
type MetricCrudHook func(project, dbAlias, col string, count int64, op OperationType)

// MetricFileHook is used to log a file operation
type MetricFileHook func(project, storeType string, op OperationType)

// MetricFunctionHook is used to log a function operation
type MetricFunctionHook func(project, service, function string)

// MetricEventingHook is used to log a eventing operation
type MetricEventingHook func(project, eventingType string)

// CreateIntentHook is used to log a create intent
type CreateIntentHook func(ctx context.Context, dbAlias, col string, req *CreateRequest) (*EventIntent, error)

// UpdateIntentHook is used to log a create intent
type UpdateIntentHook func(ctx context.Context, dbAlias, col string, req *UpdateRequest) (*EventIntent, error)

// DeleteIntentHook is used to log a create intent
type DeleteIntentHook func(ctx context.Context, dbAlias, col string, req *DeleteRequest) (*EventIntent, error)

// BatchIntentHook is used to log a create intent
type BatchIntentHook func(ctx context.Context, dbAlias string, req *BatchRequest) (*EventIntent, error)

// StageEventHook is used to stage an intended event
type StageEventHook func(ctx context.Context, intent *EventIntent, err error)

// EventingModule is the interface to mock the eventing module
type EventingModule interface {
	CreateFileIntentHook(ctx context.Context, req *CreateFileRequest) (*EventIntent, error)
	DeleteFileIntentHook(ctx context.Context, path string, meta map[string]interface{}) (*EventIntent, error)
	HookStage(ctx context.Context, intent *EventIntent, err error)
}
