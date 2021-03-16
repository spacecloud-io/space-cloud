package transport

import (
	"net/http"

	"github.com/gorilla/websocket"
)

// Transport is the object which handles all communication with the server
type Transport struct {
	// Transport variables
	sslEnabled bool
	addr       string

	// Client drivers
	httpClient *http.Client
	con        *websocket.Conn
}

type CallBackFunction func(string, interface{})

// New initialises a new transport
func New(addr string, sslEnabled bool) *Transport {

	return &Transport{
		sslEnabled: sslEnabled,
		addr:       addr,
		httpClient: &http.Client{},
	}
}

func (t *Transport) GetWebsocketConn() *websocket.Conn {
	return t.con
}