package pubsub

import (
	v1alpha1 "github.com/spacecloud-io/space-cloud/pkg/apis/core/v1alpha1"
)

type (
	EventType string
)

const (
	SubscribeEvent   EventType = "subscribe"
	UnsubscribeEvent EventType = "unsubscribe"
	MessageEvent     EventType = "message"
)

// Message defines the type of event and the associated data
type Message struct {
	Event EventType              `json:"event"`
	Data  map[string]interface{} `json:"data"`
}

// PublishMessage defines the type for publishing a message
type PublishMessage struct {
	ID       string            `json:"id"`
	MetaData map[string]string `json:"metadata,omitempty"`
	Payload  interface{}       `json:"payload"`
}

// PublishOptions defines the options for publishing a message
type PublishOptions struct {
	RequireAck bool `json:"requireAck"`
}

// SubscribeOptions defines the options for subscribing a message
type SubscribeOptions struct {
	Mode     string `json:"mode"`
	Capacity int    `json:"capacity"`
	Autoack  bool   `json:"autoack"`
	Format   string `json:"format"`
}

// ChannelsWithSchema define the channels schema and component
type ChannelsWithSchema struct {
	Channels   map[string]Channel `json:"channels,omitempty"` // key is the url
	Components *Components        `json:"components,omitempty"`
}

// Channel defines a single channel schema
type Channel struct {
	Name    string         `json:"name,omitempty"`
	Payload ChannelPayload `json:"payload,omitempty"`
}

// ChannelPayload define channel's payload
type ChannelPayload struct {
	Schema   map[string]*v1alpha1.ChannelSchema `json:"schema,omitempty"`
	Example  interface{}                        `json:"example,omitempty"`
	Examples []interface{}                      `json:"examples,omitempty"`
}

// Components stores the components for the schema refs
type Components struct {
	Schemas map[string]interface{} `json:"schemas,omitempty"`
}

// AsyncAPI defines the AsyncAPI 2.6.0 standard.
type AsyncAPI struct {
	SpecVersion string             `json:"asyncapi"` // Required
	Info        Info               `json:"info"`     // Required
	Channels    Channels           `json:"channels"` // Required
	Servers     Servers            `json:"servers,omitempty"`
	Components  AsyncAPIComponents `json:"components,omitempty"`
}

// Servers represents "servers" specified by AsyncAPI standard.
type Servers map[string]*ServerItem

// Channels represents "channels" specified by AsyncAPI standard.
type Channels map[string]*ChannelItem

// ChannelItem represents the two operations - "publish" and "subscribe"
type ChannelItem struct {
	Subscribe *Operation `json:"subscribe,omitempty"`
	Publish   *Operation `json:"publish,omitempty"`
}

// Operation represents the details of each operation
type Operation struct {
	Message MessageOneOrMany `json:"message,omitempty"`
	ID      string           `json:"operationId,omitempty"`
}

// OneOf consists of array of messages
type MessageOneOrMany struct {
	MessageEntity `json:",inline"`
	OneOf         []MessageEntity `json:"oneOf,omitempty"`
}

// MessageEntity defines the message as specified by AsyncAPI standard.
type MessageEntity struct {
	Name        string                 `json:"name,omitempty"`
	ContentType string                 `json:"contentType,omitempty"`
	Payload     map[string]interface{} `json:"payload,omitempty"`
}

// Info defines the info as specified by AsyncAPI standard.
type Info struct {
	Title       string `json:"title"`   // Required
	Version     string `json:"version"` // Required
	Description string `json:"description,omitempty"`
}

// ServerItem defines the type of server.
type ServerItem struct {
	URL         string `json:"url"`      // Required.
	Protocol    string `json:"protocol"` // Required.
	Description string `json:"description,omitempty"`
}

// AsyncAPIComponents defines the components specified by AsyncAPI standard.
type AsyncAPIComponents struct {
	Schemas map[string]interface{} `json:"schemas,omitempty"`
}
