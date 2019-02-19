package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/modules/realtime"
	"github.com/spaceuptech/space-cloud/utils"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleWebsocket(realtime *realtime.Module) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("upgrade:", err)
			return
		}
		client := utils.CreateWebsocketClient(c)
		defer client.Close()

		client.RoutineWrite()

		client.Read(func(req *model.Message) {

			switch req.Type {
			case utils.TypeRealtimeSubscribe:
				// For realtime subscribe event
				data := req.Data.(model.RealtimeRequest)
				queryID := realtime.AddLiveQuery(data.Group, client, data.Where)
				client.Write(model.Message{
					Type: req.Type,
					Data: model.RealtimeResponse{Group: data.Group, ID: queryID},
				})

			case utils.TypeRealtimeUnsubscribe:
				// For realtime unsubscribe event
				data := req.Data.(model.RealtimeRequest)
				realtime.RemoveLiveQuery(data.Group, client.ClientID(), data.ID)
				client.Write(model.Message{
					Type: req.Type,
					Data: model.RealtimeResponse{Group: data.Group, ID: data.ID},
				})
			}
		})
	}
}
