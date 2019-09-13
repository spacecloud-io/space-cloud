package eventing

import (
	"context"
	"encoding/json"
	"errors"
	"math/rand"
	"time"

	uuid "github.com/satori/go.uuid"

	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"
)

func (m *Module) batchRequests(ctx context.Context, requests []*model.QueueEventRequest) error {
	// Create the meta information
	token := rand.Intn(utils.MaxEventTokens)
	batchID := uuid.NewV1().String()

	// Create an eventDocs array
	var eventDocs []*model.EventDocument

	// Iterate over requests
	for _, req := range requests {

		// Iterate over matching rules
		rules := m.getMatchingRules(req.Type, map[string]string{})
		for _, r := range rules {
			eventDoc := m.generateQueueEventRequest(token, r.Retries, batchID, utils.EventStatusStaged, r.Service, r.Function, req)
			eventDocs = append(eventDocs, eventDoc)
		}
	}

	// Persist the events
	createRequest := &model.CreateRequest{Document: eventDocs, Operation: utils.All}
	if err := m.crud.Create(ctx, m.config.DBType, m.project, m.config.Col, createRequest); err != nil {
		return errors.New("eventing module couldn't log the request -" + err.Error())
	}

	// Broadcast the event so the concerned worker can process it immediately
	m.broadcastEvents(eventDocs)
	return nil
}

func (m *Module) generateQueueEventRequest(token, retries int, batchID, status, service, function string, event *model.QueueEventRequest) *model.EventDocument {

	timestamp := time.Now().UTC().UnixNano() / int64(time.Millisecond)

	if event.Timestamp > timestamp {
		timestamp = event.Timestamp
	}

	// Add the delay if provided
	if event.Delay > 0 {
		timestamp += event.Delay
	}

	data, _ := json.Marshal(event.Payload)

	if retries == 0 {
		retries = 3
	}

	return &model.EventDocument{
		ID:             uuid.NewV1().String(),
		BatchID:        batchID,
		Type:           event.Type,
		Token:          token,
		Timestamp:      timestamp,
		EventTimestamp: time.Now().UTC().UnixNano() / int64(time.Millisecond),
		Payload:        string(data),
		Status:         status,
		MaxRetries:     retries,
		Retries:        0,
		Delivered:      false,
		Service:        service,
		Function:       function,
	}
}

func (m *Module) generateCancelEventRequest(eventID string) *model.UpdateRequest {
	return &model.UpdateRequest{
		Find:      map[string]interface{}{utils.GetIDVariable(m.config.DBType): eventID},
		Operation: utils.One,
		Update: map[string]interface{}{
			"$set": map[string]interface{}{"status": utils.EventStatusCancelled},
		},
	}
}

func (m *Module) generateStageEventRequest(eventID string) *model.UpdateRequest {
	return &model.UpdateRequest{
		Find:      map[string]interface{}{utils.GetIDVariable(m.config.DBType): eventID},
		Operation: utils.One,
		Update: map[string]interface{}{
			"$set": map[string]interface{}{"status": utils.EventStatusStaged},
		},
	}
}

func getCreateRows(doc interface{}, op string) []interface{} {
	var rows []interface{}
	switch op {
	case utils.One:
		rows = []interface{}{doc}
	case utils.All:
		rows = doc.([]interface{})
	default:
		rows = []interface{}{}
	}

	return rows
}

func (m *Module) getMatchingRules(name string, options map[string]string) []*config.EventingRule {
	rules := []*config.EventingRule{}

	for _, rule := range m.config.Rules {
		if rule.Type == name && isOptionsValid(rule.Options, options) {
			rules = append(rules, &rule)
		}
	}

	for _, rule := range m.config.Rules {
		if rule.Type == name && isOptionsValid(rule.Options, options) {
			rules = append(rules, &rule)
		}
	}

	return rules
}

func isOptionsValid(ruleOptions, providedOptions map[string]string) bool {
	for k, v := range providedOptions {
		if v2, p := ruleOptions[k]; !p || v != v2 {
			return false
		}
	}

	return true
}

func writeSafelyToChannel(ctx context.Context, channel chan error, err error) {
	select {
	case <-ctx.Done():
	case channel <- err:
	}
}
