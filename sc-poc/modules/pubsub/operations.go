package pubsub

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/spacecloud-io/space-cloud/pkg/apis/core/v1alpha1"
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

func (a *App) createInternalChannels() {
	openapiProvisionChannel := Channel{
		Name: "openapi-provision",
		Payload: ChannelPayload{
			Schema: map[string]*v1alpha1.ChannelSchema{
				"doc": {
					Type: "string",
				},
			},
		},
	}

	asyncapiProvisionChannel := Channel{
		Name: "asyncapi-provision",
		Payload: ChannelPayload{
			Schema: map[string]*v1alpha1.ChannelSchema{
				"doc": {
					Type: "object",
					Properties: map[string]*v1alpha1.ChannelSchema{
						"name": {
							Type: "string",
						},
						"age": {
							Type: "integer",
						},
					},
					Required: []string{"name"},
				},
			},
		},
	}
	a.channels = append(a.channels, openapiProvisionChannel, asyncapiProvisionChannel)
}

// Channels return channels with their schema
func (a *App) Channels() ChannelsWithSchema {
	channels := ChannelsWithSchema{
		Channels: make(map[string]Channel),
		Components: &Components{
			Schemas: map[string]interface{}{
				"APIManMsg": map[string]interface{}{
					"type":                 "object",
					"additionalProperties": true,
				},
			},
		},
	}

	for _, topic := range a.channels {
		channelPath := getChannelPath(topic.Name)
		channels.Channels[channelPath] = topic
	}
	return channels
}

func getChannelPath(name string) string {
	if name[0] != '/' {
		name = "/" + name
	}
	return strings.ReplaceAll(name, "-", "/")
}
