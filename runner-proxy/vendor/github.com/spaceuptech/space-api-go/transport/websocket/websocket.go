package websocket

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"github.com/spaceuptech/space-api-go/config"
	"github.com/spaceuptech/space-api-go/types"
)

type websocketOptions struct {
	projectId string
	token     string
}

type Socket struct {
	url                  string
	isConnect            bool
	isConnecting         bool
	connectedOnce        bool
	options              websocketOptions
	pendingMsg           []types.WebsocketMessage
	socket               *websocket.Conn
	sendMessage          chan types.WebsocketMessage
	registerCallbackMap  map[string]func(data interface{})
	onReconnectCallbacks []func()
	mux                  sync.RWMutex
}

func Init(url string, config *config.Config) *Socket {
	url = "ws://" + url + "/v1/api/" + config.Project + "/socket/json"
	if config.IsSecure {
		url = "wss://" + url + "/v1/api/" + config.Project + "/socket/json"
	}

	s := &Socket{
		url:                 url,
		options:             websocketOptions{projectId: config.Project, token: config.Token},
		registerCallbackMap: map[string]func(data interface{}){},
		pendingMsg:          []types.WebsocketMessage{},
		mux:                 sync.RWMutex{},
	}

	writeMessage := make(chan types.WebsocketMessage)
	s.setWriterChannel(writeMessage)

	// create a websocket writer
	go s.writerRoutine()

	return s
}

func (s *Socket) connect() error {
	if !s.checkIsConnecting() {
		return nil
	}
	conn, _, err := websocket.DefaultDialer.Dial(s.url, nil)
	if err != nil {
		s.resetIsConnecting()
		return err
	}

	s.resetIsConnecting()
	s.setSocket(conn)
	s.setConnected(true)

	if s.isConnectedOnce() {
		for _, fn := range s.onReconnectCallbacks {
			go fn()
		}
	}
	s.setConnectedOnce(true)
	s.sendPendingMessages()
	return nil
}

func (s *Socket) writerRoutine() {
	var isStartReader = true
	for msg := range s.sendMessage {
		if !s.getConnected() {
			s.addPendingMsg(msg)
			continue
		}

		_ = s.socket.WriteJSON(msg)
		if isStartReader {
			go s.read()
			isStartReader = false
		}
	}
}

func (s *Socket) read() {
	for {
		msg := &types.WebsocketMessage{}
		if s.getConnected() {
			if err := s.socket.ReadJSON(msg); err != nil {
				s.setConnected(false)
				time.Sleep(5 * time.Second)
				continue
			}
		} else {
			if err := s.connect(); err != nil {
				time.Sleep(5 * time.Second)
				continue
			}
		}

		if msg != nil {
			if msg.ID != "" {
				cb, ok := s.getRegisteredCallBack(msg.ID)
				if ok {
					go cb(msg.Data)
					s.unregisterCallback(msg.ID)
					continue
				}
			}

			cb, ok := s.getRegisteredCallBack(msg.Type)
			if ok {
				go cb(msg.Data)
			}
		}
	}
}
