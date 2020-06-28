package eventing

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"time"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// HookDBCreateIntent handles the create intent request
func (m *Module) HookDBCreateIntent(ctx context.Context, dbAlias, col string, req *model.CreateRequest) (*model.EventIntent, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	// Return if eventing module isn't enabled
	if !m.config.Enabled {
		return &model.EventIntent{Invalid: true}, nil
	}

	rows := getCreateRows(req.Document, req.Operation)

	// Create the meta information
	token := rand.Intn(utils.MaxEventTokens)
	batchID := m.generateBatchID()

	// Process the documents
	eventDocs := m.processCreateDocs(token, batchID, dbAlias, col, rows)

	// Mark event as invalid if no events are generated
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

// HookDBBatchIntent handles the batch intent requests
func (m *Module) HookDBBatchIntent(ctx context.Context, dbAlias string, req *model.BatchRequest) (*model.EventIntent, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	// Return if eventing module isn't enabled
	if !m.config.Enabled {
		return &model.EventIntent{Invalid: true}, nil
	}

	// Create the meta information
	token := rand.Intn(utils.MaxEventTokens)
	batchID := m.generateBatchID()
	eventDocs := make([]*model.EventDocument, 0)

	// Iterate over all batched requests
	for _, r := range req.Requests {
		switch r.Type {
		case string(utils.Create):
			// Get the rows
			rows := getCreateRows(r.Document, r.Operation)
			docs := m.processCreateDocs(token, batchID, dbAlias, r.Col, rows)
			eventDocs = append(eventDocs, docs...)

		case string(utils.Update):
			docs, ok := m.processUpdateDeleteHook(token, utils.EventDBUpdate, batchID, dbAlias, r.Col, r.Find)
			if ok {
				eventDocs = append(eventDocs, docs...)
			}

		case string(utils.Delete):
			docs, ok := m.processUpdateDeleteHook(token, utils.EventDBDelete, batchID, dbAlias, r.Col, r.Find)
			if ok {
				eventDocs = append(eventDocs, docs...)
			}

		default:
			return nil, errors.New("invalid batch request type")
		}
	}

	// Mark event as invalid if no events are generated
	if len(eventDocs) == 0 {
		return &model.EventIntent{Invalid: true}, nil
	}

	// Persist the event intent
	createRequest := &model.CreateRequest{Document: convertToArray(eventDocs), Operation: utils.All, IsBatch: true}
	if err := m.crud.InternalCreate(ctx, m.config.DBAlias, m.project, utils.TableEventingLogs, createRequest, false); err != nil {
		return nil, errors.New("eventing module couldn't log the request -" + err.Error())
	}

	return &model.EventIntent{BatchID: batchID, Token: token, Docs: eventDocs}, nil
}

// HookDBUpdateIntent handles the update intent requests
func (m *Module) HookDBUpdateIntent(ctx context.Context, dbAlias, col string, req *model.UpdateRequest) (*model.EventIntent, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	// Return if eventing module isn't enabled
	if !m.config.Enabled {
		return &model.EventIntent{Invalid: true}, nil
	}

	return m.hookDBUpdateDeleteIntent(ctx, utils.EventDBUpdate, dbAlias, col, req.Find)
}

// HookDBDeleteIntent handles the delete intent requests
func (m *Module) HookDBDeleteIntent(ctx context.Context, dbAlias, col string, req *model.DeleteRequest) (*model.EventIntent, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	// Return if eventing module isn't enabled
	if !m.config.Enabled {
		return &model.EventIntent{Invalid: true}, nil
	}

	return m.hookDBUpdateDeleteIntent(ctx, utils.EventDBDelete, dbAlias, col, req.Find)
}

// hookDBUpdateDeleteIntent is used as the hook for update and delete events
func (m *Module) hookDBUpdateDeleteIntent(ctx context.Context, eventType, dbAlias, col string, find map[string]interface{}) (*model.EventIntent, error) {
	// Create a unique batch id and token
	batchID := m.generateBatchID()
	token := rand.Intn(utils.MaxEventTokens)

	eventDocs, ok := m.processUpdateDeleteHook(token, eventType, batchID, dbAlias, col, find)
	if ok {
		// Persist the event intent
		createRequest := &model.CreateRequest{Document: convertToArray(eventDocs), Operation: utils.All, IsBatch: true}
		if err := m.crud.InternalCreate(ctx, m.config.DBAlias, m.project, utils.TableEventingLogs, createRequest, false); err != nil {
			return nil, errors.New("eventing module couldn't log the request - " + err.Error())
		}

		return &model.EventIntent{BatchID: batchID, Token: token, Docs: eventDocs}, nil
	}

	return &model.EventIntent{Invalid: true}, nil
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
		log.Println("Eventing Error: event could not be updated", err)
		return
	}

	for _, doc := range intent.Docs {
		// Mark all docs as staged
		doc.Status = utils.EventStatusStaged

		// TODO: Optimise this step
		if doc.Type == utils.EventDBUpdate {
			dbEvent := new(model.DatabaseEventMessage)
			if err := json.Unmarshal([]byte(doc.Payload.(string)), dbEvent); err != nil {
				log.Println("Eventing Staging Error:", err)
				continue
			}

			req := &model.ReadRequest{
				Find:      dbEvent.Find.(map[string]interface{}),
				Operation: utils.One,
			}

			result, err := m.crud.Read(ctx, dbEvent.DBType, dbEvent.Col, req)
			if err != nil {
				log.Println("Eventing Staging Error:", err)
				continue
			}

			dbEvent.Doc = result

			data, err := json.Marshal(dbEvent)
			if err != nil {
				log.Println("Eventing Staging Error:", err)
				continue
			}

			doc.Payload = string(data)
			doc.Timestamp = time.Now().Format(time.RFC3339)

			updateRequest := model.UpdateRequest{
				Find:      map[string]interface{}{"_id": doc.ID},
				Operation: utils.All,
				Update:    map[string]interface{}{"$set": map[string]interface{}{"payload": doc.Payload}},
			}
			if err := m.crud.InternalUpdate(ctx, m.config.DBAlias, m.project, utils.TableEventingLogs, &updateRequest); err != nil {
				log.Println("Eventing Error: event could not be updated", err)
				return
			}
		}
	}

	// Broadcast the event so the concerned worker can process it immediately
	if !intent.Invalid {
		m.transmitEvents(intent.Token, intent.Docs)
	}
}

func (m *Module) processCreateDocs(token int, batchID, dbAlias, col string, rows []interface{}) []*model.EventDocument {
	// Get event listeners
	rules := m.getMatchingRules(utils.EventDBCreate, map[string]string{"col": col, "db": dbAlias})

	// Return if length of rules is zero
	if len(rules) == 0 {
		return nil
	}

	eventDocs := make([]*model.EventDocument, 0)
	for _, doc := range rows {

		findForCreate, possible := m.schema.CheckIfEventingIsPossible(dbAlias, col, doc.(map[string]interface{}), false)
		if !possible {
			return nil
		}

		// Iterate over all rules
		for _, rule := range rules {
			eventDoc := m.generateQueueEventRequest(token, rule.ID,
				batchID, utils.EventStatusIntent, &model.QueueEventRequest{
					Type:    utils.EventDBCreate,
					Payload: model.DatabaseEventMessage{DBType: dbAlias, Col: col, Doc: doc, Find: findForCreate},
				})
			eventDocs = append(eventDocs, eventDoc)
		}
	}

	return eventDocs
}

func (m *Module) processUpdateDeleteHook(token int, eventType, batchID, dbAlias, col string, find map[string]interface{}) ([]*model.EventDocument, bool) {
	// Get event listeners
	rules := m.getMatchingRules(eventType, map[string]string{"col": col, "db": dbAlias})

	// Return if length of rules is zero
	if len(rules) == 0 {
		return nil, false
	}

	findForUpdate, possible := m.schema.CheckIfEventingIsPossible(dbAlias, col, find, true)
	if !possible {
		return nil, false
	}

	eventDocs := make([]*model.EventDocument, len(rules))

	for i, rule := range rules {
		// Create an event doc
		eventDocs[i] = m.generateQueueEventRequest(token, rule.ID,
			batchID, utils.EventStatusIntent, &model.QueueEventRequest{
				Type:    eventType,
				Payload: model.DatabaseEventMessage{DBType: dbAlias, Col: col, Find: findForUpdate}, // The doc here contains the where clause
			})
	}

	// Mark event as invalid if no events are generated
	if len(eventDocs) == 0 {
		return nil, false
	}

	return eventDocs, true
}
