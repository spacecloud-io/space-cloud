package pubsub

import (
	"context"
	"encoding/json"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
)

// Publish publishes message to a topic
func (a *App) Publish(topic string, data interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	msg := message.NewMessage(watermill.NewUUID(), b)
	return a.pubSub.Publish(topic, msg)
}

// Subscribe subscribes to a topic
func (a *App) Subscribe(ctx context.Context, topic string) (<-chan *message.Message, error) {
	messages, err := a.pubSub.Subscribe(ctx, topic)
	if err != nil {
		return nil, err
	}
	return messages, nil
}
