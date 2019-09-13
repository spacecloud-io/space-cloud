package eventing

import (
	"context"
	"log"
	"time"

	"github.com/mitchellh/mapstructure"

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"
)

func (m *Module) processStagedEvents(t *time.Time) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	// Create a context with 5 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	start, end := m.syncMan.GetAssignedTokens()

	readRequest := model.ReadRequest{Operation: utils.All, Find: map[string]interface{}{
		"status": utils.EventStatusStaged,
		"token": map[string]interface{}{
			"$gte": start,
			"$lte": end,
		},
	}}

	dbType, col := m.config.DBType, m.config.Col

	results, err := m.crud.Read(ctx, dbType, m.project, col, &readRequest)
	if err != nil {
		log.Fatalln("Eventing stage routine error:", err)
		return
	}

	eventDocs := results.([]interface{})
	for _, temp := range eventDocs {
		eventDoc := new(model.EventDocument)
		if err := mapstructure.Decode(temp, eventDoc); err == nil {
			m.processStagedEvent(ctx, eventDoc)
		}
	}
}

func (m *Module) processStagedEvent(ctx context.Context, eventDoc *model.EventDocument) {
	// Return if the event is already being processed
	if _, loaded := m.processingEvents.LoadOrStore(eventDoc.ID, true); loaded {
		return
	}

	// Delete the event from the processing list without fail
	defer m.processingEvents.Delete(eventDoc.ID)

	// Call the function to process the event
	result, err := m.functions.CallWithContext(ctx, eventDoc.Service, eventDoc.Function, map[string]interface{}{"id": "space-cloud"}, eventDoc)
	if err != nil {
		log.Println("Eventing staged event handler could not get response from service:", err)
		return
	}

	// Return if the result is not an object
	obj, ok := result.(map[string]interface{})
	if !ok {
		return
	}

	if ackTemp, p := obj["ack"]; p {
		if ack, ok := ackTemp.(bool); ack && ok {
			// Check if response contains an event request
			var eventRequests []*model.QueueEventRequest
			if item, p := obj["event"]; p {
				req := new(model.QueueEventRequest)
				if err := mapstructure.Decode(item, req); err == nil {
					eventRequests = append(eventRequests, req)
				}
			}

			if items, p := obj["events"]; p {
				array := items.([]interface{})
				for _, item := range array {
					req := new(model.QueueEventRequest)
					if err := mapstructure.Decode(item, req); err == nil {
						eventRequests = append(eventRequests, req)
					}
				}
			}

			if len(eventRequests) > 0 {
				if err := m.batchRequests(ctx, eventRequests); err != nil {
					log.Println("Eventing: Couldn't persist events err -", err)
				}
			}
			return
		}
	}
}
