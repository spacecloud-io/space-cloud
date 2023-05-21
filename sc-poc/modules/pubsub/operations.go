package pubsub

import (
	"context"
	"encoding/json"

	"github.com/ThreeDotsLabs/watermill/message"
)

// Publish publishes message to a topic
func (a *App) Publish(topic string, msg PublishMessage, options PublishOptions) error {
	b, err := json.Marshal(msg.Payload)
	if err != nil {
		return err
	}

	watermilMsg := message.NewMessage(msg.ID, b)
	watermilMsg.Metadata = msg.MetaData
	return a.pubSub.Publish(topic, watermilMsg)
}

// Subscribe subscribes to a topic
func (a *App) Subscribe(ctx context.Context, clientID, topic string, options SubscribeOptions) (<-chan *message.Message, error) {
	messages, err := a.pubSub.Subscribe(ctx, topic)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

// Channels return channels with their schema
func (a *App) Channels() ChannelsWithSchema {
	return ChannelsWithSchema{
		Channels: map[string]Channel{
			"/sc/api": {
				Name: "api-provision",
				Payload: ChannelPayload{
					Schema: map[string]interface{}{
						"$ref": "#/components/schemas/APIManMsg",
					},
				},
			},
		},
		Components: &Components{
			Schemas: map[string]interface{}{
				"APIManMsg": map[string]interface{}{
					"type":                 "object",
					"additionalProperties": true,
				},
			},
		},
	}
}
