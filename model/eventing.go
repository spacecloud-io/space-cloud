package model

import "context"

// EventKind is the type describing the kind of event
type EventKind string

// QueueEventRequest is the payload to add a new event to the task queue
type QueueEventRequest struct {
	Name      string      `json:"name"`                // The type of the event
	Delay     int64       `json:"delay,omitempty"`     // Time in seconds
	Timestamp int64       `json:"timestamp,omitempty"` // Milliseconds from unix epoch (UTC)
	Payload   interface{} `json:"Payload,omitempty"`   // Payload contains necessary event data
	Retries   int         `json:"retries,omitempty"`
}

// EventIntent describes an intent made in the eventing system
type EventIntent struct {
	BatchID string
	Token   int
	Docs    []interface{}
	Invalid bool
}

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

// CreateMessage is the event payload for create message event
type CreateMessage struct {
	DBType string        `json:"db"`
	Col    string        `json:"col"`
	Docs   []interface{} `json:"docs"`
}

// UpdateMessage is the event payload for update message event
type UpdateMessage struct {
	DBType string `json:"db"`
	Col    string `json:"col"`
	DocID  string `json:"docId"`
}

// DeleteMessage is the event payload for delete message event
type DeleteMessage struct {
	DBType string `json:"db"`
	Col    string `json:"col"`
	DocID  string `json:"docId"`
}
