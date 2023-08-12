package typescript

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/spacecloud-io/space-cloud/cmd/spacectl/commands/client/generate/driver"
)

func (t *Typescript) GenerateAPIs(doc *openapi3.T) (string, string, error) {
	fileName := "api.ts"
	var b strings.Builder

	// imports and client class
	client := `import * as httpTypes from './types';

export class Client {
	private baseURL: string;

	constructor(baseURL: string) {
		this.baseURL = baseURL;
	}

`
	_, _ = b.WriteString(client)

	// apis
	apis := generateAPIs(doc)
	_, _ = b.WriteString(apis)

	// class closing brace
	_, _ = b.WriteString("}")
	return b.String(), fileName, nil
}

func generateAPIs(spec *openapi3.T) string {
	apis := ""
	for path, pathDef := range spec.Paths {
		apis += getFuncFromOperation(path, http.MethodGet, pathDef.Get)
		apis += getFuncFromOperation(path, http.MethodPost, pathDef.Post)
		apis += getFuncFromOperation(path, http.MethodPut, pathDef.Put)
		apis += getFuncFromOperation(path, http.MethodDelete, pathDef.Delete)
	}

	return apis
}

func getFuncFromOperation(path string, method string, operation *openapi3.Operation) string {
	if !driver.IsOperationValidForTypeGen(operation) {
		return ""
	}

	opName := getTypeName(operation.OperationID, false)
	paramsArg := ""
	queryParams := ""
	body := "null"
	pathVar := fmt.Sprintf("const path: string = this.baseURL + %q;", path)

	if method == http.MethodPost || method == http.MethodPut {
		paramsArg = "data: httpTypes." + opName + "Request"
		body = "JSON.stringify(data)"
	}

	if len(operation.Parameters) != 0 {
		paramsArg = "data: httpTypes." + opName + "Request"
		queryParams = `	const queryParams = new URLSearchParams();
        for (const key in data) {
            const value = data[key as keyof typeof data]
            if (value) queryParams.append(key, value.toString())
            
        }
`
		pathVar = fmt.Sprintf(`const path: string = this.baseURL + "%s?" + queryParams.toString();`, path)
	}

	s := ""
	s += driver.AddPadding(1)
	s += fmt.Sprintf("%s = async (%s): Promise<httpTypes.%s> => {\n", opName, paramsArg, opName+"Response")
	s += queryParams
	s += driver.AddPadding(2)
	s += pathVar
	s += fmt.Sprintf(`
		const options: RequestInit = {
			method: %q,
			headers: {
				'Content-Type': 'application/json'
			},
			body: %s
		}

		try {
			const response = await fetch(path, options);
			if (response.status != 200) {
				const errorMsg: string = "Request failed with status: " + response.status;
				throw new Error(errorMsg);
			}

			const data: httpTypes.%s = await response.json();
			return data;
		} catch (error) {
			throw error
		}
	}
	
`, method, body, opName+"Response")
	return s
}
