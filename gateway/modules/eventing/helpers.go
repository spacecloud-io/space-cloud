package eventing

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"text/template"
	"time"

	"github.com/fatih/structs"
	"github.com/segmentio/ksuid"
	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
	tmpl2 "github.com/spaceuptech/space-cloud/gateway/utils/tmpl"
)

func (m *Module) transmitEvents(eventToken int, eventDocs []*model.EventDocument) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	nodeID, err := m.syncMan.GetAssignedSpaceCloudID(ctx, m.project, eventToken)
	if err != nil {
		_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), "Eventing module could not get space-cloud url", err, nil)
		return
	}

	// Ignore if pubsub client has not been initialised
	if m.pubsubClient == nil {
		return
	}

	if err := m.pubsubClient.Send(ctx, getEventingTopic(nodeID), eventDocs); err != nil {
		_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), "Eventing module could not transmit event", err, nil)
	}
}

func (m *Module) getSpaceCloudIDFromBatchID(batchID string) string {
	return strings.Split(batchID, "--")[1]
}

func (m *Module) generateBatchID() string {
	return fmt.Sprintf("%s--%s", ksuid.New().String(), m.nodeID)
}

func (m *Module) batchRequests(ctx context.Context, requests []*model.QueueEventRequest, batchID string) error {
	return m.batchRequestsRaw(ctx, "", rand.Intn(utils.MaxEventTokens), requests, batchID)
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
		rules := m.getMatchingRules(ctx, req)
		for _, r := range rules {
			eventDoc := m.generateQueueEventRequest(ctx, token, r, batchID, utils.EventStatusStaged, req)
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

func (m *Module) generateQueueEventRequest(ctx context.Context, token int, rule *config.EventingTrigger, batchID, status string, event *model.QueueEventRequest) *model.EventDocument {
	return m.generateQueueEventRequestRaw(ctx, token, rule, "", batchID, status, event)
}

func (m *Module) generateQueueEventRequestRaw(ctx context.Context, token int, rule *config.EventingTrigger, eventDocID, batchID, status string, event *model.QueueEventRequest) *model.EventDocument {
	timestamp := time.Now()

	if eventDocID == "" {
		eventDocID = ksuid.New().String()
	}

	// Parse the timestamp provided
	eventTs, err := time.Parse(time.RFC3339Nano, event.Timestamp)
	if err != nil {
		// Log warning only if time stamp was provided in the request
		if event.Timestamp != "" {
			helpers.Logger.LogWarn(helpers.GetRequestID(ctx), fmt.Sprintf("Invalid timestamp format (%s) provided. Defaulting to current time.", event.Timestamp), nil)
		}
		eventTs = timestamp
	}

	// Add the delay if provided. Delay is always provided as milliseconds
	if event.Delay > 0 {
		eventTs = eventTs.Add(time.Duration(event.Delay) * time.Millisecond)
	}

	data, _ := json.Marshal(event.Payload)

	return &model.EventDocument{
		ID:          eventDocID,
		BatchID:     batchID,
		Type:        event.Type,
		RuleName:    rule.ID,
		Token:       token,
		Timestamp:   eventTs.Format(time.RFC3339Nano),
		Payload:     string(data),
		Status:      status,
		TriggerType: rule.TriggerType,
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
		_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), "Eventing was unable to queue dlq event to batch requests", err, map[string]interface{}{})
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

func (m *Module) getMatchingRules(ctx context.Context, req *model.QueueEventRequest) []*config.EventingTrigger {
	rules := make([]*config.EventingTrigger, 0)

	for _, rule := range m.config.Rules {
		// Skip trigger if its event type does not match incoming request
		if rule.Type != req.Type {
			continue
		}

		// Skip rule if the options do not match
		if !isOptionsValid(rule.Options, req.Options) {
			continue
		}

		// Skip rules if filter does not match
		if rule.Filter != nil && req.Payload != nil {
			if _, err := m.auth.MatchRule(ctx, m.project, rule.Filter, map[string]interface{}{"args": map[string]interface{}{"data": req.Payload}}, map[string]interface{}{}, model.ReturnWhereStub{}); err != nil {
				continue
			}
		}

		// Add rule to list of returned rules
		rule.TriggerType = "external"
		rules = append(rules, rule)
	}

	for _, rule := range m.config.InternalRules {
		// Skip trigger if its event type does not match incoming request
		if rule.Type != req.Type {
			continue
		}

		// Skip rule if the options do not match
		if !isOptionsValid(rule.Options, req.Options) {
			continue
		}

		// Skip rules if filter does not match
		if rule.Filter != nil {
			if _, err := m.auth.MatchRule(ctx, m.project, rule.Filter, map[string]interface{}{"args": map[string]interface{}{"data": req.Payload}}, map[string]interface{}{}, model.ReturnWhereStub{}); err != nil {
				continue
			}
		}

		// Add rule to list of returned rules
		rule.TriggerType = "internal"
		rules = append(rules, rule)
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
		arr := strings.Split(v, ",")
		if v2, p := providedOptions[k]; !p || !utils.StringExists(arr, v2) {
			return false
		}
	}
	return true
}

func (m *Module) selectRule(name string) (*config.EventingTrigger, error) {
	if rule, ok := m.config.Rules[name]; ok {
		return rule, nil
	}
	if rule, ok := m.config.InternalRules[name]; ok {
		return rule, nil
	}
	return &config.EventingTrigger{}, helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Could not find rule with name %s", name), nil, nil)
}

func (m *Module) validate(ctx context.Context, project, token string, event *model.QueueEventRequest) error {
	if event.Type == utils.EventDBCreate || event.Type == utils.EventDBDelete || event.Type == utils.EventDBUpdate || event.Type == utils.EventFileCreate || event.Type == utils.EventFileDelete {
		return fmt.Errorf("cannot create internal event (%s) with project token", event.Type)
	}

	if _, err := m.auth.IsEventingOpAuthorised(ctx, project, token, event); err != nil {
		return err
	}

	schema, p := m.schemas[event.Type]
	if !p {
		return nil
	}

	_, err := m.schema.SchemaValidator(ctx, event.Type, schema, event.Payload.(map[string]interface{}))
	return err
}

func (m *Module) createGoTemplate(kind, triggerName, tmpl string) error {
	key := getGoTemplateKey(kind, triggerName)

	// Create a new template object
	t := template.New(key)
	t = t.Funcs(tmpl2.CreateGoFuncMaps(nil))
	val, err := t.Parse(tmpl)
	if err != nil {
		return helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Invalid golang template provided", err, nil)
	}

	m.templates[key] = val
	return nil
}

func getGoTemplateKey(kind, triggerName string) string {
	return fmt.Sprintf("%s---%s", kind, triggerName)
}

func (m *Module) adjustReqBody(ctx context.Context, trigger, token string, endpoint *config.EventingTrigger, auth, params interface{}) (interface{}, error) {
	var req interface{}
	var err error

	switch endpoint.Tmpl {
	case config.TemplatingEngineGo:
		if tmpl, p := m.templates[getGoTemplateKey("trigger", trigger)]; p {
			req, err = tmpl2.GoTemplate(ctx, tmpl, endpoint.OpFormat, token, auth, params)
			if err != nil {
				return nil, err
			}
		}
	default:
		helpers.Logger.LogWarn(helpers.GetRequestID(ctx), fmt.Sprintf("Invalid templating engine (%s) provided. Skipping templating step.", endpoint.Tmpl), map[string]interface{}{"trigger": trigger})
		return params, nil
	}

	if req == nil {
		return params, nil
	}
	return req, nil
}

func (m *Module) generateWebhookToken(ctx context.Context, trigger *config.EventingTrigger, doc interface{}) (string, error) {
	var req interface{}
	var err error

	switch trigger.Tmpl {
	case config.TemplatingEngineGo:
		if tmpl, p := m.templates[getGoTemplateKey("claim", trigger.ID)]; p {
			req, err = tmpl2.GoTemplate(ctx, tmpl, trigger.OpFormat, "", nil, doc)
			if err != nil {
				return "", err
			}
		}
	default:
		helpers.Logger.LogWarn(helpers.GetRequestID(ctx), fmt.Sprintf("Invalid templating engine (%s) provided. Skipping templating step.", trigger.Tmpl), map[string]interface{}{"trigger": trigger})
		return m.auth.GetInternalAccessToken(ctx)
	}

	if req == nil {
		return m.auth.GetInternalAccessToken(ctx)
	}
	return m.auth.CreateToken(ctx, req.(map[string]interface{}))
}

func getEventingTopic(nodeID string) string {
	return fmt.Sprintf("eventing-%s", nodeID)
}

func getEventResponseTopic(nodeID string) string {
	return fmt.Sprintf("event-response-%s", nodeID)
}
