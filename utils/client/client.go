package client

import (
	"context"

	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/proto"
)

// Client is the inteface for the websocket and grpc sockets
type Client interface {
	Write(res *model.Message)
	Read(cb DataCallback)
	RoutineWrite()
	ClientID() string
	Close()
	Context() context.Context
}

// DataCallback is the callback invoked when data is read by the socket
type DataCallback func(data *model.Message) bool

// CreateWebsocketClient makes a client object to manage the socket
func CreateWebsocketClient(socket *websocket.Conn) *WebsocketClient {
	channel := make(chan *model.Message, 5)
	ctx, cancel := context.WithCancel(context.Background())
	id := uuid.NewV1().String()
	return &WebsocketClient{id, channel, ctx, cancel, socket}
}

// CreateGRPCRealtimeClient makes a client object to manage the grpc for realtime
func CreateGRPCRealtimeClient(stream proto.SpaceCloud_RealTimeServer) *GRPCRealtimeClient {
	channel := make(chan *model.Message, 5)
	ctx, cancel := context.WithCancel(context.Background())
	id := uuid.NewV1().String()
	return &GRPCRealtimeClient{id, channel, ctx, cancel, stream}
}

// CreateGRPCServiceClient makes a client object to manage the grpc for services
func CreateGRPCServiceClient(stream proto.SpaceCloud_ServiceServer) *GRPCServiceClient {
	channel := make(chan *model.Message, 5)
	ctx, cancel := context.WithCancel(context.Background())
	id := uuid.NewV1().String()
	return &GRPCServiceClient{id, channel, ctx, cancel, stream}
}

// CreateGRPCPubsubClient makes a client object to manage the grpc for pubsub
func CreateGRPCPubsubClient(stream proto.SpaceCloud_PubsubSubscribeServer) *GRPCPubsubClient {
	channel := make(chan *model.Message, 5)
	ctx, cancel := context.WithCancel(context.Background())
	id := uuid.NewV1().String()
	return &GRPCPubsubClient{id, channel, ctx, cancel, stream}
}