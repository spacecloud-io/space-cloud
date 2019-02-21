package utils

import (
	"context"
	"log"

	"github.com/gorilla/websocket"
	"github.com/satori/go.uuid"

	"github.com/spaceuptech/space-cloud/model"
)

// Client is the object handling all client interactions
type Client struct {
	id      string
	channel chan interface{}
	ctx     context.Context
	cancel  context.CancelFunc
	socket  *websocket.Conn
}

// DataCallback is the callback invoked when data is read by the socket
type DataCallback func(data *model.Message)

// RoutineWrite starts a json writer routine
func (c *Client) RoutineWrite() {
	go func() {
		for res := range c.channel {
			err := c.socket.WriteJSON(res)
			if err != nil {
				log.Println(err)
			}
		}
	}()
}

// Write wrties the object to the client
func (c *Client) Write(res interface{}) {
	select {
	case c.channel <- res:
	case <-c.ctx.Done():
	}
}

// Close closes the client
func (c *Client) Close() {
	c.cancel()
	c.socket.Close()
}

// Read startes a blocking reader routine
func (c *Client) Read(cb DataCallback) {
	defer c.Close()
	for {
		data := &model.Message{}
		err := c.socket.ReadJSON(data)
		if err != nil {
			log.Println(err)
			return
		}

		cb(data)
	}
}

// ClientID returns the client's id
func (c *Client) ClientID() string {
	return c.id
}

// CreateWebsocketClient makes a client object to manage the socket
func CreateWebsocketClient(socket *websocket.Conn) *Client {
	channel := make(chan interface{}, 5)
	ctx, cancel := context.WithCancel(context.Background())
	id := uuid.NewV1().String()
	return &Client{id, channel, ctx, cancel, socket}
}
