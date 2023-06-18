package openapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/spacecloud-io/space-cloud/modules/rpc"
	"github.com/spacecloud-io/space-cloud/pkg/apis/core/v1alpha1"
)

// createRPC creates a RPC from an OpenAPI operation
func (s *Source) createRPC(url, operationID, operationType, method string, respBody, reqBody *openapi3.Schema, parameters openapi3.Parameters, plugins []v1alpha1.HTTPPlugin) *rpc.RPC {
	if operationType == "" {
		operationType = "mutation"
		if method == http.MethodGet {
			operationType = "query"
		}
	}
	for _, param := range parameters {
		if reqBody == nil {
			reqBody = &openapi3.Schema{
				Type:       "object",
				Properties: make(openapi3.Schemas),
				Required:   make([]string, 0),
			}
		}
		reqBody.Properties[param.Value.Name] = param.Value.Schema
		if param.Value.Required {
			reqBody.Required = append(reqBody.Required, param.Value.Name)
		}

		contentSchema, ok := param.Value.Content["application/json"]
		if !ok {
			continue
		}
		if contentSchema.Schema != nil {
			reqBody.Properties[param.Value.Name] = contentSchema.Schema
		}

	}
	var reqSchema *openapi3.SchemaRef
	var respSchema *openapi3.SchemaRef
	if reqBody != nil {
		reqSchema = &openapi3.SchemaRef{Value: reqBody}
	}
	if respBody != nil {
		respSchema = &openapi3.SchemaRef{Value: respBody}
	}

	cleanseSchemaRef(reqSchema)
	cleanseSchemaRef(respSchema)
	return &rpc.RPC{
		Name:           operationID,
		OperationType:  operationType,
		RequestSchema:  reqSchema,
		ResponseSchema: respSchema,
		HTTPOptions: &v1alpha1.HTTPOptions{
			Method: method,
		},
		Plugins: plugins,
		Call:    s.createCall(url, operationID, method, reqBody, parameters),
	}
}

// createCall creates a callback function for an openapi operation
func (s *Source) createCall(url, operationID, method string, reqBody *openapi3.Schema, parameters openapi3.Parameters) func(ctx context.Context, vars map[string]any) (any, error) {
	return func(ctx context.Context, vars map[string]any) (any, error) {
		// Create request parameters
		urlPath, vars, headers, err := createRequestParams(url, reqBody, parameters, vars)
		if err != nil {
			return nil, err
		}
		url := fmt.Sprintf("%s%s", s.Spec.Source.URL, urlPath)

		// Create the request
		req, err := http.NewRequestWithContext(ctx, method, url, nil)
		if err != nil {
			return nil, err
		}

		// Set request body for POST method
		if method != http.MethodGet {
			data, err := json.Marshal(vars)
			if err != nil {
				return nil, err
			}
			req, err = http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(data))
			if err != nil {
				return nil, err
			}
		}

		req.Header.Set("Content-Type", "application/json")
		// Set additional headers
		for header, val := range headers {
			req.Header.Set(header, val)
		}

		// Make the request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		// Read and decode the response body
		var response interface{}
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			return nil, err
		}
		return response, nil
	}
}

// createRequestParams takes a url path, a request body schema, a list of parameters and a map of variables
// and returns a modified url path with query parameters and a modified map of variables
func createRequestParams(urlPath string, reqBodySchema *openapi3.Schema, params openapi3.Parameters, vars map[string]any) (string, map[string]any, map[string]string, error) {
	// queryString is a slice of strings to store the query parameters
	queryString := url.Values{}
	// reqBodyKeys is a slice of strings to store the keys of the request body properties
	reqBodyKeys := make([]string, 0)
	// loop over the request body properties and append their keys to reqBodyKeys
	for reqBodyKey := range reqBodySchema.Properties {
		reqBodyKeys = append(reqBodyKeys, reqBodyKey)
	}

	headers := make(map[string]string)

	// loop over the parameters
	for _, param := range params {
		// if the parameter is in the query location
		if param.Value.In == openapi3.ParameterInQuery {
			// get the value of the parameter from the vars map
			queryVal := vars[param.Value.Name]
			queryVal, ok := vars[param.Value.Name]
			if !ok {
				continue
			}
			// if the parameter schema type is object or array
			if param.Value.Schema.Value.Type == openapi3.TypeObject || param.Value.Schema.Value.Type == openapi3.TypeArray {
				// marshal the value to JSON
				queryVal, err := json.Marshal(queryVal)
				if err != nil {
					return "", nil, nil, err
				}
				// append the parameter name and value to queryString with an equal sign
				queryString.Set(param.Value.Name, fmt.Sprintf("%s", string(queryVal)))
			} else {
				// otherwise, append the parameter name and value to queryString with an equal sign and replace any spaces with plus signs
				queryString.Set(param.Value.Name, fmt.Sprintf("%v", queryVal))
			}
			// if the parameter name is not in the reqBodyKeys slice
			if !containsKey(reqBodyKeys, param.Value.Name) {
				// delete the parameter name from the vars map
				delete(vars, param.Value.Name)
			}
		} else if param.Value.In == openapi3.ParameterInPath {
			// if the parameter is in the path location
			// get the value of the parameter from the vars map
			queryVal, ok := vars[param.Value.Name]
			if !ok {
				continue
			}
			// replace the parameter placeholder in the url path with the value
			urlPath = strings.Replace(urlPath, fmt.Sprintf("{%s}", param.Value.Name), fmt.Sprintf("%v", queryVal), 1)
			// if the parameter name is not in the reqBodyKeys slice
			if !containsKey(reqBodyKeys, param.Value.Name) {
				// delete the parameter name from the vars map
				delete(vars, param.Value.Name)
			}
		} else if param.Value.In == openapi3.ParameterInHeader {
			// if the parameter is in the path location
			// get the value of the parameter from the vars map
			queryVal, ok := vars[param.Value.Name]
			if !ok {
				continue
			}
			headers[param.Value.Name] = fmt.Sprintf("%v", queryVal)
			// if the parameter name is not in the reqBodyKeys slice
			if !containsKey(reqBodyKeys, param.Value.Name) {
				// delete the parameter name from the vars map
				delete(vars, param.Value.Name)
			}
		}

	}
	// if queryString is not empty
	if len(queryString) > 0 {
		// return the url path with a question mark and join queryString with ampersands
		return fmt.Sprintf("%s?%s", urlPath, queryString.Encode()), vars, headers, nil
	}
	// otherwise, return the url path and vars as they are
	return urlPath, vars, headers, nil
}

// containsKey checks if a string is present in an array
func containsKey(array []string, str string) bool {
	// Loop through the array
	for _, element := range array {
		// Compare the element with the string
		if element == str {
			// Return true if they are equal
			return true
		}
	}
	// Return false if the string is not found
	return false
}

// cleanseSchemaRef removes ref from openapi spec recursively
func cleanseSchemaRef(schema *openapi3.SchemaRef) {
	if schema == nil {
		return
	}

	for _, p := range schema.Value.Properties {
		cleanseSchemaRef(p)
	}

	cleanseSchemaRef(schema.Value.Items)
	cleanseSchemaRef(schema.Value.AdditionalProperties.Schema)
	cleanseSchemaRef(schema.Value.Not)
	for _, s := range schema.Value.AllOf {
		cleanseSchemaRef(s)
	}
	for _, s := range schema.Value.AnyOf {
		cleanseSchemaRef(s)
	}
	for _, s := range schema.Value.OneOf {
		cleanseSchemaRef(s)
	}

	schema.Ref = ""
}
