package functions

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/spacecloud-io/space-cloud/config"
	"github.com/spacecloud-io/space-cloud/managers/apis"
	"github.com/spacecloud-io/space-cloud/model"
	"github.com/spacecloud-io/space-cloud/utils"
	"github.com/spaceuptech/helpers"
	"go.uber.org/zap"
)

// GetRoutes returns all the apis that are exposed by this app
func (a *App) GetRoutes() []*apis.API {
	resAPIs := make([]*apis.API, 0)
	for name, service := range a.Services {
		projectID := strings.Split(name, "---")[0]
		serviceName := strings.Split(name, "---")[1]
		for endpointName := range service.Endpoints {
			// Prepare schema for request body
			requestSchemaJSON, _ := json.Marshal(m{
				"type": openapi3.TypeObject,
				"properties": m{
					"operationName": m{"type": openapi3.TypeString},
					"query":         m{"type": openapi3.TypeString},
					"variables":     m{"type": openapi3.TypeObject, "additionalProperties": true},
				},
				"required": t{"query"},
			})
			requestSchema := new(openapi3.SchemaRef)
			if err := requestSchema.UnmarshalJSON(requestSchemaJSON); err != nil {
				panic(err)
			}

			// Prepare schema for graphql response
			responseSchemaJSON, _ := json.Marshal(m{
				"type": openapi3.TypeObject,
				"properties": m{
					"data": m{"type": openapi3.TypeObject, "additionalProperites": true},
					"errors": m{
						"type":  openapi3.TypeArray,
						"items": m{"type": openapi3.TypeObject, "additionalProperites": true},
					},
				},
			})
			responseSchema := new(openapi3.SchemaRef)
			if err := responseSchema.UnmarshalJSON(responseSchemaJSON); err != nil {
				panic(err)
			}
			responseSchemaDescription := "Successful response"

			pathDef := &openapi3.PathItem{
				Post: &openapi3.Operation{
					Tags:        []string{serviceName},
					Description: fmt.Sprintf("%s-%s", serviceName, endpointName),
					RequestBody: &openapi3.RequestBodyRef{
						Value: &openapi3.RequestBody{
							Required: true,
							Content: openapi3.Content{
								"application/json": &openapi3.MediaType{
									Schema: requestSchema,
								},
							},
						},
					},
					Responses: openapi3.Responses{
						"200": &openapi3.ResponseRef{
							Value: &openapi3.Response{
								Description: &responseSchemaDescription,
								Content: openapi3.Content{
									"application/json": &openapi3.MediaType{
										Schema: responseSchema,
									},
								},
							},
						},
						"401": apis.SCErrorResponseSchema,
						"403": apis.SCErrorResponseSchema,
						"500": apis.SCErrorResponseSchema,
					},
				},
			}
			url := fmt.Sprintf("/v1/api/%s/services/%s/%s", projectID, serviceName, endpointName)
			op := fmt.Sprintf("%s-%s-%s", projectID, serviceName, endpointName)
			api := apis.API{Path: url, PathDef: pathDef, Op: op}
			resAPIs = append(resAPIs, &api)
		}
	}
	return resAPIs
}

// GetHandler returns the graphql api handler
func (a *App) GetHandler(op string) (apis.HandlerFunc, error) {
	h := func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		// Prepare context object
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Minute)
		defer cancel()

		// Parse request body
		req := new(model.FunctionsRequest)
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			a.logger.Error("Invalid function request payload provided", zap.Error(err))
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusBadRequest, errors.New("invalid request payload received"))
			return
		}

		projectID, serviceName, endpoint := getRequestParamsFromURL(r.URL.Path)
		_, endpointObj, err := a.getEndpoint(projectID, serviceName, endpoint)
		if err != nil {
			a.logger.Error("Invalid function request payload provided", zap.Error(err))
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusBadRequest, errors.New("invalid request payload received"))
			return
		}

		timeout, err := a.getIntOutputFromPlugins(endpointObj, config.PluginTimeout)
		if err != nil {
			a.logger.Error("Invalid function request payload provided", zap.Error(err))
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusBadRequest, errors.New("invalid request payload received"))
			return
		}

		// Create a new context
		ctx, cancel = context.WithTimeout(r.Context(), time.Duration(timeout)*time.Second)
		defer cancel()

		// TODO: Check if function call is authorised
		// TODO: Implement caching

		reqParams := model.RequestParams{}
		status, result, err := a.handleCall(ctx, projectID, serviceName, endpoint, "", nil, reqParams, nil)
		if err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Receieved error from service call (%s:%s)", serviceName, endpoint), err, nil)
			_ = helpers.Response.SendErrorResponse(ctx, w, status, err)
			return
		}

		_ = helpers.Response.SendResponse(ctx, w, http.StatusOK, result)
	}

	return h, nil
}

func (a *App) handleCall(ctx context.Context, projectID, serviceID, endpointID, token string, auth, params interface{}, cacheInfo *config.ReadCacheOptions) (int, interface{}, error) {
	var url string
	var method string
	var ogToken string

	// Load the service rule
	service, endpoint, err := a.getEndpoint(projectID, serviceID, endpointID)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}
	serviceURL := strings.TrimSuffix(service.URL, "/")
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
		url = fmt.Sprintf("http://localhost:4122/v1/api/%s/graphql", projectID)
	default:
		return http.StatusBadRequest, nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Invalid endpoint kind (%s) provided", endpoint.Kind), nil, nil)
	}

	// TODO: Implement caching logic

	/***************** Set the request method ***************/
	// Set the default method
	if endpoint.Method == "" {
		endpoint.Method = "POST"
	}
	method = endpoint.Method

	/***************** Set the request body *****************/
	newHeaders, newParams, err := a.adjustReqBody(ctx, projectID, serviceID, endpointID, ogToken, endpoint, auth, params)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	/***************** Set the request token ****************/
	endPointToken, err := a.getStringOutputFromPlugins(endpoint, config.PluginToken)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	// Overwrite the token if provided
	if endPointToken != "" {
		token = endPointToken
	}

	/******** Fire the request and get the response ********/

	scToken := ""
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
	res, err = a.adjustResBody(ctx, projectID, serviceID, endpointID, ogToken, endpoint, auth, res)
	return status, res, err
}

func getRequestParamsFromURL(url string) (project, service, endpoint string) {
	url = strings.TrimPrefix(url, "/v1/api/")
	project = strings.Split(url, "/")[0]
	service = strings.Split(url, "/")[2]
	endpoint = strings.Split(url, "/")[3]
	return
}
