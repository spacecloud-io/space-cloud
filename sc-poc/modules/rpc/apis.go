package rpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/spacecloud-io/space-cloud/managers/apis"
	"github.com/spacecloud-io/space-cloud/utils"
	"go.uber.org/zap"
)

// GetRoutes returns all the apis that are exposed by this app
func (m *Module) GetAPIRoutes() apis.APIs {
	return m.apis
}

func (m *Module) prepareAPIs(rpcSource Source) {
	for _, rpc := range rpcSource.GetRPCs() {
		httpOpts := processHTTPOptions(rpc.Name, rpc.HTTPOptions)

		// Create an open api path definition
		operation := &openapi3.Operation{
			OperationID: rpc.Name,
			// TODO: Add tags based on the `space-cloud.io/tags` annotations
			Tags: []string{},
		}

		// Add security header is authentication was required
		if rpc.Authenticate == nil {
			operation.Security = &openapi3.SecurityRequirements{{"bearerAuth": []string{}}}
		}

		// Get request and response schemas
		requestSchema, responseSchema := rpc.RequestSchema, rpc.ResponseSchema

		if !(httpOpts.Method == http.MethodPost || httpOpts.Method == http.MethodPut) {
			// Need to add the variables as query parameters
			if requestSchema == nil {
				requestSchema = &openapi3.SchemaRef{
					Value: &openapi3.Schema{},
				}
			}
			parameters := make(openapi3.Parameters, 0, len(requestSchema.Value.Properties))
			for propertyKey, schemaRef := range requestSchema.Value.Properties {

				parameter := &openapi3.Parameter{
					In:       "query",
					Name:     propertyKey,
					Required: utils.StringExists(propertyKey, requestSchema.Value.Required...),
				}

				if utils.StringExists(schemaRef.Value.Type, "object", "array") {
					parameter.Content = openapi3.NewContentWithJSONSchemaRef(schemaRef)
				} else {
					parameter.Schema = schemaRef
				}

				parameters = append(parameters, &openapi3.ParameterRef{Value: parameter})
			}

			operation.Parameters = parameters
		} else {
			operation.RequestBody = &openapi3.RequestBodyRef{
				Value: &openapi3.RequestBody{
					Required: true,
					Content:  openapi3.NewContentWithJSONSchemaRef(requestSchema),
				},
			}
		}

		{
			// Prepare response schemas
			operation.Responses = openapi3.Responses{
				"200": &openapi3.ResponseRef{
					Value: openapi3.NewResponse().WithDescription(fmt.Sprintf("Response for %s", rpc.Name)).WithJSONSchemaRef(responseSchema),
				},
				"400": apis.SCErrorResponseSchema,
				"401": apis.SCErrorResponseSchema,
				"403": apis.SCErrorResponseSchema,
				"500": apis.SCErrorResponseSchema,
			}
		}

		// Prepare the path item
		pathItem := &openapi3.PathItem{}
		switch httpOpts.Method {
		case http.MethodPost:
			pathItem.Post = operation
		case http.MethodPut:
			pathItem.Put = operation
		case http.MethodDelete:
			pathItem.Delete = operation
		case http.MethodGet:
			pathItem.Get = operation
		}

		// Add the extensions
		operation.Extensions = map[string]interface{}{
			"x-client-gen":      true,
			"x-request-op-type": rpc.OperationType,
		}
		for k, v := range rpc.Extensions {
			operation.Extensions[k] = v
		}

		// TODO: See if we need support for schema definitions
		// Prepare the schemas definitions
		// schemas := make(openapi3.Schemas, len(compiledQuery.VariableSchema.Definitions))
		// for propertyName, schema := range compiledQuery.VariableSchema.Definitions {
		// 	schemaRef := new(openapi3.SchemaRef)
		// 	schemajSON, _ := schema.MarshalJSON()
		// 	_ = schemaRef.UnmarshalJSON([]byte(strings.ReplaceAll(string(schemajSON), "#/$defs/", "#/components/schemas/")))
		// 	schemas[propertyName] = schemaRef
		// }

		m.apis = append(m.apis, &apis.API{
			Name: fmt.Sprintf("complied_query_%s", rpc.Name),
			Path: httpOpts.URL,
			OpenAPI: &apis.OpenAPI{
				PathDef: pathItem,
				// Schemas: schemas,
			},
			Plugins: rpc.Plugins,
			Handler: m.handleRPCRequest(rpc, operation),
		})
	}
}

func (m *Module) handleRPCRequest(rpc *RPC, operation *openapi3.Operation) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Prepare context object
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Minute)
		defer cancel()

		// Extract variables
		variables := make(map[string]interface{})

		if utils.StringExists(rpc.HTTPOptions.Method, http.MethodGet, http.MethodDelete) {
			// For get and delete methods
			for _, param := range operation.Parameters {
				stringData := r.URL.Query()[param.Value.Name]

				if len(stringData) == 0 {
					continue
					// TODO: Throw an error if the variable was neither exported or injected
					// TODO: Use the default value provided in the graphql query if parameter is not provided

					// a.logger.Error("Variable parameter not provided", zap.String("param", param.Value.Name))
					// _ = utils.SendErrorResponse(w, http.StatusBadRequest, fmt.Errorf("invalid request payload received for param '%s'", param.Value.Name))
					// return
				}

				if len(param.Value.Content) > 0 {
					var data interface{}
					if err := json.Unmarshal([]byte(stringData[0]), &data); err != nil {
						m.logger.Error("Invalid graphql variable parameter provided", zap.String("param", param.Value.Name), zap.Error(err))
						_ = utils.SendErrorResponse(w, http.StatusBadRequest, fmt.Errorf("invalid request payload received for param '%s'", param.Value.Name))
						return
					}

					variables[param.Value.Name] = data
				} else {
					switch param.Value.Schema.Value.Type {
					case "string":
						variables[param.Value.Name] = stringData[0]

					case "integer":
						num, err := strconv.Atoi(stringData[0])
						if err != nil {
							m.logger.Error("Invalid graphql variable parameter provided", zap.String("param", param.Value.Name), zap.Error(err))
							_ = utils.SendErrorResponse(w, http.StatusBadRequest, fmt.Errorf("invalid request payload received for param '%s'", param.Value.Name))
							return
						}
						variables[param.Value.Name] = num

					case "number":
						num, err := strconv.ParseFloat(stringData[0], 64)
						if err != nil {
							m.logger.Error("Invalid graphql variable parameter provided", zap.String("param", param.Value.Name), zap.Error(err))
							_ = utils.SendErrorResponse(w, http.StatusBadRequest, fmt.Errorf("invalid request payload received for param '%s'", param.Value.Name))
							return
						}
						variables[param.Value.Name] = num

					case "boolean":
						b, err := strconv.ParseBool(stringData[0])
						if err != nil {
							m.logger.Error("Invalid graphql variable parameter provided", zap.String("param", param.Value.Name), zap.Error(err))
							_ = utils.SendErrorResponse(w, http.StatusBadRequest, fmt.Errorf("invalid request payload received for param '%s'", param.Value.Name))
							return
						}
						variables[param.Value.Name] = b
					}
				}

			}
		} else {
			// For all other methods
			if err := json.NewDecoder(r.Body).Decode(&variables); err != nil {
				m.logger.Error("Invalid graphql request payload provided", zap.Error(err))
				_ = utils.SendErrorResponse(w, http.StatusBadRequest, errors.New("invalid request payload received"))
				return
			}
		}

		if rpc.Authenticate != nil {
			// We need to throw an error if the request is not authenticated
			if err := rpc.Authenticate(ctx, variables); err != nil {
				// We need to throw an error if requested claim is not present in token
				m.logger.Error("Unable to authenticate request",
					zap.String("rpc", rpc.Name), zap.Error(err))
				_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err)
				return
			}
		}

		result, err := rpc.Call(ctx, variables)
		if err != nil {
			m.logger.Error("Unable to execute request",
				zap.String("rpc", rpc.Name), zap.Error(err))
			_ = utils.SendErrorResponse(w, http.StatusInternalServerError, errors.New("unable to execute request"))
			return
		}

		// Send successful response to client
		_ = utils.SendResponse(w, http.StatusOK, result)
	}
}
