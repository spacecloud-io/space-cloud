package remoteservice

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"strconv"
	"strings"
	"text/template"

	"github.com/spacecloud-io/space-cloud/config"
	"github.com/spacecloud-io/space-cloud/utils"
	tmpl2 "github.com/spacecloud-io/space-cloud/utils/tmpl"
	"github.com/spaceuptech/helpers"
)

func (a *App) getEndpoint(projectID, serviceName, endpoint string) (*config.Service, *config.Endpoint, error) {
	alias := fmt.Sprintf("%s---%s", projectID, serviceName)
	serviceObj, ok := a.Services[alias]
	if !ok {
		return nil, nil, fmt.Errorf("could not find endpoint (%s) for service (%s) in project (%s)", endpoint, serviceName, projectID)
	}

	endpointObj, ok := serviceObj.Endpoints[endpoint]
	if !ok {
		return nil, nil, fmt.Errorf("could not find endpoint (%s) for service (%s) in project (%s)", endpoint, serviceName, projectID)
	}
	return serviceObj, endpointObj, nil
}

func (a *App) adjustReqBody(ctx context.Context, projectID, serviceID, endpointID, token string, endpoint *config.Endpoint, auth, params interface{}) (config.Headers, io.Reader, error) {
	var req, graph interface{}
	var err error

	tmpl, err := a.getStringOutputFromPlugins(endpoint, config.PluginTmpl)
	if err != nil {
		return nil, nil, err
	}

	opFormat, err := a.getStringOutputFromPlugins(endpoint, config.PluginOutputFormat)
	if err != nil {
		return nil, nil, err
	}

	headers, err := a.getHeadersFromPlugins(endpoint)
	if err != nil {
		return nil, nil, err
	}

	reqPayloadFormat, err := a.getStringOutputFromPlugins(endpoint, config.PluginReqPayloadFormat)
	if err != nil {
		return nil, nil, err
	}

	switch tmpl {
	case string(config.TemplatingEngineGo):
		if tmpl, p := a.templates[getGoTemplateKey("request", projectID, serviceID, endpointID)]; p {
			req, err = tmpl2.GoTemplate(ctx, tmpl, opFormat, token, auth, params)
			if err != nil {
				return nil, nil, err
			}
		}
		if tmpl, p := a.templates[getGoTemplateKey("graph", projectID, serviceID, endpointID)]; p {
			graph, err = tmpl2.GoTemplate(ctx, tmpl, "string", token, auth, params)
			if err != nil {
				return nil, nil, err
			}
		}
	default:
		helpers.Logger.LogWarn(helpers.GetRequestID(ctx), fmt.Sprintf("Invalid templating engine (%s) provided. Skipping templating step.", tmpl), map[string]interface{}{"serviceId": serviceID, "endpointId": endpointID})
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
	requestHeader = append(requestHeader, *headers...)

	var requestBody io.Reader
	switch reqPayloadFormat {
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

func (a *App) adjustResBody(ctx context.Context, projectID, serviceID, endpointID, token string, endpoint *config.Endpoint, auth, params interface{}) (interface{}, error) {
	var res interface{}
	var err error

	tmpl, err := a.getStringOutputFromPlugins(endpoint, config.PluginTmpl)
	if err != nil {
		return nil, err
	}

	opFormat, err := a.getStringOutputFromPlugins(endpoint, config.PluginOutputFormat)
	if err != nil {
		return nil, err
	}

	switch tmpl {
	case string(config.TemplatingEngineGo):
		if tmpl, p := a.templates[getGoTemplateKey("response", projectID, serviceID, endpointID)]; p {
			res, err = tmpl2.GoTemplate(ctx, tmpl, opFormat, token, auth, params)
			if err != nil {
				return nil, err
			}
		}
	default:
		helpers.Logger.LogWarn(helpers.GetRequestID(ctx), fmt.Sprintf("Invalid templating engine (%s) provided. Skipping templating step.", tmpl), map[string]interface{}{"serviceId": serviceID, "endpointId": endpointID})
		return params, nil
	}

	if res == nil {
		return params, nil
	}
	return res, nil
}

func (a *App) createGoTemplate(kind, projectID, serviceID, endpointID, tmpl string) error {
	key := getGoTemplateKey(kind, projectID, serviceID, endpointID)

	// Create a new template object
	t := template.New(key)
	t = t.Funcs(tmpl2.CreateGoFuncMaps(nil))
	val, err := t.Parse(tmpl)
	if err != nil {
		return helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Invalid golang template provided", err, nil)
	}

	a.templates[key] = val
	return nil
}

// CombineProjectRemoteServiceName joins project id and service name
func CombineProjectRemoteServiceName(projectID, serviceName string) string {
	return fmt.Sprintf("%s---%s", projectID, serviceName)
}

// SplitPojectRemoteServiceName return project id and service name
func SplitPojectRemoteServiceName(name string) (projectID, serviceName string) {
	projectID = strings.Split(name, "---")[0]
	serviceName = strings.Split(name, "---")[1]
	return
}

func combineProjectServiceEndpointName(projectID, serviceName, endpointName string) string {
	return fmt.Sprintf("%s-%s-%s", projectID, serviceName, endpointName)
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

func getGoTemplateKey(kind, projectID, serviceID, endpointID string) string {
	return fmt.Sprintf("%s---%s---%s---%s", kind, projectID, serviceID, endpointID)
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

func getRequestParamsFromURL(url string) (project, service, endpoint string) {
	url = strings.TrimPrefix(url, "/v1/api/")
	project = strings.Split(url, "/")[0]
	service = strings.Split(url, "/")[2]
	endpoint = strings.Split(url, "/")[3]
	return
}
