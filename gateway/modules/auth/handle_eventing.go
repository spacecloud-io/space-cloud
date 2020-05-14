package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
)

// IsEventingOpAuthorised checks if the eventing operation is authorised
func (m *Module) IsEventingOpAuthorised(ctx context.Context, project, token string, event *model.QueueEventRequest) error {
	m.RLock()
	defer m.RUnlock()

	rule, err := m.getEventingRule(event.Type)
	if err != nil {
		return err
	}

	if rule.Rule == "allow" {
		if m.project == project {
			return nil
		}
		return errors.New("invalid project details provided")
	}

	var auth map[string]interface{}
	auth, err = m.parseToken(token)
	if err != nil {
		return err
	}

	if _, err = m.matchRule(ctx, project, rule, map[string]interface{}{
		"args": map[string]interface{}{"auth": auth, "params": event.Payload, "token": token},
	}, auth); err != nil {
		return err
	}

	return nil
}

func (m *Module) getEventingRule(eventType string) (*config.Rule, error) {
	if m.eventingRules == nil {
		return nil, fmt.Errorf("rule not found for given event type (%s)", eventType)
	}
	if rule, p := m.eventingRules[eventType]; p {
		return rule, nil
	}
	if rule, p := m.eventingRules["default"]; p {
		return rule, nil
	}
	return nil, fmt.Errorf("rule not found for given event type (%s)", eventType)
}
