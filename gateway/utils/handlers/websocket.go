package handlers

import (
	"errors"
	"fmt"
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
	"github.com/sirupsen/logrus"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
	"github.com/spaceuptech/space-cloud/gateway/utils/client"
	"github.com/spaceuptech/space-cloud/gateway/utils/graphql"
)

// RealtimeInterface is used to accept the realtime module
type RealtimeInterface interface {
	RemoveClient(clientID string)
	Subscribe(clientID string, data *model.RealtimeRequest, sendFeed model.SendFeed) ([]*model.FeedData, error)
	Unsubscribe(clientID string, data *model.RealtimeRequest)
}

// GraphQLWebsocketInterface is sued to accept the graphql module
type GraphQLWebsocketInterface interface {
	GetDBAlias(field *ast.Field) (string, error)
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// HandleWebsocket handles all websocket communications
func HandleWebsocket(realtime RealtimeInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		projectID := vars["project"]

		socket, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("upgrade:", err)
			return
		}

		state, err := s.projects.LoadProject(projectID)
		if err != nil {
			log.Println("Websocket error:", err)
			return
		}

		// Create a new client
		c := client.CreateWebsocketClient(socket)
		defer c.Close()

		// Unregister the client
		defer state.Realtime.RemoveClient(c.ClientID())

		go c.RoutineWrite()

		// Get client details
		clientID := c.ClientID()

		c.Read(func(req *model.Message) bool {
			switch req.Type {
			case utils.TypeRealtimeSubscribe:

				// For realtime subscribe event
				data := new(model.RealtimeRequest)
				if err := mapstructure.Decode(req.Data, data); err != nil {
					logrus.Errorf("Unable to decode incoming subscription request - %v", err)
					res := model.RealtimeResponse{Ack: false, Error: err.Error()}
					c.Write(&model.Message{ID: req.ID, Type: req.Type, Data: res})
					return true
				}
				data.Project = projectID

				// Subscribe to the realtime feed
				feedData, err := state.Realtime.Subscribe(ctx, clientID, data, func(feed *model.FeedData) {
					c.Write(&model.Message{Type: utils.TypeRealtimeFeed, Data: feed})
				})
				if err != nil {
					logrus.Errorf("Unable to process incoming subscription request - %v", err)
					res := model.RealtimeResponse{Group: data.Group, ID: data.ID, Ack: false, Error: err.Error()}
					c.Write(&model.Message{ID: req.ID, Type: req.Type, Data: res})
					return true
				}

				// Send response to c
				res := model.RealtimeResponse{Group: data.Group, ID: data.ID, Ack: true, Docs: feedData}
				c.Write(&model.Message{ID: req.ID, Type: req.Type, Data: res})

			case utils.TypeRealtimeUnsubscribe:
				// For realtime subscribe event
				data := new(model.RealtimeRequest)
				if err := mapstructure.Decode(req.Data, data); err != nil {
					logrus.Errorf("Unable to decode incoming subscription request - %v", err)
					res := model.RealtimeResponse{Ack: false, Error: err.Error()}
					c.Write(&model.Message{ID: req.ID, Type: req.Type, Data: res})
					return true
				}
				data.Project = projectID

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

// HandleGraphqlSocket handles graphql subscriptions
func HandleGraphqlSocket(projects *projects.Projects) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		projectID := mux.Vars(r)["project"]

		// Load the project state
		state, err := projects.LoadProject(projectID)
		if err != nil {
			log.Println("Websocket graphql: invalid project provided")
			return
		}

		// Create a map to store subscription ids
		var graphqlIDMapper sync.Map

		respHeader := make(http.Header)
		respHeader.Add("Sec-WebSocket-Protocol", "graphql-ws")
		socket, err := upgrader.Upgrade(w, r, respHeader)
		if err != nil {
			log.Println("upgrade:", err)
			return
		}
		defer utils.CloseTheCloser(socket)

		// Variable to store token between consecutive requests
		var token string

		// Create a new client ID that we will use to make subscriptions
		clientID := ksuid.New().String()
		defer state.Realtime.RemoveClient(clientID)

		// Make a channel to send graphql responses
		channel := make(chan *graphqlMessage)
		defer close(channel)

		// Flag to mark processing just once
		onlyOnce := true

		// Channel to close the indefinite ticker on close
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
			// Read json payload
			m := graphqlMessage{}
			if err := socket.ReadJSON(&m); err != nil {
				_ = socket.WriteJSON(&graphqlMessage{ID: m.ID, Type: utils.GqlError, Payload: payloadObject{Error: []gqlError{{Message: err.Error()}}}})
				return
			}

			switch m.Type {
			case utils.GqlConnectionInit:
				// Check if the request is authorised
				token = m.Payload.ConnectionParams.Token
				channel <- &graphqlMessage{ID: m.ID, Type: utils.GqlConnectionAck, Payload: payloadObject{}}

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
								channel <- &graphqlMessage{ID: m.ID, Type: utils.GqlConnectionKeepAlive, Payload: payloadObject{}}
							}
						}
					}()
				}

			case utils.GqlStart:

				// parse the source
				doc, err := parser.Parse(parser.ParseParams{Source: source.NewSource(&source.Source{Body: []byte(m.Payload.Query)})})
				if err != nil {
					channel <- &graphqlMessage{ID: m.ID, Type: utils.GqlError, Payload: payloadObject{Error: []gqlError{{Message: err.Error()}}}}
					continue
				}

				opDefinition, ok := doc.Definitions[0].(*ast.OperationDefinition)
				if !ok {
					channel <- &graphqlMessage{ID: m.ID, Type: utils.GqlError, Payload: payloadObject{Error: []gqlError{{Message: errors.New("erros in operation definition of schema").Error()}}}}
					continue
				}

				v, ok := opDefinition.SelectionSet.Selections[0].(*ast.Field)
				if !ok {
					channel <- &graphqlMessage{ID: m.ID, Type: utils.GqlError, Payload: payloadObject{Error: []gqlError{{Message: errors.New("error in selection set of schema").Error()}}}}
					continue
				}

				whereData, err := graphql.ExtractWhereClause(v.Arguments, utils.M{})
				if err != nil {
					channel <- &graphqlMessage{ID: m.ID, Type: utils.GqlError, Payload: payloadObject{Error: []gqlError{{Message: err.Error()}}}}
					continue
				}

				dbAlias, err := state.Graph.GetDBAlias(v)
				if err != nil {
					channel <- &graphqlMessage{ID: m.ID, Type: utils.GqlError, Payload: payloadObject{Error: []gqlError{{Message: err.Error()}}}}
					continue
				}

				data := &model.RealtimeRequest{Token: token, Where: whereData, DBType: dbAlias, Project: projectID, Group: v.Name.Value, Type: m.Type, ID: m.ID}
				for _, dirValue := range v.Arguments {
					if dirValue.Name.Value == "skipInitial" {
						if boolVal, ok := dirValue.Value.(*ast.BooleanValue); ok {
							data.Options = model.LiveQueryOptions{SkipInitial: boolVal.Value}
						}
					}
				}

				graphqlIDMapper.Store(m.ID, data.Group)

				// Subscribe to realtime feed
				feedData, err := state.Realtime.Subscribe(ctx, clientID, data, func(feed *model.FeedData) {
					feed.TypeName = "subscribe_" + feed.Group

					channel <- &graphqlMessage{ID: feed.QueryID, Type: utils.GqlData, Payload: payloadObject{Data: map[string]interface{}{feed.Group: filterGraphqlSubscriptionResults(v, feed), "find": feed.Find}}}
				})

				if err != nil {
					channel <- &graphqlMessage{ID: m.ID, Type: utils.GqlError, Payload: payloadObject{Error: []gqlError{{Message: err.Error()}}}}
					continue
				}

				for _, feed := range feedData {
					feed.TypeName = "subscribe_" + feed.Group
					channel <- &graphqlMessage{ID: m.ID, Type: utils.GqlData, Payload: payloadObject{Data: map[string]interface{}{feed.Group: filterGraphqlSubscriptionResults(v, feed)}}}
				}

			case utils.GqlStop:
				group, ok := graphqlIDMapper.Load(m.ID)
				if !ok {
					channel <- &graphqlMessage{ID: m.ID, Type: utils.GqlError, Payload: payloadObject{Error: []gqlError{{Message: errors.New("got " + utils.GqlStop + " wanted " + utils.GqlStart).Error()}}}}
					continue
				}
				data := new(model.RealtimeRequest)
				data.ID = m.ID
				data.Group = group.(string)

				state.Realtime.Unsubscribe(clientID, data)
				channel <- (&graphqlMessage{ID: m.ID, Type: utils.GQL_STOP, Payload: payloadObject{}})
				graphqlIDMapper.Delete(m.ID)

			default:
				_ = socket.WriteJSON(&graphqlMessage{ID: m.ID, Type: utils.GqlError, Payload: payloadObject{Error: []gqlError{{Message: fmt.Sprintf("invalid message type (%s) provided", m.Type)}}}})
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
