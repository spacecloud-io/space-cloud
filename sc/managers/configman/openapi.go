package configman

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"

	"github.com/spacecloud-io/space-cloud/model"
)

func addOperationToOpenAPIDoc(openapiDoc openapi3.T, app string, types model.OperationTypes) error {
	for typeName, typeDef := range types {
		// Prepare the request and response schema refs
		requestSchema := new(openapi3.SchemaRef)
		responseSchema := new(openapi3.SchemaRef)

		// Parse the schemas
		if typeDef.RequestSchema != nil {
			data, _ := json.Marshal(typeDef.RequestSchema)
			_ = requestSchema.UnmarshalJSON(data)
		}
		if typeDef.ResponseSchema != nil {
			data, _ := json.Marshal(typeDef.ResponseSchema)
			_ = responseSchema.UnmarshalJSON(data)
		}

		// Create a path item if it doesn't already exists
		path := fmt.Sprintf("/v1/operation/%s/%s", app, typeName)
		if _, p := openapiDoc.Paths[path]; !p {
			params := prepareParams(typeDef.RequiredParents, false)
			openapiDoc.Paths[path] = &openapi3.PathItem{Parameters: params}
		}
		pathItem := openapiDoc.Paths[path]

		// Create a request body if specified
		var requestBody *openapi3.RequestBodyRef
		if typeDef.RequestSchema != nil {
			// Marshal the schema
			data, _ := json.Marshal(typeDef.RequestSchema)
			schema := new(openapi3.SchemaRef)
			_ = schema.UnmarshalJSON(data)

			requestBody = &openapi3.RequestBodyRef{Value: &openapi3.RequestBody{
				Required: true,
				Content:  openapi3.Content{"application/json": &openapi3.MediaType{Schema: schema}},
			}}
		}

		// Create a response body if specified
		responses := openapi3.Responses{
			"400": errorResponseSchema,
			"401": errorResponseSchema,
			"403": errorResponseSchema,
			"500": errorResponseSchema,
		}
		if typeDef.ResponseSchema != nil {
			// Marshal the schema
			data, _ := json.Marshal(typeDef.ResponseSchema)
			schema := new(openapi3.SchemaRef)
			_ = schema.UnmarshalJSON(data)

			responses["200"] = &openapi3.ResponseRef{Value: &openapi3.Response{
				Content: openapi3.Content{"application/json": &openapi3.MediaType{Schema: schema}},
			}}
		}

		operation := &openapi3.Operation{
			Tags:        []string{strings.Title(app)},
			Description: fmt.Sprintf("Operate on '%s' from the '%s' app", typeName, app),
			OperationID: fmt.Sprintf("operate%s%s", strings.Title(app), strings.Title(typeName)),
			RequestBody: requestBody,
			Responses:   responses,
		}

		// Add the operation to the path item based on the input method
		switch typeDef.Method {
		case http.MethodPost:
			pathItem.Post = operation
		case http.MethodPut:
			pathItem.Put = operation
		case http.MethodGet:
			pathItem.Get = operation
		case http.MethodDelete:
			pathItem.Delete = operation
		default:
			return fmt.Errorf("invalid method provided for operation '%s/%s'", app, typeName)
		}
	}

	return nil
}

func addConfigToOpenAPIDoc(openapiDoc openapi3.T, app string, types model.ConfigTypes) {
	for typeName, typeDef := range types {
		// Get the schema ref
		data, _ := json.Marshal(typeDef.Schema)
		schema := new(openapi3.SchemaRef)
		_ = schema.UnmarshalJSON(data)

		// Add the single list path
		openapiDoc.Paths[fmt.Sprintf("/v1/config/%s/%s", app, typeName)] = &openapi3.PathItem{
			Parameters: prepareParams(typeDef.RequiredParents, false),
			Post: &openapi3.Operation{
				Tags:        []string{strings.Title(app)},
				Description: fmt.Sprintf("Apply '%s' from the '%s' app", typeName, app),
				OperationID: fmt.Sprintf("apply%s%s", strings.Title(app), strings.Title(typeName)),
				RequestBody: &openapi3.RequestBodyRef{Value: &openapi3.RequestBody{
					Required: true,
					Content:  openapi3.Content{"application/json": &openapi3.MediaType{Schema: getResourceObjectSchema(schema)}},
				}},
				Responses: openapi3.Responses{
					"200": okResponseSchema,
					"400": errorResponseSchema,
					"401": errorResponseSchema,
					"403": errorResponseSchema,
					"500": errorResponseSchema,
				},
			},
			Get: &openapi3.Operation{
				Tags:        []string{strings.Title(app)},
				Description: fmt.Sprintf("Read multiple '%s' from the '%s' app", typeName, app),
				OperationID: fmt.Sprintf("readMany%s%s", strings.Title(app), strings.Title(typeName)),
				Responses: openapi3.Responses{
					"200": &openapi3.ResponseRef{Value: &openapi3.Response{
						Content: openapi3.Content{
							"application/json": &openapi3.MediaType{Schema: &openapi3.SchemaRef{
								Value: &openapi3.Schema{
									Type: model.TypeObject,
									Properties: openapi3.Schemas{
										"list": &openapi3.SchemaRef{Value: &openapi3.Schema{
											Type:  openapi3.TypeArray,
											Items: getResourceObjectSchema(schema),
										}},
									},
								},
							}},
						},
					}},
					"400": errorResponseSchema,
					"401": errorResponseSchema,
					"403": errorResponseSchema,
					"500": errorResponseSchema,
				},
			},
			Delete: &openapi3.Operation{
				Tags:        []string{strings.Title(app)},
				Description: fmt.Sprintf("Delete multiple '%s' from the '%s' app", typeName, app),
				OperationID: fmt.Sprintf("deleteMany%s%s", strings.Title(app), strings.Title(typeName)),
				Responses: openapi3.Responses{
					"200": okResponseSchema,
					"400": errorResponseSchema,
					"401": errorResponseSchema,
					"403": errorResponseSchema,
					"500": errorResponseSchema,
				},
			},
		}
		openapiDoc.Paths[fmt.Sprintf("/v1/config/%s/%s/{name}", app, typeName)] = &openapi3.PathItem{
			Parameters: prepareParams(typeDef.RequiredParents, true),
			Get: &openapi3.Operation{
				Tags:        []string{strings.Title(app)},
				Description: fmt.Sprintf("Read single '%s' from the '%s' app", typeName, app),
				OperationID: fmt.Sprintf("readSingle%s%s", strings.Title(app), strings.Title(typeName)),
				Responses: openapi3.Responses{
					"200": &openapi3.ResponseRef{Value: &openapi3.Response{
						Content: openapi3.Content{
							"application/json": &openapi3.MediaType{Schema: getResourceObjectSchema(schema)},
						},
					}},
					"400": errorResponseSchema,
					"401": errorResponseSchema,
					"403": errorResponseSchema,
					"500": errorResponseSchema,
				},
			},

			Delete: &openapi3.Operation{
				Tags:        []string{strings.Title(app)},
				Description: fmt.Sprintf("Delete single '%s' from the '%s' app", typeName, app),
				OperationID: fmt.Sprintf("deleteSingle%s%s", strings.Title(app), strings.Title(typeName)),
				Responses: openapi3.Responses{
					"200": okResponseSchema,
					"400": errorResponseSchema,
					"401": errorResponseSchema,
					"403": errorResponseSchema,
					"500": errorResponseSchema,
				},
			},
		}
	}
}

func prepareParams(requiredParents []string, includePathParam bool) openapi3.Parameters {
	params := make(openapi3.Parameters, len(requiredParents))
	for i, parent := range requiredParents {
		params[i] = &openapi3.ParameterRef{
			Value: &openapi3.Parameter{
				In:     "query",
				Name:   parent,
				Schema: &openapi3.SchemaRef{Value: &openapi3.Schema{Type: openapi3.TypeString}},
			},
		}
	}

	// Add the name path parameter
	if includePathParam {
		params = append(params, &openapi3.ParameterRef{Value: &openapi3.Parameter{
			In:     "path",
			Name:   "name",
			Schema: &openapi3.SchemaRef{Value: &openapi3.Schema{Type: openapi3.TypeString}},
		}})
	}
	return params
}

func getResourceObjectSchema(spec *openapi3.SchemaRef) *openapi3.SchemaRef {
	return &openapi3.SchemaRef{Value: &openapi3.Schema{
		Type: model.TypeObject,
		Properties: openapi3.Schemas{
			"meta": &openapi3.SchemaRef{Value: &openapi3.Schema{
				Type: model.TypeObject,
				Properties: openapi3.Schemas{
					"module": &openapi3.SchemaRef{Value: &openapi3.Schema{Type: model.TypeString}},
					"type":   &openapi3.SchemaRef{Value: &openapi3.Schema{Type: model.TypeString}},
					"name":   &openapi3.SchemaRef{Value: &openapi3.Schema{Type: model.TypeString}},
					"parents": &openapi3.SchemaRef{Value: &openapi3.Schema{
						Type:                 model.TypeObject,
						AdditionalProperties: &openapi3.SchemaRef{Value: &openapi3.Schema{Type: model.TypeString}},
					}},
				},
			}},
			"spec": spec,
		},
		Required: []string{"meta", "spec"},
	}}
}
