package eventing

import (
	"context"
	"encoding/json"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

func (m *Module) processIntents(t *time.Time) {

	// Return if module is not enabled
	if !m.IsEnabled() {
		return
	}
	m.lock.RLock()
	dbAlias, col := m.config.DBAlias, utils.TableEventingLogs
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

	results, err := m.crud.Read(ctx, dbAlias, col, &readRequest)
	if err != nil {
		logrus.Errorf("Eventing intent routine error - %s", err.Error())
		return
	}

	eventDocs := results.([]interface{})
	for _, doc := range eventDocs {

		// Parse event doc to EventDocument
		eventDoc := new(model.EventDocument)
		if err := mapstructure.Decode(doc, eventDoc); err != nil {
			logrus.Errorf("Could not covert object (%v) as intent event doc - %s", doc, err.Error())
			continue
		}

		timestamp, err := time.Parse(time.RFC3339, eventDoc.EventTimestamp) // We are using event timestamp since intent are processed wrt the time the event was created
		if err != nil {
			logrus.Errorf("Could not parse (%s) in intent event doc (%s) as time - %s", eventDoc.EventTimestamp, eventDoc.ID, err.Error())
			continue
		}

		if t.After(timestamp.Add(30 * time.Second)) {
			go m.processIntent(eventDoc)
		}
	}
}

// TODO: potential bug for prematurely processing an intent when operation is still underway e.g -> uploading a large file
func (m *Module) processIntent(eventDoc *model.EventDocument) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	// Create a context with 5 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Get the eventID
	eventID := eventDoc.ID

	switch eventDoc.Type {

	case utils.EventDBCreate:
		// Unmarshal the payload
		createEvent := model.DatabaseEventMessage{}
		_ = json.Unmarshal([]byte(eventDoc.Payload.(string)), &createEvent)

		// Check if document exists in database
		readRequest := &model.ReadRequest{Operation: utils.One, Find: createEvent.Find.(map[string]interface{})}
		if _, err := m.crud.Read(ctx, createEvent.DBType, createEvent.Col, readRequest); err != nil {

			// Mark event as cancelled if it document doesn't exist
			if err := m.crud.InternalUpdate(ctx, m.config.DBAlias, m.project, utils.TableEventingLogs, m.generateCancelEventRequest(eventID)); err != nil {
				logrus.Errorf("Eventing: Couldn't cancel intent - %s", err.Error())
			}
			return
		}

		// Mark event as staged if document does exist
		if err := m.crud.InternalUpdate(ctx, m.config.DBAlias, m.project, utils.TableEventingLogs, m.generateStageEventRequest(eventID)); err != nil {
			logrus.Errorf("Eventing: Couldn't update intent to staged - %s", err.Error())
			return
		}

		// Broadcast the event so the concerned worker can process it immediately
		eventDoc.Status = utils.EventStatusStaged
		m.transmitEvents(eventDoc.Token, []*model.EventDocument{eventDoc})

	case utils.EventDBUpdate:
		// Unmarshal the payload
		updateEvent := model.DatabaseEventMessage{}
		_ = json.Unmarshal([]byte(eventDoc.Payload.(string)), &updateEvent)

		// Get the document from the database
		timestamp := time.Now()
		readRequest := &model.ReadRequest{Operation: utils.One, Find: updateEvent.Find.(map[string]interface{})}
		result, err := m.crud.Read(ctx, updateEvent.DBType, updateEvent.Col, readRequest)
		if err != nil {
			// Do nothing if there is an error while reading
			return
		}

		// Update the payload and mark event as staged
		updateEvent.Doc = result
		data, _ := json.Marshal(updateEvent)
		if err := m.crud.InternalUpdate(ctx, m.config.DBAlias, m.project, utils.TableEventingLogs, &model.UpdateRequest{
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
			eventDoc.Status = utils.EventStatusStaged
			eventDoc.Payload = string(data)
			eventDoc.Timestamp = timestamp.Format(time.RFC3339)
			m.transmitEvents(eventDoc.Token, []*model.EventDocument{eventDoc})
		}

	case utils.EventDBDelete:
		// Unmarshal the payload
		deleteEvent := model.DatabaseEventMessage{}
		_ = json.Unmarshal([]byte(eventDoc.Payload.(string)), &deleteEvent)

		// Check if document exists in database
		readRequest := &model.ReadRequest{Operation: utils.One, Find: deleteEvent.Find.(map[string]interface{})}
		if _, err := m.crud.Read(ctx, deleteEvent.DBType, deleteEvent.Col, readRequest); err == nil {

			// Mark the event as cancelled if the document still exists
			_ = m.crud.InternalUpdate(ctx, m.config.DBAlias, m.project, utils.TableEventingLogs, m.generateCancelEventRequest(eventID))
			return
		}

		// Mark the event as staged if the document doesn't exist
		if err := m.crud.InternalUpdate(ctx, m.config.DBAlias, m.project, utils.TableEventingLogs, m.generateStageEventRequest(eventID)); err == nil {
			// Broadcast the event so the concerned worker can process it immediately
			eventDoc.Status = utils.EventStatusStaged
			m.transmitEvents(eventDoc.Token, []*model.EventDocument{eventDoc})
		}
	case utils.EventFileCreate:

		filePayload := model.FilePayload{}
		_ = json.Unmarshal([]byte(eventDoc.Payload.(string)), &filePayload)
		// Check if document exists in database

		token, err := m.auth.GetInternalAccessToken()
		if err != nil {
			logrus.Errorf("Eventing: Error generating token in intent staging - %s", err.Error())
			return
		}

		if err := m.fileStore.DoesExists(ctx, m.project, token, filePayload.Path); err != nil {

			// Mark event as cancelled if it document doesn't exist
			if err := m.crud.InternalUpdate(ctx, m.config.DBAlias, m.project, utils.TableEventingLogs, m.generateCancelEventRequest(eventID)); err != nil {
				logrus.Errorf("Eventing: Couldn't cancel intent - %s", err.Error())
			}
			return
		}

		// Mark event as staged if document does exist
		if err := m.crud.InternalUpdate(ctx, m.config.DBAlias, m.project, utils.TableEventingLogs, m.generateStageEventRequest(eventID)); err != nil {
			logrus.Errorf("Eventing: Couldn't update intent to staged - %s", err.Error())
			return
		}

		// Broadcast the event so the concerned worker can process it immediately
		eventDoc.Status = utils.EventStatusStaged
		m.transmitEvents(eventDoc.Token, []*model.EventDocument{eventDoc})

	case utils.EventFileDelete:
		filePayload := model.FilePayload{}
		_ = json.Unmarshal([]byte(eventDoc.Payload.(string)), &filePayload)

		token, err := m.auth.GetInternalAccessToken()
		if err != nil {
			logrus.Errorf("Eventing: Error generating token in intent staging - %s", err.Error())
			return
		}

		if err := m.fileStore.DoesExists(ctx, m.project, token, filePayload.Path); err == nil {
			// Mark the event as cancelled if the object still exists
			_ = m.crud.InternalUpdate(ctx, m.config.DBAlias, m.project, utils.TableEventingLogs, m.generateCancelEventRequest(eventID))
			return
		}

		// Mark the event as staged if the object doesn't exist
		if err := m.crud.InternalUpdate(ctx, m.config.DBAlias, m.project, utils.TableEventingLogs, m.generateStageEventRequest(eventID)); err == nil {
			// Broadcast the event so the concerned worker can process it immediately
			eventDoc.Status = utils.EventStatusStaged
			m.transmitEvents(eventDoc.Token, []*model.EventDocument{eventDoc})
		}

	}
}
