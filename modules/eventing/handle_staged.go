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

	// Return if module is not enabled
	if !m.config.Enabled {
		return
	}

	// Create a context with 5 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
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
		log.Println("Eventing stage routine error:", err)
		return
	}

	eventDocs := results.([]interface{})
	for _, temp := range eventDocs {
		eventDoc := new(model.EventDocument)
		if err := mapstructure.Decode(temp, eventDoc); err == nil {
			timestamp := eventDoc.Timestamp
			currentTimestamp := t.UTC().UnixNano() / int64(time.Millisecond)

			if currentTimestamp > timestamp {
				go m.processStagedEvent(eventDoc)
			}
		}
	}
}

func (m *Module) processStagedEvent(eventDoc *model.EventDocument) {
	// Create a context with 5 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Return if the event is already being processed
	if _, loaded := m.processingEvents.LoadOrStore(eventDoc.ID, true); loaded {
		return
	}

	// Delete the event from the processing list without fail
	defer m.processingEvents.Delete(eventDoc.ID)

	// Call the function to process the event
	ctxLocal, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Create a variable to track retries
	retries := 0

	for {
		result, err := m.functions.CallWithContext(ctxLocal, eventDoc.Service, eventDoc.Function, map[string]interface{}{"id": "space-cloud-eventing"}, eventDoc)
		if err == nil {
			// Check if the result is an object
			obj, ok := result.(map[string]interface{})
			if ok {
				// Check if ack is present in response
				if ackTemp, p := obj["ack"]; p {

					// Check if the ack is true
					if ack, ok := ackTemp.(bool); ok && ack {

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

						m.crud.InternalUpdate(ctx, m.config.DBType, m.project, m.config.Col, m.generateProcessedEventRequest(eventDoc.ID))
						return
					}
				}
			}
		}

		log.Println("Eventing staged event handler could not get response from service:", err)

		// Increment the retries. Exit the loop if max retries reached.
		retries++
		if retries >= eventDoc.Retries {
			// Mark event as failed
			break
		}

		// Sleep for 5 seconds
		time.Sleep(5 * time.Second)
	}

	if err := m.crud.InternalUpdate(context.TODO(), m.config.DBType, m.project, m.config.Col, m.generateFailedEventRequest(eventDoc.ID, "Max retires limit reached")); err != nil {
		log.Println("Eventing staged event handler could not update event doc:", err)
	}
}
