package client

import (
	"context"

	"github.com/gorilla/websocket"
	"github.com/segmentio/ksuid"

	"github.com/spaceuptech/space-cloud/model"
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
	id := ksuid.New().String()
	return &WebsocketClient{id, channel, ctx, cancel, socket}
}
