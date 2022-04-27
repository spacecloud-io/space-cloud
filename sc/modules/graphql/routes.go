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
	"github.com/spaceuptech/helpers"
	"go.uber.org/zap"

	"github.com/spacecloud-io/space-cloud/managers/apis"
)

// GetRoutes returns all the apis that are exposed by this app
func (a *App) GetRoutes() []*apis.API {
	projects := make([]interface{}, 0, len(a.dbSchemas))
	for project := range a.dbSchemas {
		projects = append(projects, project)
	}

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
		Parameters: openapi3.Parameters{{
			Value: &openapi3.Parameter{
				Name:     "project",
				In:       "path",
				Required: true,
				Schema: &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type: openapi3.TypeString,
						Enum: projects,
					},
				},
			},
		}},
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

	return []*apis.API{{Path: "/v1/api/{project}/graphql", PathDef: pathDef, Op: "graphql"}}
}

// GetHandler returns the graphql api handler
func (a *App) GetHandler(op string) (apis.HandlerFunc, error) {
	type request struct {
		OperationName string                 `json:"operationName"`
		Query         string                 `json:"query"`
		Variables     map[string]interface{} `json:"variables"`
	}
	h := func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		// Prepare context object
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Minute)
		defer cancel()

		// Parse request body
		req := new(request)
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			a.logger.Error("Invalid graphql request payload provided", zap.Error(err))
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusBadRequest, errors.New("invalid request payload received"))
			return
		}

		// Get the appropriate graphql schema for the project
		project := pathParams["project"]
		schema, p := a.schemas[project]
		if !p {
			a.logger.Error("Invalid project provided in graphql request", zap.String("project", project))
			_ = helpers.Response.SendErrorResponse(ctx, w, http.StatusBadRequest, fmt.Errorf("invalid project '%s' provided", project))
			return
		}

		// Parse the graphql request
		source := source.NewSource(&source.Source{
			Body: []byte(req.Query),
			Name: "GraphQL request",
		})
		rawAst, err := parser.Parse(parser.ParseParams{Source: source})
		if err != nil {
			gErrs := gqlerrors.FormatErrors(err)
			a.logger.Error("Invalid graphql query provided", zap.String("query", req.Query), zap.Any("error", gErrs))
			_ = helpers.Response.SendResponse(ctx, w, http.StatusOK, map[string]interface{}{"errors": gErrs})
			return
		}

		// Execute the query
		result := graphql.Execute(graphql.ExecuteParams{
			Context:       ctx,
			Schema:        *schema,
			AST:           rawAst,
			OperationName: req.OperationName,
			Args:          req.Variables,
			Root:          newStore(),
		})
		_ = helpers.Response.SendResponse(ctx, w, http.StatusOK, result)
	}

	return h, nil
}
