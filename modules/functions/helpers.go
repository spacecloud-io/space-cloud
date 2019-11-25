package functions

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/utils"
)

func (m *Module) handleCall(ctx context.Context, service, endpoint, token string, params interface{}) (interface{}, error) {

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

	method := s.Endpoints[endpoint].Method
	if method == "" {
		method = "POST"
	}

	var res interface{}
	if err := m.manager.MakeHTTPRequest(ctx, method, url, token, params, &res); err != nil {
		return nil, err
	}

	return res, nil
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
