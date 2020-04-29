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
	"github.com/graphql-go/graphql/language/ast"
	"github.com/stretchr/testify/mock"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/modules"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

func TestHandleWebsocket(t *testing.T) {
	t.Parallel()
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
					args:          []interface{}{mock.Anything, mock.Anything, mock.Anything},
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
					args:          []interface{}{mock.Anything, &model.RealtimeRequest{Group: "col", DBType: "db", ID: "q1", Where: map[string]interface{}{"foo": "bar"}}, mock.Anything},
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
					args:          []interface{}{mock.Anything, &model.RealtimeRequest{Group: "col", DBType: "db", ID: "q1", Where: map[string]interface{}{"foo": "bar"}}, mock.Anything},
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
					args:          []interface{}{mock.Anything, &model.RealtimeRequest{Group: "col", DBType: "db", ID: "q1", Where: map[string]interface{}{"foo": "bar"}}, mock.Anything},
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
					args:          []interface{}{mock.Anything, &model.RealtimeRequest{}, mock.Anything},
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
			s := httptest.NewServer(HandleWebsocket(&mockWebsocketModules{realtime: &realtime}))
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
			time.Sleep(100 * time.Millisecond)
			realtime.AssertExpectations(t)
		})
	}
}

func TestHandleGraphqlSocket(t *testing.T) {
	t.Parallel()
	type mockArg struct {
		method        string
		args          []interface{}
		paramReturned []interface{}
	}
	tests := []struct {
		name             string
		realtimeMockArgs []mockArg
		graphMockArgs    []mockArg
		send             []*graphqlMessage
		rcv              []*graphqlMessage
		push             []*model.FeedData
	}{
		{
			name: "valid init ack",
			realtimeMockArgs: []mockArg{
				{
					method:        "RemoveClient",
					args:          []interface{}{mock.Anything},
					paramReturned: []interface{}{},
				},
			},
			graphMockArgs: []mockArg{},
			push:          []*model.FeedData{},
			send:          []*graphqlMessage{{Type: utils.GqlConnectionInit, ID: "1"}},
			rcv:           []*graphqlMessage{{Type: utils.GqlConnectionAck, ID: "1"}},
		},
		{
			name: "valid start",
			realtimeMockArgs: []mockArg{
				{
					method:        "RemoveClient",
					args:          []interface{}{mock.Anything},
					paramReturned: []interface{}{},
				},
				{
					method:        "Subscribe",
					args:          []interface{}{mock.Anything, &model.RealtimeRequest{Type: "start", Token: "abc", Group: "col", DBType: "db", ID: "2", Where: map[string]interface{}{"foo": "bar"}}, mock.Anything},
					paramReturned: []interface{}{[]*model.FeedData{{Group: "col", Payload: map[string]interface{}{"f1": "1", "f2": 2}, Find: map[string]interface{}{"foo": "bar"}}}, nil},
				},
			},
			graphMockArgs: []mockArg{{method: "GetDBAlias", args: []interface{}{mock.Anything}, paramReturned: []interface{}{"db", nil}}},
			push:          []*model.FeedData{},
			send: []*graphqlMessage{
				{Type: utils.GqlStart, ID: "2", Payload: payloadObject{Query: `
subscription {
	col(where: {foo: bar}) @db {
    payload {
			f1
		}
		find
  }
}
`, Token: "abc"}},
			},
			rcv: []*graphqlMessage{
				{Type: utils.GqlData, ID: "2", Payload: payloadObject{Data: map[string]interface{}{"col": map[string]interface{}{"payload": map[string]interface{}{"f1": "1", "__typename": "col"}, "find": map[string]interface{}{"foo": "bar"}}}}},
			},
		},
		{
			name: "valid start with payload without selection",
			realtimeMockArgs: []mockArg{
				{
					method:        "RemoveClient",
					args:          []interface{}{mock.Anything},
					paramReturned: []interface{}{},
				},
				{
					method:        "Subscribe",
					args:          []interface{}{mock.Anything, &model.RealtimeRequest{Type: "start", Token: "abc", Group: "col", DBType: "db", ID: "2", Where: map[string]interface{}{"foo": "bar"}}, mock.Anything},
					paramReturned: []interface{}{[]*model.FeedData{{Group: "col", Payload: map[string]interface{}{"f1": "1", "f2": 2}, Find: map[string]interface{}{"foo": "bar"}}}, nil},
				},
			},
			graphMockArgs: []mockArg{{method: "GetDBAlias", args: []interface{}{mock.Anything}, paramReturned: []interface{}{"db", nil}}},
			push:          []*model.FeedData{},
			send: []*graphqlMessage{
				{Type: utils.GqlStart, ID: "2", Payload: payloadObject{Token: "abc", Query: `
subscription {
	col(where: {foo: bar}) @db {
    payload
		find
  }
}
`}},
			},
			rcv: []*graphqlMessage{
				{Type: utils.GqlData, ID: "2", Payload: payloadObject{Data: map[string]interface{}{"col": map[string]interface{}{"payload": map[string]interface{}{"f1": "1", "f2": float64(2), "__typename": "col"}, "find": map[string]interface{}{"foo": "bar"}}}}},
			},
		},
		{
			name: "valid start with payload without selection and push delete feed",
			realtimeMockArgs: []mockArg{
				{
					method:        "RemoveClient",
					args:          []interface{}{mock.Anything},
					paramReturned: []interface{}{},
				},
				{
					method:        "Subscribe",
					args:          []interface{}{mock.Anything, &model.RealtimeRequest{Type: "start", Token: "abc", Group: "col", DBType: "db", ID: "2", Where: map[string]interface{}{"foo": "bar"}}, mock.Anything},
					paramReturned: []interface{}{[]*model.FeedData{{Type: "initial", Group: "col", Payload: map[string]interface{}{"f1": "1", "f2": "2"}, Find: map[string]interface{}{"foo": "bar"}}}, nil},
				},
			},
			graphMockArgs: []mockArg{{method: "GetDBAlias", args: []interface{}{mock.Anything}, paramReturned: []interface{}{"db", nil}}},
			push:          []*model.FeedData{{Group: "col", Type: utils.RealtimeDelete, Find: map[string]interface{}{"foo": "bar"}}},
			send: []*graphqlMessage{
				{Type: utils.GqlConnectionInit, ID: "1"},
				{Type: utils.GqlStart, ID: "2", Payload: payloadObject{Token: "abc", Query: `
subscription {
	col(where: {foo: bar}) @db {
    payload
		find
		type
  }
}
`}},
			},
			rcv: []*graphqlMessage{
				{Type: utils.GqlConnectionAck, ID: "1"},
				{Type: utils.GqlData, ID: "2", Payload: payloadObject{Data: map[string]interface{}{"col": map[string]interface{}{"payload": map[string]interface{}{"f1": "1", "f2": "2", "__typename": "col"}, "find": map[string]interface{}{"foo": "bar"}, "type": "initial"}}}},
				{Type: utils.GqlData, ID: "2", Payload: payloadObject{Data: map[string]interface{}{"col": map[string]interface{}{"payload": map[string]interface{}{"__typename": "col", "foo": "bar"}, "find": map[string]interface{}{"foo": "bar"}, "type": "delete"}}}},
			},
		},
		{
			name: "valid start and query with variables",
			realtimeMockArgs: []mockArg{
				{
					method:        "RemoveClient",
					args:          []interface{}{mock.Anything},
					paramReturned: []interface{}{},
				},
				{
					method:        "Subscribe",
					args:          []interface{}{mock.Anything, &model.RealtimeRequest{Type: "start", Group: "col", DBType: "db", ID: "2", Where: map[string]interface{}{"foo": "bar"}}, mock.Anything},
					paramReturned: []interface{}{[]*model.FeedData{{Group: "col", Payload: map[string]interface{}{"f1": "1", "f2": 2}, Find: map[string]interface{}{"foo": "bar"}}}, nil},
				},
			},
			graphMockArgs: []mockArg{{method: "GetDBAlias", args: []interface{}{mock.Anything}, paramReturned: []interface{}{"db", nil}}},
			push:          []*model.FeedData{},
			send: []*graphqlMessage{
				{Type: utils.GqlStart, ID: "2", Payload: payloadObject{Query: `
subscription {
	col(where: $where) @db {
    payload {
			f1
		}
		find
  }
}
`, Variables: map[string]interface{}{"where": map[string]string{"foo": "bar"}}}},
			},
			rcv: []*graphqlMessage{
				{Type: utils.GqlData, ID: "2", Payload: payloadObject{Data: map[string]interface{}{"col": map[string]interface{}{"payload": map[string]interface{}{"f1": "1", "__typename": "col"}, "find": map[string]interface{}{"foo": "bar"}}}}},
			},
		},
		{
			name: "valid start with skip initial",
			realtimeMockArgs: []mockArg{
				{
					method:        "RemoveClient",
					args:          []interface{}{mock.Anything},
					paramReturned: []interface{}{},
				},
				{
					method:        "Subscribe",
					args:          []interface{}{mock.Anything, &model.RealtimeRequest{Type: "start", Group: "col", DBType: "db", ID: "2", Where: map[string]interface{}{"foo": "bar"}, Options: model.LiveQueryOptions{SkipInitial: true}}, mock.Anything},
					paramReturned: []interface{}{[]*model.FeedData{{Group: "col", Payload: map[string]interface{}{"f1": "1", "f2": 2}, Find: map[string]interface{}{"foo": "bar"}}}, nil},
				},
			},
			graphMockArgs: []mockArg{{method: "GetDBAlias", args: []interface{}{mock.Anything}, paramReturned: []interface{}{"db", nil}}},
			push:          []*model.FeedData{},
			send: []*graphqlMessage{
				{Type: utils.GqlStart, ID: "2", Payload: payloadObject{Query: `
subscription {
	col(where: {foo: bar}, skipInitial: true) @db {
    payload {
			f1
		}
		find
  }
}
`}},
			},
			rcv: []*graphqlMessage{
				{Type: utils.GqlData, ID: "2", Payload: payloadObject{Data: map[string]interface{}{"col": map[string]interface{}{"payload": map[string]interface{}{"f1": "1", "__typename": "col"}, "find": map[string]interface{}{"foo": "bar"}}}}},
			},
		},
		{
			name: "valid start with invalid skip initial",
			realtimeMockArgs: []mockArg{
				{
					method:        "RemoveClient",
					args:          []interface{}{mock.Anything},
					paramReturned: []interface{}{},
				},
				{
					method:        "Subscribe",
					args:          []interface{}{mock.Anything, &model.RealtimeRequest{Type: "start", Token: "abc", Group: "col", DBType: "db", ID: "2", Where: map[string]interface{}{"foo": "bar"}, Options: model.LiveQueryOptions{SkipInitial: false}}, mock.Anything},
					paramReturned: []interface{}{[]*model.FeedData{{Group: "col", Payload: map[string]interface{}{"f1": "1", "f2": 2}, Find: map[string]interface{}{"foo": "bar"}}}, nil},
				},
			},
			graphMockArgs: []mockArg{{method: "GetDBAlias", args: []interface{}{mock.Anything}, paramReturned: []interface{}{"db", nil}}},
			push:          []*model.FeedData{},
			send: []*graphqlMessage{
				{Type: utils.GqlStart, ID: "2", Payload: payloadObject{Token: "abc", Query: `subscription {	col(where: {foo: bar}, skipInitial: "bad value") @db {find}}`}},
			},
			rcv: []*graphqlMessage{
				{Type: utils.GqlData, ID: "2", Payload: payloadObject{Data: map[string]interface{}{"col": map[string]interface{}{"find": map[string]interface{}{"foo": "bar"}}}}},
			},
		},
		{
			name: "valid  start with invalid db alias",
			realtimeMockArgs: []mockArg{
				{
					method:        "RemoveClient",
					args:          []interface{}{mock.Anything},
					paramReturned: []interface{}{},
				},
			},
			graphMockArgs: []mockArg{{method: "GetDBAlias", args: []interface{}{mock.Anything}, paramReturned: []interface{}{"", errors.New("forced error")}}},
			push:          []*model.FeedData{},
			send: []*graphqlMessage{
				{Type: utils.GqlStart, ID: "2", Payload: payloadObject{Query: `subscription {	col(where: {foo: bar}, skipInitial: "bad value") @db {find}}`}},
			},
			rcv: []*graphqlMessage{
				{Type: utils.GqlError, ID: "2"},
			},
		},
		{
			name: "invalid query string",
			realtimeMockArgs: []mockArg{
				{
					method:        "RemoveClient",
					args:          []interface{}{mock.Anything},
					paramReturned: []interface{}{},
				},
			},
			graphMockArgs: []mockArg{},
			push:          []*model.FeedData{},
			send: []*graphqlMessage{
				{Type: utils.GqlStart, ID: "1", Payload: payloadObject{Query: `bad string`}},
			},
			rcv: []*graphqlMessage{
				{Type: utils.GqlError, ID: "1"},
			},
		},
		{
			name: "stop without start",
			realtimeMockArgs: []mockArg{
				{
					method:        "RemoveClient",
					args:          []interface{}{mock.Anything},
					paramReturned: []interface{}{},
				},
			},
			graphMockArgs: []mockArg{},
			push:          []*model.FeedData{},
			send: []*graphqlMessage{
				{Type: utils.GqlStop, ID: "2"},
			},
			rcv: []*graphqlMessage{
				{Type: utils.GqlError, ID: "2"},
			},
		},
		{
			name: "valid stop",
			realtimeMockArgs: []mockArg{
				{
					method:        "RemoveClient",
					args:          []interface{}{mock.Anything},
					paramReturned: []interface{}{},
				},
				{
					method:        "Subscribe",
					args:          []interface{}{mock.Anything, &model.RealtimeRequest{Type: "start", Group: "col", DBType: "db", ID: "2", Where: map[string]interface{}{"foo": "bar"}}, mock.Anything},
					paramReturned: []interface{}{[]*model.FeedData{}, nil},
				},
				{
					method:        "Unsubscribe",
					args:          []interface{}{mock.Anything, &model.RealtimeRequest{DBType: "db", Group: "col", ID: "2"}},
					paramReturned: []interface{}{[]*model.FeedData{}, nil},
				},
			},
			graphMockArgs: []mockArg{{method: "GetDBAlias", args: []interface{}{mock.Anything}, paramReturned: []interface{}{"db", nil}}},
			push:          []*model.FeedData{},
			send: []*graphqlMessage{
				{Type: utils.GqlStart, ID: "2", Payload: payloadObject{Query: `
subscription {
	col(where: {foo: bar}) @db {
    payload {
			f1
		}
		find
  }
}
`}},
				{Type: utils.GqlStop, ID: "2"},
			},
			rcv: []*graphqlMessage{
				{Type: utils.GqlComplete, ID: "2"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create the mocked struct
			realtime := mockRealtimeModule{push: tt.push}
			graph := mockGraphQLModule{}

			// Create the expectations
			for _, m := range tt.realtimeMockArgs {
				realtime.On(m.method, m.args...).Return(m.paramReturned...)
			}
			for _, m := range tt.graphMockArgs {
				graph.On(m.method, m.args...).Return(m.paramReturned...)
			}

			// Create the mock server
			s := httptest.NewServer(HandleGraphqlSocket(&mockWebsocketModules{&realtime, &graph}))
			defer s.Close()

			// Convert http://127.0.0.1 to ws://127.0.0.
			u := "ws" + strings.TrimPrefix(s.URL, "http")

			// Connect to the server
			ws, _, err := websocket.DefaultDialer.Dial(u, nil)
			if err != nil {
				t.Fatalf("HandleGraphQLSocket() = Unable to connect to server - %v", err)
			}
			defer utils.CloseTheCloser(ws)

			for _, msg := range tt.send {
				if err := ws.WriteJSON(msg); err != nil {
					t.Fatalf("HandleGraphQLSocket() = Unable to send message to server - %v", err)
					return
				}
			}

			for _, m := range tt.rcv {
				res := new(graphqlMessage)
				if err := ws.ReadJSON(res); err != nil {
					t.Fatalf("HandleGraphQLSocket() = Unable to read message to server - %v", err)
					return
				}

				if m.Type == utils.GqlError && res.Type == utils.GqlError {
					continue
				}

				if !reflect.DeepEqual(m, res) {
					t.Fatalf("HandleWebsocket() = got - %v; wanted - %v", res, m)
				}
			}

			_ = ws.Close()
			time.Sleep(10 * time.Millisecond)
			realtime.AssertExpectations(t)
			graph.AssertExpectations(t)
		})
	}
}

type mockWebsocketModules struct {
	realtime modules.RealtimeInterface
	graphql  modules.GraphQLInterface
}

func (m *mockWebsocketModules) Realtime() modules.RealtimeInterface {
	return m.realtime
}

func (m *mockWebsocketModules) GraphQL() modules.GraphQLInterface {
	return m.graphql
}

// Create all the mock interfaces
type mockRealtimeModule struct {
	mock.Mock

	push []*model.FeedData
}

func (m *mockRealtimeModule) RemoveClient(clientID string) {
	m.Called(clientID)
}

func (m *mockRealtimeModule) Subscribe(clientID string, data *model.RealtimeRequest, sendFeed model.SendFeed) ([]*model.FeedData, error) {
	c := m.Called(clientID, data, sendFeed)
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

func (m *mockRealtimeModule) Unsubscribe(clientID string, data *model.RealtimeRequest) error {
	m.Called(clientID, data)
	return nil
}

func (m *mockRealtimeModule) HandleRealtimeEvent(ctx context.Context, eventDoc *model.CloudEventPayload) error {
	return m.Called(ctx, eventDoc).Error(0)
}

func (m *mockRealtimeModule) ProcessRealtimeRequests(eventDoc *model.CloudEventPayload) error {
	return m.Called(eventDoc).Error(0)
}

type mockGraphQLModule struct {
	mock.Mock
}

func (m *mockGraphQLModule) GetDBAlias(field *ast.Field) (string, error) {
	c := m.Called(field)
	return c.String(0), c.Error(1)
}

func (m *mockGraphQLModule) ExecGraphQLQuery(ctx context.Context, req *model.GraphQLRequest, token string, cb model.GraphQLCallback) {
	m.Called(ctx, req, token, cb)
}
