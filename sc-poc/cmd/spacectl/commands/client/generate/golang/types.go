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

	// Prepare the response schema
	responseSchema := operation.Responses["200"].Value.Content["application/json"].Schema

	opName := getTypeName(operation.OperationID, false)
	var s string
	if len(requestSchema.Value.Properties) > 0 {
		s += generateTypeDef(requestSchema, opName+"Request", method, opName)
	}

	s += generateTypeDef(responseSchema, opName+"Response", method, opName)
	return s
}

func generateTypeDef(schema *openapi3.SchemaRef, name string, method string, opName string) string {
	var s string
	pendingTypes := make(map[string]*openapi3.SchemaRef)

	s += fmt.Sprintf("// %s\n", name)
	s += fmt.Sprintf("type %s struct {\n", name)
	for k, nestedSchema := range schema.Value.Properties {
		typeName := getTypeName(k, false)
		s += fmt.Sprintf("	%s ", typeName)

		switch nestedSchema.Value.Type {
		case "object":
			typeName = method + "_" + typeName + "_" + opName
			s += fmt.Sprintf("%s `json:%q`\n", typeName, k)
			pendingTypes[typeName] = nestedSchema

		case "array":
			typeName = method + "_" + typeName + "_" + opName
			s += fmt.Sprintf("[]%s `json:%q`\n", typeName, k)
			pendingTypes[typeName] = nestedSchema.Value.Items

		case "integer", "number":
			s += fmt.Sprintf("int32 `json:%q`\n", k)

		case "string":
			s += fmt.Sprintf("string `json:%q`\n", k)

		case "boolean":
			s += fmt.Sprintf("bool `json:%q`\n", k)

		}
	}
	s += "}\n\n"
	for name, schema := range pendingTypes {
		s += generateTypeDef(schema, name, method, opName)
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
