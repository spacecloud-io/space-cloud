package eventing

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/fatih/structs"
	"github.com/mitchellh/mapstructure"
	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

func (m *Module) processStagedEvents(t *time.Time) {

	// Return if module is not enabled
	if !m.IsEnabled() {
		return
	}
	m.lock.RLock()
	dbAlias, col := m.config.DBAlias, utils.TableEventingLogs
	m.lock.RUnlock()

	// Create a context with 5 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	start, end := m.syncMan.GetAssignedTokens()

	readRequest := model.ReadRequest{Operation: utils.All, Options: &model.ReadOptions{Sort: []string{"ts"}, Limit: &limit}, Find: map[string]interface{}{
		"status": utils.EventStatusStaged,
		"token": map[string]interface{}{
			"$gte": start,
			"$lte": end,
		},
	}}

	attr := map[string]string{"project": m.project, "db": dbAlias, "col": col}
	reqParams := model.RequestParams{Resource: "db-read", Op: "access", Attributes: attr}
	results, err := m.crud.Read(ctx, dbAlias, col, &readRequest, reqParams)
	if err != nil {
		_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), "Eventing stage routine error", err, nil)
		return
	}

	eventDocs := results.([]interface{})
	for _, temp := range eventDocs {
		eventDoc := new(model.EventDocument)
		if err := mapstructure.Decode(temp, eventDoc); err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Could not covert object (%v) as staged event doc", temp), err, nil)
			continue
		}

		timestamp, err := time.Parse(time.RFC3339Nano, eventDoc.Timestamp) // We are using event timestamp since intent are processed wrt the time the event was created
		if err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Could not parse time (%s) in staged event doc (%s) as time", eventDoc.Timestamp, eventDoc.ID), err, nil)
			continue
		}

		timestamp = timestamp.Add(15 * time.Second)

		if t.After(timestamp) || t.Equal(timestamp) {
			go m.processStagedEvent(eventDoc)
		}
	}
}

func (m *Module) processStagedEvent(eventDoc *model.EventDocument) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	// Return if the event is already being processed
	if _, loaded := m.processingEvents.LoadOrStore(eventDoc.ID, true); loaded {
		return
	}

	// Delete the event from the processing list without fail
	defer m.processingEvents.Delete(eventDoc.ID)

	eventType, triggerName := eventDoc.Type, eventDoc.RuleName

	rule, err := m.selectRule(triggerName)
	if err != nil {
		_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Error processing staged event", err, nil)
		return
	}

	var maxRetries int
	if rule.Retries > 0 {
		maxRetries = rule.Retries
	} else {
		maxRetries = 3
	}

	if rule.Timeout == 0 {
		rule.Timeout = 5000
	}

	// Call the function to process the event
	timeoutLocal := time.Duration(5000*maxRetries*rule.Timeout) * time.Millisecond

	ctx, cancel := context.WithTimeout(context.Background(), timeoutLocal)
	defer cancel()

	// Create a variable to track retries
	retries := 0

	// Payload will be of type json. Unmarshal it before sending
	var doc interface{}
	_ = json.Unmarshal([]byte(eventDoc.Payload.(string)), &doc)
	eventDoc.Payload = doc

	cloudEvent := model.CloudEventPayload{SpecVersion: "1.0", Type: eventType, Source: m.syncMan.GetEventSource(), ID: eventDoc.ID,
		Time: eventDoc.Timestamp, Data: eventDoc.Payload}

	doc = structs.Map(&cloudEvent)
	newDoc, err := m.adjustReqBody(ctx, triggerName, "", rule, nil, doc)
	if err != nil {
		if err := m.logInvocation(ctx, eventDoc.ID, []byte("{}"), 0, "", err.Error()); err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), "eventing module couldn't log the invocation ", err, nil)
			return
		}
		if err := m.crud.InternalUpdate(context.Background(), m.config.DBAlias, m.project, utils.TableEventingLogs, m.generateFailedEventRequest(eventDoc.ID, "Max retires limit reached")); err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), "Eventing staged event handler could not update event doc ", err, nil)
		}
		_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to adjust request body according to template for trigger (%s)", triggerName), err, nil)
		return
	}

	// Generate the token
	token, err := m.generateWebhookToken(ctx, rule, doc, newDoc)
	if err != nil {
		if err := m.logInvocation(ctx, eventDoc.ID, []byte("{}"), 0, "", err.Error()); err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), "eventing module couldn't log the invocation ", err, nil)
			return
		}
		if err := m.crud.InternalUpdate(context.Background(), m.config.DBAlias, m.project, utils.TableEventingLogs, m.generateFailedEventRequest(eventDoc.ID, "Max retires limit reached")); err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), "Eventing staged event handler could not update event doc ", err, nil)
		}
		_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), "error invoking web hook in eventing unable to get internal access token", err, nil)
		return
	}

	for {
		if err := m.invokeWebhook(ctx, token, &http.Client{}, rule, eventDoc, newDoc); err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), "Eventing staged event handler could not get response from service", err, nil)

			// Increment the retries. Exit the loop if max retries reached.
			retries++
			if retries >= maxRetries+1 {
				// Mark event as failed
				break
			}

			// Sleep for 5 seconds
			time.Sleep(5 * time.Second)
			continue
		}

		// Reaching here means the event was successfully processed. Let's simply return
		return
	}

	if err := m.triggerDLQEvent(ctx, eventDoc); err != nil {
		_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Couldn't create DLQ event for event id %v", eventDoc.ID), err, nil)
	}
	if err := m.crud.InternalUpdate(context.Background(), m.config.DBAlias, m.project, utils.TableEventingLogs, m.generateFailedEventRequest(eventDoc.ID, "Max retires limit reached")); err != nil {
		_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), "Eventing staged event handler could not update event doc", err, nil)
	}
}

func (m *Module) invokeWebhook(ctx context.Context, token string, client model.HTTPEventingInterface, rule *config.EventingTrigger, eventDoc *model.EventDocument, params interface{}) error {
	ctxLocal, cancel := context.WithTimeout(ctx, time.Duration(rule.Timeout)*time.Millisecond)
	defer cancel()

	scToken, err := m.auth.GetSCAccessToken(ctx)
	if err != nil {
		return helpers.Logger.LogError(helpers.GetRequestID(ctx), "error invoking web hook in eventing unable to get sc access token", err, nil)
	}

	var eventResponse model.EventResponse
	if err := m.MakeInvocationHTTPRequest(ctxLocal, client, http.MethodPost, rule.URL, eventDoc.ID, token, scToken, params, &eventResponse); err != nil {
		return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("error invoking web hook in eventing unable to send http request to url %s", rule.URL), err, nil)
	}

	// Check if response contains an error
	if eventResponse.Error != "" {
		return errors.New(eventResponse.Error)
	}

	var eventRequests []*model.QueueEventRequest
	// Check if response contains an event request
	if eventResponse.Event != nil {
		eventRequests = append(eventRequests, eventResponse.Event)
	}

	if eventResponse.Events != nil {
		eventRequests = append(eventRequests, eventResponse.Events...)
	}

	if eventResponse.Response != nil {
		if m.pubsubClient != nil {
			sendTopic := getEventResponseTopic(m.getSpaceCloudIDFromBatchID(eventDoc.BatchID))
			err = m.pubsubClient.Send(ctxLocal, sendTopic, model.EventResponseMessage{BatchID: eventDoc.BatchID, Response: eventResponse.Response})
			if err != nil {
				return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("error invoking web hook in eventing unable to send http request for synchronous response to node %s", sendTopic), err, nil)
			}
		}
	}

	if len(eventRequests) > 0 {
		if err := m.batchRequests(ctx, eventRequests, eventDoc.BatchID); err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), "error invoking web hook in eventing unable to persist events off", err, nil)
		}
	}

	_ = m.crud.InternalUpdate(ctxLocal, m.config.DBAlias, m.project, utils.TableEventingLogs, m.generateProcessedEventRequest(eventDoc.ID))
	return nil
}
