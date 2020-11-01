package modules

import (
	"context"

	"github.com/graphql-go/graphql/language/ast"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// RealtimeInterface is used to mock the realtime module
type RealtimeInterface interface {
	RemoveClient(clientID string)
	Subscribe(clientID string, data *model.RealtimeRequest, sendFeed model.SendFeed) ([]*model.FeedData, error)
	Unsubscribe(ctx context.Context, data *model.RealtimeRequest, clientID string) error

	HandleRealtimeEvent(ctxRoot context.Context, eventDoc *model.CloudEventPayload) error
	ProcessRealtimeRequests(ctx context.Context, eventDoc *model.CloudEventPayload) error
}

// GraphQLInterface is used to mock the graphql module
type GraphQLInterface interface {
	GetDBAlias(ctx context.Context, field *ast.Field, token string, store utils.M) (string, error)
	ExecGraphQLQuery(ctx context.Context, req *model.GraphQLRequest, token string, cb model.GraphQLCallback)
}
