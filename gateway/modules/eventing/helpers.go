package eventing

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/fatih/structs"
	"github.com/segmentio/ksuid"
	"github.com/sirupsen/logrus"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

func (m *Module) transmitEvents(eventToken int, eventDocs []*model.EventDocument) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	url, err := m.syncMan.GetAssignedSpaceCloudURL(ctx, m.project, eventToken)
	if err != nil {
		logrus.Errorln("Eventing module could not get space-cloud url:", err)
		return
	}

	token, err := m.adminMan.GetInternalAccessToken()
	if err != nil {
		logrus.Errorln("Eventing module could not transmit event:", err)
		return
	}

	scToken, err := m.auth.GetSCAccessToken()
	if err != nil {
		logrus.Errorln("Eventing module could not transmit event:", err)
		return
	}

	var res interface{}
	if err := m.syncMan.MakeHTTPRequest(ctx, "POST", url, token, scToken, eventDocs, &res); err != nil {
		logrus.Errorln("Eventing module could not transmit event:", err)
		log.Println(res)
	}
}

func (m *Module) getSpaceCloudIDFromBatchID(batchID string) string {
	return strings.Split(batchID, "--")[1]
}

func (m *Module) generateBatchID() string {
	return fmt.Sprintf("%s--%s", ksuid.New().String(), m.syncMan.GetNodeID())
}

func (m *Module) batchRequests(ctx context.Context, requests []*model.QueueEventRequest, batchID string) error {
	return m.batchRequestsRaw(ctx, "", 0, requests, batchID)
}
func (m *Module) batchRequestsRaw(ctx context.Context, eventDocID string, token int, requests []*model.QueueEventRequest, batchID string) error {
	// Create the meta information
	if token == 0 {
		token = rand.Intn(utils.MaxEventTokens)
	}

	// Create an eventDocs array
	var eventDocs []*model.EventDocument

	// Iterate over requests
	for _, req := range requests {

		// Iterate over matching rules
		rules := m.getMatchingRules(req.Type, map[string]string{})
		for _, r := range rules {
			eventDoc := m.generateQueueEventRequest(token, r.ID, batchID, utils.EventStatusStaged, req)
			eventDocs = append(eventDocs, eventDoc)
		}
	}

	// Return if no docs are to be queued
	if len(eventDocs) == 0 {
		return nil
	}

	// Persist the events
	createRequest := &model.CreateRequest{Document: convertToArray(eventDocs), Operation: utils.All, IsBatch: true}
	if err := m.crud.InternalCreate(ctx, m.config.DBAlias, m.project, utils.TableEventingLogs, createRequest, false); err != nil {
		return errors.New("eventing module couldn't log the request -" + err.Error())
	}

	// Broadcast the event so the concerned worker can process it immediately
	m.transmitEvents(token, eventDocs)
	return nil
}

func (m *Module) generateQueueEventRequest(token int, name, batchID, status string, event *model.QueueEventRequest) *model.EventDocument {
	return m.generateQueueEventRequestRaw(token, name, "", batchID, status, event)
}

func (m *Module) generateQueueEventRequestRaw(token int, name, eventDocID, batchID, status string, event *model.QueueEventRequest) *model.EventDocument {
	timestamp := time.Now()

	if eventDocID == "" {
		eventDocID = ksuid.New().String()
	}

	// Parse the timestamp provided
	eventTs, err := time.Parse(time.RFC3339, event.Timestamp)
	if err != nil {
		// Log warning only if time stamp was provided in the request
		if event.Timestamp != "" {
			logrus.Warningf("Invalid timestamp format (%s) provided. Defaulting to current time.", event.Timestamp)
		}
		eventTs = timestamp
	}

	if eventTs.After(timestamp) {
		timestamp = eventTs
	}

	// Add the delay if provided. Delay is always provided as milliseconds
	if event.Delay > 0 {
		timestamp = timestamp.Add(time.Duration(event.Delay) * time.Millisecond)
	}

	data, _ := json.Marshal(event.Payload)

	return &model.EventDocument{
		ID:        eventDocID,
		BatchID:   batchID,
		Type:      event.Type,
		RuleName:  name,
		Token:     token,
		Timestamp: timestamp.Format(time.RFC3339),
		Payload:   string(data),
		Status:    status,
	}
}

func (m *Module) generateCancelEventRequest(eventID string) *model.UpdateRequest {
	return &model.UpdateRequest{
		Find:      map[string]interface{}{"_id": eventID},
		Operation: utils.All,
		Update: map[string]interface{}{
			"$set": map[string]interface{}{"status": utils.EventStatusCancelled},
		},
	}
}

func (m *Module) generateStageEventRequest(eventID string) *model.UpdateRequest {
	return &model.UpdateRequest{
		Find:      map[string]interface{}{"_id": eventID},
		Operation: utils.All,
		Update: map[string]interface{}{
			"$set": map[string]interface{}{"status": utils.EventStatusStaged},
		},
	}
}

func (m *Module) generateFailedEventRequest(eventID, remark string) *model.UpdateRequest {
	return &model.UpdateRequest{
		Find:      map[string]interface{}{"_id": eventID},
		Operation: utils.All,
		Update: map[string]interface{}{
			"$set": map[string]interface{}{"status": utils.EventStatusFailed, "remark": remark},
		},
	}
}

func (m *Module) generateProcessedEventRequest(eventID string) *model.UpdateRequest {
	return &model.UpdateRequest{
		Find:      map[string]interface{}{"_id": eventID},
		Operation: utils.All,
		Update: map[string]interface{}{
			"$set": map[string]interface{}{"status": utils.EventStatusProcessed},
		},
	}
}

func (m *Module) triggerDLQEvent(ctx context.Context, eventDoc *model.EventDocument) error {
	req := &model.QueueEventRequest{
		Type: fmt.Sprintf("%s%s", utils.DLQEventTriggerPrefix, eventDoc.RuleName),
		Payload: map[string]interface{}{
			"event_id":        eventDoc.ID,
			"event_type":      eventDoc.Type,
			"event_payload":   eventDoc.Payload,
			"event_timestamp": eventDoc.Timestamp,
			"event_name":      eventDoc.RuleName,
		},
	}

	if err := m.batchRequests(ctx, []*model.QueueEventRequest{req}, m.generateBatchID()); err != nil {
		_ = utils.LogError(fmt.Sprintf("error queueing dlq event in eventing unable to batch requests - %s", err.Error()), "eventing", "triggerDLQEvent", err)
		return err
	}

	m.metricHook(m.project, req.Type)
	return nil
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

func (m *Module) getMatchingRules(name string, options map[string]string) []config.EventingRule {
	rules := make([]config.EventingRule, 0)

	for n, rule := range m.config.Rules {
		if rule.Type == name && isOptionsValid(rule.Options, options) {
			rule.ID = n
			rules = append(rules, rule)
		}
	}

	for n, rule := range m.config.InternalRules {
		if rule.Type == name && isOptionsValid(rule.Options, options) {
			rule.ID = n
			rules = append(rules, rule)
		}
	}
	return rules
}

func convertToArray(eventDocs []*model.EventDocument) []interface{} {
	docs := make([]interface{}, len(eventDocs))

	for i, doc := range eventDocs {
		docs[i] = structs.Map(doc)
	}

	return docs
}

func isOptionsValid(ruleOptions, providedOptions map[string]string) bool {
	for k, v := range ruleOptions {
		if v2, p := providedOptions[k]; !p || v != v2 {
			return false
		}
	}
	return true
}

func (m *Module) selectRule(name string) (config.EventingRule, error) {
	if rule, ok := m.config.Rules[name]; ok {
		return rule, nil
	}
	if rule, ok := m.config.InternalRules[name]; ok {
		return rule, nil
	}
	return config.EventingRule{}, fmt.Errorf("could not find rule with name %s", name)
}

func (m *Module) validate(ctx context.Context, project, token string, event *model.QueueEventRequest) error {
	if event.Type == utils.EventDBCreate || event.Type == utils.EventDBDelete || event.Type == utils.EventDBUpdate || event.Type == utils.EventFileCreate || event.Type == utils.EventFileDelete {
		return nil
	}

	if err := m.auth.IsEventingOpAuthorised(ctx, project, token, event); err != nil {
		return err
	}

	schema, p := m.schemas[event.Type]
	if !p {
		return nil
	}

	_, err := m.schema.SchemaValidator(event.Type, schema, event.Payload.(map[string]interface{}))
	return err
}
