package eventing

import (
	"context"
	"encoding/json"
	"time"

	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/model"
)

func (m *Module) createProcessUpdateEventsRoutine() {
	m.updateEventC = make(chan *queueUpdateEvent, 1000)
	for i := 0; i < 50; i++ {
		go m.routineProcessUpdateEvents()
	}
}

func (m *Module) routineProcessUpdateEvents() {
	for ev := range m.updateEventC {
		m.queueUpdateEvent(ev)
	}
}

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
	ch, err := m.pubsubClient.Subscribe(context.Background(), getEventingTopic(m.nodeID))
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

func (m *Module) routineHandleEventResponseMessages() {
	ch, err := m.pubsubClient.Subscribe(context.Background(), getEventResponseTopic(m.nodeID))
	if err != nil {
		panic(err)
	}

	for msg := range ch {
		pubsubMsg := new(model.PubSubMessage)
		if err := json.Unmarshal([]byte(msg.Payload), pubsubMsg); err != nil {
			_ = helpers.Logger.LogError("event-response-process", "Unable to marshal incoming process event response message", err, map[string]interface{}{"payload": msg.Payload})
			continue
		}

		m.handleEventResponseMessage(pubsubMsg)
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

func (m *Module) handleEventResponseMessage(msg *model.PubSubMessage) {
	// Create a context
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Unmarshal the incoming message
	eventResponse, ok := msg.Payload.(map[string]interface{})
	if !ok {
		_ = helpers.Logger.LogError("event-response-process", "Unable to extract event response message from incoming request", nil, map[string]interface{}{"payload": msg.Payload})
		_ = m.pubsubClient.SendAck(ctx, msg.ReplyTo, false)
		return
	}

	m.ProcessEventResponseMessage(ctx, eventResponse["batchId"].(string), eventResponse["response"])
	_ = m.pubsubClient.SendAck(ctx, msg.ReplyTo, true)
}
