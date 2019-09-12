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
	token := rand.Intn(m.maxTokens)
	batchID := uuid.NewV1().String()

	// Process the documents
	eventDocs := m.processCreateDocs(token, batchID, rows)

	// Persist the event intent
	createRequest := &model.CreateRequest{Document: eventDocs, Operation: utils.All}
	if err := m.crud.Create(ctx, m.config.DBType, m.project, m.config.Col, createRequest); err != nil {
		return nil, errors.New("eventing module couldn't log the request -" + err.Error())
	}

	return &model.EventIntent{BatchID: batchID, Token: token, Docs: eventDocs}, nil
}

// HandleBatchIntent handles the batch intent requests
func (m *Module) HandleBatchIntent(ctx context.Context, dbType string, req *model.BatchRequest) (*model.EventIntent, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	// Create the meta information
	token := rand.Intn(m.maxTokens)
	batchID := uuid.NewV1().String()
	eventDocs := []interface{}{}

	// Iterate over all batched requests
	for _, r := range req.Requests {
		switch r.Type {
		case string(utils.Create):
			// Get the rows
			rows := getCreateRows(r.Document, r.Operation)
			docs := m.processCreateDocs(token, batchID, rows)
			eventDocs = append(eventDocs, docs...)

		case string(utils.Update):
			eventDoc, ok := m.processUpdateDeleteHook(token, utils.EventUpdate, batchID, dbType, r.Find)
			if ok {
				eventDocs = append(eventDocs, eventDoc)
			}

		case string(utils.Delete):
			eventDoc, ok := m.processUpdateDeleteHook(token, utils.EventDelete, batchID, dbType, r.Find)
			if ok {
				eventDocs = append(eventDocs, eventDoc)
			}

		default:
			return nil, errors.New("invalid batch request type")
		}
	}

	// Persist the event intent
	createRequest := &model.CreateRequest{Document: eventDocs, Operation: utils.All}
	if err := m.crud.Create(ctx, m.config.DBType, m.project, m.config.Col, createRequest); err != nil {
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

func (m *Module) handleUpdateDeleteIntent(ctx context.Context, eventName, dbType, col string, find map[string]interface{}) (*model.EventIntent, error) {
	// Create a unique batch id and token
	batchID := uuid.NewV1().String()
	token := rand.Intn(m.maxTokens)

	eventDoc, ok := m.processUpdateDeleteHook(token, eventName, batchID, dbType, find)
	if ok {
		// Persist the event intent
		createRequest := &model.CreateRequest{Document: eventDoc, Operation: utils.One}
		if err := m.crud.Create(ctx, m.config.DBType, m.project, m.config.Col, createRequest); err != nil {
			return nil, errors.New("eventing module couldn't log the request -" + err.Error())
		}

		return &model.EventIntent{BatchID: batchID, Token: token, Docs: []interface{}{eventDoc}}, nil
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
	} else {
		// Set the status to stagged if no error occurred
		set["status"] = utils.EventStatusStagged
	}

	// Create the find and update clauses
	find := map[string]interface{}{"batchid": intent.BatchID}
	update := map[string]interface{}{"$set": set}

	updateRequest := model.UpdateRequest{Find: find, Operation: utils.All, Update: update}
	if err := m.crud.Update(ctx, m.config.DBType, m.project, m.config.Col, &updateRequest); err != nil {
		log.Println("Eventing Error: event could not be updated", err)
		return
	}

	// Broadcast the event so the concerned worker can process the events
}
