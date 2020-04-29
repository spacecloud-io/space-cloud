package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
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
	"github.com/spaceuptech/space-cloud/gateway/modules"
	"github.com/spaceuptech/space-cloud/gateway/utils"
	"github.com/spaceuptech/space-cloud/gateway/utils/client"
	"github.com/spaceuptech/space-cloud/gateway/utils/graphql"
)

// WebsocketModulesInterface is used to accept the modules object
type WebsocketModulesInterface interface {
	Realtime() modules.RealtimeInterface
	GraphQL() modules.GraphQLInterface
}

// RealtimeInterface is used to accept the realtime module
type RealtimeInterface interface {
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// HandleWebsocket handles all websocket communications
func HandleWebsocket(modules WebsocketModulesInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		projectID := vars["project"]

		socket, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("upgrade:", err)
			return
		}

		c := client.CreateWebsocketClient(socket)
		defer c.Close()

		realtime := modules.Realtime()

		defer realtime.RemoveClient(c.ClientID())

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

				// Subscribe to realtime feed
				feedData, err := realtime.Subscribe(clientID, data, func(feed *model.FeedData) {
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

				if err := realtime.Unsubscribe(clientID, data); err != nil {
					res := model.RealtimeResponse{Group: data.Group, ID: data.ID, Ack: false}
					c.Write(&model.Message{ID: req.ID, Type: req.Type, Data: res})
				}

				// Send response to client
				res := model.RealtimeResponse{Group: data.Group, ID: data.ID, Ack: true}
				c.Write(&model.Message{ID: req.ID, Type: req.Type, Data: res})
			default:
				c.Write(&model.Message{ID: req.ID, Type: req.Type, Data: map[string]string{"error": "Invalid message type"}})
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
	Query     string                 `json:"query,omitempty"`
	Token     string                 `json:"authToken"`
	Variables map[string]interface{} `json:"variables"`
	Error     []gqlError             `json:"errors,omitempty"`
	Data      interface{}            `json:"data,omitempty"`
}

type gqlError struct {
	Message string `json:"message"`
}

// HandleGraphqlSocket handles graphql subscriptions
func HandleGraphqlSocket(modules WebsocketModulesInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		projectID := mux.Vars(r)["project"]

		realtime := modules.Realtime()
		graph := modules.GraphQL()

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

		// Create a new client ID that we will use to make subscriptions
		clientID := ksuid.New().String()
		defer realtime.RemoveClient(clientID)

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

				whereData, err := graphql.ExtractWhereClause(v.Arguments, utils.M{"vars": m.Payload.Variables})
				if err != nil {
					channel <- &graphqlMessage{ID: m.ID, Type: utils.GqlError, Payload: payloadObject{Error: []gqlError{{Message: err.Error()}}}}
					continue
				}

				dbAlias, err := graph.GetDBAlias(v)
				if err != nil {
					channel <- &graphqlMessage{ID: m.ID, Type: utils.GqlError, Payload: payloadObject{Error: []gqlError{{Message: err.Error()}}}}
					continue
				}

				data := &model.RealtimeRequest{Token: m.Payload.Token, Where: whereData, DBType: dbAlias, Project: projectID, Group: v.Name.Value, Type: m.Type, ID: m.ID}
				for _, dirValue := range v.Arguments {
					if dirValue.Name.Value == "skipInitial" {
						if boolVal, ok := dirValue.Value.(*ast.BooleanValue); ok {
							data.Options = model.LiveQueryOptions{SkipInitial: boolVal.Value}
						}
					}
				}

				graphqlIDMapper.Store(m.ID, getGraphQLMapKey(data.DBType, data.Group))

				// Subscribe to realtime feed
				feedData, err := realtime.Subscribe(clientID, data, func(feed *model.FeedData) {
					feed.TypeName = "subscribe_" + feed.Group
					if feed.Type == utils.RealtimeDelete {
						// Make a new map
						find := feed.Find.(map[string]interface{})
						payload := make(map[string]interface{}, len(find))

						// Copy the kev value pairs of find in this new ma
						for k, v := range find {
							payload[k] = v
						}

						// Set the payload
						feed.Payload = payload
					}
					channel <- &graphqlMessage{ID: m.ID, Type: utils.GqlData, Payload: payloadObject{Data: map[string]interface{}{feed.Group: filterGraphqlSubscriptionResults(v, feed)}}}
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
				key, ok := graphqlIDMapper.Load(m.ID)
				if !ok {
					channel <- &graphqlMessage{ID: m.ID, Type: utils.GqlError, Payload: payloadObject{Error: []gqlError{{Message: errors.New("got " + utils.GqlStop + " wanted " + utils.GqlStart).Error()}}}}
					continue
				}
				data := new(model.RealtimeRequest)
				data.ID = m.ID
				data.DBType, data.Group = getValuesFromGraphQLKey(key.(string))

				if err := realtime.Unsubscribe(clientID, data); err != nil {
					channel <- &graphqlMessage{ID: m.ID, Type: utils.GqlError, Payload: payloadObject{Error: []gqlError{{Message: err.Error()}}}}
					continue
				}

				channel <- &graphqlMessage{ID: m.ID, Type: utils.GqlComplete}
				graphqlIDMapper.Delete(m.ID)

			default:
				_ = socket.WriteJSON(&graphqlMessage{ID: m.ID, Type: utils.GqlError, Payload: payloadObject{Error: []gqlError{{Message: fmt.Sprintf("invalid message type (%s) provided", m.Type)}}}})
			}
		}
	}
}

func getGraphQLMapKey(dbAlias, col string) string {
	return fmt.Sprintf("%s--%s", dbAlias, col)
}

func getValuesFromGraphQLKey(key string) (dbAlias, col string) {
	arr := strings.Split(key, "--")
	return arr[0], arr[1]
}

func filterGraphqlSubscriptionResults(field *ast.Field, feed *model.FeedData) map[string]interface{} {

	filteredResults := map[string]interface{}{}
	feedMap := structs.Map(feed)

	for _, returnField := range field.SelectionSet.Selections {
		returnFieldName := returnField.(*ast.Field).Name.Value

		if returnFieldName == "payload" {
			if returnField.GetSelectionSet() == nil {
				filteredResults[returnFieldName] = feedMap[returnFieldName]
				filteredResults[returnFieldName].(map[string]interface{})["__typename"] = feed.Group
				continue
			}

			v, ok := feedMap[returnFieldName]
			if !ok {
				continue
			}

			result := graphql.Filter(returnField.(*ast.Field), v)
			result.(map[string]interface{})["__typename"] = feed.Group
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
