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

		// Create an empty project variable
		var project string

		// Create a new client
		c := client.CreateWebsocketClient(socket)

		defer func() {
			// Unregister service if project could be loaded
			state, err := s.projects.LoadProject(project)
			if err == nil {
				// Unregister the service
				state.Realtime.RemoveClient(c.ClientID())
				state.Functions.UnregisterService(c.ClientID())
			}
		}()

		// Close the client to free up resources
		defer c.Close()

		// Start the writer routine
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

				// Set the clients project
				project = data.Project

				// Load the project state
				state, err := s.projects.LoadProject(project)
				if err != nil {
					res := model.RealtimeResponse{Group: data.Group, ID: data.ID, Ack: false, Error: err.Error()}
					c.Write(&model.Message{ID: req.ID, Type: utils.TypeRealtimeSubscribe, Data: res})
					return true
				}
				// Subscribe to the realtime feed
				feedData, err := state.Realtime.Subscribe(ctx, clientID, state.Auth, state.Crud, data, func(feed *model.FeedData) {
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

				// Load the project state
				state, err := s.projects.LoadProject(project)
				if err != nil {
					res := model.RealtimeResponse{Group: data.Group, ID: data.ID, Ack: false, Error: err.Error()}
					c.Write(&model.Message{ID: req.ID, Type: req.Type, Data: res})
					return true
				}

				state.Realtime.Unsubscribe(clientID, data)

				// Send response to c
				res := model.RealtimeResponse{Group: data.Group, ID: data.ID, Ack: true}
				c.Write(&model.Message{ID: req.ID, Type: req.Type, Data: res})

			case utils.TypeServiceRegister:
				// TODO add security rule for functions registered as well
				data := new(model.ServiceRegisterRequest)
				mapstructure.Decode(req.Data, data)

				// Set the clients project
				project = data.Project

				state, err := s.projects.LoadProject(project)
				if err != nil {
					c.Write(&model.Message{ID: req.ID, Type: req.Type, Data: map[string]interface{}{"ack": false}})
					return true
				}
				state.Functions.RegisterService(clientID, data, func(payload *model.FunctionsPayload) {
					c.Write(&model.Message{Type: utils.TypeServiceRequest, Data: payload})
				})

				c.Write(&model.Message{ID: req.ID, Type: req.Type, Data: map[string]interface{}{"ack": true}})

			case utils.TypeServiceRequest:
				data := new(model.FunctionsPayload)
				mapstructure.Decode(req.Data, data)

				// Handle response if project could be loaded
				state, err := s.projects.LoadProject(project)
				if err == nil {
					state.Functions.HandleServiceResponse(data)
				}
			}
			return true
		})
	}
}
