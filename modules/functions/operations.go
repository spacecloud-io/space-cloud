package functions

import (
	"context"
	"time"

	"github.com/spaceuptech/space-cloud/config"
)

// Call simply calls a function on a service
func (m *Module) Call(service, function, token string, params interface{}, timeout int) (interface{}, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	return m.handleCall(ctx, service, function, token, params)
}

// CallWithContext invokes function on a service. The response from the function is returned back along with
// any errors if they occurred.
func (m *Module) CallWithContext(ctx context.Context, service, function, token string, params interface{}) (interface{}, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	return m.handleCall(ctx, service, function, token, params)
}

func (m *Module) AddInternalRule(service string, rule *config.Service) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.config.InternalServices[service] = rule
}
