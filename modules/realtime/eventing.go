package realtime

import (
	"encoding/json"
	"log"
	"sync/atomic"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/nats-io/nats.go"
	uuid "github.com/satori/go.uuid"

	"github.com/spaceuptech/space-cloud/model"
)

func (m *Module) registerEventHandlerService() error {
	m.RLock()
	defer m.RUnlock()

	// Get the internal access token
	token, err := m.auth.GetInternalAccessToken()
	if err != nil {
		return err
	}

	// Register a function
	req := &model.ServiceRegisterRequest{
		Service: serviceName,
		Project: m.project,
		Token:   token,
	}

	m.functions.RegisterService(m.nodeID, req, m.handleFunctionCall)
	return nil
}

func (m *Module) handleFunctionCall(payload *model.FunctionsPayload) {

	eventDoc := new(model.EventDocument)
	if err := mapstructure.Decode(payload.Params, eventDoc); err != nil {
		log.Println("Realtime Module Events Handler Error:", err)

		// Create and send a -ve response
		m.sendFunctionsAck(payload.ID, false)
		return
	}

	// Make a receive channel
	channel := make(chan struct{}, 1)

	// Get alive member count
	count := int32(m.syncMan.GetAliveNodeCount())

	// Listen for a single response
	reply := uuid.NewV1().String()
	sub, err := m.ec.Subscribe(reply, func(res *handlerAck) {
		if res.Ack {
			count = atomic.AddInt32(&count, -1)
			if atomic.LoadInt32(&count) == 0 {
				channel <- struct{}{}
			}
		}
	})
	if err != nil {
		log.Println("Realtime Module Events Handler Error:", err)

		// Create and send a -ve response
		m.sendFunctionsAck(payload.ID, false)
		return
	}

	defer sub.Unsubscribe()

	if err := m.ec.PublishRequest(pubSubTopic, reply, eventDoc); err != nil {
		log.Println("Realtime Module Events Handler Error:", err)

		// Create and send a -ve response
		m.sendFunctionsAck(payload.ID, false)
		return
	}

	select {
	case <-time.After(10 * time.Second):
		log.Println("RealtimeHandler: expired")
		m.sendFunctionsAck(payload.ID, false)
		return

	case <-channel:
		// Create and send a +ve response
		m.sendFunctionsAck(payload.ID, true)
		return
	}
}

func (m *Module) handleRealtimeRequests(msg *nats.Msg) {
	eventDoc := new(model.EventDocument)
	if err := json.Unmarshal(msg.Data, eventDoc); err != nil {
		log.Println("Realtime Module Request Handler Error:", err)

		// Create and send a -ve response
		m.sendRealtimeAck(msg.Reply, false)
		return
	}

	dbEvent := new(model.DatabaseEventMessage)
	if err := json.Unmarshal([]byte(eventDoc.Payload), dbEvent); err != nil {
		log.Println("Realtime Module Request Handler Error:", err)

		// Create and send a -ve response
		m.sendRealtimeAck(msg.Reply, false)
		return
	}

	feedData := &model.FeedData{
		DocID:     dbEvent.DocID,
		Type:      eventingToRealtimeEvent(eventDoc.Type),
		Payload:   dbEvent.Doc,
		TimeStamp: eventDoc.Timestamp,
		Group:     dbEvent.Col,
		DBType:    dbEvent.DBType,
	}

	m.helperSendFeed(feedData)
	// Create and send a -ve response
	m.sendRealtimeAck(msg.Reply, true)
}

func (m *Module) sendFunctionsAck(id string, ack bool) {
	res := model.FunctionsPayload{
		ID:     id,
		Auth:   nil,
		Params: map[string]bool{"ack": ack},
	}
	m.functions.HandleServiceResponse(&res)
}

func (m *Module) sendRealtimeAck(reply string, ack bool) {
	res := handlerAck{Ack: ack}
	if err := m.ec.Publish(reply, res); err != nil {
		log.Println("Realtime request handler could not publish response:", err)
	}
}
