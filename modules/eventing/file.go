package eventing

import (
	"context"
	"errors"
	"log"
	"math/rand"

	"github.com/segmentio/ksuid"

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"
)

// HookDBCreateIntent handles the create intent request
func (m *Module) HookFileCreateIntent(ctx context.Context, req *model.CreateFileRequest) (*model.EventIntent, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	// Return if eventing module isn't enabled
	if !m.config.Enabled {
		return &model.EventIntent{Invalid: true}, nil
	}

	// Create the meta information
	token := rand.Intn(utils.MaxEventTokens)
	batchID := ksuid.New().String()

	rules := m.getMatchingRules(utils.EventFileCreate, map[string]string{})

	// Process the documents
	eventDocs := make([]*model.EventDocument, 0)
	for _, rule := range rules {
		eventDocs = append(eventDocs, m.generateQueueEventRequest(token, rule.Retries,
			batchID, utils.EventStatusIntent, rule.Url, &model.QueueEventRequest{
				Type:    utils.EventFileCreate,
				Payload: req.Path,
			}))
	}

	// Persist the event intent
	createRequest := &model.CreateRequest{Document: convertToArray(eventDocs), Operation: utils.All}
	if err := m.crud.InternalCreate(ctx, m.config.DBType, m.project, m.config.Col, createRequest); err != nil {
		return nil, errors.New("eventing module couldn't log the request - " + err.Error())
	}

	return &model.EventIntent{BatchID: batchID, Token: token, Docs: eventDocs}, nil
}

// HookDBDeleteIntent handles the delete intent requests
func (m *Module) HookFileDeleteIntent(ctx context.Context, path string) (*model.EventIntent, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	// Return if eventing module isn't enabled
	if !m.config.Enabled {
		return &model.EventIntent{Invalid: true}, nil
	}

	// Create a unique batch id and token
	batchID := ksuid.New().String()
	token := rand.Intn(utils.MaxEventTokens)

	rules := m.getMatchingRules(utils.EventFileDelete, map[string]string{})

	// Process the documents
	eventDocs := make([]*model.EventDocument, 0)
	for _, rule := range rules {
		eventDocs = append(eventDocs, m.generateQueueEventRequest(token, rule.Retries,
			batchID, utils.EventStatusIntent, rule.Url, &model.QueueEventRequest{
				Type:    utils.EventFileDelete,
				Payload: path,
			}))
	}

	// Persist the event intent
	createRequest := &model.CreateRequest{Document: convertToArray(eventDocs), Operation: utils.All}
	if err := m.crud.InternalCreate(ctx, m.config.DBType, m.project, m.config.Col, createRequest); err != nil {
		return nil, errors.New("eventing module couldn't log the request - " + err.Error())
	}

	return &model.EventIntent{BatchID: batchID, Token: token, Docs: eventDocs}, nil
}

// HookDBStage stages the event so that it can be processed
func (m *Module) HookFileStage(ctx context.Context, intent *model.EventIntent, err error) {
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

	for _, doc := range intent.Docs {
		// Mark all docs as staged
		doc.Status = utils.EventStatusStaged
	}

	// Broadcast the event so the concerned worker can process it immediately
	if !intent.Invalid {
		m.transmitEvents(intent.Token, intent.Docs)
	}
}
