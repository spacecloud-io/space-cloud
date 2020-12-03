package eventing

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/spaceuptech/helpers"

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

	attr := map[string]string{"project": m.project, "db": dbAlias, "col": col}
	reqParams := model.RequestParams{Resource: "db-read", Op: "access", Attributes: attr}
	results, err := m.crud.Read(ctx, dbAlias, col, &readRequest, reqParams)
	if err != nil {
		_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), "Eventing intent routine error", err, nil)
		return
	}

	eventDocs := results.([]interface{})
	for _, doc := range eventDocs {

		// Parse event doc to EventDocument
		eventDoc := new(model.EventDocument)
		if err := mapstructure.Decode(doc, eventDoc); err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Could not covert object (%v) as intent event doc", doc), err, nil)
			continue
		}

		timestamp, err := time.Parse(time.RFC3339Nano, eventDoc.EventTimestamp) // We are using event timestamp since intent are processed wrt the time the event was created
		if err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Could not parse (%s) in intent event doc (%s) as time", eventDoc.EventTimestamp, eventDoc.ID), err, nil)
			continue
		}

		if t.After(timestamp.Add(5 * time.Minute)) {
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
	case utils.EventFileCreate:

		filePayload := model.FilePayload{}
		_ = json.Unmarshal([]byte(eventDoc.Payload.(string)), &filePayload)
		// Check if document exists in database

		token, err := m.auth.GetInternalAccessToken(ctx)
		if err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), "Eventing: Error generating token in intent staging", err, nil)
			return
		}

		if err := m.fileStore.DoesExists(ctx, m.project, token, filePayload.Path); err != nil {

			// Mark event as cancelled if it document doesn't exist
			m.updateEventC <- &queueUpdateEvent{
				project: m.project,
				db:      m.config.DBAlias,
				col:     utils.TableEventingLogs,
				req:     m.generateCancelEventRequest(eventID),
				err:     "Eventing: Couldn't cancel intent",
			}
			return
		}

		// Mark event as staged if document does exist
		if err := m.crud.InternalUpdate(ctx, m.config.DBAlias, m.project, utils.TableEventingLogs, m.generateStageEventRequest(eventID)); err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), "Eventing: Couldn't update intent to staged", err, nil)
			return
		}

		// Broadcast the event so the concerned worker can process it immediately
		eventDoc.Status = utils.EventStatusStaged
		m.transmitEvents(eventDoc.Token, []*model.EventDocument{eventDoc})

	case utils.EventFileDelete:
		filePayload := model.FilePayload{}
		_ = json.Unmarshal([]byte(eventDoc.Payload.(string)), &filePayload)

		token, err := m.auth.GetInternalAccessToken(ctx)
		if err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), "Eventing: Error generating token in intent staging", err, nil)
			return
		}

		if err := m.fileStore.DoesExists(ctx, m.project, token, filePayload.Path); err == nil {
			// Mark the event as cancelled if the object still exists
			m.updateEventC <- &queueUpdateEvent{
				project: m.project,
				db:      m.config.DBAlias,
				col:     utils.TableEventingLogs,
				req:     m.generateCancelEventRequest(eventID),
				err:     "Eventing: Couldn't update intent to cancelled",
			}
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
