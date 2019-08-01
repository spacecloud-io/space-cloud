package auth

import (
	"errors"

	"github.com/spaceuptech/space-cloud/config"
)

// IsFuncCallAuthorised checks if the func call is authorised
func (m *Module) IsFuncCallAuthorised(project, service, function, token string, params interface{}) (TokenClaims, error) {
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

	auth, err := m.parseToken(token)
	if err != nil {
		return nil, err
	}

	if err = m.matchRule(project, rule, map[string]interface{}{
		"args": map[string]interface{}{"auth": auth, "params": params},
	}, auth); err != nil {
		return nil, err
	}

	return auth, nil
}

func (m *Module) getFunctionRule(service, function string) (*config.Rule, error) {
	if m.funcRules == nil {
		return nil, ErrRuleNotFound
	}

	if service, p := m.funcRules[service]; p && service.Functions != nil {
		if funcStub, p := service.Functions[function]; p {
			return funcStub.Rule, nil
		}
	}

	return nil, ErrRuleNotFound
}
