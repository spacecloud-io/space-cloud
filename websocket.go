package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/modules/auth"
	"github.com/spaceuptech/space-cloud/modules/realtime"
	"github.com/spaceuptech/space-cloud/utils"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleWebsocket(realtime *realtime.Module, auth *auth.Module) http.HandlerFunc {
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
				data := new(model.RealtimeRequest)
				mapstructure.Decode(req.Data, data)

				// Check if the user is authicated
				authObj, err := auth.IsAuthenticated(data.Token, data.DBType, data.Group, utils.Read)
				if err != nil {
					client.Write(model.Message{
						ID:   req.ID,
						Type: req.Type,
						Data: model.RealtimeResponse{Group: data.Group, ID: data.ID, Ack: false, Error: err.Error()},
					})
					return
				}

				// Create an args object
				args := map[string]interface{}{
					"args":    map[string]interface{}{"find": data.Where, "op": utils.All, "auth": authObj},
					"project": data.Project, // Don't forget to do this for every request
				}

				// Check if user is authorized to make this request
				err = auth.IsAuthorized(data.DBType, data.Group, utils.Read, args)
				if err != nil {
					client.Write(model.Message{
						ID:   req.ID,
						Type: req.Type,
						Data: model.RealtimeResponse{Group: data.Group, ID: data.ID, Ack: false, Error: err.Error()},
					})
					return
				}

				realtime.AddLiveQuery(data.ID, data.Group, client, data.Where)
				client.Write(model.Message{
					ID:   req.ID,
					Type: req.Type,
					Data: model.RealtimeResponse{Group: data.Group, ID: data.ID, Ack: true},
				})

			case utils.TypeRealtimeUnsubscribe:
				// For realtime unsubscribe event
				data := new(model.RealtimeRequest)
				mapstructure.Decode(req.Data, data)

				realtime.RemoveLiveQuery(data.Group, client.ClientID(), data.ID)
				client.Write(model.Message{
					ID:   req.ID,
					Type: req.Type,
					Data: model.RealtimeResponse{Group: data.Group, ID: data.ID, Ack: true},
				})
			}
		})
	}
}
