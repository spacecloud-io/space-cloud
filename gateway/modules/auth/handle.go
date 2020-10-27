package auth

import (
	"context"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
)

// AuthorizeRequest authorizes a request using the rule provided
func (m *Module) AuthorizeRequest(ctx context.Context, rule *config.Rule, project, token string, args map[string]interface{}) (map[string]interface{}, error) {
	m.RLock()
	defer m.RUnlock()

	// Return if rule is allow
	if rule.Rule == "allow" {
		return map[string]interface{}{}, nil
	}

	// Parse token
	auth, err := m.jwt.ParseToken(ctx, token)
	if err != nil {
		return nil, err
	}

	args["auth"] = auth
	args["token"] = token
	if _, err := m.matchRule(ctx, project, rule, map[string]interface{}{"args": args}, auth, model.ReturnWhereStub{}); err != nil {
		return nil, err
	}

	return auth, err
}
