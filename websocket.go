package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"
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

		client := utils.CreateWebsocketClient(c)
		defer s.realtime.RemoveClient(client.ClientID())
		defer s.functions.UnregisterService(client.ClientID())

		defer client.Close()
		go client.RoutineWrite()

		client.Read(func(req *model.Message) {
			switch req.Type {
			case utils.TypeRealtimeSubscribe:
				s.realtime.Subscribe(client, s.auth, s.crud, req)

			case utils.TypeRealtimeUnsubscribe:
				s.realtime.Unsubscribe(client, req)

			case utils.TypeServiceRegister:
				// TODO add security rule for functions registered as well
				data := new(model.ServiceRegisterRequest)
				mapstructure.Decode(req.Data, data)

				s.functions.RegisterService(client, data)

			case utils.TypeServiceRequest:
				data := new(model.FunctionsPayload)
				mapstructure.Decode(req.Data, data)

				s.functions.HandleServiceResponse(req.ID, data)
			}

		})
	}
}
