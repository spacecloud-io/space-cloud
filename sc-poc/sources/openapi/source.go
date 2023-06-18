package openapi

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"sync"

	"github.com/caddyserver/caddy/v2"
	"github.com/getkin/kin-openapi/openapi3"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/mitchellh/mapstructure"
	"github.com/spacecloud-io/space-cloud/managers/source"
	"github.com/spacecloud-io/space-cloud/modules/rpc"
	"github.com/spacecloud-io/space-cloud/pkg/apis/core/v1alpha1"
)

var openapisourcesResource = schema.GroupVersionResource{Group: "core.space-cloud.io", Version: "v1alpha1", Resource: "openapisources"}

func init() {
	source.RegisterSource(Source{}, openapisourcesResource)
}

// Source describes the compiled openapi
type Source struct {
	v1alpha1.OpenAPISource

	// Internal stuff
	wg     *sync.WaitGroup
	logger *zap.Logger
}

// CaddyModule returns the Caddy module information.
func (Source) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  caddy.ModuleID(source.GetModuleName(openapisourcesResource)),
		New: func() caddy.Module { return new(Source) },
	}
}

// Provision provisions the source
func (s *Source) Provision(ctx caddy.Context) error {
	s.logger = ctx.Logger(s)
	s.wg = &sync.WaitGroup{}
	s.wg.Add(1)
	return nil
}

// GetPriority returns the priority of the source. Higher
func (s *Source) GetPriority() int {
	return 0
}

// GetProviders returns the providers this source is applicable for
func (s *Source) GetProviders() []string {
	return []string{"rpc"}
}

// GetRPCs returns the rpcs generated by the openapi spec
func (s *Source) GetRPCs() rpc.RPCs {
	loader := openapi3.NewLoader()
	var schema *openapi3.T

	if s.Spec.OpenAPI.Ref != nil {
		// Load the OpenAPI specification from url
		url, err := url.Parse(s.Spec.OpenAPI.Ref.URL)
		if err != nil {
			s.logger.Error("Unable to parse url to load openapi spec", zap.Error(err))
			return nil
		}
		// Load the OpenAPI specification from URL data
		schema, err = loader.LoadFromURI(url)
		if err != nil {
			s.logger.Error("Unable to load openapi spec", zap.Any("url", url), zap.Error(err))
			return nil
		}
	} else {
		// Marshal the OpenAPI specification to JSON
		data, err := json.Marshal(s.Spec.OpenAPI.Value)
		if err != nil {
			s.logger.Error("Unable to marshal openapi spec", zap.Error(err))
			return nil
		}
		// Load the OpenAPI specification from JSON data
		schema, err = loader.LoadFromData(data)
		if err != nil {
			s.logger.Error("Unable to load openapi spec", zap.Error(err))
			return nil
		}
	}

	rpcs := make([]*rpc.RPC, 0)
	for url, pathSpec := range schema.Paths {
		// Process GET requests
		if getSpec := pathSpec.Get; getSpec != nil {
			var respBody *openapi3.Schema
			var reqBody *openapi3.Schema
			var operationType string
			var plugins []v1alpha1.HTTPPlugin

			// Check if the response has a 200 status code
			if res, ok := getSpec.Responses[strconv.FormatInt(int64(200), 10)]; ok {
				// Check if the response has JSON content
				if content := res.Value.Content.Get("application/json"); content != nil {
					respBody = content.Schema.Value
				}
			}
			// Check if operation type is present in extensions
			if ext, ok := getSpec.Extensions["x-request-op-type"]; ok {
				operationType = ext.(string)
			}
			// Check if plugins is present in extensions
			if ext, ok := getSpec.Extensions["x-sc-plugins"]; ok {
				// export value of ext to plugins variable
				if err := mapstructure.Decode(ext, &plugins); err != nil {
					s.logger.Error("Unable to decode x-sc-plugins", zap.Error(err))
				}
			}
			rpcs = append(rpcs, s.createRPC(url, getSpec.OperationID, operationType, http.MethodGet, respBody, reqBody, getSpec.Parameters, plugins))
		}

		// Process POST requests
		if postSpec := pathSpec.Post; postSpec != nil {
			var respBody *openapi3.Schema
			var reqBody *openapi3.Schema
			var operationType string
			var plugins []v1alpha1.HTTPPlugin

			// Check if the response has a 200 status code
			if res, ok := postSpec.Responses[strconv.FormatInt(int64(200), 10)]; ok {
				// Check if the response has JSON content
				if content := res.Value.Content.Get("application/json"); content != nil {
					respBody = content.Schema.Value
				}
			}
			// Check if the response has a 204 status code
			if res, ok := postSpec.Responses[strconv.FormatInt(int64(204), 10)]; ok {
				// Check if the response has JSON content
				if content := res.Value.Content.Get("application/json"); content != nil {
					respBody = content.Schema.Value
				}
			}
			// Get the request body schema
			if postSpec.RequestBody != nil {
				if reqContent := postSpec.RequestBody.Value.Content.Get("application/json"); reqContent != nil {
					reqBody = reqContent.Schema.Value
				}
			}
			// Check if operation type is present in extensions
			if ext, ok := postSpec.Extensions["x-request-op-type"]; ok {
				operationType = ext.(string)
			}
			// Check if plugins is present in extensions
			if ext, ok := postSpec.Extensions["x-sc-plugins"]; ok {
				// export value of ext to plugins variable
				if err := mapstructure.Decode(ext, &plugins); err != nil {
					s.logger.Error("Unable to decode x-sc-plugins", zap.Error(err))
				}
			}
			rpcs = append(rpcs, s.createRPC(url, postSpec.OperationID, operationType, http.MethodPost, respBody, reqBody, postSpec.Parameters, plugins))
		}

		// Process PUT requests
		if putSpec := pathSpec.Put; putSpec != nil {
			var respBody *openapi3.Schema
			var reqBody *openapi3.Schema
			var operationType string
			var plugins []v1alpha1.HTTPPlugin

			// Check if the response has a 200 status code
			if res, ok := putSpec.Responses[strconv.FormatInt(int64(200), 10)]; ok {
				// Check if the response has JSON content
				if content := res.Value.Content.Get("application/json"); content != nil {
					respBody = content.Schema.Value
				}
			}
			// Check if the response has a 204 status code
			if res, ok := putSpec.Responses[strconv.FormatInt(int64(204), 10)]; ok {
				// Check if the response has JSON content
				if content := res.Value.Content.Get("application/json"); content != nil {
					respBody = content.Schema.Value
				}
			}
			// Get the request body schema
			if putSpec.RequestBody != nil {
				if reqContent := putSpec.RequestBody.Value.Content.Get("application/json"); reqContent != nil {
					reqBody = reqContent.Schema.Value
				}
			}
			// Check if operation type is present in extensions
			if ext, ok := putSpec.Extensions["x-request-op-type"]; ok {
				operationType = ext.(string)
			}
			// Check if plugins is present in extensions
			if ext, ok := putSpec.Extensions["x-sc-plugins"]; ok {
				// export value of ext to plugins variable
				if err := mapstructure.Decode(ext, &plugins); err != nil {
					s.logger.Error("Unable to decode x-sc-plugins", zap.Error(err))
				}
			}
			rpcs = append(rpcs, s.createRPC(url, putSpec.OperationID, operationType, http.MethodPut, respBody, reqBody, putSpec.Parameters, plugins))
		}

		// Process DELETE requests
		if deleteSpec := pathSpec.Delete; deleteSpec != nil {
			var respBody *openapi3.Schema
			var reqBody *openapi3.Schema
			var operationType string
			var plugins []v1alpha1.HTTPPlugin

			// Check if the response has a 200 status code
			if res, ok := deleteSpec.Responses[strconv.FormatInt(int64(200), 10)]; ok {
				// Check if the response has JSON content
				if content := res.Value.Content.Get("application/json"); content != nil {
					respBody = content.Schema.Value
				}
			}
			// Check if the response has a 204 status code
			if res, ok := deleteSpec.Responses[strconv.FormatInt(int64(204), 10)]; ok {
				// Check if the response has JSON content
				if content := res.Value.Content.Get("application/json"); content != nil {
					respBody = content.Schema.Value
				}
			}
			// Get the request body schema
			if deleteSpec.RequestBody != nil {
				if reqContent := deleteSpec.RequestBody.Value.Content.Get("application/json"); reqContent != nil {
					reqBody = reqContent.Schema.Value
				}
			}
			// Check if operation type is present in extensions
			if ext, ok := deleteSpec.Extensions["x-request-op-type"]; ok {
				operationType = ext.(string)
			}
			// Check if plugins is present in extensions
			if ext, ok := deleteSpec.Extensions["x-sc-plugins"]; ok {
				// export value of ext to plugins variable
				if err := mapstructure.Decode(ext, &plugins); err != nil {
					s.logger.Error("Unable to decode x-sc-plugins", zap.Error(err))
				}
			}
			rpcs = append(rpcs, s.createRPC(url, deleteSpec.OperationID, operationType, http.MethodDelete, respBody, reqBody, deleteSpec.Parameters, plugins))
		}
	}

	return rpcs
}

// Interface guards
var (
	_ caddy.Provisioner = (*Source)(nil)
	_ source.Source     = (*Source)(nil)
	_ rpc.Source        = (*Source)(nil)
)
