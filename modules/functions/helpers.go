package functions

import (
	"context"

	"github.com/spaceuptech/space-cloud/config"
)

func (m *Module) handleCall(ctx context.Context, service, function, token string, params interface{}) (interface{}, error) {

	// Load the service rule
	s := m.loadService(service)

	// Generate the url
	url := s.URL
	url = url + s.Endpoints[function].Path

	var res interface{}
	if err := m.manager.MakeHTTPRequest(ctx, "POST", url, token, params, &res); err != nil {
		return nil, err
	}

	return res, nil
}

func (m *Module) loadService(service string) *config.Service {
	if s, p := m.config.InternalServices[service]; p {
		return s
	}

	return m.config.Services[service]
}
