package client

import (
	"context"
	"log"

	"github.com/gorilla/websocket"

	"github.com/spaceuptech/space-cloud/model"
)

// WebsocketClient is the websocket client
type WebsocketClient struct {
	id      string
	channel chan *model.Message
	ctx     context.Context
	cancel  context.CancelFunc
	socket  *websocket.Conn
}

// RoutineWrite starts a json writer routine
func (c *WebsocketClient) RoutineWrite() {
	for res := range c.channel {
		err := c.socket.WriteJSON(res)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

// Write writes the object to the client
func (c *WebsocketClient) Write(res *model.Message) {
	select {
	case <-c.ctx.Done():
	case c.channel <- res:
	}
}

// Close closes the client
func (c *WebsocketClient) Close() {
	c.cancel()
	close(c.channel)
	c.socket.Close()
}

// Read starts a blocking reader routine
func (c *WebsocketClient) Read(cb DataCallback) {
	for {
		data := &model.Message{}
		err := c.socket.ReadJSON(data)
		if err != nil {
			return
		}

		// Close the reader if callback returned false
		next := cb(data)
		if !next {
			return
		}
	}
}

// ClientID returns the client's id
func (c *WebsocketClient) ClientID() string {
	return c.id
}

// Context returns the client's context
func (c *WebsocketClient) Context() context.Context {
	return c.ctx
}
