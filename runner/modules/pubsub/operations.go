package pubsub

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/segmentio/ksuid"

	"github.com/spaceuptech/space-cloud/runner/model"
	"github.com/spaceuptech/space-cloud/runner/utils"
)

// PublishString delivers a message in a fire and forget fashion
func (m *Module) PublishString(ctx context.Context, topic, value string) error {
	// Send the message
	return m.client.Publish(ctx, m.getTopicName(topic), value).Err()
}

// Send delivers a message reliably
func (m *Module) Send(ctx context.Context, topic string, value interface{}) error {
	// Create a new subscription on reply to channel
	replyTo := m.getTopicName(ksuid.New().String())

	// Create a subscription
	pubsub := m.client.Subscribe(context.TODO(), replyTo)
	defer utils.CloseTheCloser(pubsub)

	// Check if the subscription is active
	if _, err := pubsub.Receive(ctx); err != nil {
		return err
	}

	// Send the message
	data, err := json.Marshal(model.PubSubMessage{ReplyTo: replyTo, Payload: value})
	if err != nil {
		return err
	}

	if err := m.client.Publish(ctx, m.getTopicName(topic), string(data)).Err(); err != nil {
		return err
	}

	// Wait for the message to come back
	msg, err := pubsub.ReceiveMessage(ctx)
	if err != nil {
		return err
	}

	if msg.Payload != "ACK" {
		return fmt.Errorf("invalid response received in redis send - %s", msg.Payload)
	}

	return nil
}

// SendAck acknowledges the receipt of a message
func (m *Module) SendAck(ctx context.Context, replyTo string, ack bool) error {
	// Prepare response message
	msg := "ACK"
	if !ack {
		msg = "NACK"
	}

	// Send the acknowledgement
	return m.client.Publish(ctx, replyTo, msg).Err()
}

// Subscribe creates a subscription on a topic
func (m *Module) Subscribe(ctx context.Context, topic string) (<-chan *redis.Message, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	// Check if subscription already exists. Return corresponding channel if it does.
	if sub, p := m.mapping[topic]; p {
		return sub.ch, nil
	}

	// Make a redis subscription
	pubsub := m.client.Subscribe(context.TODO(), m.getTopicName(topic))
	if _, err := pubsub.Receive(ctx); err != nil {
		return nil, err
	}

	// Make a channel to listen for subscriptions
	ch := pubsub.Channel()
	m.mapping[topic] = &subscription{ch, pubsub}
	return ch, nil
}

// CancelSubscription cancels an active subscription
func (m *Module) CancelSubscription(topic string) {
	m.lock.Lock()
	defer m.lock.Unlock()

	// Close the subscription if it exists
	if sub, p := m.mapping[topic]; p {
		_ = sub.pubsub.Close()
	}

	// Remove it from the mapping
	delete(m.mapping, topic)
}
