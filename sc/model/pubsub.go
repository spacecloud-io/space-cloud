package model

import (
	"errors"

	"github.com/mitchellh/mapstructure"
)

// PubSubMessage describes the format of pubsub send message
type PubSubMessage struct {
	ReplyTo string      `json:"replyTo"`
	Payload interface{} `json:"payload"`
}

// Unmarshal parses the payload into the object provided
func (m *PubSubMessage) Unmarshal(ptr interface{}) error {
	if m.Payload == nil {
		return errors.New("no payload has been provided")
	}

	return mapstructure.Decode(m.Payload, ptr)
}
