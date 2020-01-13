package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/fatih/structs"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/parser"
	"github.com/graphql-go/graphql/language/source"
	"github.com/mitchellh/mapstructure"
	"github.com/segmentio/ksuid"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
	"github.com/spaceuptech/space-cloud/gateway/utils/client"
	"github.com/spaceuptech/space-cloud/gateway/utils/graphql"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (s *Server) handleWebsocket() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		project := vars["project"]

		socket, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("upgrade:", err)
			return
		}

		// Create a new client
		c := client.CreateWebsocketClient(socket)

		state, err := s.projects.LoadProject(project)
		if err != nil {
			log.Println("Websocket error:", err)
			return
		}

		// Unregister the client
		defer state.Realtime.RemoveClient(c.ClientID())

		// Close the client to free up resources
		defer c.Close()
		go c.RoutineWrite()

		// Get c details
		ctx := c.Context()
		clientID := c.ClientID()

		c.Read(func(req *model.Message) bool {
			switch req.Type {
			case utils.TypeRealtimeSubscribe:

				// For realtime subscribe event
				data := new(model.RealtimeRequest)
				_ = mapstructure.Decode(req.Data, data)
				data.Project = project

				// Subscribe to the realtime feed
				feedData, err := state.Realtime.Subscribe(ctx, clientID, data, func(feed *model.FeedData) {
					c.Write(&model.Message{Type: utils.TypeRealtimeFeed, Data: feed})
				})
				if err != nil {
					res := model.RealtimeResponse{Group: data.Group, ID: data.ID, Ack: false, Error: err.Error()}
					c.Write(&model.Message{ID: req.ID, Type: req.Type, Data: res})
					return false
				}

				// Send response to c
				res := model.RealtimeResponse{Group: data.Group, ID: data.ID, Ack: true, Docs: feedData}
				c.Write(&model.Message{ID: req.ID, Type: req.Type, Data: res})

			case utils.TypeRealtimeUnsubscribe:
				// For realtime subscribe event
				data := new(model.RealtimeRequest)
				_ = mapstructure.Decode(req.Data, data)
				data.Project = project

				state.Realtime.Unsubscribe(clientID, data)

				// Send response to c
				res := model.RealtimeResponse{Group: data.Group, ID: data.ID, Ack: true}
				c.Write(&model.Message{ID: req.ID, Type: req.Type, Data: res})
			default:
				c.Write(&model.Message{ID: req.ID, Type: req.Type, Data: map[string]string{"error": "Invalid message type"}})
				return false
			}
			return true
		})
	}
}

type graphqlMessage struct {
	Payload payloadObject `json:"payload"`
	ID      string        `json:"id"`
	Type    string        `json:"type"`
}

type payloadObject struct {
	ConnectionParams token       `json:"connectionParams,omitempty"`
	Query            string      `json:"query,omitempty"`
	Error            []gqlError  `json:"error,omitempty"`
	Data             interface{} `json:"data,omitempty"`
}

type token struct {
	Token string `json:"token,omitempty"`
}

type gqlError struct {
	Message string `json:"message"`
}

var graphqlIDMapper sync.Map

func (s *Server) handleGraphqlSocket() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		project := mux.Vars(r)["project"]

		// Load the project state
		state, err := s.projects.LoadProject(project)
		if err != nil {
			log.Println("Websocket graphql: invalid project provided")
			return
		}

		respHeader := make(http.Header)
		respHeader.Add("Sec-WebSocket-Protocol", "graphql-ws")
		socket, err := upgrader.Upgrade(w, r, respHeader)
		if err != nil {
			log.Println("upgrade:", err)
			return
		}
		defer utils.CloseTheCloser(socket)

		clientID := ksuid.New().String()
		defer state.Realtime.RemoveClient(clientID)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		channel := make(chan *graphqlMessage)
		defer close(channel)
		onlyOnce := true

		closeConnAliveRoutine := make(chan bool)
		defer close(closeConnAliveRoutine)

		go func() {
			for res := range channel {
				err := socket.WriteJSON(res)
				if err != nil {
					log.Println(err)
				}
			}
		}()

		for {

			m := graphqlMessage{}
			if err := socket.ReadJSON(&m); err != nil {
				return
			}

			var token string
			switch m.Type {
			case utils.GQL_CONNECTION_INIT:
				// Check if the request is authorised
				token = m.Payload.ConnectionParams.Token
				if err := s.adminMan.IsTokenValid(token); err != nil {
					channel <- &graphqlMessage{ID: m.ID, Type: utils.GQL_CONNECTION_ERROR, Payload: payloadObject{Error: []gqlError{{Message: err.Error()}}}}
					return
				}

				channel <- (&graphqlMessage{ID: m.ID, Type: utils.GQL_CONNECTION_ACK, Payload: payloadObject{}})

				if onlyOnce {
					onlyOnce = false
					go func() {
						ticker := time.NewTicker(20 * time.Second)
						defer ticker.Stop()
						for {
							select {
							case <-closeConnAliveRoutine:
								return
							case <-ticker.C:
								channel <- &graphqlMessage{ID: m.ID, Type: utils.GQL_CONNECTION_KEEP_ALIVE, Payload: payloadObject{}}
							}
						}
					}()
				}

			case utils.GQL_START:

				parserSource := source.NewSource(&source.Source{
					Body: []byte(m.Payload.Query),
				})
				// parse the source
				doc, err := parser.Parse(parser.ParseParams{Source: parserSource})
				if err != nil {
					channel <- &graphqlMessage{ID: m.ID, Type: utils.GQL_ERROR, Payload: payloadObject{Error: []gqlError{{Message: err.Error()}}}}
					closeConnAliveRoutine <- true
					return
				}

				opDefinition, ok := doc.Definitions[0].(*ast.OperationDefinition)
				if !ok {
					channel <- &graphqlMessage{ID: m.ID, Type: utils.GQL_ERROR, Payload: payloadObject{Error: []gqlError{{Message: errors.New("erros in operation definition of schema").Error()}}}}
					closeConnAliveRoutine <- true
					return
				}

				v, ok := opDefinition.SelectionSet.Selections[0].(*ast.Field)
				if !ok {
					channel <- &graphqlMessage{ID: m.ID, Type: utils.GQL_ERROR, Payload: payloadObject{Error: []gqlError{{Message: errors.New("error in selection set of schema").Error()}}}}
					closeConnAliveRoutine <- true
					return
				}

				whereData, err := graphql.ExtractWhereClause(v.Arguments, utils.M{})
				if err != nil {
					channel <- &graphqlMessage{ID: m.ID, Type: utils.GQL_ERROR, Payload: payloadObject{Error: []gqlError{{Message: err.Error()}}}}
					closeConnAliveRoutine <- true
					return
				}
				dbType, err := state.Graph.GetDBAlias(v)
				if err != nil {
					channel <- &graphqlMessage{ID: m.ID, Type: utils.GQL_ERROR, Payload: payloadObject{Error: []gqlError{{Message: err.Error()}}}}
					closeConnAliveRoutine <- true
					return
				}

				data := &model.RealtimeRequest{Token: token, Where: whereData, DBType: dbType, Project: project, Group: v.Name.Value, Type: m.Type, ID: m.ID}
				for _, dirValue := range v.Arguments {
					if dirValue.Name.Value == "skipInitial" {
						if dirValue.Value.(*ast.BooleanValue).Value {
							data.Options = model.LiveQueryOptions{SkipInitial: true}
						}
					}
				}

				graphqlIDMapper.Store(m.ID, data.Group)

				// Subscribe to realtime feed
				feedData, err := state.Realtime.Subscribe(ctx, clientID, data, func(feed *model.FeedData) {
					feed.TypeName = "subscribe_" + feed.Group

					channel <- &graphqlMessage{ID: feed.QueryID, Type: utils.GQL_DATA, Payload: payloadObject{Data: map[string]interface{}{feed.Group: filterGraphqlSubscriptionResults(v, feed), "find": feed.Find}}}
				})

				if err != nil {
					channel <- &graphqlMessage{ID: m.ID, Type: utils.GQL_ERROR, Payload: payloadObject{Error: []gqlError{{Message: err.Error()}}}}
					closeConnAliveRoutine <- true
					return
				}

				for _, feed := range feedData {
					feed.TypeName = "subscribe_" + feed.Group

					channel <- &graphqlMessage{ID: m.ID, Type: utils.GQL_DATA, Payload: payloadObject{Data: map[string]interface{}{feed.Group: filterGraphqlSubscriptionResults(v, feed)}}}
				}

			case utils.GQL_STOP:
				group, ok := graphqlIDMapper.Load(m.ID)
				if !ok {
					channel <- &graphqlMessage{ID: m.ID, Type: utils.GQL_ERROR, Payload: payloadObject{Error: []gqlError{{Message: errors.New("got " + utils.GQL_STOP + " wanted " + utils.GQL_START).Error()}}}}
					closeConnAliveRoutine <- true
					return
				}
				data := new(model.RealtimeRequest)
				data.ID = m.ID
				data.Group = group.(string)

				state.Realtime.Unsubscribe(clientID, data)
				channel <- (&graphqlMessage{ID: m.ID, Type: utils.GQL_STOP, Payload: payloadObject{}})
				graphqlIDMapper.Delete(m.ID)
			}
		}
	}
}

func filterGraphqlSubscriptionResults(field *ast.Field, feed *model.FeedData) map[string]interface{} {

	filteredResults := map[string]interface{}{}
	feedMap := structs.Map(feed)

	for _, returnField := range field.SelectionSet.Selections {
		returnFieldName := returnField.(*ast.Field).Name.Value

		if returnFieldName == "payload" {
			result := map[string]interface{}{}
			for _, value := range returnField.GetSelectionSet().Selections {
				valueName := value.(*ast.Field).Name.Value
				v, ok := feedMap[returnFieldName]
				if !ok {
					continue
				}
				val, ok := v.(map[string]interface{})
				if !ok {
					continue
				}
				a, ok := val[valueName]
				if ok {
					result[valueName] = a
				}
			}

			result["__typename"] = feed.Group
			filteredResults[returnFieldName] = result
			continue
		}

		value, ok := feedMap[returnFieldName]
		if ok {
			filteredResults[returnFieldName] = value
		}

	}

	return filteredResults
}
