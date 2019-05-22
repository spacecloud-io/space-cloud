package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/modules/auth"
	"github.com/spaceuptech/space-cloud/modules/crud"
	"github.com/spaceuptech/space-cloud/modules/realtime"
	"github.com/spaceuptech/space-cloud/utils"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleWebsocket(realtime *realtime.Module, auth *auth.Module, crud *crud.Module) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("upgrade:", err)
			return
		}

		client := utils.CreateWebsocketClient(c)
		defer realtime.RemoveClient(client.ClientID())
		defer client.Close()
		client.RoutineWrite()

		client.Read(func(req *model.Message) {
			switch req.Type {
			case utils.TypeRealtimeSubscribe:
				realtime.Subscribe(client, auth, crud, req)

			case utils.TypeRealtimeUnsubscribe:
				realtime.Unsubscribe(client, req)
			}
		})
	}
}
