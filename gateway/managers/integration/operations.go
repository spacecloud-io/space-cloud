package integration

import (
	"context"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
)

// HandleConfigAuth handles the authentication of the config requests
func (m *Manager) HandleConfigAuth(ctx context.Context, resource, op string, claims map[string]interface{}, attr map[string]string) config.IntegrationAuthResponse {
	m.lock.RLock()
	defer m.lock.RUnlock()

	res := authResponse{checkResponse: false, err: nil}

	// Return if the request is not made by an integration
	if !isIntegrationRequest(claims) {
		return res
	}

	// Set the value of the result
	res.checkResponse = true
	res.err = m.checkPermissions(ctx, "config", resource, op, claims, attr)
	return res
}

// InvokeHook invokes all the hooks registered for the given request
func (m *Manager) InvokeHook(ctx context.Context, params model.RequestParams) config.IntegrationAuthResponse {
	m.lock.RLock()
	defer m.lock.RUnlock()

	// Don't invoke hook if request is internal
	if role, p := params.Claims["role"]; p && role == "sc-internal" {
		return authResponse{}
	}

	return m.invokeHooks(ctx, params)
}
