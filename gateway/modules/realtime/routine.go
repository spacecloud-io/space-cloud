package realtime

import (
	"context"
	"fmt"
	"sync"

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
