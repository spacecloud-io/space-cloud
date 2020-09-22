package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// IsEventingOpAuthorised checks if the eventing operation is authorised
func (m *Module) IsEventingOpAuthorised(ctx context.Context, project, token string, event *model.QueueEventRequest) (model.RequestParams, error) {
	m.RLock()
	defer m.RUnlock()

	rule, err := m.getEventingRule(ctx, event.Type)
	if err != nil {
		return model.RequestParams{}, err
	}

	if m.project != project {
		return model.RequestParams{}, errors.New("invalid project details provided")
	}

	var auth map[string]interface{}
	if rule.Rule != "allow" {
		auth, err = m.jwt.ParseToken(ctx, token)
		if err != nil {
			return model.RequestParams{}, err
		}
	}

	// Check if internal token
	if auth != nil {
		if id, p := auth["id"]; p && id == utils.InternalUserID {
			hookResponse := m.integrationMan.InvokeHook(ctx, model.RequestParams{
				Claims:     auth,
				Resource:   "internal-api-access",
				Op:         "eventing-queue",
				Attributes: map[string]string{"project": project},
			})
			if hookResponse.CheckResponse() {
				attr := map[string]string{"project": project, "type": event.Type}
				return model.RequestParams{Claims: auth, Resource: "eventing-queue", Op: "access", Attributes: attr}, hookResponse.Error()
			}
		}
	}

	if _, err = m.matchRule(ctx, project, rule, map[string]interface{}{
		"args": map[string]interface{}{"auth": auth, "params": event.Payload, "token": token},
	}, auth); err != nil {
		return model.RequestParams{}, err
	}

	attr := map[string]string{"project": project, "type": event.Type}
	return model.RequestParams{Claims: auth, Resource: "eventing-queue", Op: "access", Attributes: attr}, nil
}

func (m *Module) getEventingRule(ctx context.Context, eventType string) (*config.Rule, error) {
	if m.eventingRules == nil {
		return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Security rules not initialized for event type (%s)", eventType), nil, nil)
	}
	if rule, p := m.eventingRules[eventType]; p {
		return rule, nil
	}
	if rule, p := m.eventingRules["default"]; p {
		return rule, nil
	}
	return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("No security rule provided for event type (%s)", eventType), nil, nil)
}
