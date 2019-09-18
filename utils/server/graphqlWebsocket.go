package server

import (
	"context"
	"log"

	"github.com/gorilla/websocket"
)

type graphqlWebsocket interface {
	write(res *graphqlMessage)
	writeRoutine()
	close()
}

type graphqlWebsocketClient struct {
	channel chan *graphqlMessage
	ctx     context.Context
	cancel  context.CancelFunc
	socket  *websocket.Conn
}

func (g *graphqlWebsocketClient) writeRoutine() {
	for res := range g.channel {
		err := g.socket.WriteJSON(res)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func (g *graphqlWebsocketClient) write(res *graphqlMessage) {
	select {
	case g.channel <- res:
	}
}

func (g *graphqlWebsocketClient) close() {
	g.cancel()
	close(g.channel)
	g.socket.Close()
}
