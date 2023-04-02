package pubsub

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/message"
)

type (
	// Source describes the implementation of source from the pubsub module
	Source interface {
		Publish(topic string, data interface{}) error
		Subscribe(ctx context.Context, topic string) (<-chan *message.Message, error)
	}
)
