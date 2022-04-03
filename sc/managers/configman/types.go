package configman

import (
	"context"

	"github.com/getkin/kin-openapi/openapi3"

	"github.com/spacecloud-io/space-cloud/model"
)

type (
	loadApp func(appName string) (interface{}, error)

	// HookImpl is a controller which implements the hook interface
	HookImpl interface {
		Hook(ctx context.Context, obj *model.ResourceObject) error
	}
)

var (
	okResponseDescription = "SpaceCloud config/operation ok response"
	okResponseSchema      = &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: &okResponseDescription,
			Content: openapi3.Content{
				"application/json": &openapi3.MediaType{
					Schema: &openapi3.SchemaRef{
						Value: &openapi3.Schema{
							Type:       openapi3.TypeObject,
							Properties: openapi3.Schemas{},
						},
					},
				},
			},
		},
	}

	errorResponseDescription = "SpaceCloud config/operation error response"
	errorResponseSchema      = &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: &errorResponseDescription,
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
								"schemaErrors": &openapi3.SchemaRef{
									Value: &openapi3.Schema{
										Type:  openapi3.TypeArray,
										Items: &openapi3.SchemaRef{Value: &openapi3.Schema{Type: openapi3.TypeString}},
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
