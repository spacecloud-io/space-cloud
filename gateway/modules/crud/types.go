package crud

import (
	"context"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
)

type integrationManagerInterface interface {
	InvokeHook(ctx context.Context, params model.RequestParams) config.IntegrationAuthResponse
}
