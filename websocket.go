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

		// Create an empty project variable
		var project string

		// Create a new client
		client := client.CreateWebsocketClient(c)

		defer func() {
			// Unregister service if project could be loaded
			state, err := s.projects.LoadProject(project)
			if err == nil {
				// Unregister the service
				state.Realtime.RemoveClient(client.ClientID())
				state.Functions.UnregisterService(client.ClientID())
			}
		}()

		// Close the client to free up resources
		defer client.Close()

		// Start the writer routine
		go client.RoutineWrite()

		// Get client details
		ctx := client.Context()
		clientID := client.ClientID()

		client.Read(func(req *model.Message) {
			switch req.Type {
			case utils.TypeRealtimeSubscribe:

				// For realtime subscribe event
				data := new(model.RealtimeRequest)
				mapstructure.Decode(req.Data, data)

				// Set the clients project
				project = data.Project

				// Load the project state
				state, err := s.projects.LoadProject(project)
				if err != nil {
					res := model.RealtimeResponse{Group: data.Group, ID: data.ID, Ack: false, Error: err.Error()}
					client.Write(&model.Message{ID: req.ID, Type: utils.TypeRealtimeSubscribe, Data: res})
					return
				}
				// Subscribe to the realtime feed
				feedData, err := state.Realtime.Subscribe(ctx, clientID, state.Auth, state.Crud, data, func(feed *model.FeedData) {
					client.Write(&model.Message{Type: utils.TypeRealtimeFeed, Data: feed})
				})
				if err != nil {
					res := model.RealtimeResponse{Group: data.Group, ID: data.ID, Ack: false, Error: err.Error()}
					client.Write(&model.Message{ID: req.ID, Type: req.Type, Data: res})
					return
				}

				// Send response to client
				res := model.RealtimeResponse{Group: data.Group, ID: data.ID, Ack: true, Docs: feedData}
				client.Write(&model.Message{ID: req.ID, Type: req.Type, Data: res})

			case utils.TypeRealtimeUnsubscribe:
				// For realtime subscribe event
				data := new(model.RealtimeRequest)
				mapstructure.Decode(req.Data, data)

				// Load the project state
				state, err := s.projects.LoadProject(project)
				if err != nil {
					res := model.RealtimeResponse{Group: data.Group, ID: data.ID, Ack: false, Error: err.Error()}
					client.Write(&model.Message{ID: req.ID, Type: req.Type, Data: res})
					return
				}

				state.Realtime.Unsubscribe(clientID, data)

				// Send response to client
				res := model.RealtimeResponse{Group: data.Group, ID: data.ID, Ack: true}
				client.Write(&model.Message{ID: req.ID, Type: req.Type, Data: res})

			case utils.TypeServiceRegister:
				// TODO add security rule for functions registered as well
				data := new(model.ServiceRegisterRequest)
				mapstructure.Decode(req.Data, data)

				// Set the clients project
				project = data.Project

				state, err := s.projects.LoadProject(project)
				if err != nil {
					client.Write(&model.Message{ID: req.ID, Type: req.Type, Data: map[string]interface{}{"ack": false}})
					return
				}
				state.Functions.RegisterService(clientID, data, func(payload *model.FunctionsPayload) {
					client.Write(&model.Message{Type: utils.TypeServiceRequest, Data: payload})
				})

				client.Write(&model.Message{ID: req.ID, Type: req.Type, Data: map[string]interface{}{"ack": true}})

			case utils.TypeServiceRequest:
				data := new(model.FunctionsPayload)
				mapstructure.Decode(req.Data, data)

				// Handle response if project could be loaded
				state, err := s.projects.LoadProject(project)
				if err == nil {
					state.Functions.HandleServiceResponse(data)
				}
			}
		})
	}
}
