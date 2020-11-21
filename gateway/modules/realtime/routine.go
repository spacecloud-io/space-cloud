package realtime

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

func (m *Module) helperSendFeed(ctx context.Context, data *model.FeedData) {
	clientsTemp, ok := m.groups.Load(createGroupKey(data.DBType, data.Group))
	if !ok {
		// This should be on the debug level
		helpers.Logger.LogDebug(helpers.GetRequestID(ctx), fmt.Sprintf("Realtime handler could not find key (%s)", createGroupKey(data.DBType, data.Group)), nil)
		return
	}

	clients := clientsTemp.(*clientsStub)
	clients.clients.Range(func(key interface{}, value interface{}) bool {
		queries := value.(*sync.Map)
		queries.Range(func(id interface{}, value interface{}) bool {
			query := value.(*queryStub)

			dataPoint := &model.FeedData{
				QueryID: id.(string), Group: data.Group, Payload: data.Payload, Find: data.Find,
				TimeStamp: data.TimeStamp, Type: data.Type, DBType: data.DBType,
			}

			switch data.Type {
			case utils.RealtimeDelete:
				_ = m.auth.PostProcessMethod(ctx, query.actions, dataPoint.Payload)
				query.sendFeed(dataPoint)
				m.metrics.AddDBOperation(m.project, data.DBType, data.Group, 1, model.Read)

			case utils.RealtimeInsert, utils.RealtimeUpdate:
				if utils.Validate(query.whereObj, data.Payload) {
					_ = m.auth.PostProcessMethod(ctx, query.actions, dataPoint.Payload)
					query.sendFeed(dataPoint)
					m.metrics.AddDBOperation(m.project, data.DBType, data.Group, 1, model.Read)
				}

			default:
				helpers.Logger.LogInfo(helpers.GetRequestID(ctx), "Realtime Module Error: Invalid event type received", map[string]interface{}{"dataType": data.Type})
			}
			return true
		})
		return true
	})
}

func (m *Module) routineHandleMessages() {
	ch, err := m.pubsubClient.Subscribe(context.Background(), getSendTopic(m.nodeID))
	if err != nil {
		panic(err)
	}

	for msg := range ch {
		pubsubMsg := new(model.PubSubMessage)
		if err := json.Unmarshal([]byte(msg.Payload), pubsubMsg); err != nil {
			_ = helpers.Logger.LogError("realtime-process", "Unable to marshal incoming realtime process request", err, map[string]interface{}{"payload": msg.Payload})
			continue
		}

		go m.handlePubSubMessage(pubsubMsg)
	}
}

func (m *Module) handlePubSubMessage(msg *model.PubSubMessage) {
	// Create a context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Unmarshal the incoming message
	event := new(model.CloudEventPayload)
	if err := msg.Unmarshal(event); err != nil {
		_ = helpers.Logger.LogError("realtime-process", "Unable to extract cloud event doc from incoming process realtime event request", err, map[string]interface{}{"payload": msg.Payload})
		_ = m.pubsubClient.SendAck(ctx, msg.ReplyTo, false)
		return
	}

	// Process the cloud event doc
	if err := m.ProcessRealtimeRequests(ctx, event); err != nil {
		_ = helpers.Logger.LogError("realtime-process", "Unable to extract event docs from incoming process event request", err, map[string]interface{}{"payload": msg.Payload})
		_ = m.pubsubClient.SendAck(ctx, msg.ReplyTo, false)
		return
	}

	_ = m.pubsubClient.SendAck(ctx, msg.ReplyTo, true)
}
