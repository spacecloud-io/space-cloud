package functions

import (
	"context"

	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/utils"
)

func (m *Module) handleCall(ctx context.Context, service, function, token string, params interface{}) (interface{}, error) {

	// Load the service rule
	s := m.loadService(service)

	// Generate the url
	url := m.manager.ResolveURL(s.Kind, s.URL, s.Scheme)
	url = url + s.Functions[function].Path

	var res interface{}
	if err := utils.MakeHTTPRequest(ctx, "POST", url, token, params, &res); err != nil {
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
