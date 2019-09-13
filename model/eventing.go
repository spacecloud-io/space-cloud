package model

import "context"

// EventKind is the type describing the kind of event
type EventKind string

// EventDocument is the format in which the event is persistent on disk
type EventDocument struct {
	ID             string `json:"_id"`
	BatchID        string `json:"batchid"`
	Type           string `json:"type"`
	Token          int    `json:"token"`
	Timestamp      int64  `json:"timestamp"`
	EventTimestamp int64  `json:"event_timestamp"`
	Payload        string `json:"payload"`
	Status         string `json:"status"`
	MaxRetries     int    `json:"max_retries"`
	Retries        int    `json:"retries"`
	Delivered      bool   `json:"delivered"`
	Service        string `json:"service"`
	Function       string `json:"function"`
}

// QueueEventRequest is the payload to add a new event to the task queue
type QueueEventRequest struct {
	Type      string      `json:"type"`                // The type of the event
	Delay     int64       `json:"delay,omitempty"`     // Time in seconds
	Timestamp int64       `json:"timestamp,omitempty"` // Milliseconds from unix epoch (UTC)
	Payload   interface{} `json:"Payload,omitempty"`   // Payload contains necessary event dat
}

// EventIntent describes an intent made in the eventing system
type EventIntent struct {
	BatchID string
	Token   int
	Docs    []*EventDocument
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

// DatabaseEventMessage is the event payload for create, update and delete events
type DatabaseEventMessage struct {
	DBType string      `json:"db"`
	Col    string      `json:"col"`
	DocID  string      `json:"docId"`
	Doc    interface{} `json:"doc"`
}
