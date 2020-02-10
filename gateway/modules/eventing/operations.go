package eventing

import (
	"context"

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
func (m *Module) QueueEvent(ctx context.Context, project, token string, req *model.QueueEventRequest) error {
	m.lock.RLock()
	defer m.lock.RUnlock()

	err := m.validate(ctx, project, token, req)
	if err != nil {
		return err
	}

	return m.batchRequests(ctx, []*model.QueueEventRequest{req})
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
