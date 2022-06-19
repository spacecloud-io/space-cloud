package remoteservice

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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
		projectID, serviceName := SplitPojectRemoteServiceName(name)
		for endpointName := range service.Endpoints {
			// Prepare schema for request body
			requestSchemaJSON, _ := json.Marshal(m{
				"type": openapi3.TypeObject,
				"properties": m{
					"params":  m{"type": openapi3.TypeObject, "additionalProperties": true},
					"timeout": m{"type": openapi3.TypeInteger},
					"cache": m{
						"type": openapi3.TypeObject,
						"properties": m{
							"ttl":               m{"type": openapi3.TypeInteger},
							"instantInvalidate": m{"type": openapi3.TypeBoolean},
						},
					},
				},
				"required": t{},
			})
			requestSchema := new(openapi3.SchemaRef)
			if err := requestSchema.UnmarshalJSON(requestSchemaJSON); err != nil {
				panic(err)
			}

			// Prepare schema for graphql response
			responseSchemaJSON, _ := json.Marshal(m{
				"type":                 openapi3.TypeObject,
				"additionalProperties": true,
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
			op := combineProjectServiceEndpointName(projectID, serviceName, endpointName)
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

		// Get the JWT token from header
		token := utils.GetTokenFromHeader(r)

		// Get path params
		projectID, serviceName, endpoint := getRequestParamsFromURL(r.URL.Path)

		// Get request timeout
		_, endpointObj, err := a.getEndpoint(projectID, serviceName, endpoint)
		if err != nil {
			a.logger.Error("Unable to get endpoint", zap.Error(err))
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusBadRequest, errors.New("unable to get endpoint"))
			return
		}
		timeout, err := a.getTimeoutPlugins(endpointObj)
		if err != nil {
			a.logger.Error("Unable to get endpoint timeout", zap.Error(err))
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusBadRequest, errors.New("unable to get endpoint timeout"))
			return
		}

		// Create a new context
		ctx, cancel = context.WithTimeout(r.Context(), time.Duration(timeout)*time.Second)
		defer cancel()

		// TODO: Check if function call is authorised
		// actions, reqParams, err := auth.IsFuncCallAuthorised(ctx, projectID, serviceName, endpoint, token, req.Params)
		// if err != nil {
		// 	_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusForbidden, err)
		// 	return
		// }

		reqParams := utils.ExtractRequestParams(r, model.RequestParams{}, req)
		status, result, err := a.CallWithContext(ctx, projectID, serviceName, endpoint, token, reqParams, req)
		if err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Receieved error from service call (%s:%s)", serviceName, endpoint), err, nil)
			_ = helpers.Response.SendErrorResponse(ctx, w, status, err)
			return
		}

		// TODO: Implement post process method
		// _ = authHelpers.PostProcessMethod(ctx, auth.GetAESKey(), actions, result)

		_ = helpers.Response.SendResponse(ctx, w, http.StatusOK, result)
	}

	return h, nil
}

// CallWithContext handles remote service request
func (a *App) CallWithContext(ctx context.Context, projectID, serviceID, endpointID, token string, reqParams model.RequestParams, req *model.FunctionsRequest) (int, interface{}, error) {
	reqParams.HTTPParams.Payload = map[string]interface{}{
		"service":  serviceID,
		"endpoint": endpointID,
		"params":   req.Params,
	}

	// TODO: Implement integration hooks
	// hookResponse := m.integrationMan.InvokeHook(ctx, reqParams)
	// if hookResponse.CheckResponse() {
	// 	// Check if an error occurred
	// 	if err := hookResponse.Error(); err != nil {
	// 		return hookResponse.Status(), nil, err
	// 	}

	// 	// Gracefully return
	// 	return hookResponse.Status(), hookResponse.Result(), nil
	// }

	// TODO: Add metric hook for cache
	status, result, err := a.handleCall(ctx, projectID, serviceID, endpointID, token, reqParams.Claims, req.Params, req.Cache)
	if err != nil {
		return status, result, err
	}

	// TODO: Implement metric hook
	// m.metricHook(m.project, service, function)

	return status, result, nil
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
	// var redisKey string
	// var cacheResponse *caching.CacheResult
	// if endpoint.Method == http.MethodGet && cacheInfo != nil {
	// 	// First step is to load all the options
	// 	cacheOptionsArray := make([]interface{}, len(endpoint.CacheOptions))
	// 	for index, key := range endpoint.CacheOptions {
	// 		value, err := utils.LoadValue(key, map[string]interface{}{"args": map[string]interface{}{"auth": auth, "token": ogToken, "url": url}})
	// 		if err != nil {
	// 			return http.StatusBadRequest, nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to extract value for cache option key (%s)", key), err, nil)
	// 		}
	// 		cacheOptionsArray[index] = value
	// 	}

	// 	// Check if response is present in the cache
	// 	cacheResponse, err = m.caching.GetRemoteService(ctx, m.project, serviceID, endpointID, cacheInfo, cacheOptionsArray)
	// 	if err != nil {
	// 		return http.StatusInternalServerError, nil, err
	// 	}

	// 	// Return if cache is present in response
	// 	if cacheResponse.IsCacheHit() {
	// 		return http.StatusOK, cacheResponse.GetResult(), nil
	// 	}

	// 	// Store the redis key for future use
	// 	redisKey = cacheResponse.Key()
	// }

	/***************** Set the request method ***************/
	// Set the default method
	if endpoint.Method == "" {
		endpoint.Method = "POST"
	}
	method = endpoint.Method

	/***************** Set the request body *****************/
	var newParams io.Reader
	var newHeaders config.Headers
	// Make a request object
	if !(method == http.MethodGet || method == http.MethodDelete) {
		newHeaders, newParams, err = a.applyRequestBodyPlugin(ctx, projectID, serviceID, endpointID, ogToken, endpoint, auth, params)
		if err != nil {
			return http.StatusBadRequest, nil, err
		}
	}

	/***************** Set the request token ****************/
	endPointToken, err := a.getTokenPlugin(endpoint)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	// Overwrite the token if provided
	if endPointToken != "" {
		token = endPointToken
	}

	// Create a new token if claims are provided
	// claims, _, err := a.getTemplateOutputFromPlugins(endpoint, config.PluginClaims)
	// if err != nil {
	// 	return http.StatusInternalServerError, nil, err
	// }
	// if claims != "" {
	// TODO: Implement generate webhook token
	// newToken, err := m.generateWebhookToken(ctx, serviceID, endpointID, token, endpoint, auth, params)
	// if err != nil {
	// 	return http.StatusInternalServerError, nil, err
	// }
	// token = newToken
	// }

	/******** Fire the request and get the response ********/

	scToken := ""
	// TODO: Implement get sc access token
	// scToken, err := m.auth.GetSCAccessToken(ctx)
	// if err != nil {
	// 	return http.StatusInternalServerError, nil, err
	// }

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
	res, err = a.applyResponseTemplatePlugin(ctx, projectID, serviceID, endpointID, ogToken, endpoint, auth, res)
	return status, res, err
}
