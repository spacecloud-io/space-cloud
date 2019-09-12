package eventing

import (
	"encoding/json"
	"time"

	uuid "github.com/satori/go.uuid"

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"
)

func (m *Module) generateQueueEventRequest(token int, batchID, status string, event *model.QueueEventRequest) map[string]interface{} {

	timestamp := time.Now().UTC().UnixNano() / int64(time.Millisecond)

	if event.Timestamp > timestamp {
		timestamp = event.Timestamp
	}

	// Add the delay if provided
	if event.Delay > 0 {
		timestamp += event.Delay
	}

	data, _ := json.Marshal(event.Payload)

	retries := 3
	if event.Retries > 0 {
		retries = event.Retries
	}

	// Get the id field name
	id := utils.GetIDVariable(utils.DBType(m.config.DBType))
	doc := map[string]interface{}{
		id:                uuid.NewV1().String(),
		"batchid":         batchID,
		"name":            event.Name,
		"token":           token,
		"timestamp":       timestamp,
		"event_timestamp": time.Now().UTC().UnixNano() / int64(time.Millisecond),
		"payload":         string(data),
		"status":          status,
		"max_retries":     retries,
		"retries":         0,
	}

	return doc
}

func getCreateRows(doc interface{}, op string) []interface{} {
	var rows []interface{}
	switch op {
	case utils.One:
		rows = []interface{}{doc}
	case utils.All:
		rows = doc.([]interface{})
	default:
		rows = []interface{}{}
	}

	return rows
}

func (m *Module) processCreateDocs(token int, batchID string, rows []interface{}) []interface{} {
	eventDocs := make([]interface{}, len(rows))
	for i, doc := range rows {
		eventDocs[i] = m.generateQueueEventRequest(token, batchID, utils.EventStatusIntent, &model.QueueEventRequest{
			Name:    utils.EventCreate,
			Payload: map[string]interface{}{"doc": doc},
		})
	}

	return eventDocs
}

func (m *Module) processUpdateDeleteHook(token int, eventName, batchID, dbType string, find map[string]interface{}) (interface{}, bool) {
	// Check if id field is valid
	if idTemp, p := find[utils.GetIDVariable(utils.DBType(dbType))]; p {
		if id, ok := utils.AcceptableIDType(idTemp); ok {

			// Create an event doc
			event := m.generateQueueEventRequest(token, batchID, utils.EventStatusIntent, &model.QueueEventRequest{
				Name:    eventName,
				Payload: map[string]string{"id": id},
			})

			return event, true
		}
	}

	return nil, false
}
