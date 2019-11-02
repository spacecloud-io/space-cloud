package model

// EventKind is the type describing the kind of event
type EventKind string

// EventDocument is the format in which the event is persistent on disk
type EventDocument struct {
	ID             string      `structs:"_id" json:"_id" bson:"_id" mapstructure:"_id"`
	BatchID        string      `structs:"batchid" json:"batchid" bson:"batchid" mapstructure:"batchid"`
	Type           string      `structs:"type" json:"type" bson:"type" mapstructure:"type"`
	Token          int         `structs:"token" json:"token" bson:"token" mapstructure:"token"`
	Timestamp      int64       `structs:"timestamp" json:"timestamp" bson:"timestamp" mapstructure:"timestamp"`                         // The timestamp of when the event should get executed
	EventTimestamp int64       `structs:"event_timestamp" json:"event_timestamp" bson:"event_timestamp" mapstructure:"event_timestamp"` // The time stamp of when the event was logged
	Payload        interface{} `structs:"payload" json:"payload" bson:"payload" mapstructure:"payload"`
	Status         string      `structs:"status" json:"status" bson:"status" mapstructure:"status"`
	Retries        int         `structs:"retries" json:"retries" bson:"retries" mapstructure:"retries"`
	Url            string      `structs:"url" json:"url" bson:"url" mapstructure:"url"`
	Remark         string      `structs:"remark" json:"remark" bson:"remark" mapstructure:"remark"`
}

// CloudEventPayload is the the JSON event spec by Cloud Events Specification
type CloudEventPayload struct {
	SpecVersion string      `json:"specversion"`
	Type        string      `json:"type"`
	Source      string      `json:"source"`
	Id          string      `json:"id"`
	Time        string      `json:"time"`
	Data        interface{} `json:"data"`
}

type EventResponse struct {
	Event  *QueueEventRequest   `json:"event"`
	Events []*QueueEventRequest `json:"events"`
	Error  string               `json:"error"`
}

// QueueEventRequest is the payload to add a new event to the task queue
type QueueEventRequest struct {
	Type      string            `json:"type"`                // The type of the event
	Delay     int64             `json:"delay,omitempty"`     // Time in seconds
	Timestamp int64             `json:"timestamp,omitempty"` // Milliseconds from unix epoch (UTC)
	Payload   interface{}       `json:"payload,omitempty"`   // Payload contains necessary event dat
	Options   map[string]string `json:"options"`
}

// EventIntent describes an intent made in the eventing system
type EventIntent struct {
	BatchID string
	Token   int
	Docs    []*EventDocument
	Invalid bool
}

// DatabaseEventMessage is the event payload for create, update and delete events
type DatabaseEventMessage struct {
	DBType string      `json:"db" mapstructure:"db"`
	Col    string      `json:"col" mapstructure:"col"`
	DocID  string      `json:"docId" mapstructure:"docId"`
	Doc    interface{} `json:"doc" mapstructure:"doc"`
}
