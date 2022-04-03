package configman

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"

	"github.com/spacecloud-io/space-cloud/model"
)

func loadTypeDefinition(module, typeName string) (*model.TypeDefinition, error) {
	controllerLock.RLock()
	defer controllerLock.RUnlock()

	defs, p := controllerDefinitions[module]
	if !p {
		return nil, fmt.Errorf("provided module '%s' does not exist", module)
	}

	typeDef, p := defs[typeName]
	if !p {
		return nil, fmt.Errorf("type '%s' does not exist in module '%s'", typeName, module)
	}

	return typeDef, nil
}

func loadHook(module string, typeDef *model.TypeDefinition, phase model.HookPhase, loadApp loadApp) (HookImpl, error) {
	controllerLock.RLock()
	defer controllerLock.RUnlock()

	// Check if hooks are defined for that phase
	if typeDef.Hooks == nil {
		return nil, nil
	}
	if _, p := typeDef.Hooks[phase]; !p {
		return nil, nil
	}

	ctrl, err := loadApp(module)
	if err != nil {
		return nil, err
	}

	hookImpl, ok := ctrl.(HookImpl)
	if !ok {
		return nil, fmt.Errorf("controller '%s' doesn't implement hook functionality", module)
	}

	return hookImpl, nil
}

func addOpenAPIPath(app string, types model.Types) {
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
					Content:  openapi3.Content{"application/json": &openapi3.MediaType{Schema: schema}},
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
								Value: &openapi3.Schema{Type: openapi3.TypeArray, Items: schema}}},
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
							"application/json": &openapi3.MediaType{Schema: schema},
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
