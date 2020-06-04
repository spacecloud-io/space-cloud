package functions

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"text/template"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

func (m *Module) handleCall(ctx context.Context, serviceID, endpointID, token string, auth, params interface{}) (interface{}, error) {
	var url string
	var method string

	// Load the service rule
	service := m.loadService(serviceID)
	serviceURL := strings.TrimSuffix(service.URL, "/")
	endpoint := service.Endpoints[endpointID]

	/***************** Set the request url ******************/

	// Adjust the endpointPath to account for variables
	endpointPath, err := adjustPath(endpoint.Path, params)
	if err != nil {
		return nil, err
	}

	switch endpoint.Kind {
	case config.EndpointKindInternal:
		if !strings.HasPrefix(endpointPath, "/") {
			endpointPath = "/" + endpointPath
		}

		url = serviceURL + endpointPath

	case config.EndpointKindExternal:
		url = endpointPath

	case config.EndpointKindPrepared:
		url = fmt.Sprintf("http://localhost:4122/v1/api/%s/graphql", m.project)

	default:
		return nil, utils.LogError(fmt.Sprintf("Invalid endpoint kind (%s) provided", endpoint.Kind), module, segmentCall, nil)
	}

	/***************** Set the request method ***************/

	// Set the default method
	if endpoint.Method == "" {
		endpoint.Method = "POST"
	}
	method = endpoint.Method

	/***************** Set the request token ****************/

	// Overwrite the token if provided
	if endpoint.Token != "" {
		token = endpoint.Token
	}

	/***************** Set the request body *****************/

	params, err = m.adjustReqBody(serviceID, endpointID, endpoint, auth, params)
	if err != nil {
		return nil, err
	}

	/******** Fire the request and get the response ********/

	scToken, err := m.auth.GetSCAccessToken()
	if err != nil {
		return nil, err
	}

	var res interface{}
	if err := m.manager.MakeHTTPRequest(ctx, method, url, token, scToken, params, &res); err != nil {
		return nil, err
	}

	/**************** Return the response body ****************/

	return m.adjustResBody(serviceID, endpointID, endpoint, auth, res)
}

func (m *Module) adjustReqBody(serviceID, endpointID string, endpoint config.Endpoint, auth, params interface{}) (interface{}, error) {
	var req, graph interface{}
	var err error

	switch endpoint.Tmpl {
	case config.EndpointTemplatingEngineGo:
		if tmpl, p := m.templates[getGoTemplateKey("request", serviceID, endpointID)]; p {
			req, err = goTemplate(tmpl, endpoint.OpFormat, auth, params)
			if err != nil {
				return nil, err
			}
		}
		if tmpl, p := m.templates[getGoTemplateKey("graph", serviceID, endpointID)]; p {
			graph, err = goTemplate(tmpl, "string", auth, params)
			if err != nil {
				return nil, err
			}
		}
	default:
		utils.LogWarn(fmt.Sprintf("Invalid templating engine (%s) provided. Skipping templating step.", endpoint.Tmpl), module, "adjust-req")
		return params, nil
	}

	switch endpoint.Kind {
	case config.EndpointKindInternal, config.EndpointKindExternal:
		if req == nil {
			return params, nil
		}
		return req, nil
	case config.EndpointKindPrepared:
		return map[string]interface{}{"query": graph, "variables": req}, nil
	default:
		return nil, utils.LogError(fmt.Sprintf("Invalid endpoint kind (%s) provided", endpoint.Kind), module, "adjust-req", nil)
	}
}

func (m *Module) adjustResBody(serviceID, endpointID string, endpoint config.Endpoint, auth, params interface{}) (interface{}, error) {
	var res interface{}
	var err error

	switch endpoint.Tmpl {
	case config.EndpointTemplatingEngineGo:
		if tmpl, p := m.templates[getGoTemplateKey("response", serviceID, endpointID)]; p {
			res, err = goTemplate(tmpl, endpoint.OpFormat, auth, params)
			if err != nil {
				return nil, err
			}
		}
	default:
		utils.LogWarn(fmt.Sprintf("Invalid templating engine (%s) provided. Skipping templating step.", endpoint.Tmpl), module, "adjust-res")
		return params, nil
	}

	if res == nil {
		return params, nil
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

func (m *Module) createGoTemplate(kind, serviceID, endpointID, tmpl string) error {
	key := getGoTemplateKey(kind, serviceID, endpointID)

	// Create a new template object
	t := template.New(key)
	t = t.Funcs(m.createGoFuncMaps())
	val, err := t.Parse(tmpl)
	if err != nil {
		return utils.LogError("Invalid golang template provided", module, segmentSetConfig, err)
	}

	m.templates[key] = val
	return nil
}

func getGoTemplateKey(kind, serviceID, endpointID string) string {
	return fmt.Sprintf("%s---%s---%s", kind, serviceID, endpointID)
}
