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
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("upgrade:", err)
			return
		}

		wsClient := client.CreateWebsocketClient(c)
		defer s.realtime.RemoveClient(wsClient.ClientID())
		defer s.functions.UnregisterService(wsClient.ClientID())

		defer wsClient.Close()
		go wsClient.RoutineWrite()

		// Get wsClient details
		ctx := wsClient.Context()
		clientID := wsClient.ClientID()

		wsClient.Read(func(req *model.Message) {
			switch req.Type {
			case utils.TypeRealtimeSubscribe:
				// For realtime subscribe event
				data := new(model.RealtimeRequest)
				mapstructure.Decode(req.Data, data)

				// Subscribe to realttme feed
				feedData, err := s.realtime.Subscribe(ctx, clientID, s.auth, s.crud, data, func(feed *model.FeedData) {
					wsClient.Write(&model.Message{Type: utils.TypeRealtimeFeed, Data: feed})
				})
				if err != nil {
					res := model.RealtimeResponse{Group: data.Group, ID: data.ID, Ack: false, Error: err.Error()}
					wsClient.Write(&model.Message{ID: req.ID, Type: req.Type, Data: res})
					return
				}

				// Send response to wsClient
				res := model.RealtimeResponse{Group: data.Group, ID: data.ID, Ack: true, Docs: feedData}
				wsClient.Write(&model.Message{ID: req.ID, Type: req.Type, Data: res})

			case utils.TypeRealtimeUnsubscribe:
				// For realtime subscribe event
				data := new(model.RealtimeRequest)
				mapstructure.Decode(req.Data, data)

				s.realtime.Unsubscribe(clientID, data)

				// Send response to wsClient
				res := model.RealtimeResponse{Group: data.Group, ID: data.ID, Ack: true}
				wsClient.Write(&model.Message{ID: req.ID, Type: req.Type, Data: res})

			case utils.TypeServiceRegister:
				// TODO add security rule for functions registered as well
				data := new(model.ServiceRegisterRequest)
				mapstructure.Decode(req.Data, data)

				s.functions.RegisterService(clientID, data, func(payload *model.FunctionsPayload) {
					wsClient.Write(&model.Message{Type: utils.TypeServiceRequest, Data: payload})
				})

				wsClient.Write(&model.Message{ID: req.ID, Type: req.Type, Data: map[string]interface{}{"ack": true}})

			case utils.TypeServiceRequest:
				data := new(model.FunctionsPayload)
				mapstructure.Decode(req.Data, data)

				s.functions.HandleServiceResponse(data)
			}
		})
	}
}
