package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/spacecloud-io/space-cloud/pkg/apis/core/v1alpha1"
)

// Publish publishes message to a topic
func (a *Module) Publish(topic string, msg PublishMessage, options PublishOptions) error {
	b, err := json.Marshal(msg.Payload)
	if err != nil {
		return err
	}

	watermilMsg := message.NewMessage(msg.ID, b)
	watermilMsg.Metadata = msg.MetaData
	return a.pubSub.Publish(topic, watermilMsg)
}

// Subscribe subscribes to a topic
func (a *Module) Subscribe(ctx context.Context, clientID, topic string, options SubscribeOptions) (<-chan *message.Message, error) {
	messages, err := a.pubSub.Subscribe(ctx, topic)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (a *Module) createInternalChannels() {
	openapiProvisionChannel := v1alpha1.PubsubChannelSpec{
		Channel: "/openapi/provision",
		Payload: &v1alpha1.ChannelSchema{
			Type:                 "object",
			AdditionalProperties: json.RawMessage(fmt.Sprintf(`%t`, true)),
		},
		ProducerOptions: &v1alpha1.ChannelOptions{
			Plugins: []v1alpha1.HTTPPlugin{
				{
					Name:   "",
					Driver: "deny_user",
				},
			},
		},
		ConsumerOptions: &v1alpha1.ChannelOptions{
			Plugins: []v1alpha1.HTTPPlugin{
				{
					Name:   "",
					Driver: "authenticate_sc_user",
				},
			},
		},
	}

	asyncapiProvisionChannel := v1alpha1.PubsubChannelSpec{
		Channel: "/asyncapi/provision",
		Payload: &v1alpha1.ChannelSchema{
			Type:                 "object",
			AdditionalProperties: json.RawMessage(fmt.Sprintf(`%t`, true)),
		},
		ProducerOptions: &v1alpha1.ChannelOptions{
			Plugins: []v1alpha1.HTTPPlugin{
				{
					Name:   "",
					Driver: "deny_user",
				},
			},
		},
		ConsumerOptions: &v1alpha1.ChannelOptions{
			Plugins: []v1alpha1.HTTPPlugin{
				{
					Name:   "",
					Driver: "authenticate_sc_user",
				},
			},
		},
	}

	a.channels = append(a.channels, openapiProvisionChannel, asyncapiProvisionChannel)
}

// Channels return channels with their schema
func (a *Module) Channels() ChannelsWithSchema {
	channels := ChannelsWithSchema{
		Channels: make(map[string]v1alpha1.PubsubChannelSpec),
	}

	for _, channel := range a.channels {
		channelPath := getChannelPath(channel.Channel)
		channels.Channels[channelPath] = channel
	}
	return channels
}

func getChannelPath(name string) string {
	if name[0] != '/' {
		name = "/" + name
	}
	return strings.ReplaceAll(name, "-", "/")
}
