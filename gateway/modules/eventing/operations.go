package eventing

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/segmentio/ksuid"

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

// QueueEvent queues a new event
func (m *Module) QueueEvent(ctx context.Context, project, token string, req *model.QueueEventRequest) (interface{}, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	err := m.validate(ctx, project, token, req)
	if err != nil {
		logrus.Errorf("error queueing event in eventing unable to validate - %s", err.Error())
		return nil, err
	}

	batchID := m.generateBatchID()
	responseChan := make(chan interface{}, 1)
	defer close(responseChan) // close channel
	m.eventChanMap.Store(batchID, eventResponse{time: time.Now(), response: responseChan})
	defer m.eventChanMap.Delete(batchID)

	if err = m.batchRequests(ctx, []*model.QueueEventRequest{req}, batchID); err != nil {
		logrus.Errorf("error queueing event in eventing unable to batch requests - %s", err.Error())
		return nil, err
	}

	// if true then wait for event response
	if req.IsSynchronous {
		for {
			select {
			case <-ctx.Done():
				// clear channel
				return nil, ctx.Err()
			case result := <-responseChan:
				return result, nil
			}
		}
	}
	return nil, nil
}

// SendEventResponse sends response to client via channel
func (m *Module) SendEventResponse(batchID string, payload interface{}) {
	// get channel from map
	value, ok := m.eventChanMap.Load(batchID)
	if !ok {
		logrus.Errorf("error sending synchronous event response to client unable to find channel in map for batch %s", batchID)
		return
	}
	result := value.(eventResponse)

	// send response to client
	result.response <- payload
}

// AddInternalRules adds triggers which are used for space cloud internally
func (m *Module) AddInternalRules(eventingRules []config.EventingRule) {
	m.lock.Lock()
	defer m.lock.Unlock()

	for _, incomingRule := range eventingRules {
		isPresent := false
		for _, storedRule := range m.config.InternalRules {

			// Add the rule for the only if it doesn't already exist
			if isRulesMatching(&storedRule, &incomingRule) {
				isPresent = true
				break
			}
		}

		if !isPresent {
			key := ksuid.New().String()
			m.config.InternalRules[key] = incomingRule
		}
	}
}
