package eventing

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/model"
)

func (m *Module) createProcessUpdateEventsRoutine() {
	m.updateEventC = make(chan *queueUpdateEvent, 1000)
	for i := 0; i < 25; i++ {
		go m.routineProcessUpdateEvents()
	}
}

func (m *Module) routineProcessUpdateEvents() {
	for {
		select {
		case <-m.globalCloserChannel:
			return
		case ev := <-m.updateEventC:
			m.dbQueryForEventStatusUpdate(ev)
			m.deleteEventFromProcessingEventsMapChannel <- []string{ev.req.Find["_id"].(string)}
		}
	}
}

func (m *Module) routineUpdateEventsStatusInDB(updateChan chan *queueUpdateEvent) {
	duration := 250 * time.Millisecond
	t := time.NewTimer(duration)
	defer t.Stop()
	arr := make([]*queueUpdateEvent, 0)

	flush := func() {
		if len(arr) == 0 {
			t.Reset(duration)
			return
		}
		eventIDs, updateRequest := m.generateInOperatorUpdateRequest(arr)
		m.dbQueryForEventStatusUpdate(updateRequest)
		m.deleteEventFromProcessingEventsMapChannel <- eventIDs
		arr = make([]*queueUpdateEvent, 0)
		t.Reset(duration)
	}

	for {
		select {
		case <-m.globalCloserChannel:
			return
		case ev := <-updateChan:
			arr = append(arr, ev)
			if len(arr) > 200 {
				flush()
			}

		case <-t.C:
			flush()
		}
	}
}

func (m *Module) routineDeleteEventsFromSyncMap() {
	duration := 5 * time.Second
	t := time.NewTimer(duration)
	defer t.Stop()

	passiveArr, activeArr := make([]string, 0), make([]string, 0)
	for {
		select {
		case <-m.globalCloserChannel:
			return
		case eventIDs := <-m.deleteEventFromProcessingEventsMapChannel:
			activeArr = append(activeArr, eventIDs...)
		case <-t.C:
			for _, eventID := range passiveArr {
				// Delete the event from the processing list without fail
				m.processingEvents.Delete(eventID)
			}
			passiveArr = activeArr
			activeArr = make([]string, 0)
			t.Reset(duration)
		}
	}
}

func (m *Module) routineProcessIntents() {
	duration := 10 * time.Second
	t := time.NewTimer(duration)
	defer t.Stop()

	for {
		select {
		case <-m.globalCloserChannel:
			fmt.Println("Closing routineProcessIntents")
			return
		case ct := <-t.C:
			m.processIntents(&ct)
			t.Reset(duration)
		}
	}
}

func (m *Module) routineProcessStaged() {
	duration := 10 * time.Second
	t := time.NewTimer(duration)
	defer t.Stop()

	for {
		select {
		case <-m.globalCloserChannel:
			return
		case ct := <-t.C:
			m.processStagedEvents(&ct)
			t.Reset(duration)
		}
	}
}

func (m *Module) routineProcessEventsWithBuffering() {
	for {
		select {
		case <-m.globalCloserChannel:
			// Before closing the routine, delete all the un processed events in the buffer
			m.processingEvents = sync.Map{}
			return

		case eventDoc := <-m.bufferedEventProcessingChannel:
			go m.processStagedEvent(eventDoc)
		}
	}
}

func (m *Module) routineHandleMessages() {
	ch, err := m.pubsubClient.Subscribe(context.Background(), getEventingTopic(m.nodeID))
	if err != nil {
		panic(err)
	}

	for {
		select {
		case <-m.globalCloserChannel:
			return
		case msg := <-ch:
			if msg == nil {
				continue
			}
			pubsubMsg := new(model.PubSubMessage)
			if err := json.Unmarshal([]byte(msg.Payload), pubsubMsg); err != nil {
				_ = helpers.Logger.LogError("event-process", "Unable to marshal incoming process event request", err, map[string]interface{}{"payload": msg.Payload})
				continue
			}

			m.handlePubSubMessage(pubsubMsg)
		}
	}
}

func (m *Module) routineHandleEventResponseMessages() {
	ch, err := m.pubsubClient.Subscribe(context.Background(), getEventResponseTopic(m.nodeID))
	if err != nil {
		panic(err)
	}

	for {
		select {
		case <-m.globalCloserChannel:
			return
		case msg := <-ch:
			if msg == nil {
				continue
			}
			pubsubMsg := new(model.PubSubMessage)
			if err := json.Unmarshal([]byte(msg.Payload), pubsubMsg); err != nil {
				_ = helpers.Logger.LogError("event-response-process", "Unable to marshal incoming process event response message", err, map[string]interface{}{"payload": msg.Payload})
				continue
			}

			m.handleEventResponseMessage(pubsubMsg)
		}
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
