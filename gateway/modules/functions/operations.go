package functions

import (
	"context"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
)

// CallWithContext invokes function on a service. The response from the function is returned back along with
// any errors if they occurred.
func (m *Module) CallWithContext(ctx context.Context, service, function, token string, reqParams model.RequestParams, params interface{}) (int, interface{}, error) {
	status, result, err := m.handleCall(ctx, service, function, token, reqParams.Claims, params)
	if err != nil {
		return status, result, err
	}

	m.metricHook(m.project, service, function)
	return status, result, nil
}

// AddInternalRule add an internal rule to internal service
func (m *Module) AddInternalRule(service string, rule *config.Service) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.config.InternalServices[service] = rule
}
