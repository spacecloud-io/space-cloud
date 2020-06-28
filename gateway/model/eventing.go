package model

// EventKind is the type describing the kind of event
type EventKind string

// EventDocument is the format in which the event is persistent on disk
type EventDocument struct {
	ID             string      `structs:"_id" json:"_id" bson:"_id" mapstructure:"_id"`
	BatchID        string      `structs:"batchid" json:"batchid" bson:"batchid" mapstructure:"batchid"`
	Type           string      `structs:"type" json:"type" bson:"type" mapstructure:"type"`
	RuleName       string      `structs:"rule_name" json:"rule_name" bson:"rule_name" mapstructure:"rule_name"`
	Token          int         `structs:"token" json:"token" bson:"token" mapstructure:"token"`
	Timestamp      string      `structs:"ts" json:"ts" bson:"ts" mapstructure:"ts"`                         // The timestamp of when the event should get executed
	EventTimestamp string      `structs:"event_ts" json:"event_ts" bson:"event_ts" mapstructure:"event_ts"` // The time stamp of when the event was logged
	Payload        interface{} `structs:"payload" json:"payload" bson:"payload" mapstructure:"payload"`
	Status         string      `structs:"status" json:"status" bson:"status" mapstructure:"status"`
	Remark         string      `structs:"remark" json:"remark" bson:"remark" mapstructure:"remark"`
}

// InvocationDocument is the format in which the invocation are persistent on disk
type InvocationDocument struct {
	ID                 string `struct:"_id" json:"_id" bson:"_id" mapstructure:"_id"`
	EventID            string `struct:"event_id" json:"event_id" bson:"event_id" mapstructure:"event_id"`
	InvocationTime     string `struct:"invocation_time" json:"invocation_time" bson:"invocation_time" mapstructure:"invocation_time"`
	RequestPayload     string `struct:"request_payload" json:"request_payload" bson:"request_payload" mapstructure:"request_payload"`
	ResponseStatusCode int    `struct:"response_status_code" json:"response_status_code" bson:"response_status_code" mapstructure:"response_status_code"`
	ResponseBody       string `struct:"response_body" json:"response_body" bson:"response_body" mapstructure:"response_body"`
	ErrorMessage       string `struct:"error_msg" json:"error_msg" bson:"error_msg" mapstructure:"error_msg"`
	Remark             string `struct:"remark" json:"remark" bson:"remark" mapstructure:"remark"`
}

// CloudEventPayload is the the JSON event spec by Cloud Events Specification
type CloudEventPayload struct {
	SpecVersion string      `json:"specversion"`
	Type        string      `json:"type"`
	Source      string      `json:"source"`
	ID          string      `json:"id"`
	Time        string      `json:"time"`
	Data        interface{} `json:"data"`
}

// EventResponse is struct response of events
type EventResponse struct {
	Event    *QueueEventRequest   `json:"event,omitempty"`
	Events   []*QueueEventRequest `json:"events,omitempty"`
	Response interface{}          `json:"response,omitempty"` // for getting response of synchronous events
	Error    string               `json:"error,omitempty"`
}

// QueueEventRequest is the payload to add a new event to the task queue
type QueueEventRequest struct {
	Type          string            `json:"type"`                // The type of the event
	Delay         int64             `json:"delay,omitempty"`     // Time in seconds
	Timestamp     string            `json:"timestamp,omitempty"` // Milliseconds from unix epoch (UTC)
	Payload       interface{}       `json:"payload,omitempty"`   // Payload contains necessary event dat
	Options       map[string]string `json:"options"`
	IsSynchronous bool              `json:"isSynchronous"` // if true then client will wait for response of event
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
	Doc    interface{} `json:"doc" mapstructure:"doc"`
	Find   interface{} `json:"find" mapstructure:"find"`
}
