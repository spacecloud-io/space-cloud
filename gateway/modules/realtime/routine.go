package realtime

import (
	"log"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

func (m *Module) helperSendFeed(data *model.FeedData) {
	clientsTemp, ok := m.groups.Load(createGroupKey(data.DBType, data.Group))
	if !ok {
		// This should be on the debug level
		logrus.Debugln("Realtime handler could not find key:", createGroupKey(data.DBType, data.Group))
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
				_ = m.auth.PostProcessMethod(query.actions, dataPoint.Payload)
				query.sendFeed(dataPoint)
				m.metrics.AddDBOperation(m.project, data.DBType, data.Group, 1, utils.Read)

			case utils.RealtimeInsert, utils.RealtimeUpdate:
				if utils.Validate(query.whereObj, data.Payload) {
					_ = m.auth.PostProcessMethod(query.actions, dataPoint.Payload)
					query.sendFeed(dataPoint)
					m.metrics.AddDBOperation(m.project, data.DBType, data.Group, 1, utils.Read)
				}

			default:
				log.Println("Realtime Module Error: Invalid event type received -", data.Type)
			}
			return true
		})
		return true
	})
}
