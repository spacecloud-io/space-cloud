package eventing

import (
	"context"
	"errors"
	"log"
	"math/rand"

	uuid "github.com/satori/go.uuid"

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"
)

// HandleCreateIntent handles the create intent request
func (m *Module) HandleCreateIntent(ctx context.Context, dbType, col string, req *model.CreateRequest) (*model.EventIntent, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	rows := getCreateRows(req.Document, req.Operation)

	// Create the meta information
	token := rand.Intn(utils.MaxEventTokens)
	batchID := uuid.NewV1().String()

	// Process the documents
	eventDocs := m.processCreateDocs(token, batchID, dbType, col, rows)

	// Persist the event intent
	createRequest := &model.CreateRequest{Document: eventDocs, Operation: utils.All}
	if err := m.crud.InternalCreate(ctx, m.config.DBType, m.project, m.config.Col, createRequest); err != nil {
		return nil, errors.New("eventing module couldn't log the request -" + err.Error())
	}

	return &model.EventIntent{BatchID: batchID, Token: token, Docs: eventDocs}, nil
}

// HandleBatchIntent handles the batch intent requests
func (m *Module) HandleBatchIntent(ctx context.Context, dbType string, req *model.BatchRequest) (*model.EventIntent, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	// Create the meta information
	token := rand.Intn(utils.MaxEventTokens)
	batchID := uuid.NewV1().String()
	var eventDocs []*model.EventDocument

	// Iterate over all batched requests
	for _, r := range req.Requests {
		switch r.Type {
		case string(utils.Create):
			// Get the rows
			rows := getCreateRows(r.Document, r.Operation)
			eventDocs := m.processCreateDocs(token, batchID, dbType, r.Col, rows)
			eventDocs = append(eventDocs, eventDocs...)

		case string(utils.Update):
			eventDocs, ok := m.processUpdateDeleteHook(token, utils.EventUpdate, batchID, dbType, r.Col, r.Find)
			if ok {
				eventDocs = append(eventDocs, eventDocs...)
			}

		case string(utils.Delete):
			eventDocs, ok := m.processUpdateDeleteHook(token, utils.EventDelete, batchID, dbType, r.Col, r.Find)
			if ok {
				eventDocs = append(eventDocs, eventDocs...)
			}

		default:
			return nil, errors.New("invalid batch request type")
		}
	}

	// Persist the event intent
	createRequest := &model.CreateRequest{Document: eventDocs, Operation: utils.All}
	if err := m.crud.InternalCreate(ctx, m.config.DBType, m.project, m.config.Col, createRequest); err != nil {
		return nil, errors.New("eventing module couldn't log the request -" + err.Error())
	}

	return &model.EventIntent{BatchID: batchID, Token: token, Docs: eventDocs}, nil
}

// HandleUpdateIntent handles the update intent requests
func (m *Module) HandleUpdateIntent(ctx context.Context, dbType, col string, req *model.UpdateRequest) (*model.EventIntent, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	return m.handleUpdateDeleteIntent(ctx, utils.EventUpdate, dbType, col, req.Find)
}

// HandleDeleteIntent handles the delete intent requests
func (m *Module) HandleDeleteIntent(ctx context.Context, dbType, col string, req *model.DeleteRequest) (*model.EventIntent, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	return m.handleUpdateDeleteIntent(ctx, utils.EventDelete, dbType, col, req.Find)
}

func (m *Module) handleUpdateDeleteIntent(ctx context.Context, eventType, dbType, col string, find map[string]interface{}) (*model.EventIntent, error) {
	// Create a unique batch id and token
	batchID := uuid.NewV1().String()
	token := rand.Intn(utils.MaxEventTokens)

	eventDocs, ok := m.processUpdateDeleteHook(token, eventType, batchID, dbType, col, find)
	if ok {
		// Persist the event intent
		createRequest := &model.CreateRequest{Document: eventDocs, Operation: utils.All}
		if err := m.crud.InternalCreate(ctx, m.config.DBType, m.project, m.config.Col, createRequest); err != nil {
			return nil, errors.New("eventing module couldn't log the request -" + err.Error())
		}

		return &model.EventIntent{BatchID: batchID, Token: token, Docs: eventDocs}, nil
	}

	return &model.EventIntent{Invalid: true}, nil
}

// HandleStage stages the event so that it can be processed
func (m *Module) HandleStage(ctx context.Context, intent *model.EventIntent, err error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	// Return if the intent is invalid
	if intent.Invalid {
		return
	}

	set := map[string]interface{}{}
	if err != nil {
		// Set the status to cancelled if error occurred
		set["status"] = utils.EventStatusCancelled
		set["remark"] = err.Error()
		intent.Invalid = true
	} else {
		// Set the status to staged if no error occurred
		set["status"] = utils.EventStatusStaged
	}

	// Create the find and update clauses
	find := map[string]interface{}{"batchid": intent.BatchID}
	update := map[string]interface{}{"$set": set}

	updateRequest := model.UpdateRequest{Find: find, Operation: utils.All, Update: update}
	if err := m.crud.InternalUpdate(ctx, m.config.DBType, m.project, m.config.Col, &updateRequest); err != nil {
		log.Println("Eventing Error: event could not be updated", err)
		return
	}

	// Broadcast the event so the concerned worker can process it immediately
	if !intent.Invalid {
		m.broadcastEvents(intent.Docs)
	}
}

func (m *Module) processCreateDocs(token int, batchID, dbType, col string, rows []interface{}) []*model.EventDocument {
	// Get event listeners
	rules := m.getMatchingRules(utils.EventCreate, map[string]string{"col": col, "db": dbType})

	var eventDocs []*model.EventDocument
	for _, doc := range rows {

		// Skip the doc if id isn't present
		idTemp, p := doc.(map[string]interface{})[utils.GetIDVariable(dbType)]
		if !p {
			continue
		}

		// Skip the doc if id isn't of type string
		docID, ok := idTemp.(string)
		if !ok {
			continue
		}
		// Iterate over all rules
		for _, rule := range rules {
			eventDocs = append(eventDocs, m.generateQueueEventRequest(token, rule.Retries,
				batchID, utils.EventStatusIntent, rule.Service, rule.Function, &model.QueueEventRequest{
					Type:    utils.EventCreate,
					Payload: model.DatabaseEventMessage{DBType: dbType, Col: col, Doc: doc, DocID: docID},
				}))
		}
	}

	return eventDocs
}

func (m *Module) processUpdateDeleteHook(token int, eventType, batchID, dbType, col string, find map[string]interface{}) ([]*model.EventDocument, bool) {
	// Get event listeners
	rules := m.getMatchingRules(eventType, map[string]string{"col": col, "db": dbType})

	// Check if id field is valid
	if idTemp, p := find[utils.GetIDVariable(dbType)]; p {
		if id, ok := utils.AcceptableIDType(idTemp); ok {

			eventDocs := make([]*model.EventDocument, len(rules))

			for _, rule := range rules {
				// Create an event doc
				eventDocs = append(eventDocs, m.generateQueueEventRequest(token, rule.Retries,
					batchID, utils.EventStatusIntent, rule.Service, rule.Function, &model.QueueEventRequest{
						Type:    eventType,
						Payload: model.DatabaseEventMessage{DBType: dbType, Col: col, DocID: id},
					}))
			}

			return eventDocs, true
		}
	}

	return nil, false
}
