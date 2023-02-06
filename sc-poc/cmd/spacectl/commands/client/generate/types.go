package generate

import (
	"fmt"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/spacecloud-io/space-cloud/utils"
)

func generateTypes(doc *openapi3.T) string {
	var b strings.Builder
	for _, pathDef := range doc.Paths {
		if isOperationValidForTypeGen(pathDef.Get) {
			s := generateTypesFromOperation(pathDef.Get)
			_, _ = b.WriteString(s)
		}

		if isOperationValidForTypeGen(pathDef.Delete) {
			s := generateTypesFromOperation(pathDef.Delete)
			_, _ = b.WriteString(s)
		}

		if isOperationValidForTypeGen(pathDef.Post) {
			s := generateTypesFromOperation(pathDef.Post)
			_, _ = b.WriteString(s)
		}

		if isOperationValidForTypeGen(pathDef.Put) {
			s := generateTypesFromOperation(pathDef.Put)
			_, _ = b.WriteString(s)
		}
	}
	return b.String()
}

func isOperationValidForTypeGen(operation *openapi3.Operation) bool {
	if operation == nil {
		return false
	}

	if _, p := operation.Extensions["x-client-gen"]; p {
		return true
	}

	return false
}

func generateTypesFromOperation(operation *openapi3.Operation) string {
	// Prepare the request schema
	requestSchema := &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "object", Properties: openapi3.Schemas{}}}
	if operation.RequestBody != nil {
		requestSchema = operation.RequestBody.Value.Content["application/json"].Schema
	}
	for _, param := range operation.Parameters {
		if len(param.Value.Content) > 0 {
			requestSchema.Value.Properties[param.Value.Name] = param.Value.Content["application/json"].Schema
		} else {
			requestSchema.Value.Properties[param.Value.Name] = param.Value.Schema
		}

		if param.Value.Required {
			requestSchema.Value.Required = append(requestSchema.Value.Required, param.Value.Name)
		}
	}

	// Prepare the response schema
	responseSchema := operation.Responses["200"].Value.Content["application/json"].Schema

	var s string
	if len(requestSchema.Value.Properties) > 0 {
		s = fmt.Sprintf("export interface %sRequest ", getTypeName(operation.OperationID, false))
		s += generateTypeDef(requestSchema, 0)
		s += ";\n\n"
	}

	s += fmt.Sprintf("export interface %sResponse ", getTypeName(operation.OperationID, false))
	s += generateTypeDef(responseSchema, 0)
	s += ";\n\n"

	return s
}

func generateTypeDef(schema *openapi3.SchemaRef, depth int) string {
	switch schema.Value.Type {
	case "object":
		s := "{\n"
		for k, nestedSchema := range schema.Value.Properties {
			s += addPadding(depth + 1)
			s += k
			if !utils.StringExists(k, schema.Value.Required...) {
				s += "?"
			}
			s += ": "
			s += generateTypeDef(nestedSchema, depth+1)
			s += ";\n"
		}
		s += addPadding(depth)
		s += "}"

		return s

	case "array":
		s := generateTypeDef(schema.Value.Items, depth)
		s += "[]"
		return s

	case "string", "number", "boolean":
		s := schema.Value.Type
		return s
	case "integer":
		s := "number"
		return s
	}

	return ""
}

func getTypeName(name string, skipFirst bool) string {
	arr := strings.Split(name, "-")
	for i, item := range arr {
		if i == 0 && skipFirst {
			arr[i] = item
			continue
		}

		arr[i] = strings.Title(item)
	}

	return strings.Join(arr, "")
}

func addPadding(depth int) string {
	s := ""
	for i := 0; i < depth; i++ {
		s += "  "
	}

	return s
}
