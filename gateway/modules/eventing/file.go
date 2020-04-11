package eventing

import (
	"context"
	"errors"
	"math/rand"

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

	rules := m.getMatchingRules(utils.EventFileCreate, map[string]string{})

	// Process the documents
	eventDocs := make([]*model.EventDocument, 0)
	for _, rule := range rules {
		eventDoc := m.generateQueueEventRequest(token, rule.ID,
			batchID, utils.EventStatusIntent, &model.QueueEventRequest{
				Type: utils.EventFileCreate,
				Payload: &model.FilePayload{
					Meta: req.Meta,
					Path: req.Path,
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

	rules := m.getMatchingRules(utils.EventFileDelete, map[string]string{})
	// Process the documents
	eventDocs := make([]*model.EventDocument, 0)
	for _, rule := range rules {
		eventDoc := m.generateQueueEventRequest(token, rule.ID,
			batchID, utils.EventStatusIntent, &model.QueueEventRequest{
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
