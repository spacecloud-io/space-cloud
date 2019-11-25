package eventing

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/mitchellh/mapstructure"

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"
)

func (m *Module) processStagedEvents(t *time.Time) {

	// Return if module is not enabled
	if !m.IsEnabled() {
		return
	}
	m.lock.RLock()
	project := m.project
	dbType, col := m.config.DBType, m.config.Col
	m.lock.RUnlock()

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

	results, err := m.crud.Read(ctx, dbType, project, col, &readRequest)
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
	m.lock.RLock()
	defer m.lock.RUnlock()

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

	// Payload will be of type json. Unmarshal it before sending
	var doc interface{}
	json.Unmarshal([]byte(eventDoc.Payload.(string)), &doc)
	eventDoc.Payload = doc

	cloudEvent := model.CloudEventPayload{SpecVersion: "1.0-rc1", Type: eventDoc.Type, Source: m.syncMan.GetEventSource(), Id: eventDoc.ID,
		Time: time.Unix(0, eventDoc.Timestamp*int64(time.Millisecond)).Format(time.RFC3339), Data: eventDoc.Payload}

	for {
		token, err := m.auth.GetInternalAccessToken()
		if err != nil {
			log.Println("Eventing: Couldn't trigger functions -", err)
			return
		}

		var eventResponse model.EventResponse
		err = m.syncMan.MakeHTTPRequest(ctxLocal, "POST", eventDoc.Url, token, cloudEvent, &eventResponse)
		if err == nil {
			var eventRequests []*model.QueueEventRequest

			// Check if response contains an event request
			if eventResponse.Event != nil {
				eventRequests = append(eventRequests, eventResponse.Event)
			}

			if eventResponse.Events != nil {
				eventRequests = append(eventRequests, eventResponse.Events...)
			}

			if len(eventRequests) > 0 {
				if err := m.batchRequests(ctx, eventRequests); err != nil {
					log.Println("Eventing: Couldn't persist events err -", err)
				}
			}

			m.crud.InternalUpdate(ctxLocal, m.config.DBType, m.project, m.config.Col, m.generateProcessedEventRequest(eventDoc.ID))
			return
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

	if err := m.crud.InternalUpdate(ctx, m.config.DBType, m.project, m.config.Col, m.generateFailedEventRequest(eventDoc.ID, "Max retires limit reached")); err != nil {
		log.Println("Eventing staged event handler could not update event doc:", err)
	}
}
