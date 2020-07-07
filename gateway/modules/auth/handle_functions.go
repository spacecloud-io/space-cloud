package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
)

// IsFuncCallAuthorised checks if the func call is authorised
func (m *Module) IsFuncCallAuthorised(ctx context.Context, project, service, function, token string, params interface{}) (*model.PostProcess, model.RequestParams, error) {
	m.RLock()
	defer m.RUnlock()

	rule, err := m.getFunctionRule(service, function)
	if err != nil {
		return nil, model.RequestParams{}, err
	}

	if m.project != project {
		return nil, model.RequestParams{}, errors.New("invalid project details provided")
	}

	var auth map[string]interface{}
	if rule.Rule != "allow" {
		auth, err = m.parseToken(token)
		if err != nil {
			return nil, model.RequestParams{}, err
		}
	}

	actions, err := m.matchRule(ctx, project, rule, map[string]interface{}{
		"args": map[string]interface{}{"auth": auth, "params": params, "token": token},
	}, auth)
	if err != nil {
		return nil, model.RequestParams{}, err
	}

	attr := map[string]string{"project": project, "service": service, "endpoint": function}
	reqParams := model.RequestParams{Resource: "service-call", Op: "access", Claims: auth, Attributes: attr}
	return actions, reqParams, nil
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
