package model

import (
	"context"

	"github.com/spaceuptech/space-cloud/utils"
)

// MetricCrudHook is used to log a database operation
type MetricCrudHook func(project, dbType, col string, count int64, op utils.OperationType)

// CreateIntentHook is used to log a create intent
type CreateIntentHook func(ctx context.Context, dbType, col string, req *CreateRequest) (*EventIntent, error)

// UpdateIntentHook is used to log a create intent
type UpdateIntentHook func(ctx context.Context, dbType, col string, req *UpdateRequest) (*EventIntent, error)

// DeleteIntentHook is used to log a create intent
type DeleteIntentHook func(ctx context.Context, dbType, col string, req *DeleteRequest) (*EventIntent, error)

// BatchIntentHook is used to log a create intent
type BatchIntentHook func(ctx context.Context, dbType string, req *BatchRequest) (*EventIntent, error)

// StageEventHook is used to stage an intended event
type StageEventHook func(ctx context.Context, intent *EventIntent, err error)

// CrudHooks is the struct to store the hooks related to the crud module
type CrudHooks struct {
	Create CreateIntentHook
	Update UpdateIntentHook
	Delete DeleteIntentHook
	Batch  BatchIntentHook
	Stage  StageEventHook
}
