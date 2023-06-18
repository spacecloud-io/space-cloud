package golang

import (
	"fmt"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"golang.org/x/tools/imports"
)

func GenerateTypes(spec *openapi3.T, pkgName string) (string, error) {
	var b strings.Builder

	// package name and imports
	pkgOut := fmt.Sprintf("package %s\n\n", pkgName)
	_, _ = b.WriteString(pkgOut)

	// types
	typeDefinitionsOut := generateTypes(spec)
	_, _ = b.WriteString(typeDefinitionsOut)

	// The generation code produces unindented horrors. Use the Go Imports
	// to make it all pretty.
	outBytes, err := imports.Process(pkgName+".go", []byte(b.String()), nil)
	if err != nil {
		return "", fmt.Errorf("error formatting Go code: %w", err)
	}
	return string(outBytes), nil
}

func generateTypes(doc *openapi3.T) string {
	var b strings.Builder

	for _, pathDef := range doc.Paths {
		if isOperationValidForTypeGen(pathDef.Get) {
			s := generateTypesFromOperation(pathDef.Get, "Get")
			_, _ = b.WriteString(s)
		}

		if isOperationValidForTypeGen(pathDef.Delete) {
			s := generateTypesFromOperation(pathDef.Delete, "Delete")
			_, _ = b.WriteString(s)
		}

		if isOperationValidForTypeGen(pathDef.Post) {
			s := generateTypesFromOperation(pathDef.Post, "Post")
			_, _ = b.WriteString(s)
		}

		if isOperationValidForTypeGen(pathDef.Put) {
			s := generateTypesFromOperation(pathDef.Put, "Put")
			_, _ = b.WriteString(s)
		}
	}
	return b.String()
}

func generateTypesFromOperation(operation *openapi3.Operation, method string) string {
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

	// Prepare the request and response schema
	var s string
	opName := getTypeName(operation.OperationID, false)

	if len(requestSchema.Value.Properties) > 0 {
		s += generateTypeDef(requestSchema, opName+"Request")
	}

	for statusCode, content := range operation.Responses {
		responseSchema := content.Value.Content["application/json"].Schema
		s += generateTypeDef(responseSchema, opName+"Response"+statusCode)
	}

	// Prepare the result struct
	structName := opName + "Result"
	s += fmt.Sprintf(`
	// %s produces API's response object
	type %s struct {
		httpResponse *http.Response
	}

	func (r *%s) Raw() io.Reader {
		return r.httpResponse.Body
	}

	func (r *%s) StatusCode() int {
		return r.httpResponse.StatusCode
	}
	`, structName, structName, structName, structName)

	for statusCode := range operation.Responses {
		resultTypeName := opName + "Response" + statusCode
		fnName := "Result"
		if statusCode != "200" {
			fnName += statusCode
		}
		s += fmt.Sprintf(`
		func (r *%s) %s() (*%s, error) {
			if r.StatusCode() != %s {
				return nil, nil
			}
		
			var obj %s
			err := json.NewDecoder(r.Raw()).Decode(&obj)
			if err != nil {
				return nil, err
			}
		
			return &obj, nil
		}
		`, structName, fnName, resultTypeName, statusCode, resultTypeName)
	}

	return s
}

func generateTypeDef(schema *openapi3.SchemaRef, name string) string {
	var s string
	pendingTypes := make(map[string]*openapi3.SchemaRef)

	s += fmt.Sprintf("// %s\n", name)
	s += fmt.Sprintf("type %s struct {\n", name)
	for k, nestedSchema := range schema.Value.Properties {
		required := isRequired(schema.Value.Required, k)
		primitiveGoType := getGoTypes(nestedSchema.Value.Type)
		child := getTypeName(k, false)
		s += fmt.Sprintf("	%s ", child)

		switch nestedSchema.Value.Type {
		case "object":
			typeName := name + "_" + child
			if !required {
				k += ",omitempty"
				s += "*"
			}
			s += fmt.Sprintf("%s `json:%q`\n", typeName, k)
			pendingTypes[typeName] = nestedSchema

		case "array":
			typeName := name + "_" + child
			s += "[]"
			if !required {
				k += ",omitempty"
			}

			arrayType := nestedSchema.Value.Items.Value.Type
			if arrayType == "object" {
				s += fmt.Sprintf("%s `json:%q`\n", typeName, k)
				pendingTypes[typeName] = nestedSchema.Value.Items
			} else if arrayType == "array" {
				s += nestedArray(nestedSchema.Value.Items)
				s += fmt.Sprintf(" `json:%q`\n", k)
			} else {
				s += fmt.Sprintf("%s `json:%q`\n", getGoTypes(arrayType), k)
			}

		case "integer", "number", "string", "boolean":
			if !required {
				k += ",omitempty"
			}
			s += primitiveGoType
			s += fmt.Sprintf(" `json:%q`\n", k)
		}
	}
	s += "}\n\n"
	for name, schema := range pendingTypes {
		s += generateTypeDef(schema, name)
	}
	return s
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
	s1 := strings.Join(arr, "")

	arr = strings.Split(s1, "_")
	for i, item := range arr {
		if i == 0 && skipFirst {
			arr[i] = item
			continue
		}

		arr[i] = strings.Title(item)
	}

	return strings.Join(arr, "")
}

func isRequired(required []string, name string) bool {
	for _, r := range required {
		if r == name {
			return true
		}
	}
	return false
}

func getGoTypes(jsonType string) string {
	switch jsonType {
	case "string":
		return "string"
	case "boolean":
		return "bool"
	case "integer":
		return "int32"
	case "number":
		return "float32"
	}
	return "any"
}

func nestedArray(schema *openapi3.SchemaRef) string {
	s := ""
	if schema.Value.Type == "array" {
		s = "[]"
		s += nestedArray(schema.Value.Items)
	} else {
		s += getGoTypes(schema.Value.Type)
	}
	return s
}
