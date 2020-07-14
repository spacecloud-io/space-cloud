package routing

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

func (r *Routing) modifyRequest(ctx context.Context, modules modulesInterface, route *config.Route, req *http.Request) (string, interface{}, int, error) {
	// Extract the token
	token := utils.GetTokenFromHeader(req)

	// Extract the params only if content-type is `application/json`
	var params interface{}
	var data []byte
	var err error
	if req.Header.Get("Content-Type") == "application/json" {
		data, err = ioutil.ReadAll(req.Body)
		if err != nil {
			return "", nil, http.StatusBadRequest, err
		}

		if err := json.Unmarshal(data, &params); err != nil {
			utils.LogWarn("Unable to unmarshal body to JSON", module, handleRequest)
			req.Body = ioutil.NopCloser(bytes.NewBuffer(data))
		}
	}

	// Finally we authorize the request
	a := modules.Auth()
	args := map[string]interface{}{"params": params, "query": makeQueryArguments(req)}
	auth, err := a.AuthorizeRequest(ctx, route.Rule, route.Project, token, args)
	if err != nil {
		return "", nil, http.StatusForbidden, err
	}

	// Set the headers
	state := map[string]interface{}{"args": params, "auth": auth}
	headers := append(r.globalConfig.RequestHeaders, route.Modify.RequestHeaders...)
	prepareHeaders(headers, state).UpdateHeader(req.Header)

	// Don't forget to reset the body
	if params != nil {
		// Generate new request body if template was provided
		newParams, err := r.adjustBody("request", route.Project, token, route, auth, params)
		if err != nil {
			return "", nil, http.StatusBadRequest, err
		}

		// Marshal it then set it
		data, _ = json.Marshal(newParams)
		req.Body = ioutil.NopCloser(bytes.NewBuffer(data))
		req.Header.Set("Content-Length", strconv.Itoa(len(data)))
		req.ContentLength = int64(len(data))
	}

	return token, auth, http.StatusOK, err
}

func (r *Routing) modifyResponse(res *http.Response, route *config.Route, token string, auth interface{}) error {
	// Extract the params only if content-type is `application/json` and a response template is provided
	var params interface{}
	var data []byte
	var err error

	if res.Header.Get("Content-Type") == "application/json" && route.Modify.ResTmpl != "" {
		data, err = ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(data, &params); err != nil {
			utils.LogWarn("Unable to unmarshal response body to JSON", module, handleResponse)
			res.Body = ioutil.NopCloser(bytes.NewBuffer(data))
		}
	}

	// Set the headers
	state := map[string]interface{}{"args": params, "auth": auth}
	headers := append(r.globalConfig.ResponseHeaders, route.Modify.ResponseHeaders...)
	prepareHeaders(headers, state).UpdateHeader(res.Header)

	// If params is not nil we need to template the response
	if params != nil {
		newParams, err := r.adjustBody("response", route.Project, token, route, auth, params)
		if err != nil {
			return err
		}

		// Marshal it then set it
		data, _ = json.Marshal(newParams)
		res.Body = ioutil.NopCloser(bytes.NewBuffer(data))
		res.Header.Set("Content-Length", strconv.Itoa(len(data)))
		res.ContentLength = int64(len(data))
	}

	return nil
}

func makeQueryArguments(r *http.Request) map[string]interface{} {
	// Prepare the query parameters
	queryParams := r.URL.Query()
	params := make(map[string]interface{}, len(queryParams))
	for k, v := range queryParams {
		params[k] = v[0]
	}

	// Prepare the headers
	headers := make(map[string]interface{}, len(r.Header))
	for k, v := range r.Header {
		headers[k] = v[0]
	}

	// Prepare path array
	pathArrayTemp := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
	pathArray := make([]interface{}, len(pathArrayTemp))
	for i, v := range pathArrayTemp {
		pathArray[i] = v
	}

	// Finally we return the entire query object
	return map[string]interface{}{
		"path":      r.URL.Path,
		"pathArray": pathArray,
		"params":    params,
		"headers":   headers,
	}
}
