package graphql

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/gqlerrors"
	"github.com/graphql-go/graphql/language/parser"
	"github.com/graphql-go/graphql/language/source"
	"go.uber.org/zap"

	"github.com/spacecloud-io/space-cloud/managers/apis"
	"github.com/spacecloud-io/space-cloud/modules/auth"
	"github.com/spacecloud-io/space-cloud/modules/graphql/rootvalue"
	"github.com/spacecloud-io/space-cloud/utils"
)

// GetRoutes returns all the apis that are exposed by this app
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
		Name:    "graphql",
		Path:    "/v1/graphql",
		PathDef: pathDef,
		Handler: func(w http.ResponseWriter, r *http.Request) {
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

			// Parse the graphql request
			source := source.NewSource(&source.Source{
				Body: []byte(req.Query),
				Name: "GraphQL request",
			})
			rawAst, err := parser.Parse(parser.ParseParams{Source: source})
			if err != nil {
				gErrs := gqlerrors.FormatErrors(err)
				a.logger.Error("Invalid graphql query provided", zap.String("query", req.Query), zap.Error(err))
				_ = utils.SendResponse(w, http.StatusOK, map[string]interface{}{"errors": gErrs})
				return
			}

			// Preprocess request for auth
			isAuthRequired, injectedClaims := preprocessForAuth(rawAst)
			if isAuthRequired {
				authResult, p := auth.GetAuthenticationResult(r)
				if !p || !authResult.IsAuthenticated {
					// Send an error if request is unauthenticated
					auth.SendUnauthenticatedResponse(w, r, a.logger, nil)
					return
				}

				// Inject the claims in the variables
				for claim, variable := range injectedClaims {
					v, p := authResult.Claims[claim]
					if !p {
						// We need to throw an error if requested claim is not present in token
						gErrs := gqlerrors.FormatErrors(fmt.Errorf("token does not contain required claim '%s'", claim))
						a.logger.Error("Token does not contain required claim",
							zap.String("query", req.Query), zap.String("claim", claim))
						_ = utils.SendResponse(w, http.StatusOK, map[string]interface{}{"errors": gErrs})
						return
					}

					_ = utils.StoreValue(variable, v, req.Variables)
				}
			}

			// Execute the query
			rootValue := rootvalue.New(rawAst)
			result := graphql.Execute(graphql.ExecuteParams{
				Context:       ctx,
				Schema:        a.schema,
				AST:           rawAst,
				OperationName: req.OperationName,
				Args:          req.Variables,
				Root:          rootValue,
			})

			// Append errors returned from remote sources
			if rootValue.HasErrors() {
				result.Errors = append(result.Errors, rootValue.GetFormatedErrors()...)
			}

			// Send response to client
			_ = utils.SendResponse(w, http.StatusOK, result)
		},
	}}
}
