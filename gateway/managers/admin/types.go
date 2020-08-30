package admin

import (
	"context"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
)

// IntegrationInterface s used to describe the features of integration manager we need.
type IntegrationInterface interface {
	HandleConfigAuth(ctx context.Context, resource, op string, claims map[string]interface{}, attr map[string]string) config.IntegrationAuthResponse
	InvokeHook(ctx context.Context, params model.RequestParams) config.IntegrationAuthResponse
}
