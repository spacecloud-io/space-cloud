package eventing

import (
	"context"
	"errors"
	"math/rand"

	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// CreateFileIntentHook handles the create file intent request
func (m *Module) CreateFileIntentHook(ctx context.Context, req *model.CreateFileRequest) (*model.EventIntent, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	// Return if eventing module isn't enabled
	if !m.config.Enabled {
		return &model.EventIntent{Invalid: true}, nil
	}

	// Create the meta information
	token := rand.Intn(utils.MaxEventTokens)
	batchID := m.generateBatchID()

	rules := m.getMatchingRules(ctx, &model.QueueEventRequest{Type: utils.EventFileCreate, Options: map[string]string{}})

	// Process the documents
	eventDocs := make([]*model.EventDocument, 0)
	path := req.Path
	if req.Type == "file" {
		path += "/" + req.Name
	}
	for _, rule := range rules {
		eventDoc := m.generateQueueEventRequest(ctx, token, rule.ID, batchID, utils.EventStatusIntent, &model.QueueEventRequest{
			Type: utils.EventFileCreate,
			Payload: &model.FilePayload{
				Meta: req.Meta,
				Path: path,
				Type: req.Type,
			},
		})
		eventDocs = append(eventDocs, eventDoc)
	}

	if len(eventDocs) == 0 {
		return &model.EventIntent{Invalid: true}, nil
	}

	// Persist the event intent
	createRequest := &model.CreateRequest{Document: convertToArray(eventDocs), Operation: utils.All, IsBatch: true}
	if err := m.crud.InternalCreate(ctx, m.config.DBAlias, m.project, utils.TableEventingLogs, createRequest, false); err != nil {
		return nil, errors.New("eventing module couldn't log the request - " + err.Error())
	}

	return &model.EventIntent{BatchID: batchID, Token: token, Docs: eventDocs}, nil
}

// DeleteFileIntentHook handles the delete file intent requests
func (m *Module) DeleteFileIntentHook(ctx context.Context, path string, meta map[string]interface{}) (*model.EventIntent, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	// Return if eventing module isn't enabled
	if !m.config.Enabled {
		return &model.EventIntent{Invalid: true}, nil
	}

	// Create a unique batch id and token
	batchID := m.generateBatchID()
	token := rand.Intn(utils.MaxEventTokens)

	rules := m.getMatchingRules(ctx, &model.QueueEventRequest{Type: utils.EventFileDelete, Options: map[string]string{}})
	// Process the documents
	eventDocs := make([]*model.EventDocument, 0)
	for _, rule := range rules {
		eventDoc := m.generateQueueEventRequest(ctx, token, rule.ID, batchID, utils.EventStatusIntent, &model.QueueEventRequest{
			Type: utils.EventFileDelete,
			Payload: &model.FilePayload{
				Path: path,
				Meta: meta,
			},
		})
		eventDocs = append(eventDocs, eventDoc)
	}

	if len(eventDocs) == 0 {
		return &model.EventIntent{Invalid: true}, nil
	}

	// Persist the event intent
	createRequest := &model.CreateRequest{Document: convertToArray(eventDocs), Operation: utils.All, IsBatch: true}
	if err := m.crud.InternalCreate(ctx, m.config.DBAlias, m.project, utils.TableEventingLogs, createRequest, false); err != nil {
		return nil, errors.New("eventing module couldn't log the request - " + err.Error())
	}

	return &model.EventIntent{BatchID: batchID, Token: token, Docs: eventDocs}, nil
}

// HookStage stages the event so that it can be processed
func (m *Module) HookStage(ctx context.Context, intent *model.EventIntent, err error) {
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
	if err := m.crud.InternalUpdate(ctx, m.config.DBAlias, m.project, utils.TableEventingLogs, &updateRequest); err != nil {
		helpers.Logger.LogInfo(helpers.GetRequestID(ctx), "Eventing Error: event could not be updated", map[string]interface{}{"error": err})
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
