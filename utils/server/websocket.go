package server

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"
	"github.com/spaceuptech/space-cloud/utils/client"
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
