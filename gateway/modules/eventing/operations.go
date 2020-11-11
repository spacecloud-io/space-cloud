package eventing

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
)

// IsEnabled returns whether the eventing module is enabled or not
func (m *Module) IsEnabled() bool {
	m.lock.RLock()
	defer m.lock.RUnlock()

	// Return false if config isn't present
	if m.config == nil {
		return false
	}

	return m.config.Enabled
}

// QueueAdminEvent queues a new event created by the admin. This does no validation and hence must be used cautiously.
// For most use cases, consider using QueueEvent instead.
func (m *Module) QueueAdminEvent(ctx context.Context, req *model.QueueEventRequest) error {
	batchID := m.generateBatchID()

	// Prepare the find object for update and delete events
	if err := m.prepareFindObject(req); err != nil {
		return err
	}

	if err := m.batchRequests(ctx, []*model.QueueEventRequest{req}, batchID); err != nil {
		return helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to queue admin event cannot batch requests", err, nil)
	}

	m.metricHook(m.project, req.Type)
	return nil
}

// QueueEvent queues a new event
func (m *Module) QueueEvent(ctx context.Context, project, token string, req *model.QueueEventRequest) (interface{}, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	if err := m.validate(ctx, project, token, req); err != nil {
		return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to queue event validation failed", err, nil)
	}

	batchID := m.generateBatchID()

	responseChan := make(chan interface{}, 1)
	defer close(responseChan) // close channel

	m.eventChanMap.Store(batchID, eventResponse{time: time.Now(), response: responseChan})
	defer m.eventChanMap.Delete(batchID)

	if err := m.batchRequests(ctx, []*model.QueueEventRequest{req}, batchID); err != nil {
		return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to queue event cannot batch requests", err, nil)
	}

	// if true then wait for event response
	if req.IsSynchronous {
		for {
			select {
			case <-ctx.Done():
				// clear channel
				return nil, ctx.Err()
			case result := <-responseChan:
				m.metricHook(m.project, req.Type)
				return result, nil
			}
		}
	}

	m.metricHook(m.project, req.Type)
	return nil, nil
}

// ProcessEventResponseMessage sends response to client via channel
func (m *Module) ProcessEventResponseMessage(ctx context.Context, batchID string, payload interface{}) {
	// get channel from map
	value, ok := m.eventChanMap.Load(batchID)
	if !ok {
		helpers.Logger.LogWarn(helpers.GetRequestID(ctx), fmt.Sprintf("Event source (%s) not accepting any responses", batchID), nil)
		return
	}
	result := value.(eventResponse)

	// send response to client
	result.response <- payload
}

// SetRealtimeTriggers adds triggers which are used for space cloud internally
func (m *Module) SetRealtimeTriggers(eventingRules []*config.EventingTrigger) {
	m.lock.Lock()
	defer m.lock.Unlock()

	for key := range m.config.InternalRules {
		if strings.HasPrefix(key, "realtime") {
			delete(m.config.InternalRules, key)
		}
	}

	for _, incomingRule := range eventingRules {
		key := strings.Join([]string{"realtime", incomingRule.Options["db"], incomingRule.Options["col"], incomingRule.Type}, "-")
		m.config.InternalRules[key] = incomingRule
	}
}
