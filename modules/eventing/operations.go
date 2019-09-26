package eventing

import (
	uuid "github.com/satori/go.uuid"

	"github.com/spaceuptech/space-cloud/config"
)

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
