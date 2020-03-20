package handlers

import (
	"context"
	"errors"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/mock"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

type mockRealtimeModule struct {
	mock.Mock

	push []*model.FeedData
}

func TestHandleWebsocket(t *testing.T) {
	type mockArg struct {
		method        string
		args          []interface{}
		paramReturned []interface{}
	}
	tests := []struct {
		name     string
		mockArgs []mockArg
		send     interface{}
		rcv      []*model.Message
		push     []*model.FeedData
	}{
		{
			name: "invalid request",
			mockArgs: []mockArg{
				{
					method:        "RemoveClient",
					args:          []interface{}{mock.Anything},
					paramReturned: []interface{}{},
				},
			},
			send: "abc",
			rcv:  []*model.Message{{Type: "unknown", ID: "0", Data: map[string]interface{}{"error": "invalid request sent"}}},
		},
		{
			name: "invalid message type request",
			mockArgs: []mockArg{
				{
					method:        "RemoveClient",
					args:          []interface{}{mock.Anything},
					paramReturned: []interface{}{},
				},
			},
			send: model.Message{Type: "stupid type", ID: "1"},
			rcv:  []*model.Message{{Type: "stupid type", ID: "1", Data: map[string]interface{}{"error": "Invalid message type"}}},
		},
		{
			name: "subscription request - no data",
			mockArgs: []mockArg{
				{
					method:        "RemoveClient",
					args:          []interface{}{mock.Anything},
					paramReturned: []interface{}{},
				},
				{
					method:        "Subscribe",
					args:          []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything},
					paramReturned: []interface{}{[]*model.FeedData{}, nil},
				},
			},
			send: model.Message{Type: utils.TypeRealtimeSubscribe, ID: "1"},
			rcv:  []*model.Message{{Type: utils.TypeRealtimeSubscribe, ID: "1", Data: map[string]interface{}{"ack": true}}},
		},
		{
			name: "subscription request - with data",
			mockArgs: []mockArg{
				{
					method:        "RemoveClient",
					args:          []interface{}{mock.Anything},
					paramReturned: []interface{}{},
				},
				{
					method:        "Subscribe",
					args:          []interface{}{mock.Anything, mock.Anything, &model.RealtimeRequest{Group: "col", DBType: "db", ID: "q1", Where: map[string]interface{}{"foo": "bar"}}, mock.Anything},
					paramReturned: []interface{}{[]*model.FeedData{}, nil},
				},
			},
			send: model.Message{Type: utils.TypeRealtimeSubscribe, ID: "1", Data: model.RealtimeRequest{Group: "col", DBType: "db", ID: "q1", Where: map[string]interface{}{"foo": "bar"}}},
			rcv:  []*model.Message{{Type: utils.TypeRealtimeSubscribe, ID: "1", Data: map[string]interface{}{"ack": true, "id": "q1", "group": "col"}}},
		},
		{
			name: "subscription request - with data and push",
			mockArgs: []mockArg{
				{
					method:        "RemoveClient",
					args:          []interface{}{mock.Anything},
					paramReturned: []interface{}{},
				},
				{
					method:        "Subscribe",
					args:          []interface{}{mock.Anything, mock.Anything, &model.RealtimeRequest{Group: "col", DBType: "db", ID: "q1", Where: map[string]interface{}{"foo": "bar"}}, mock.Anything},
					paramReturned: []interface{}{[]*model.FeedData{}, nil},
				},
			},
			send: model.Message{Type: utils.TypeRealtimeSubscribe, ID: "1", Data: model.RealtimeRequest{Group: "col", DBType: "db", ID: "q1", Where: map[string]interface{}{"foo": "bar"}}},
			rcv: []*model.Message{
				{Type: utils.TypeRealtimeSubscribe, ID: "1", Data: map[string]interface{}{"ack": true, "id": "q1", "group": "col"}},
				{Type: utils.TypeRealtimeFeed, Data: map[string]interface{}{"group": "col1", "dbType": "dbType1"}},
				{Type: utils.TypeRealtimeFeed, Data: map[string]interface{}{"group": "col2", "dbType": "dbType2"}},
			},
			push: []*model.FeedData{{Group: "col1", DBType: "dbType1"}, {Group: "col2", DBType: "dbType2"}},
		},
		{
			name: "subscription request - with data and response array",
			mockArgs: []mockArg{
				{
					method:        "RemoveClient",
					args:          []interface{}{mock.Anything},
					paramReturned: []interface{}{},
				},
				{
					method:        "Subscribe",
					args:          []interface{}{mock.Anything, mock.Anything, &model.RealtimeRequest{Group: "col", DBType: "db", ID: "q1", Where: map[string]interface{}{"foo": "bar"}}, mock.Anything},
					paramReturned: []interface{}{[]*model.FeedData{{Payload: map[string]interface{}{"foo": "bar"}}}, nil},
				},
			},
			send: model.Message{Type: utils.TypeRealtimeSubscribe, ID: "1", Data: model.RealtimeRequest{Group: "col", DBType: "db", ID: "q1", Where: map[string]interface{}{"foo": "bar"}}},
			rcv:  []*model.Message{{Type: utils.TypeRealtimeSubscribe, ID: "1", Data: map[string]interface{}{"ack": true, "id": "q1", "group": "col", "docs": []interface{}{map[string]interface{}{"payload": map[string]interface{}{"foo": "bar"}}}}}},
		},
		{
			name: "subscription request - incorrect data",
			mockArgs: []mockArg{
				{
					method:        "RemoveClient",
					args:          []interface{}{mock.Anything},
					paramReturned: []interface{}{},
				},
			},
			send: model.Message{Type: utils.TypeRealtimeSubscribe, ID: "1", Data: map[string]int{"group": 98}},
			rcv:  []*model.Message{{Type: utils.TypeRealtimeSubscribe, ID: "1", Data: map[string]interface{}{"ack": false}}},
		},
		{
			name: "subscription request - subscribe throws an error",
			mockArgs: []mockArg{
				{
					method:        "RemoveClient",
					args:          []interface{}{mock.Anything},
					paramReturned: []interface{}{},
				},
				{
					method:        "Subscribe",
					args:          []interface{}{mock.Anything, mock.Anything, &model.RealtimeRequest{}, mock.Anything},
					paramReturned: []interface{}{nil, errors.New("some stupid error")},
				},
			},
			send: model.Message{Type: utils.TypeRealtimeSubscribe, ID: "1", Data: model.RealtimeRequest{}},
			rcv:  []*model.Message{{Type: utils.TypeRealtimeSubscribe, ID: "1", Data: map[string]interface{}{"ack": false}}},
		},
		{
			name: "unsubscribe request - invalid request",
			mockArgs: []mockArg{
				{
					method:        "RemoveClient",
					args:          []interface{}{mock.Anything},
					paramReturned: []interface{}{},
				},
			},
			send: model.Message{Type: utils.TypeRealtimeUnsubscribe, ID: "1", Data: ""},
			rcv:  []*model.Message{{Type: utils.TypeRealtimeUnsubscribe, ID: "1", Data: map[string]interface{}{"ack": false}}},
		},
		{
			name: "unsubscribe request - valid request",
			mockArgs: []mockArg{
				{
					method:        "RemoveClient",
					args:          []interface{}{mock.Anything},
					paramReturned: []interface{}{},
				},
				{
					method:        "Unsubscribe",
					args:          []interface{}{mock.Anything, mock.Anything},
					paramReturned: []interface{}{},
				},
			},
			send: model.Message{Type: utils.TypeRealtimeUnsubscribe, ID: "1", Data: model.RealtimeRequest{}},
			rcv:  []*model.Message{{Type: utils.TypeRealtimeUnsubscribe, ID: "1", Data: map[string]interface{}{"ack": true}}},
		},
		{
			name: "unsubscribe request - valid request with some data",
			mockArgs: []mockArg{
				{
					method:        "RemoveClient",
					args:          []interface{}{mock.Anything},
					paramReturned: []interface{}{},
				},
				{
					method:        "Unsubscribe",
					args:          []interface{}{mock.Anything, &model.RealtimeRequest{ID: "q1", Group: "col"}},
					paramReturned: []interface{}{},
				},
			},
			send: model.Message{Type: utils.TypeRealtimeUnsubscribe, ID: "1", Data: model.RealtimeRequest{ID: "q1", Group: "col"}},
			rcv:  []*model.Message{{Type: utils.TypeRealtimeUnsubscribe, ID: "1", Data: map[string]interface{}{"ack": true, "id": "q1", "group": "col"}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create the mocked struct
			realtime := mockRealtimeModule{push: tt.push}

			// Create the expectations
			for _, m := range tt.mockArgs {
				realtime.On(m.method, m.args...).Return(m.paramReturned...)
			}

			// Create the mock server
			s := httptest.NewServer(HandleWebsocket(realtime))
			defer s.Close()

			// Convert http://127.0.0.1 to ws://127.0.0.
			u := "ws" + strings.TrimPrefix(s.URL, "http")

			// Connect to the server
			ws, _, err := websocket.DefaultDialer.Dial(u, nil)
			if err != nil {
				t.Fatalf("HandleWebsocket() = Unable to connect to server - %v", err)
			}
			defer utils.CloseTheCloser(ws)

			if err := ws.WriteJSON(tt.send); err != nil {
				t.Fatalf("HandleWebsocket() = Unable to send message to server - %v", err)
				return
			}

			for _, m := range tt.rcv {
				res := new(model.Message)
				if err := ws.ReadJSON(res); err != nil {
					t.Fatalf("HandleWebsocket() = Unable to read message to server - %v", err)
					return
				}

				// Check try check if ack if false
				if _, p := m.Data.(map[string]interface{})["ack"]; p {
					if !res.Data.(map[string]interface{})["ack"].(bool) && !m.Data.(map[string]interface{})["ack"].(bool) {
						break
					}
				}

				if !reflect.DeepEqual(m, res) {
					t.Fatalf("HandleWebsocket() = got - %v; wanted - %v", res, m)
				}
			}

			_ = ws.Close()
			time.Sleep(10 * time.Millisecond)
			realtime.AssertExpectations(t)
		})
	}
}

func (m mockRealtimeModule) RemoveClient(clientID string) {
	m.Called(clientID)
}

func (m mockRealtimeModule) Subscribe(ctx context.Context, clientID string, data *model.RealtimeRequest, sendFeed model.SendFeed) ([]*model.FeedData, error) {
	c := m.Called(ctx, clientID, data, sendFeed)
	if err := c.Error(1); err != nil {
		return nil, err
	}

	if m.push != nil {
		go func() {
			time.Sleep(400 * time.Millisecond)
			for _, data := range m.push {
				sendFeed(data)
			}
		}()
	}

	return c.Get(0).([]*model.FeedData), nil
}

func (m mockRealtimeModule) Unsubscribe(clientID string, data *model.RealtimeRequest) {
	m.Called(clientID, data)
}
