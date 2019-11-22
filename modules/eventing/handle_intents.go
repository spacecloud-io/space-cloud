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

func (m *Module) processIntents(t *time.Time) {

	// Return if module is not enabled
	if !m.IsEnabled() {
		return
	}
	m.lock.RLock()
	project := m.project
	dbType, col := m.config.DBType, m.config.Col
	m.lock.RUnlock()

	// Create a context with 5 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	start, end := m.syncMan.GetAssignedTokens()

	readRequest := model.ReadRequest{Operation: utils.All, Find: map[string]interface{}{
		"status": utils.EventStatusIntent,
		"token": map[string]interface{}{
			"$gte": start,
			"$lte": end,
		},
	}}

	results, err := m.crud.Read(ctx, dbType, project, col, &readRequest)
	if err != nil {
		log.Println("Eventing intent routine error:", err)
		return
	}

	eventDocs := results.([]interface{})
	for _, doc := range eventDocs {

		// Parse event doc to EventDocument
		eventDoc := new(model.EventDocument)
		if err := mapstructure.Decode(doc, eventDoc); err != nil {
			continue
		}

		timestamp := eventDoc.EventTimestamp
		currentTimestamp := t.UTC().UnixNano() / int64(time.Millisecond)

		if currentTimestamp > timestamp+(30*1000) {
			go m.processIntent(eventDoc)
		}
	}
}

func (m *Module) processIntent(eventDoc *model.EventDocument) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	// Create a context with 5 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Get the eventID
	eventID := eventDoc.ID

	switch eventDoc.Type {

	case utils.EventCreate:
		// Unmarshal the payload
		createEvent := model.DatabaseEventMessage{}
		if err := json.Unmarshal([]byte(eventDoc.Payload.(string)), &createEvent); err != nil {
			return
		}

		idVar := utils.GetIDVariable(createEvent.DBType)

		// Check if document exists in database
		readRequest := &model.ReadRequest{Operation: utils.One, Find: map[string]interface{}{idVar: createEvent.DocID}}
		if _, err := m.crud.Read(ctx, createEvent.DBType, m.project, createEvent.Col, readRequest); err != nil {

			// Mark event as cancelled if it document doesn't exist
			if err := m.crud.InternalUpdate(ctx, m.config.DBType, m.project, m.config.Col, m.generateCancelEventRequest(eventID)); err != nil {
				log.Println("Eventing: Couldn't cancel intent -", err)
			}
			return
		}

		// Mark event as staged if document does exist
		if err := m.crud.InternalUpdate(ctx, m.config.DBType, m.project, m.config.Col, m.generateStageEventRequest(eventID)); err != nil {
			log.Println("Eventing: Couldn't update intent to staged -", err)
			return
		}

		// Broadcast the event so the concerned worker can process it immediately
		eventDoc.Status = utils.EventStatusProcessed
		m.transmitEvents(eventDoc.Token, []*model.EventDocument{eventDoc})

	case utils.EventUpdate:
		// Unmarshal the payload
		updateEvent := model.DatabaseEventMessage{}
		json.Unmarshal([]byte(eventDoc.Payload.(string)), &updateEvent)
		idVar := utils.GetIDVariable(updateEvent.DBType)

		// Get the document from the database
		timestamp := time.Now().UTC().UnixNano() / int64(time.Millisecond)
		readRequest := &model.ReadRequest{Operation: utils.One, Find: map[string]interface{}{idVar: updateEvent.DocID}}
		result, err := m.crud.Read(ctx, updateEvent.DBType, m.project, updateEvent.Col, readRequest)
		if err != nil {
			// Do nothing if there is an error while reading
			return
		}

		// Update the payload and mark event as staged
		updateEvent.Doc = result
		data, _ := json.Marshal(updateEvent)
		if err := m.crud.InternalUpdate(ctx, m.config.DBType, m.project, m.config.Col, &model.UpdateRequest{
			Find: map[string]interface{}{"_id": eventID},
			Update: map[string]interface{}{
				"$set": map[string]interface{}{
					"status":    utils.EventStatusStaged,
					"payload":   string(data),
					"timestamp": timestamp,
				},
			},
		}); err == nil {
			// Broadcast the event so the concerned worker can process it immediately
			eventDoc.Status = utils.EventStatusProcessed
			eventDoc.Payload = string(data)
			eventDoc.Timestamp = timestamp
			m.transmitEvents(eventDoc.Token, []*model.EventDocument{eventDoc})
		}

	case utils.EventDelete:
		// Unmarshal the payload
		deleteEvent := model.DatabaseEventMessage{}
		json.Unmarshal([]byte(eventDoc.Payload.(string)), &deleteEvent)
		idVar := utils.GetIDVariable(deleteEvent.DBType)

		// Check if document exists in database
		readRequest := &model.ReadRequest{Operation: utils.One, Find: map[string]interface{}{idVar: deleteEvent.DocID}}
		if _, err := m.crud.Read(ctx, deleteEvent.DBType, m.project, deleteEvent.Col, readRequest); err == nil {

			// Mark the event as cancelled if the document still exists
			m.crud.InternalUpdate(ctx, m.config.DBType, m.project, m.config.Col, m.generateCancelEventRequest(eventID))
			return
		}

		// Mark the event as staged if the document doesn't exist
		if err := m.crud.InternalUpdate(ctx, m.config.DBType, m.project, m.config.Col, m.generateStageEventRequest(eventID)); err == nil {
			// Broadcast the event so the concerned worker can process it immediately
			eventDoc.Status = utils.EventStatusProcessed
			m.transmitEvents(eventDoc.Token, []*model.EventDocument{eventDoc})
		}

	}
}
