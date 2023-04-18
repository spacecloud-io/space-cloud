package pubsub

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
	MetaData map[string]string `json:"metadata"`
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
	Components Components         `json:"components,omitempty"`
}

// Channel defines a single channel schema
type Channel struct {
	Name    string         `json:"name,omitempty"`
	Payload ChannelPayload `json:"payload,omitempty"`
}

// ChannelPayload define channel's payload
type ChannelPayload struct {
	Schema   map[string]string `json:"schema,omitempty"`
	Example  interface{}       `json:"example,omitempty"`
	Examples []interface{}     `json:"examples,omitempty"`
}

// Components stores the components for the schema refs
type Components struct {
	Messages map[string]interface{} `json:"messages,omitempty"`
}
