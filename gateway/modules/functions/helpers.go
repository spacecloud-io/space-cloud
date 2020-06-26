package functions

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"text/template"

	tmpl2 "github.com/spaceuptech/space-cloud/gateway/utils/tmpl"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

func (m *Module) handleCall(ctx context.Context, serviceID, endpointID, token string, auth, params interface{}) (interface{}, error) {
	var url string
	var method string
	var ogToken string

	// Load the service rule
	service := m.loadService(serviceID)
	serviceURL := strings.TrimSuffix(service.URL, "/")
	endpoint := service.Endpoints[endpointID]
	ogToken = token

	/***************** Set the request url ******************/

	// Adjust the endpointPath to account for variables
	endpointPath, err := adjustPath(endpoint.Path, auth, params)
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

	// Overwrite the token if provided
	if endpoint.Token != "" {
		token = endpoint.Token
	}

	/***************** Set the request body *****************/

	newParams, err := m.adjustReqBody(serviceID, endpointID, ogToken, endpoint, auth, params)
	if err != nil {
		return nil, err
	}

	/******** Fire the request and get the response ********/

	scToken, err := m.auth.GetSCAccessToken()
	if err != nil {
		return nil, err
	}

	var res interface{}
	req := &utils.HTTPRequest{
		Params: newParams,
		Method: method, URL: url,
		Token: token, SCToken: scToken,
		Headers: prepareHeaders(endpoint, ogToken, auth, params),
	}
	if err := utils.MakeHTTPRequest(ctx, req, &res); err != nil {
		return nil, err
	}

	/**************** Return the response body ****************/

	return m.adjustResBody(serviceID, endpointID, ogToken, endpoint, auth, res)
}

func prepareHeaders(endpoint *config.Endpoint, token string, claims, params interface{}) map[string]string {
	headers := make(map[string]string, len(endpoint.Headers))
	state := map[string]interface{}{"args": params, "auth": claims, "token": token}
	for _, header := range endpoint.Headers {
		// Load the string if it exists
		value, err := utils.LoadValue(header.Value, state)
		if err == nil {
			if temp, ok := value.(string); ok {
				header.Value = temp
			} else {
				d, _ := json.Marshal(value)
				header.Value = string(d)
			}
		}

		headers[header.Key] = header.Value
	}
	return headers
}

func (m *Module) adjustReqBody(serviceID, endpointID, token string, endpoint *config.Endpoint, auth, params interface{}) (interface{}, error) {
	var req, graph interface{}
	var err error

	switch endpoint.Tmpl {
	case config.EndpointTemplatingEngineGo:
		if tmpl, p := m.templates[getGoTemplateKey("request", serviceID, endpointID)]; p {
			req, err = tmpl2.GoTemplate(module, segmentGoTemplate, tmpl, endpoint.OpFormat, token, auth, params)
			if err != nil {
				return nil, err
			}
		}
		if tmpl, p := m.templates[getGoTemplateKey("graph", serviceID, endpointID)]; p {
			graph, err = tmpl2.GoTemplate(module, segmentGoTemplate, tmpl, "string", token, auth, params)
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

func (m *Module) adjustResBody(serviceID, endpointID, token string, endpoint *config.Endpoint, auth, params interface{}) (interface{}, error) {
	var res interface{}
	var err error

	switch endpoint.Tmpl {
	case config.EndpointTemplatingEngineGo:
		if tmpl, p := m.templates[getGoTemplateKey("response", serviceID, endpointID)]; p {
			res, err = tmpl2.GoTemplate(module, segmentGoTemplate, tmpl, endpoint.OpFormat, token, auth, params)
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

func adjustPath(path string, claims, params interface{}) (string, error) {
	newPath := path
	for {
		pre := strings.IndexRune(newPath, '{')
		if pre < 0 {
			return newPath, nil
		}
		post := strings.IndexRune(newPath, '}')

		key := strings.TrimSuffix(strings.TrimPrefix(newPath[pre:post], "{"), "}")
		value, err := loadParam(key, claims, params)
		if err != nil {
			return "", err
		}

		newPath = newPath[:pre] + value + newPath[post+1:]
	}
}

func loadParam(key string, claims, params interface{}) (string, error) {
	val, err := utils.LoadValue(key, map[string]interface{}{"args": params, "auth": claims})
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
	t = t.Funcs(tmpl2.CreateGoFuncMaps(m.auth))
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
