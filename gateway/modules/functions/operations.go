package functions

import (
	"context"
	"fmt"

	"github.com/spaceuptech/helpers"

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

// GetEndpointContextTimeout returns the endpoint timeout of particular remote-service
func (m *Module) GetEndpointContextTimeout(ctx context.Context, projectID, service, function string) (int, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	resourceID := config.GenerateResourceID(m.clusterID, projectID, config.ResourceRemoteService, service)
	if serviceVal, ok := m.config[resourceID]; ok {
		if endpoint, ok := serviceVal.Endpoints[function]; ok {
			return endpoint.Timeout, nil
		}
	}
	return 0, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Could not find endpoint (%s) for service (%s)", function, service), nil, nil)
}
