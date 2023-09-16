package apis

import (
	"net/http"
	"sync"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/invopop/jsonschema"
)

// The necesary global objects to hold all registered apps
var (
	appsLock sync.RWMutex

	registeredApps apps
)

// RegisterApp marks the app which have routers
func RegisterApp(name string, priority int) {
	appsLock.Lock()
	defer appsLock.Unlock()

	registeredApps = append(registeredApps, app{name, priority})
	registeredApps.sort()
}

// Vars returns the route variables for the current request, if any.
func Vars(r *http.Request) map[string]string {
	if rv := r.Context().Value(pathParamsKey); rv != nil {
		return rv.(map[string]string)
	}

	return nil
}

// PrepareOpenAPIRequest modifies the provided operation to add a request payload based
// on struct provided
func PrepareOpenAPIRequest(operation *openapi3.Operation, modifier OpenAPIPayloadModifier) *openapi3.Operation {
	r := new(jsonschema.Reflector)
	r.ExpandedStruct = true
	r.DoNotReference = true

	jsonSchema := r.Reflect(modifier.Ptr)
	jsonSchema.ID = ""
	jsonSchema.Version = ""
	jsRawRequest, _ := jsonSchema.MarshalJSON()

	openapiRequestJS := &openapi3.SchemaRef{}
	_ = openapiRequestJS.UnmarshalJSON(jsRawRequest)

	operation.RequestBody = &openapi3.RequestBodyRef{Value: &openapi3.RequestBody{
		Required:    true,
		Description: modifier.Description,
		Content:     openapi3.NewContentWithJSONSchemaRef(openapiRequestJS),
	}}

	return operation
}

// PrepareOpenAPIResponse modifies the provided operation to add a standard response payload
// based on struct provided.
func PrepareOpenAPIResponse(operation *openapi3.Operation, modifier OpenAPIPayloadModifier) *openapi3.Operation {
	operation.Responses = openapi3.Responses{
		"400": SCErrorResponseSchema,
		"401": SCErrorResponseSchema,
		"403": SCErrorResponseSchema,
		"500": SCErrorResponseSchema,
	}

	// We want to put a 204 - "No Content" status code if response struct is empty
	if modifier.Ptr == nil {
		operation.Responses["204"] = &openapi3.ResponseRef{
			Value: openapi3.NewResponse().WithDescription(modifier.Description),
		}

		return operation
	}

	r := new(jsonschema.Reflector)
	r.ExpandedStruct = true
	r.DoNotReference = true

	jsonSchema := r.Reflect(modifier.Ptr)
	jsonSchema.ID = ""
	jsonSchema.Version = ""
	jsRawResponse, _ := jsonSchema.MarshalJSON()

	openapiResponseJS := &openapi3.SchemaRef{}
	_ = openapiResponseJS.UnmarshalJSON(jsRawResponse)

	operation.Responses["200"] = &openapi3.ResponseRef{
		Value: openapi3.NewResponse().WithDescription(modifier.Description).WithJSONSchemaRef(openapiResponseJS),
	}

	return operation
}
