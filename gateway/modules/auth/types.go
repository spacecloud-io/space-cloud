package auth

import (
	"context"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
)

type adminMan interface {
	GetSecret() string
}
type integrationManagerInterface interface {
	InvokeHook(ctx context.Context, params model.RequestParams) config.IntegrationAuthResponse
}
