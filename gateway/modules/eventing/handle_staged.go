package eventing

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"

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

	readRequest := model.ReadRequest{Operation: utils.All, Find: map[string]interface{}{
		"status": utils.EventStatusStaged,
		"token": map[string]interface{}{
			"$gte": start,
			"$lte": end,
		},
	}}

	results, err := m.crud.Read(ctx, dbAlias, col, &readRequest)
	if err != nil {
		logrus.Errorf("Eventing stage routine error - %s", err.Error())
		return
	}

	eventDocs := results.([]interface{})
	for _, temp := range eventDocs {
		eventDoc := new(model.EventDocument)
		if err := mapstructure.Decode(temp, eventDoc); err != nil {
			logrus.Errorf("Could not covert object (%v) as staged event doc - %s", temp, err.Error())
			continue
		}

		timestamp, err := time.Parse(time.RFC3339, eventDoc.Timestamp) // We are using event timestamp since intent are processed wrt the time the event was created
		if err != nil {
			logrus.Errorf("Could not parse time (%s) in staged event doc (%s) as time - %s", eventDoc.Timestamp, eventDoc.ID, err.Error())
			continue
		}

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

	evType, name := eventDoc.Type, eventDoc.RuleName

	rule, err := m.selectRule(name)
	if err != nil {
		logrus.Errorln("Error processing staged event:", err)
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

	cloudEvent := model.CloudEventPayload{SpecVersion: "1.0-rc1", Type: evType, Source: m.syncMan.GetEventSource(), ID: eventDoc.ID,
		Time: eventDoc.Timestamp, Data: eventDoc.Payload}

	for {
		if err := m.invokeWebhook(ctx, &http.Client{}, rule, eventDoc, &cloudEvent); err != nil {
			logrus.Errorf("Eventing staged event handler could not get response from service -%s", err.Error())

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
		_ = utils.LogError(fmt.Sprintf("Couldn't create DLQ event for event id %v", eventDoc.ID), "eventing", "triggerDLQEvent", err)
	}
	if err := m.crud.InternalUpdate(context.Background(), m.config.DBAlias, m.project, utils.TableEventingLogs, m.generateFailedEventRequest(eventDoc.ID, "Max retires limit reached")); err != nil {
		logrus.Errorf("Eventing staged event handler could not update event doc - %s", err.Error())
	}
}

func (m *Module) invokeWebhook(ctx context.Context, client model.HTTPEventingInterface, rule config.EventingRule, eventDoc *model.EventDocument, cloudEvent *model.CloudEventPayload) error {
	ctxLocal, cancel := context.WithTimeout(ctx, time.Duration(rule.Timeout)*time.Millisecond)
	defer cancel()
	internalToken, err := m.auth.GetInternalAccessToken()
	if err != nil {
		logrus.Errorf("error invoking web hook in eventing unable to get internal access token - %s", err.Error())
		return err
	}

	scToken, err := m.auth.GetSCAccessToken()
	if err != nil {
		logrus.Errorf("error invoking web hook in eventing unable to get sc access token - %s", err.Error())
		return err
	}

	var eventResponse model.EventResponse
	if err := m.MakeInvocationHTTPRequest(ctxLocal, client, http.MethodPost, rule.URL, eventDoc.ID, internalToken, scToken, cloudEvent, &eventResponse); err != nil {
		logrus.Errorf("error invoking web hook in eventing unable to send http request to url %s - %s", rule.URL, err.Error())
		return err
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
		url, err := m.syncMan.GetSpaceCloudURLFromID(m.getSpaceCloudIDFromBatchID(eventDoc.BatchID))
		if err != nil {
			logrus.Errorf("error invoking web hook in eventing unable to get sc addr from batchID %s - %s", eventDoc.BatchID, err)
			return err
		}
		url = fmt.Sprintf("http://%s/v1/api/%s/eventing/process-event-response", url, m.project)
		if err := m.syncMan.MakeHTTPRequest(ctxLocal, http.MethodPost, url, internalToken, scToken, map[string]interface{}{"batchID": eventDoc.BatchID, "response": eventResponse.Response}, &map[string]interface{}{}); err != nil {
			logrus.Errorf("error invoking web hook in eventing unable to send http request for synchronous response to url %s - %s", url, err.Error())
			return err
		}
	}

	if len(eventRequests) > 0 {
		if err := m.batchRequests(ctx, eventRequests, eventDoc.BatchID); err != nil {
			logrus.Errorf("error invoking web hook in eventing unable to persist events off - %s", err.Error())
		}
	}

	_ = m.crud.InternalUpdate(ctxLocal, m.config.DBAlias, m.project, utils.TableEventingLogs, m.generateProcessedEventRequest(eventDoc.ID))
	return nil
}
