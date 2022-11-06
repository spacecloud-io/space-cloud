package rest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/invopop/jsonschema"
	"github.com/spacecloud-io/space-cloud/managers/apis"
	"github.com/spacecloud-io/space-cloud/modules/graphql"
	"github.com/spacecloud-io/space-cloud/utils"
	"go.uber.org/zap"
)

func (a *App) prepareCompilesGraphqlEndpoints() error {
	arr := make(apis.APIs, len(a.CompiledGraphqlQueries))
	for i, q := range a.CompiledGraphqlQueries {
		// Prepare the path & method
		path := fmt.Sprintf("/v1/api/complied-query/%s", q.Name)
		method := http.MethodPost
		if q.Spec.HTTP != nil {
			if q.Spec.HTTP.URL != "" {
				path = q.Spec.HTTP.URL
			}

			if q.Spec.HTTP.Method != "" {
				method = q.Spec.HTTP.Method
			}
		}

		// Get the compiled graphql query
		var defaultVariables map[string]interface{}
		if len(q.Spec.Graphql.DefaultVariables.Raw) > 0 {
			if err := json.Unmarshal(q.Spec.Graphql.DefaultVariables.Raw, &defaultVariables); err != nil {
				return err
			}
		}
		compiledQuery, err := a.graphqlApp.Compile(q.Spec.Graphql.Query, q.Spec.Graphql.OperationName, defaultVariables, true)
		if err != nil {
			return err
		}

		// Create an open api path definition
		operation := &openapi3.Operation{
			OperationID: q.Name,
			// TODO: Add tags based on the `space-cloud.io/tags` annotations
			Tags: []string{},
		}

		// Add security header is authentication was required
		if compiledQuery.IsAuthRequired {
			operation.Security = &openapi3.SecurityRequirements{{"bearerAuth": []string{}}}
		}

		if q.Spec.HTTP != nil && !(q.Spec.HTTP.Method == http.MethodPost || q.Spec.HTTP.Method == http.MethodPut) {
			// Need to add the variables as query parameters
			parameters := make(openapi3.Parameters, len(compiledQuery.VariableSchema.Properties.Keys()))
			for i, propertyKey := range compiledQuery.VariableSchema.Properties.Keys() {
				t, _ := compiledQuery.VariableSchema.Properties.Get(propertyKey)

				// Prepare the schema object
				schema := t.(*jsonschema.Schema)
				schemaRef := new(openapi3.SchemaRef)
				schemajSON, _ := schema.MarshalJSON()
				_ = schemaRef.UnmarshalJSON([]byte(strings.ReplaceAll(string(schemajSON), "#/$defs/", "#/components/schemas/")))

				parameters[i] = &openapi3.ParameterRef{
					Value: &openapi3.Parameter{
						In:     "query",
						Name:   propertyKey,
						Schema: schemaRef,
					},
				}
			}

			operation.Parameters = parameters
		} else {
			schema := compiledQuery.VariableSchema
			schemaRef := new(openapi3.SchemaRef)
			schemajSON, _ := schema.MarshalJSON()
			_ = schemaRef.UnmarshalJSON([]byte(strings.ReplaceAll(string(schemajSON), "#/$defs/", "#/components/schemas/")))

			operation.RequestBody = &openapi3.RequestBodyRef{
				Value: &openapi3.RequestBody{
					Required: true,
					Content:  openapi3.NewContentWithJSONSchemaRef(schemaRef),
				},
			}
		}

		{
			// Prepare the response body
			schema := compiledQuery.ResponseSchema
			schemaRef := new(openapi3.SchemaRef)
			schemajSON, _ := schema.MarshalJSON()
			_ = schemaRef.UnmarshalJSON(schemajSON)

			operation.Responses = openapi3.Responses{
				"200": &openapi3.ResponseRef{
					Value: openapi3.NewResponse().WithDescription(fmt.Sprintf("Response for %s", q.Name)).WithJSONSchemaRef(schemaRef),
				},
				"400": apis.SCErrorResponseSchema,
				"401": apis.SCErrorResponseSchema,
				"403": apis.SCErrorResponseSchema,
				"500": apis.SCErrorResponseSchema,
			}
		}

		// Prepare the path item
		pathItem := &openapi3.PathItem{}
		switch method {
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
			"x-request-op-type": compiledQuery.OperationType,
		}
		for k, v := range compiledQuery.Extensions {
			operation.Extensions[k] = v
		}

		// Prepare the schemas
		schemas := make(openapi3.Schemas, len(compiledQuery.VariableSchema.Definitions))
		for propertyName, schema := range compiledQuery.VariableSchema.Definitions {
			schemaRef := new(openapi3.SchemaRef)
			schemajSON, _ := schema.MarshalJSON()
			_ = schemaRef.UnmarshalJSON([]byte(strings.ReplaceAll(string(schemajSON), "#/$defs/", "#/components/schemas/")))
			schemas[propertyName] = schemaRef
		}

		arr[i] = &apis.API{
			Name: fmt.Sprintf("complied_query_%s", q.Name),
			Path: path,
			OpenAPI: &apis.OpenAPI{
				PathDef: pathItem,
				Schemas: schemas,
			},
			Handler: a.handleGraphqlRequest(method, compiledQuery, operation),
		}
	}

	a.apis = append(a.apis, arr...)
	return nil
}

func (a *App) handleGraphqlRequest(method string, compiledQuery *graphql.CompiledQuery, operation *openapi3.Operation) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Prepare context object
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Minute)
		defer cancel()

		// Extract variables
		variables := make(map[string]interface{})

		if utils.StringExists(method, http.MethodGet, http.MethodDelete) {
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

				switch param.Value.Schema.Value.Type {
				case "object", "array":
					var data interface{}
					if err := json.Unmarshal([]byte(stringData[0]), &data); err != nil {
						a.logger.Error("Invalid graphql variable parameter provided", zap.String("param", param.Value.Name), zap.Error(err))
						_ = utils.SendErrorResponse(w, http.StatusBadRequest, fmt.Errorf("invalid request payload received for param '%s'", param.Value.Name))
						return
					}

					variables[param.Value.Name] = data
				case "string":
					variables[param.Value.Name] = stringData[0]

				case "integer":
					num, err := strconv.Atoi(stringData[0])
					if err != nil {
						a.logger.Error("Invalid graphql variable parameter provided", zap.String("param", param.Value.Name), zap.Error(err))
						_ = utils.SendErrorResponse(w, http.StatusBadRequest, fmt.Errorf("invalid request payload received for param '%s'", param.Value.Name))
						return
					}
					variables[param.Value.Name] = num

				case "number":
					num, err := strconv.ParseFloat(stringData[0], 64)
					if err != nil {
						a.logger.Error("Invalid graphql variable parameter provided", zap.String("param", param.Value.Name), zap.Error(err))
						_ = utils.SendErrorResponse(w, http.StatusBadRequest, fmt.Errorf("invalid request payload received for param '%s'", param.Value.Name))
						return
					}
					variables[param.Value.Name] = num

				case "boolean":
					b, err := strconv.ParseBool(stringData[0])
					if err != nil {
						a.logger.Error("Invalid graphql variable parameter provided", zap.String("param", param.Value.Name), zap.Error(err))
						_ = utils.SendErrorResponse(w, http.StatusBadRequest, fmt.Errorf("invalid request payload received for param '%s'", param.Value.Name))
						return
					}
					variables[param.Value.Name] = b
				}
			}
		} else {
			// For all other methods
			if err := json.NewDecoder(r.Body).Decode(&variables); err != nil {
				a.logger.Error("Invalid graphql request payload provided", zap.Error(err))
				_ = utils.SendErrorResponse(w, http.StatusBadRequest, errors.New("invalid request payload received"))
				return
			}
		}

		// We need to throw an error if the request is not authenticated
		if err := graphql.AuthenticateRequest(r, compiledQuery, variables); err != nil {
			// We need to throw an error if requested claim is not present in token
			a.logger.Error("Unable to authenticate request",
				zap.String("query", compiledQuery.Query), zap.Error(err))
			_ = utils.SendErrorResponse(w, http.StatusUnauthorized, err)
			return
		}

		// Execute the query
		result := a.graphqlApp.Execute(ctx, compiledQuery, variables)
		if result.HasErrors() {
			a.logger.Error("Unable to execute request",
				zap.String("query", compiledQuery.Query), zap.Reflect("errors", result.Errors))
			_ = utils.SendErrorResponse(w, http.StatusInternalServerError, errors.New("unable to execute request"))
			return
		}

		// Send successful response to client
		_ = utils.SendResponse(w, http.StatusOK, result.Data)
	}
}
