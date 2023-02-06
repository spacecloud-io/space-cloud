package generate

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

func generateAPI(name, baseURL string, doc *openapi3.T) string {
	var s string

	s += apiPrefix
	s += "\n"
	s += "// Define a service using a base URL and expected endpoints\n"
	s += "export const api = createApi({\n"
	s += fmt.Sprintf("  reducerPath: \"%s\",\n", name)
	s += fmt.Sprintf("  baseQuery: fetchBaseQuery({ baseUrl: \"%s\" }),\n", baseURL)

	tagsJSON, _ := json.Marshal(getAllTags(doc))
	s += fmt.Sprintf("  tagTypes: %s,\n", string(tagsJSON))
	s += "  endpoints: (builder) => ({\n"
	s += getQueries(doc)
	s += "  }),\n"
	s += "});\n\n"
	s += generateAPIExport(doc)
	return s
}

func generateAPIExport(doc *openapi3.T) string {
	var ops string

	ops += "export const { "

	for _, pathDef := range doc.Paths {
		ops += getExportedFuncFromOperation(pathDef.Get)
		ops += getExportedFuncFromOperation(pathDef.Delete)
		ops += getExportedFuncFromOperation(pathDef.Post)
		ops += getExportedFuncFromOperation(pathDef.Put)
	}
	ops += "} = api;\n"
	return ops
}

func getQueries(doc *openapi3.T) string {
	var s string
	for path, pathDef := range doc.Paths {
		s += getQueryFromOperation(path, http.MethodGet, pathDef.Get)
		s += getQueryFromOperation(path, http.MethodDelete, pathDef.Delete)
		s += getQueryFromOperation(path, http.MethodPost, pathDef.Post)
		s += getQueryFromOperation(path, http.MethodPut, pathDef.Put)
	}
	return s
}

func getExportedFuncFromOperation(operation *openapi3.Operation) string {
	if !isOperationValidForTypeGen(operation) {
		return ""
	}
	ops := "use"
	ops += getTypeName(operation.OperationID, false)
	ops += strings.Title(getQueryType(operation))
	ops += ", "
	return ops
}

func getQueryFromOperation(path, method string, operation *openapi3.Operation) string {
	if !isOperationValidForTypeGen(operation) {
		return ""
	}

	s := addPadding(2)
	s += getTypeName(operation.OperationID, true)
	s += ": "
	s += fmt.Sprintf("builder.%s%s({\n",
		getQueryType(operation), getAPITypes(operation))
	s += fmt.Sprintf("      query: helpers.prepareQuery(\"%s\", \"%s\"),\n", path, method)

	// Check if caching tags are provided
	if tags, p := operation.Extensions["x-tags"]; p {
		if getQueryType(operation) == "query" {
			s += fmt.Sprintf("      providesTags: helpers.getTags(%s),\n", tags.(json.RawMessage))
		} else {
			s += fmt.Sprintf("      invalidatesTags: helpers.getTags(%s),\n", tags.(json.RawMessage))
		}
	}
	s += "    }),\n"
	return s
}

func getQueryType(operation *openapi3.Operation) string {
	opType := operation.Extensions["x-request-op-type"].(json.RawMessage)

	return string(opType)[1 : len(opType)-1]
}

func getAPITypes(operation *openapi3.Operation) string {
	s := "<"
	s += fmt.Sprintf("httpTypes.%sResponse", getTypeName(operation.OperationID, false))
	s += ", "

	// Prepare the request schema
	requestSchema := &openapi3.SchemaRef{Value: &openapi3.Schema{Properties: openapi3.Schemas{}}}
	if operation.RequestBody != nil {
		requestSchema = operation.RequestBody.Value.Content["application/json"].Schema
	}
	for _, param := range operation.Parameters {
		requestSchema.Value.Properties[param.Value.Name] = param.Value.Schema
	}

	if len(requestSchema.Value.Properties) > 0 {
		s += fmt.Sprintf("httpTypes.%sRequest", getTypeName(operation.OperationID, false))
	} else {
		s += "void"
	}

	s += ">"

	return s
}

func getAllTags(doc *openapi3.T) []string {
	tags := map[string]struct{}{}
	for _, pathDef := range doc.Paths {
		for k := range getTagsFromOperation(pathDef.Get) {
			tags[k] = struct{}{}
		}
		for k := range getTagsFromOperation(pathDef.Post) {
			tags[k] = struct{}{}
		}
		for k := range getTagsFromOperation(pathDef.Delete) {
			tags[k] = struct{}{}
		}
		for k := range getTagsFromOperation(pathDef.Put) {
			tags[k] = struct{}{}
		}
	}

	arr := make([]string, 0, len(tags))
	for k := range tags {
		arr = append(arr, k)
	}
	return arr
}

func getTagsFromOperation(operation *openapi3.Operation) map[string]struct{} {
	if operation == nil {
		return nil
	}

	t, p := operation.Extensions["x-tags"]
	if !p {
		return nil
	}

	var tags []struct {
		Key  string `json:"key"`
		Type string `json:"type"`
	}

	_ = json.Unmarshal(t.(json.RawMessage), &tags)

	finalTags := make(map[string]struct{}, len(tags))
	for _, t := range tags {
		finalTags[t.Type] = struct{}{}
	}
	return finalTags
}

var apiPrefix string = `
// Need to use the React-specific entry point to allow generating React hooks
import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react';
import * as helpers from './helpers';
import * as httpTypes from "./types";
`
