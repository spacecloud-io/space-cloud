package rpc

import (
	"context"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/spacecloud-io/space-cloud/pkg/apis/core/v1alpha1"
)

type (
	// Source describes the implementation of source from the rpc module
	Source interface {
		GetRPCs() RPCs
	}

	// RPCs is an array of RPC objects
	RPCs []*RPC

	// RPC describes the meta information required by the rpc provider
	RPC struct {
		Name          string
		OperationType string
		Extensions    map[string]any

		HTTPOptions *v1alpha1.HTTPOptions
		Plugins     []v1alpha1.HTTPPlugin

		RequestSchema  *openapi3.SchemaRef
		ResponseSchema *openapi3.SchemaRef

		Call         Call
		Authenticate Authenticate
	}

	// Call describes the function to execute when the rpc is invoked
	Call func(ctx context.Context, vars map[string]any) (data any, err error)

	// Authenticate is the function which gets executed for authenticating client request
	Authenticate func(ctx context.Context, vars map[string]any) error
)
