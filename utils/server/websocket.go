package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/parser"
	"github.com/graphql-go/graphql/language/source"
	uuid "github.com/satori/go.uuid"

	"github.com/mitchellh/mapstructure"

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"
	"github.com/spaceuptech/space-cloud/utils/admin"
	"github.com/spaceuptech/space-cloud/utils/client"
	"github.com/spaceuptech/space-cloud/utils/graphql"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (s *Server) handleWebsocket() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		socket, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("upgrade:", err)
			return
		}

		c := client.CreateWebsocketClient(socket)
		defer s.realtime.RemoveClient(c.ClientID())
		defer s.functions.UnregisterService(c.ClientID())
		defer s.pubsub.UnsubscribeAll(c.ClientID())

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
				mapstructure.Decode(req.Data, data)

				// Subscribe to realtime feed
				feedData, err := s.realtime.Subscribe(ctx, clientID, s.auth, s.crud, data, func(feed *model.FeedData) {
					c.Write(&model.Message{Type: utils.TypeRealtimeFeed, Data: feed})
				})
				if err != nil {
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
				mapstructure.Decode(req.Data, data)

				s.realtime.Unsubscribe(clientID, data)

				// Send response to c
				res := model.RealtimeResponse{Group: data.Group, ID: data.ID, Ack: true}
				c.Write(&model.Message{ID: req.ID, Type: req.Type, Data: res})

			case utils.TypeServiceRegister:
				// TODO add security rule for functions registered as well
				data := new(model.ServiceRegisterRequest)
				mapstructure.Decode(req.Data, data)

				s.functions.RegisterService(clientID, data, func(payload *model.FunctionsPayload) {
					c.Write(&model.Message{Type: utils.TypeServiceRequest, Data: payload})
				})

				c.Write(&model.Message{ID: req.ID, Type: req.Type, Data: map[string]interface{}{"ack": true}})

			case utils.TypeServiceRequest:
				data := new(model.FunctionsPayload)
				mapstructure.Decode(req.Data, data)

				s.functions.HandleServiceResponse(data)
			case utils.TypePubsubSubscribe:
				// For pubsub subscribe event
				data := new(model.PubsubSubscribeRequest)
				mapstructure.Decode(req.Data, data)

				// Subscribe to pubsub feed
				var status int
				var err error
				if data.Queue == "" {
					status, err = s.pubsub.Subscribe(data.Project, data.Token, clientID, data.Subject, func(msg *model.PubsubMsg) {
						c.Write(&model.Message{ID: req.ID, Type: utils.TypePubsubSubscribeFeed, Data: msg})
					})
				} else {
					status, err = s.pubsub.QueueSubscribe(data.Project, data.Token, clientID, data.Subject, data.Queue, func(msg *model.PubsubMsg) {
						c.Write(&model.Message{ID: req.ID, Type: utils.TypePubsubSubscribeFeed, Data: msg})
					})
				}
				if err != nil {
					res := model.PubsubMsgResponse{Status: int32(status), Error: err.Error()}
					c.Write(&model.Message{ID: req.ID, Type: utils.TypePubsubSubscribe, Data: res})
					return true
				}

				// Send response to c
				res := model.PubsubMsgResponse{Status: int32(status)}
				c.Write(&model.Message{ID: req.ID, Type: utils.TypePubsubSubscribe, Data: res})

			case utils.TypePubsubUnsubscribe:
				// For pubsub unsubscribe event
				data := new(model.PubsubSubscribeRequest)
				mapstructure.Decode(req.Data, data)

				status, err := s.pubsub.Unsubscribe(clientID, data.Subject)

				// Send response to c
				if err != nil {
					res := model.PubsubMsgResponse{Status: int32(status), Error: err.Error()}
					c.Write(&model.Message{ID: req.ID, Type: utils.TypePubsubUnsubscribe, Data: res})
					return true
				}

				// Send response to c
				res := model.PubsubMsgResponse{Status: int32(status)}
				c.Write(&model.Message{ID: req.ID, Type: utils.TypePubsubUnsubscribe, Data: res})
			case utils.TypePubsubUnsubscribeAll:
				// For pubsub unsubscribe event
				data := new(model.PubsubSubscribeRequest)
				mapstructure.Decode(req.Data, data)

				status, err := s.pubsub.UnsubscribeAll(clientID)

				// Send response to c
				if err != nil {
					res := model.PubsubMsgResponse{Status: int32(status), Error: err.Error()}
					c.Write(&model.Message{ID: req.ID, Type: utils.TypePubsubUnsubscribeAll, Data: res})
					return true
				}

				// Send response to c
				res := model.PubsubMsgResponse{Status: int32(status)}
				c.Write(&model.Message{ID: req.ID, Type: utils.TypePubsubUnsubscribeAll, Data: res})
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

func (s *Server) handleGraphqlSocket(adminMan *admin.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		project := mux.Vars(r)["project"]

		respHeader := make(http.Header)
		respHeader.Add("Sec-WebSocket-Protocol", "graphql-ws")
		socket, err := upgrader.Upgrade(w, r, respHeader)
		if err != nil {
			log.Println("upgrade:", err)
			return
		}
		defer socket.Close()

		clientID := uuid.NewV1().String()
		defer s.realtime.RemoveClient(clientID)

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
				//				log.Println("Graphql websocket handler error:", err)
				return
			}

			var token string
			switch m.Type {
			case utils.GQL_CONNECTION_INIT:
				// Check if the request is authorised
				token = m.Payload.ConnectionParams.Token
				if err := adminMan.IsTokenValid(token); err != nil {
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
				dbType, err := graphql.GetDBType(v)
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
				feedData, err := s.realtime.Subscribe(ctx, clientID, s.auth, s.crud, data, func(feed *model.FeedData) {
					feed.TypeName = feed.Group
					channel <- &graphqlMessage{ID: feed.QueryID, Type: utils.GQL_DATA, Payload: payloadObject{Data: map[string]interface{}{feed.Group: feed, "__typename": feed.Group}}}
				})

				if err != nil {
					channel <- &graphqlMessage{ID: m.ID, Type: utils.GQL_ERROR, Payload: payloadObject{Error: []gqlError{{Message: err.Error()}}}}
					closeConnAliveRoutine <- true
					return
				}

				for _, feed := range feedData {
					feed.TypeName = feed.Group
					channel <- &graphqlMessage{ID: m.ID, Type: utils.GQL_DATA, Payload: payloadObject{Data: map[string]interface{}{feed.Group: feed, "__typename": feed.Group}}}
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

				s.realtime.Unsubscribe(clientID, data)
				channel <- (&graphqlMessage{ID: m.ID, Type: utils.GQL_STOP, Payload: payloadObject{}})
				graphqlIDMapper.Delete(m.ID)
			}
		}
	}
}
