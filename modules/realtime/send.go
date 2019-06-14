package realtime

import (
	"context"
	"time"

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/modules/crud"
	"github.com/spaceuptech/space-cloud/utils"
)

// SendCreate broadcasts a realtime create datapoint to the concerned clients
func (m *Module) SendCreate(dbType string, col string, req *model.CreateRequest) {
	var rows []interface{}
	switch req.Operation {
	case utils.One:
		rows = []interface{}{req.Document}
	case utils.All:
		rows = req.Document.([]interface{})
	default:
		rows = []interface{}{}
	}

	for _, t := range rows {
		data := t.(map[string]interface{})

		idVar := "id"
		if dbType == string(utils.Mongo) {
			idVar = "_id"
		}

		// Send realtime message if id fields exists
		if idTemp, p := data[idVar]; p {
			if id, ok := AcceptableIdType(idTemp); ok {
				m.send(&model.FeedData{
					Group:     col,
					DBType:    dbType,
					Type:      utils.RealtimeWrite,
					TimeStamp: time.Now().Unix(),
					DocID:     id,
					Payload:   data,
				})
			}
		}
	}
}

// SendUpdate broadcasts a realtime update datapoint to the concerned clients
func (m *Module) SendUpdate(ctx context.Context, project, dbType, col string, req *model.UpdateRequest, crud *crud.Module) {
	idVar := "id"
	if dbType == string(utils.Mongo) {
		idVar = "_id"
	}

	if idTemp, p := req.Find[idVar]; p {
		if id, ok := AcceptableIdType(idTemp); ok {
			// Create the find object
			find := map[string]interface{}{idVar: id}

			data, err := crud.Read(ctx, dbType, project, col, &model.ReadRequest{Find: find, Operation: utils.One})
			if err == nil {
				m.send(&model.FeedData{
					Group:     col,
					Type:      utils.RealtimeWrite,
					TimeStamp: time.Now().Unix(),
					DocID:     id,
					DBType:    dbType,
					Payload:   data.(map[string]interface{}),
				})
			}
		}
	}
}

// SendDelete broadcasts a realtime delete datapoint to the concerned clients
func (m *Module) SendDelete(dbType string, col string, req *model.DeleteRequest) {
	idVar := "id"
	if dbType == string(utils.Mongo) {
		idVar = "_id"
	}

	if idTemp, p := req.Find[idVar]; p {
		if id, ok := AcceptableIdType(idTemp); ok {
			m.send(&model.FeedData{
				Group:     col,
				Type:      utils.RealtimeDelete,
				TimeStamp: time.Now().Unix(),
				DocID:     id,
				DBType:    dbType,
			})
		}
	}
}

func (m *Module) send(data *model.FeedData) {
	m.RLock()
	defer m.RUnlock()

	if m.enabled {
		m.feed <- data
	}
}
