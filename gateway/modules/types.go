package modules

import (
	"context"

	"github.com/graphql-go/graphql/language/ast"

	"github.com/spaceuptech/space-cloud/gateway/model"
)

// RealtimeInterface is used to mock the realtime module
type RealtimeInterface interface {
	RemoveClient(clientID string)
	Subscribe(clientID string, data *model.RealtimeRequest, sendFeed model.SendFeed) ([]*model.FeedData, error)
	Unsubscribe(clientID string, data *model.RealtimeRequest) error

	HandleRealtimeEvent(ctxRoot context.Context, eventDoc *model.CloudEventPayload) error
	ProcessRealtimeRequests(eventDoc *model.CloudEventPayload) error
}

// GraphQLInterface is used to mock the graphql module
type GraphQLInterface interface {
	GetDBAlias(field *ast.Field) (string, error)
	ExecGraphQLQuery(ctx context.Context, req *model.GraphQLRequest, token string, cb model.GraphQLCallback)
}
