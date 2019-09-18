package server

import (
	"context"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/parser"
	"github.com/graphql-go/graphql/language/source"

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
	ConnectionParams token       `json:"connectionParams"`
	Query            string      `json:"query"`
	Error            string      `json:"error"`
	Status           string      `json:"status"`
	Data             interface{} `json:"data"`
}

type token struct {
	Token string `json:"token"`
}

type graphqlSucessData struct {
	Type  string                 `json:"type"`
	Doc   map[string]interface{} `json:"doc"`
	DocID string                 `json:"docId`
}

var graphqlIDMapper sync.Map

func (s *Server) handleGraphqlSocket(adminMan *admin.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		project := mux.Vars(r)["project"]
		socket, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("upgrade:", err)
			return
		}

		ctx, cancel := context.WithCancel(context.Background())
		channel := make(chan *graphqlMessage)
		graphqlWs := &graphqlWebsocketClient{channel: channel, ctx: ctx, cancel: cancel, socket: socket}
		defer graphqlWs.close()

		for {

			m := graphqlMessage{}
			if err := socket.ReadJSON(m); err != nil {
				return
			}
			var token string
			switch m.Type {
			case utils.GqlConnectionInit:
				// Check if the request is authorised
				token = m.Payload.ConnectionParams.Token
				if err := adminMan.IsTokenValid(token); err != nil {
					graphqlWs.write(&graphqlMessage{ID: m.ID, Payload: payloadObject{Status: utils.GqlConnectionError, Error: err.Error()}})
					return
				}
				graphqlWs.write(&graphqlMessage{ID: m.ID, Payload: payloadObject{Status: utils.GqlConnectionAck}})

			case utils.GqlStart:
				data := new(model.RealtimeRequest)

				source := source.NewSource(&source.Source{
					Body: []byte(m.Payload.Query),
				})
				// parse the source
				doc, err := parser.Parse(parser.ParseParams{Source: source})
				if err != nil {
					graphqlWs.write(&graphqlMessage{ID: m.ID, Payload: payloadObject{Status: utils.GqlConnectionError, Error: err.Error()}})
					return
				}
				v := doc.Definitions[0].(*ast.OperationDefinition).SelectionSet.Selections[0].(*ast.Field)

				for _, dirValue := range v.Arguments {
					if dirValue.Name.Value == "where" {
						whereData, err := graphql.ExtractWhereClause(v.Arguments, utils.M{})
						if err != nil {
							log.Println(err)
							return
						}
						data.Where = whereData
					} else if dirValue.Name.Value == "skipInitial" {
						if dirValue.Value.(*ast.BooleanValue).Value {
							data.Options = model.LiveQueryOptions{SkipInitial: true}
						}
					}
				}
				data.Token = token
				data.DBType = v.Directives[0].Name.Value
				data.Project = project
				data.Group = v.Name.Value
				data.Type = m.Type
				data.ID = m.ID
				graphqlIDMapper.Store(m.ID, data.Group)

				// Subscribe to realtime feed
				feedData, err := s.realtime.Subscribe(ctx, m.ID, s.auth, s.crud, data, func(feed *model.FeedData) {
					graphqlWs.write(&graphqlMessage{ID: feed.QueryID, Payload: payloadObject{Status: utils.GqlData, Data: graphqlSucessData{Type: feed.Type, Doc: feed.Payload, DocID: feed.DocID}}})
				})

				if err != nil {
					graphqlWs.write(&graphqlMessage{ID: m.ID, Payload: payloadObject{Status: utils.GqlError, Error: err.Error()}})
					return
				}

				graphqlWs.write(&graphqlMessage{ID: m.ID, Payload: payloadObject{Status: utils.GqlStart, Data: feedData}})

			case utils.GqlStop:
				// For realtime subscribe event
				group, ok := graphqlIDMapper.Load(m.ID)
				if !ok {
					return
				}
				data := new(model.RealtimeRequest)
				data.ID = m.ID
				data.Group = group.(string)

				s.realtime.Unsubscribe(m.ID, data)
				graphqlWs.write(&graphqlMessage{ID: m.ID, Payload: payloadObject{Status: utils.GqlStop}})

			}
		}
	}
}
