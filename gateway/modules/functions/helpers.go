package functions

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

func (m *Module) handleCall(ctx context.Context, service, endpoint, token string, auth, params interface{}) (interface{}, error) {

	// Load the service rule
	s := m.loadService(service)

	// Generate the url
	s.URL = strings.TrimSuffix(s.URL, "/")
	path, err := adjustPath(s.Endpoints[endpoint].Path, params)
	if err != nil {
		return nil, err
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	url := s.URL + path

	e := s.Endpoints[endpoint]

	// Set the default method
	if e.Method == "" {
		e.Method = "POST"
	}

	// Overwrite the token if provided
	if e.Token != "" {
		token = e.Token
	}

	params, err = m.adjustBody(service, endpoint, e, auth, params)
	if err != nil {
		return nil, err
	}

	scToken, err := m.auth.GetSCAccessToken()
	if err != nil {
		return nil, err
	}

	var res interface{}
	if err := m.manager.MakeHTTPRequest(ctx, e.Method, url, token, scToken, params, &res); err != nil {
		return nil, err
	}

	return res, nil
}

func (m *Module) adjustBody(serviceID, endpointID string, endpoint config.Endpoint, auth, params interface{}) (interface{}, error) {
	// Set the default endpoint kind
	if endpoint.Kind == "" {
		endpoint.Kind = config.EndpointKindSimple
	}

	switch endpoint.Kind {
	case config.EndpointKindSimple:
		return params, nil
	case config.EndpointKindTransform:
		return goTemplate(m.templates[getGoTemplateKey(serviceID, endpointID)], endpoint.OpFormat, auth, params)
	default:
		return nil, utils.LogError(fmt.Sprintf("Invalid endpoint kind (%s) provided", endpoint.Kind), module, segmentCall, nil)
	}
}

func adjustPath(path string, params interface{}) (string, error) {
	newPath := path
	for {
		pre := strings.IndexRune(newPath, '{')
		if pre < 0 {
			return newPath, nil
		}
		post := strings.IndexRune(newPath, '}')

		key := strings.TrimSuffix(strings.TrimPrefix(newPath[pre:post], "{"), "}")
		value, err := loadParam(key, params)
		if err != nil {
			return "", err
		}

		newPath = newPath[:pre] + value + newPath[post+1:]
	}
}

func loadParam(key string, params interface{}) (string, error) {
	val, err := utils.LoadValue(key, map[string]interface{}{"args": params})
	if err != nil {
		return "", err
	}

	switch value := val.(type) {
	case float64:
		return strconv.FormatFloat(value, 'f', -1, 64), nil
	case int64, int:
		return fmt.Sprintf("%d", value), nil
	case string:
		return value, nil
	}

	return "", errors.New("invalid parameter type")
}

func (m *Module) loadService(service string) *config.Service {
	if s, p := m.config.InternalServices[service]; p {
		return s
	}

	return m.config.Services[service]
}

func getGoTemplateKey(serviceID, endpointID string) string {
	return fmt.Sprintf("%s---%s", serviceID, endpointID)
}
