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
	// Create the meta information
	token := rand.Intn(utils.MaxEventTokens)

	// Create an eventDocs array
	var eventDocs []*model.EventDocument

	// Iterate over requests
	for _, req := range requests {

		// Iterate over matching rules
		rules := m.getMatchingRules(req.Type, map[string]string{})
		for _, r := range rules {
			eventDoc := m.generateQueueEventRequest(token, r.Retries, r.Name, batchID, utils.EventStatusStaged, r.URL, req)
			eventDocs = append(eventDocs, eventDoc)
		}
	}

	// Persist the events
	createRequest := &model.CreateRequest{Document: convertToArray(eventDocs), Operation: utils.All, IsBatch: true}
	if err := m.crud.InternalCreate(ctx, m.config.DBType, m.project, utils.TableEventingLogs, createRequest, false); err != nil {
		return errors.New("eventing module couldn't log the request -" + err.Error())
	}

	// Broadcast the event so the concerned worker can process it immediately
	m.transmitEvents(token, eventDocs)
	return nil
}

func (m *Module) generateQueueEventRequest(token, retries int, name string, batchID, status, url string, event *model.QueueEventRequest) *model.EventDocument {

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
		ID:             ksuid.New().String(),
		BatchID:        batchID,
		Type:           event.Type,
		RuleName:       name,
		Token:          token,
		Timestamp:      timestamp,
		EventTimestamp: time.Now().UTC().UnixNano() / int64(time.Millisecond),
		Payload:        string(data),
		Status:         status,
		Retries:        retries,
		URL:            url,
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
	var rules []config.EventingRule

	for n, rule := range m.config.Rules {
		if rule.Type == name && isOptionsValid(rule.Options, options) {
			rule.Name = n
			rules = append(rules, rule)
		}
	}

	for n, rule := range m.config.InternalRules {
		if rule.Type == name && isOptionsValid(rule.Options, options) {
			rule.Name = n
			rules = append(rules, rule)
		}
	}
	return rules
}

func isRulesMatching(rule1 *config.EventingRule, rule2 *config.EventingRule) bool {

	if rule1.Type != rule2.Type {
		return false
	}

	if !isOptionsValid(rule1.Options, rule2.Options) {
		return false
	}

	if rule1.URL != rule2.URL {
		return false
	}

	return true
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

func (m *Module) selectRule(name, evType string) (config.EventingRule, error) {
	if evType == utils.EventDBCreate || evType == utils.EventDBDelete || evType == utils.EventDBUpdate || evType == utils.EventFileCreate || evType == utils.EventFileDelete {
		return config.EventingRule{Timeout: 5000, Type: evType, Retries: 3}, nil
	}

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
