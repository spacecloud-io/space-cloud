package apis

import (
	"net/http"
	"sort"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/spacecloud-io/space-cloud/pkg/apis/core/v1alpha1"
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
	// App returns the paths it intends to expose
	App interface {
		GetAPIRoutes() APIs
	}

	// API describes how to handle a particular path
	API struct {
		Name    string
		Path    string
		OpenAPI *OpenAPI
		Handler http.Handler
		Plugins []v1alpha1.HTTPPlugin

		app string
	}

	OpenAPI struct {
		PathDef *openapi3.PathItem
		Schemas openapi3.Schemas
	}

	// APIs is a collection of APIs
	APIs []*API
)

const (
	pathParamsKey contextKey = iota
)

type (
	contextKey int
	apps       []app
	app        struct {
		name     string
		priority int
	}
)

func (a apps) sort() {
	sort.SliceStable(a, func(i, j int) bool {
		return a[i].priority > a[j].priority
	})
}
