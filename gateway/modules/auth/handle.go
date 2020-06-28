package auth

import (
	"context"

	"github.com/spaceuptech/space-cloud/gateway/config"
)

// AuthorizeRequest authorizes a request using the rule provided
func (m *Module) AuthorizeRequest(ctx context.Context, rule *config.Rule, project, token string, params interface{}) (map[string]interface{}, error) {
	m.RLock()
	defer m.RUnlock()

	// Return if rule is allow
	if rule.Rule == "allow" {
		return map[string]interface{}{}, nil
	}

	// Parse token
	auth, err := m.parseToken(token)
	if err != nil {
		return nil, err
	}

	args := map[string]interface{}{"auth": auth, "token": token, "params": params}
	if _, err := m.matchRule(ctx, project, rule, map[string]interface{}{"args": args}, auth); err != nil {
		return nil, err
	}

	return auth, err
}
