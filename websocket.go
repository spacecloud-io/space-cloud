package main

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

func (s *server) handleWebsocket() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("upgrade:", err)
			return
		}

		client := client.CreateWebsocketClient(c)
		defer s.realtime.RemoveClient(client.ClientID())
		defer s.functions.UnregisterService(client.ClientID())

		defer client.Close()
		go client.RoutineWrite()

		client.Read(func(req *model.Message) {
			switch req.Type {
			case utils.TypeRealtimeSubscribe:
				// For realtime subscribe event
				data := new(model.RealtimeRequest)
				mapstructure.Decode(req.Data, data)

				s.realtime.Subscribe(req.ID, client, s.auth, s.crud, data)

			case utils.TypeRealtimeUnsubscribe:
				// For realtime subscribe event
				data := new(model.RealtimeRequest)
				mapstructure.Decode(req.Data, data)

				s.realtime.Unsubscribe(req.ID, client, data)

			case utils.TypeServiceRegister:
				// TODO add security rule for functions registered as well
				data := new(model.ServiceRegisterRequest)
				mapstructure.Decode(req.Data, data)

				s.functions.RegisterService(req.ID, client, data)

			case utils.TypeServiceRequest:
				data := new(model.FunctionsPayload)
				mapstructure.Decode(req.Data, data)

				s.functions.HandleServiceResponse(data)
			}

		})
	}
}
