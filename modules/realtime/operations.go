package realtime

import (
	"time"

	"github.com/mitchellh/mapstructure"

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/modules/auth"
	"github.com/spaceuptech/space-cloud/modules/crud"
	"github.com/spaceuptech/space-cloud/utils"
)

// Operation performs the RealTime operations common between Websocket and GRPC.
func (m *Module) Operation(client *utils.Client, auth *auth.Module, crud *crud.Module) {

	defer client.Close()
	defer m.RemoveClient(client.ClientID())

	client.RoutineWrite()

	client.Read(func(req *model.Message) {
		switch req.Type {
		case utils.TypeRealtimeSubscribe:
			// For realtime subscribe event
			data := new(model.RealtimeRequest)
			mapstructure.Decode(req.Data, data)

			// Check if the user is authenticated
			authObj, err := auth.IsAuthenticated(data.Token, data.DBType, data.Group, utils.Read)
			if err != nil {
				client.Write(&model.Message{
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
				client.Write(&model.Message{
					ID:   req.ID,
					Type: req.Type,
					Data: model.RealtimeResponse{Group: data.Group, ID: data.ID, Ack: false, Error: err.Error()},
				})
				return
			}

			readReq := model.ReadRequest{Find: data.Where, Operation: utils.All}
			result, err := crud.Read(client.Context, data.DBType, data.Project, data.Group, &readReq)
			if err != nil {
				client.Write(&model.Message{
					ID:   req.ID,
					Type: req.Type,
					Data: model.RealtimeResponse{Group: data.Group, ID: data.ID, Ack: false, Error: err.Error()},
				})
				return
			}

			feedData := []*model.FeedData{}
			array, ok := result.([]interface{})
			if ok {
				timeStamp := time.Now().Unix()
				for _, row := range array {
					payload := row.(map[string]interface{})
					idVar := "id"
					if data.DBType == string(utils.Mongo) {
						idVar = "_id"
					}
					if docID, ok := payload[idVar].(string); ok {
						feedData = append(feedData, &model.FeedData{
							Group:     data.Group,
							Type:      utils.RealtimeWrite,
							TimeStamp: timeStamp,
							DocID:     docID,
							DBType:    data.DBType,
							Payload:   payload,
							QueryID:   data.ID,
						})
					}
				}
			}
			// Add the live query
			m.AddLiveQuery(data.ID, data.Group, client, data.Where)
			client.Write(&model.Message{
				ID:   req.ID,
				Type: req.Type,
				Data: model.RealtimeResponse{Group: data.Group, ID: data.ID, Ack: true, Docs: feedData},
			})

		case utils.TypeRealtimeUnsubscribe:
			// For realtime unsubscribe event
			data := new(model.RealtimeRequest)
			mapstructure.Decode(req.Data, data)

			m.RemoveLiveQuery(data.Group, client.ClientID(), data.ID)
			client.Write(&model.Message{
				ID:   req.ID,
				Type: req.Type,
				Data: model.RealtimeResponse{Group: data.Group, ID: data.ID, Ack: true},
			})
		}
	})
}
