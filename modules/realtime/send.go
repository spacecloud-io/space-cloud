package realtime

import (
	"encoding/json"
	"log"
	"time"

	uuid "github.com/satori/go.uuid"

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"
)

// SendCreateIntent broadcasts a realtime create datapoint to the concerned clients
func (m *Module) SendCreateIntent(project, dbType, col string, req *model.CreateRequest) string {
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
			if id, ok := idTemp.(string); ok {
				msgID := uuid.NewV1().String()
				feed := &model.FeedData{Group: col, DBType: dbType, Type: utils.RealtimeWrite, TimeStamp: time.Now().Unix(), DocID: id, Payload: data}
				m.send(project, col, &Message{ID: msgID, Data: feed, Type: typeIntent})
				return msgID
			}
		}
	}
	return ""
}

// SendUpdateIntent broadcasts a realtime update datapoint to the concerned clients
func (m *Module) SendUpdateIntent(project, dbType, col string, req *model.UpdateRequest) string {
	idVar := "id"
	if dbType == string(utils.Mongo) {
		idVar = "_id"
	}

	if idTemp, p := req.Find[idVar]; p {
		if id, ok := idTemp.(string); ok {
			msgID := uuid.NewV1().String()
			feed := &model.FeedData{Group: col, DBType: dbType, Type: utils.RealtimeUpdate, TimeStamp: time.Now().Unix(), DocID: id}
			m.send(project, col, &Message{ID: msgID, Data: feed, Type: typeIntent})
			return msgID
		}
	}
	return ""
}

// SendDeleteIntent broadcasts a realtime delete datapoint to the concerned clients
func (m *Module) SendDeleteIntent(project, dbType, col string, req *model.DeleteRequest) string {
	idVar := "id"
	if dbType == string(utils.Mongo) {
		idVar = "_id"
	}

	if idTemp, p := req.Find[idVar]; p {
		if id, ok := idTemp.(string); ok {
			msgID := uuid.NewV1().String()
			feed := &model.FeedData{Group: col, Type: utils.RealtimeDelete, TimeStamp: time.Now().Unix(), DocID: id, DBType: dbType}
			m.send(project, col, &Message{ID: msgID, Data: feed, Type: typeIntent})
			return msgID
		}
	}
	return ""
}

// SendAck send an ack for the intent
func (m *Module) SendAck(msgID, project, col string, ack bool) {

	// Don't do anything if the msgID is empty
	if msgID == "" {
		return
	}

	m.send(project, col, &Message{ID: msgID, Type: typeAck, Ack: ack})
}

func (m *Module) send(project, col string, msg *Message) {
	m.RLock()
	defer m.RUnlock()

	if !m.enabled {
		return
	}

	bytes, _ := json.Marshal(msg)
	err := m.nc.Publish(getSubjectName(project, col), bytes)
	if err != nil {
		log.Println("Realtime Error:", err)
	}
}
