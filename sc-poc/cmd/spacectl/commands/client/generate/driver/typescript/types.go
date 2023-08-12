package typescript

import (
	"fmt"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/spacecloud-io/space-cloud/cmd/spacectl/commands/client/generate/driver"
	"github.com/spacecloud-io/space-cloud/utils"
)

func (t *Typescript) GenerateTypes(doc *openapi3.T) (string, string, error) {
	fileName := "types.ts"
	var b strings.Builder
	for _, pathDef := range doc.Paths {
		if driver.IsOperationValidForTypeGen(pathDef.Get) {
			s := generateTypesFromOperation(pathDef.Get)
			_, _ = b.WriteString(s)
		}

		if driver.IsOperationValidForTypeGen(pathDef.Delete) {
			s := generateTypesFromOperation(pathDef.Delete)
			_, _ = b.WriteString(s)
		}

		if driver.IsOperationValidForTypeGen(pathDef.Post) {
			s := generateTypesFromOperation(pathDef.Post)
			_, _ = b.WriteString(s)
		}

		if driver.IsOperationValidForTypeGen(pathDef.Put) {
			s := generateTypesFromOperation(pathDef.Put)
			_, _ = b.WriteString(s)
		}
	}
	return b.String(), fileName, nil
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
			if nestedSchema.Value.Type == "null" {
				continue
			}

			s += driver.AddPadding(depth + 1)
			s += k
			if !utils.StringExists(k, schema.Value.Required...) {
				s += "?"
			}
			s += ": "
			s += generateTypeDef(nestedSchema, depth+1)
			s += ";\n"
		}
		s += driver.AddPadding(depth)
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
	case "undefined":
		s := "any"
		return s
	}

	return ""
}
