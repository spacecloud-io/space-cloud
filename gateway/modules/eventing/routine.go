package eventing

import (
	"context"
	"encoding/json"
	"time"

	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/model"
)

func (m *Module) routineProcessIntents() {
	m.tickerIntent = time.NewTicker(10 * time.Second)
	for t := range m.tickerIntent.C {
		m.processIntents(&t)
	}
}

func (m *Module) routineProcessStaged() {
	m.tickerStaged = time.NewTicker(10 * time.Second)
	for t := range m.tickerStaged.C {
		m.processStagedEvents(&t)
	}
}

func (m *Module) routineHandleMessages() {
	ch, err := m.pubsubClient.Subscribe(context.Background(), getSendTopic(m.nodeID))
	if err != nil {
		panic(err)
	}

	for msg := range ch {
		pubsubMsg := new(model.PubSubMessage)
		if err := json.Unmarshal([]byte(msg.Payload), pubsubMsg); err != nil {
			_ = helpers.Logger.LogError("event-process", "Unable to marshal incoming process event request", err, map[string]interface{}{"payload": msg.Payload})
			continue
		}

		m.handlePubSubMessage(pubsubMsg)
	}
}

func (m *Module) handlePubSubMessage(msg *model.PubSubMessage) {
	// Create a context
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Unmarshal the incoming message
	var eventDocs []*model.EventDocument
	if err := msg.Unmarshal(&eventDocs); err != nil {
		_ = helpers.Logger.LogError("event-process", "Unable to extract event docs from incoming process event request", err, map[string]interface{}{"payload": msg.Payload})
		_ = m.pubsubClient.SendAck(ctx, msg.ReplyTo, false)
		return
	}

	m.ProcessTransmittedEvents(eventDocs)
	_ = m.pubsubClient.SendAck(ctx, msg.ReplyTo, true)
}
