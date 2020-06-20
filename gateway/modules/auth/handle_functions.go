package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/spaceuptech/space-cloud/gateway/config"
)

// IsFuncCallAuthorised checks if the func call is authorised
func (m *Module) IsFuncCallAuthorised(ctx context.Context, project, service, function, token string, params interface{}) (map[string]interface{}, error) {
	m.RLock()
	defer m.RUnlock()

	rule, err := m.getFunctionRule(service, function)
	if err != nil {
		return nil, err
	}
	if rule.Rule == "allow" {
		if m.project == project {
			return map[string]interface{}{}, nil
		}
		return map[string]interface{}{}, errors.New("invalid project details provided")
	}

	var auth map[string]interface{}
	auth, err = m.parseToken(token)
	if err != nil {
		return nil, err
	}

	if _, err = m.matchRule(ctx, project, rule, map[string]interface{}{
		"args": map[string]interface{}{"auth": auth, "params": params, "token": token},
	}, auth); err != nil {
		return nil, err
	}
	return auth, nil
}

func (m *Module) getFunctionRule(service, function string) (*config.Rule, error) {
	if m.funcRules == nil {
		return nil, fmt.Errorf("no rule has been provided for endpoint (%s) in service (%s)", function, service)
	}
	if serviceStub, p := m.funcRules.InternalServices[service]; p && serviceStub.Endpoints != nil {
		if funcStub, p := serviceStub.Endpoints[function]; p && funcStub.Rule != nil {
			return funcStub.Rule, nil
		}
	} else if defaultServiceStub, p := m.funcRules.Services[service]; p && defaultServiceStub.Endpoints != nil {
		if funcStub, p := defaultServiceStub.Endpoints[function]; p && funcStub.Rule != nil {
			return funcStub.Rule, nil
		}
	}

	return nil, fmt.Errorf("no rule has been provided for endpoint (%s) in service (%s)", function, service)
}
