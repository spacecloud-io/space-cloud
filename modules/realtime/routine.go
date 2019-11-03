package realtime

import (
	"log"
	"sync"

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"
)

func (m *Module) helperSendFeed(data *model.FeedData) {
	clientsTemp, ok := m.groups.Load(createGroupKey(data.DBType, data.Group))
	if !ok {
		log.Println("Realtime hanlder could not find key:", createGroupKey(data.DBType, data.Group))
		return
	}

	clients := clientsTemp.(*clientsStub)
	clients.clients.Range(func(key interface{}, value interface{}) bool {
		queries := value.(*sync.Map)
		queries.Range(func(id interface{}, value interface{}) bool {
			query := value.(*queryStub)

			dataPoint := &model.FeedData{
				QueryID: id.(string), DocID: data.DocID, Group: data.Group, Payload: data.Payload,
				TimeStamp: data.TimeStamp, Type: data.Type, DBType: data.DBType,
			}

			switch data.Type {
			case utils.RealtimeDelete:
				query.sendFeed(dataPoint)
				m.metrics.AddDBOperation(m.project, data.DBType, data.Group, 1, utils.Read)

			case utils.RealtimeInsert, utils.RealtimeUpdate:
				if utils.Validate(query.whereObj, data.Payload) {
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
