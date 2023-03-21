package golang

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"golang.org/x/tools/imports"
)

func GenerateAPI(spec *openapi3.T, pkgName string) (string, error) {
	var b strings.Builder

	// package name and imports
	importsOut := generateImports(pkgName)
	_, _ = b.WriteString(importsOut)

	// client
	clientOut := generateClient(spec)
	_, _ = b.WriteString(clientOut)

	// The generation code produces unindented horrors. Use the Go Imports
	// to make it all pretty.
	outBytes, err := imports.Process(pkgName+".go", []byte(b.String()), nil)
	if err != nil {
		return "", fmt.Errorf("error formatting Go code: %w", err)
	}
	return string(outBytes), nil
}

func generateImports(pkgName string) string {
	s := fmt.Sprintf("package %s\n\n", pkgName)

	imports := []string{"context", "net/http", "net/url", "bytes", "encoding/json"}
	s += "import (\n"
	for _, imp := range imports {
		s += fmt.Sprintf("	%q\n", imp)
	}
	s += ")\n\n"
	return s
}

func generateClient(doc *openapi3.T) string {
	s := `
// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The endpoint of the server. All the paths in
	// the swagger spec will be appended to the server.
	Server string

	// Client for performing requests.
	Client *http.Client
}

// Creates a new Client, with reasonable defaults
func NewClient(server string) (*Client, error) {
	// create a client with sane default values
	client := Client{
		Server: server,
	}

	// create httpClient, if not already present
	if client.Client == nil {
		client.Client = &http.Client{}
	}
	return &client, nil
}`

	// Generate operation methods
	s += "\n\n"
	for path, pathDef := range doc.Paths {
		s += getFuncFromOperation(path, http.MethodGet, pathDef.Get)
		s += getFuncFromOperation(path, http.MethodPost, pathDef.Post)
		s += getFuncFromOperation(path, http.MethodPut, pathDef.Put)
		s += getFuncFromOperation(path, http.MethodDelete, pathDef.Delete)
	}
	return s
}

func getFuncFromOperation(path, method string, operation *openapi3.Operation) string {
	if !isOperationValidForTypeGen(operation) {
		return ""
	}

	opName := getTypeName(operation.OperationID, false)
	paramsOut := getOpParams(operation.Parameters)

	paramsArg := ""
	if paramsOut != "" {
		paramsArg = "params " + opName + "Request"
	}
	var s string
	switch method {
	case "GET":
		s += fmt.Sprintf("func (c *Client) %s(ctx context.Context, %s) (*%s, error) {\n", opName, paramsArg, opName+"Response")
		s += addPadding(1)
		s += fmt.Sprintf("path := c.Server + %q\n", path)
		s += fmt.Sprintf(`
	url, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	%s

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	var obj %s
	err = json.NewDecoder(resp.Body).Decode(&obj)
	if err != nil {
		return nil, err
	}

	return &obj, nil
}`, paramsOut, opName+"Response")
		s += "\n\n"

	case "POST":
		s += fmt.Sprintf("func (c *Client) %s(ctx context.Context, body %s) (*%s, error) {\n", opName, opName+"Request", opName+"Response")
		s += fmt.Sprintf("path := c.Server + %q\n", path)
		s += fmt.Sprintf(`
	url, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader := bytes.NewReader(buf)

	req, err := http.NewRequest("POST", url.String(), bodyReader)
	if err != nil {
		return nil, err
	}

	resp, err :=  c.Client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	var obj %s
	err = json.NewDecoder(resp.Body).Decode(&obj)
	if err != nil {
		return nil, err
	}

	return &obj, nil
}`, opName+"Response")
		s += "\n\n"

	case "DELETE":
		s += fmt.Sprintf("func (c *Client) %s(ctx context.Context, %s) (*%s, error) {\n", opName, paramsArg, opName+"Response")
		s += addPadding(1)
		s += fmt.Sprintf("path := c.Server + %q\n", path)
		s += fmt.Sprintf(`
	url, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	%s

	req, err := http.NewRequest("DELETE", url.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	var obj %s
	err = json.NewDecoder(resp.Body).Decode(&obj)
	if err != nil {
		return nil, err
	}

	return &obj, nil
}`, paramsOut, opName+"Response")
		s += "\n\n"

	case "PUT":
		s += fmt.Sprintf("func (c *Client) %s(ctx context.Context, body %s) (*%s, error) {\n", opName, opName+"Request", opName+"Response")
		s += fmt.Sprintf("path := c.Server + %q\n", path)
		s += fmt.Sprintf(`
	url, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader := bytes.NewReader(buf)

	req, err := http.NewRequest("PUT", url.String(), bodyReader)
	if err != nil {
		return nil, err
	}

	resp, err :=  c.Client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	var obj %s
	err = json.NewDecoder(resp.Body).Decode(&obj)
	if err != nil {
		return nil, err
	}

	return &obj, nil
}`, opName+"Response")
		s += "\n\n"
	}

	return s
}

func getOpParams(params openapi3.Parameters) string {
	if len(params) == 0 {
		return ""
	}

	s := "queryValues := url.Query()\n"
	for _, p := range params {
		s += fmt.Sprintf(`
		queryValues.Add(%q, fmt.Sprint(params.%s))
		`, p.Value.Name, getTypeName(p.Value.Name, false))
	}
	s += "url.RawQuery = queryValues.Encode()\n"
	return s
}
