package apis

import (
	"net/http"
	"sort"

	"github.com/getkin/kin-openapi/openapi3"
)

var (
	// SCErrorResponseDescription describes the description of standard SpaceCloud error response
	SCErrorResponseDescription = "SpaceCloud error response"

	// SCErrorResponseSchema is the standard format for all error response in SpaceCloud
	SCErrorResponseSchema *openapi3.ResponseRef = &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: &SCErrorResponseDescription,
			Content: openapi3.Content{
				"application/json": &openapi3.MediaType{
					Schema: &openapi3.SchemaRef{
						Value: &openapi3.Schema{
							Type: openapi3.TypeObject,
							Properties: openapi3.Schemas{
								"error": &openapi3.SchemaRef{
									Value: &openapi3.Schema{
										Type: openapi3.TypeString,
									},
								},
							},
							Required: []string{"error"},
						},
					},
				},
			},
		},
	}
)

type (
	// HandlerFunc stores a handler
	HandlerFunc func(w http.ResponseWriter, r *http.Request, pathParams map[string]string)
	// App returns the paths it intends to expose
	App interface {
		GetRoutes() []*API
		GetHandler(op string) (HandlerFunc, error)
	}

	// API describes how to handle a particular path
	API struct {
		Path    string             `json:"path"`
		PathDef *openapi3.PathItem `json:"pathDef"`
		Op      string             `json:"op"`

		app string
	}
)

type (
	apps []app
	app  struct {
		name     string
		priority int
	}
)

func (a apps) sort() {
	sort.SliceStable(a, func(i, j int) bool {
		return a[i].priority > a[j].priority
	})
}
