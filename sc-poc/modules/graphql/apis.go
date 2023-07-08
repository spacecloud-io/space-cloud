package graphql

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/graphql-go/graphql/gqlerrors"
	"go.uber.org/zap"

	"github.com/spacecloud-io/space-cloud/managers/apis"
	"github.com/spacecloud-io/space-cloud/utils"
)

// GetAPIRoutes returns all the apis that are exposed by this app
func (a *App) GetAPIRoutes() apis.APIs {

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
	responseSchemaDescription := "Successful GraphQL response"

	pathDef := &openapi3.PathItem{
		Post: &openapi3.Operation{
			Tags:        []string{"GraphQL"},
			Description: "SpaceCloud's GraphQL API",
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

	return []*apis.API{{
		Name: "graphql",
		Path: "/v1/graphql",
		OpenAPI: &apis.OpenAPI{
			PathDef: pathDef,
		},
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Prepare context object
			ctx, cancel := context.WithTimeout(r.Context(), 10*time.Minute)
			defer cancel()

			// Parse request body
			req := new(request)
			if err := json.NewDecoder(r.Body).Decode(req); err != nil {
				a.logger.Error("Invalid graphql request payload provided", zap.Error(err))
				_ = utils.SendErrorResponse(w, http.StatusBadRequest, errors.New("invalid request payload received"))
				return
			}

			// Put an empty variable map if it isn't already defined
			if req.Variables == nil {
				req.Variables = make(map[string]interface{})
			}

			// Compile the graphql query
			compiledQuery, err := a.Compile(req.Query, req.OperationName, nil, false)
			if err != nil {
				gErrs := gqlerrors.FormatErrors(err)
				a.logger.Error("Invalid graphql query provided", zap.String("query", req.Query), zap.Error(err))
				_ = utils.SendResponse(w, http.StatusOK, map[string]interface{}{"errors": gErrs})
				return
			}

			// We need to throw an error if the request is not authenticated
			if err := compiledQuery.AuthenticateRequest(ctx, req.Variables); err != nil {
				// We need to throw an error if requested claim is not present in token
				gErrs := gqlerrors.FormatErrors(err)
				a.logger.Error("Unable to authenticate request",
					zap.String("query", req.Query), zap.Error(err))
				_ = utils.SendResponse(w, http.StatusOK, map[string]interface{}{"errors": gErrs})
				return
			}

			// Execute the query
			result := compiledQuery.Execute(ctx, req.Variables)

			// Send response to client
			_ = utils.SendResponse(w, http.StatusOK, result)
		}),
	}}
}
