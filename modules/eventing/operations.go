package eventing

import (
	"context"

	uuid "github.com/satori/go.uuid"

	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/model"
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
func (m *Module) QueueEvent(ctx context.Context, req *model.QueueEventRequest) error {
	m.lock.RLock()
	defer m.lock.RUnlock()

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
			key := uuid.NewV1().String()
			m.config.InternalRules[key] = incomingRule
		}
	}

}
