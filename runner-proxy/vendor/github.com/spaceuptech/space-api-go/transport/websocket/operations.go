package websocket

import (
	"errors"
	"time"

	"github.com/segmentio/ksuid"

	"github.com/spaceuptech/space-api-go/types"
)

// RegisterOnReconnectCallback registers a callback
func (s *Socket) RegisterOnReconnectCallback(function func()) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.onReconnectCallbacks = append(s.onReconnectCallbacks, function)
}

// RegisterCallback registers a callback
func (s *Socket) RegisterCallback(evType string, function func(data interface{})) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.registerCallbackMap[evType] = function
}

func (s *Socket) DeregisterCallback(evType string) {
	s.mux.Lock()
	defer s.mux.Unlock()
	delete(s.registerCallbackMap, evType)
}

// Request sends a message to the server expecting a response in a request-response pattern
func (s *Socket) Request(msgType string, data interface{}) (interface{}, error) {
	if !s.getConnected() {
		// connect to server
		if err := s.connect(); err != nil {
			return false, err
		}
	}

	id := s.Send(msgType, data)

	timer1 := time.NewTimer(10 * time.Second)
	defer timer1.Stop()

	// channel for receiving service register acknowledgement
	ch := make(chan interface{})
	defer close(ch)

	s.RegisterCallback(id, func(data interface{}) {
		ch <- data
	})
	defer s.DeregisterCallback(id)

	select {
	case <-timer1.C:
		return false, errors.New("response time elapsed")
	case msg := <-ch:
		return msg, nil
	}
}

// Send sends a message to server over websocket protocol
func (s *Socket) Send(Type string, data interface{}) string {
	id := ksuid.New().String()
	s.sendMessage <- types.WebsocketMessage{ID: id, Type: Type, Data: data}
	return id
}
