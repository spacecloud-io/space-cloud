package functions

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"text/template"

	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/modules/global/caching"
	tmpl2 "github.com/spaceuptech/space-cloud/gateway/utils/tmpl"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

func (m *Module) handleCall(ctx context.Context, serviceID, endpointID, token string, auth, params interface{}, cacheInfo *config.ReadCacheOptions) (int, interface{}, error) {
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
	endpointPath, err := adjustPath(ctx, endpoint.Path, auth, params)
	if err != nil {
		return http.StatusBadRequest, nil, err
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
		url = fmt.Sprintf("http://localhost:4122/v1/api/%s/graphql", m.getProject())

	default:
		return http.StatusBadRequest, nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Invalid endpoint kind (%s) provided", endpoint.Kind), nil, nil)
	}

	var redisKey string
	var cacheResponse *caching.CacheResult
	if endpoint.Method == http.MethodGet && cacheInfo != nil {
		// First step is to load all the options
		cacheOptionsArray := make([]interface{}, len(endpoint.CacheOptions))
		for index, key := range endpoint.CacheOptions {
			value, err := utils.LoadValue(key, map[string]interface{}{"args": map[string]interface{}{"auth": auth, "token": ogToken, "url": url}})
			if err != nil {
				return http.StatusBadRequest, nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to extract value for cache option key (%s)", key), err, nil)
			}
			cacheOptionsArray[index] = value
		}

		// Check if response is present in the cache
		cacheResponse, err = m.caching.GetRemoteService(ctx, m.project, serviceID, endpointID, cacheInfo, cacheOptionsArray)
		if err != nil {
			return http.StatusInternalServerError, nil, err
		}

		// Return if cache is present in response
		if cacheResponse.IsCacheHit() {
			return http.StatusOK, cacheResponse.GetResult(), nil
		}

		// Store the redis key for future use
		redisKey = cacheResponse.Key()
	}

	/***************** Set the request method ***************/

	// Set the default method
	if endpoint.Method == "" {
		endpoint.Method = "POST"
	}
	method = endpoint.Method

	/***************** Set the request body *****************/

	newHeaders, newParams, err := m.adjustReqBody(ctx, serviceID, endpointID, ogToken, endpoint, auth, params)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	/***************** Set the request token ****************/

	// Overwrite the token if provided
	if endpoint.Token != "" {
		token = endpoint.Token
	}

	// Create a new token if claims are provided
	if endpoint.Claims != "" {
		newToken, err := m.generateWebhookToken(ctx, serviceID, endpointID, token, endpoint, auth, params)
		if err != nil {
			return http.StatusInternalServerError, nil, err
		}
		token = newToken
	}

	/******** Fire the request and get the response ********/

	scToken, err := m.auth.GetSCAccessToken(ctx)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	// Prepare the state object
	state := map[string]interface{}{"args": params, "auth": auth, "token": ogToken}

	var res interface{}
	req := &utils.HTTPRequest{
		Params: newParams,
		Method: method, URL: url,
		Token: token, SCToken: scToken,
		Headers: prepareHeaders(ctx, newHeaders, state),
	}
	status, err := utils.MakeHTTPRequest(ctx, req, &res)
	if err != nil {
		return status, nil, err
	}

	/**************** Return the response body ****************/

	res, err = m.adjustResBody(ctx, serviceID, endpointID, ogToken, endpoint, auth, res)

	if cacheInfo != nil && redisKey != "" {
		// Store the adjusted body in cache for future use
		if err := m.caching.SetRemoteServiceKey(ctx, redisKey, cacheResponse, cacheInfo, res); err != nil {
			return 0, nil, err
		}
	}
	return status, res, err
}

func prepareHeaders(ctx context.Context, headers config.Headers, state map[string]interface{}) config.Headers {
	out := make([]config.Header, len(headers))
	for i, header := range headers {
		// First create a new header object
		h := config.Header{Key: header.Key, Value: header.Value, Op: header.Op}

		// Load the string if it exists
		value, err := utils.LoadValue(header.Value, state)
		if err == nil {
			if temp, ok := value.(string); ok {
				h.Value = temp
			} else {
				d, _ := json.Marshal(value)
				h.Value = string(d)
			}
		}

		out[i] = h
	}
	return out
}

func (m *Module) adjustReqBody(ctx context.Context, serviceID, endpointID, token string, endpoint *config.Endpoint, auth, params interface{}) (config.Headers, io.Reader, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	var req, graph interface{}
	var err error

	switch endpoint.Tmpl {
	case config.TemplatingEngineGo:
		if tmpl, p := m.templates[getGoTemplateKey("request", serviceID, endpointID)]; p {
			req, err = tmpl2.GoTemplate(ctx, tmpl, endpoint.OpFormat, token, auth, params)
			if err != nil {
				return nil, nil, err
			}
		}
		if tmpl, p := m.templates[getGoTemplateKey("graph", serviceID, endpointID)]; p {
			graph, err = tmpl2.GoTemplate(ctx, tmpl, "string", token, auth, params)
			if err != nil {
				return nil, nil, err
			}
		}
	default:
		helpers.Logger.LogWarn(helpers.GetRequestID(ctx), fmt.Sprintf("Invalid templating engine (%s) provided. Skipping templating step.", endpoint.Tmpl), map[string]interface{}{"serviceId": serviceID, "endpointId": endpointID})
	}

	var body interface{}
	switch endpoint.Kind {
	case config.EndpointKindInternal, config.EndpointKindExternal:
		if req == nil {
			body = params
		} else {
			body = req
		}
	case config.EndpointKindPrepared:
		body = map[string]interface{}{"query": graph, "variables": req}
	default:
		return nil, nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Invalid endpoint kind (%s) provided", endpoint.Kind), nil, nil)
	}

	var requestHeader config.Headers
	requestHeader = append(requestHeader, endpoint.Headers...)

	var requestBody io.Reader
	switch endpoint.ReqPayloadFormat {
	case "", config.EndpointRequestPayloadFormatJSON:
		// Marshal json into byte array
		data, err := json.Marshal(body)
		if err != nil {
			return nil, nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Cannot marshal provided data for graphQL API endpoint (%s)", endpointID), err, map[string]interface{}{"serviceId": serviceID})
		}
		requestBody = bytes.NewReader(data)
		requestHeader = append(requestHeader, config.Header{Key: "Content-Type", Value: "application/json", Op: "set"})
	case config.EndpointRequestPayloadFormatFormData:
		buff := new(bytes.Buffer)
		writer := multipart.NewWriter(buff)

		for key, val := range body.(map[string]interface{}) {
			value, ok := val.(string)
			if !ok {
				return nil, nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Invalid type of value provided for arg (%s) expecting string as endpoint (%s) has request payload of (form-data) type ", endpointID, key), err, map[string]interface{}{"serviceId": serviceID})
			}
			_ = writer.WriteField(key, value)
		}
		err = writer.Close()
		if err != nil {
			return nil, nil, err
		}
		requestBody = bytes.NewReader(buff.Bytes())
		requestHeader = append(requestHeader, config.Header{Key: "Content-Type", Value: writer.FormDataContentType(), Op: "set"})
	}
	return requestHeader, requestBody, err
}

func (m *Module) adjustResBody(ctx context.Context, serviceID, endpointID, token string, endpoint *config.Endpoint, auth, params interface{}) (interface{}, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	var res interface{}
	var err error

	switch endpoint.Tmpl {
	case config.TemplatingEngineGo:
		if tmpl, p := m.templates[getGoTemplateKey("response", serviceID, endpointID)]; p {
			res, err = tmpl2.GoTemplate(ctx, tmpl, endpoint.OpFormat, token, auth, params)
			if err != nil {
				return nil, err
			}
		}
	default:
		helpers.Logger.LogWarn(helpers.GetRequestID(ctx), fmt.Sprintf("Invalid templating engine (%s) provided. Skipping templating step.", endpoint.Tmpl), map[string]interface{}{"serviceId": serviceID, "endpointId": endpointID})
		return params, nil
	}

	if res == nil {
		return params, nil
	}
	return res, nil
}

func (m *Module) generateWebhookToken(ctx context.Context, serviceID, endpointID, token string, endpoint *config.Endpoint, auth, params interface{}) (string, error) {
	var req interface{}
	var err error

	switch endpoint.Tmpl {
	case config.TemplatingEngineGo:
		if tmpl, p := m.templates[getGoTemplateKey("claim", serviceID, endpointID)]; p {
			req, err = tmpl2.GoTemplate(ctx, tmpl, endpoint.OpFormat, token, auth, params)
			if err != nil {
				return "", err
			}
		}
	default:
		helpers.Logger.LogWarn(helpers.GetRequestID(ctx), fmt.Sprintf("Invalid templating engine (%s) provided. Skipping templating step.", endpoint.Tmpl), map[string]interface{}{"serviceId": serviceID, "endpointId": endpointID})
		return token, nil
	}

	if req == nil {
		return token, nil
	}
	return m.auth.CreateToken(ctx, req.(map[string]interface{}))
}

func adjustPath(ctx context.Context, path string, claims, params interface{}) (string, error) {
	newPath := path
	for {
		pre := strings.IndexRune(newPath, '{')
		if pre < 0 {
			return newPath, nil
		}
		post := strings.IndexRune(newPath, '}')

		key := strings.TrimSuffix(strings.TrimPrefix(newPath[pre:post], "{"), "}")
		value, err := loadParam(ctx, key, claims, params)
		if err != nil {
			return "", err
		}

		newPath = newPath[:pre] + value + newPath[post+1:]
	}
}

func loadParam(ctx context.Context, key string, claims, params interface{}) (string, error) {
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
	m.lock.RLock()
	defer m.lock.RUnlock()

	for _, s := range m.config {
		if s.ID == service {
			return s
		}
	}

	return nil
}

func (m *Module) getProject() string {
	m.lock.RLock()
	defer m.lock.RUnlock()

	return m.project
}

func (m *Module) createGoTemplate(kind, serviceID, endpointID, tmpl string) error {
	key := getGoTemplateKey(kind, serviceID, endpointID)

	// Create a new template object
	t := template.New(key)
	t = t.Funcs(tmpl2.CreateGoFuncMaps(m.auth))
	val, err := t.Parse(tmpl)
	if err != nil {
		return helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Invalid golang template provided", err, nil)
	}

	m.templates[key] = val
	return nil
}

func getGoTemplateKey(kind, serviceID, endpointID string) string {
	return fmt.Sprintf("%s---%s---%s", kind, serviceID, endpointID)
}
