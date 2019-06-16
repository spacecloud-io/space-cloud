package realtime

import (
	"context"
	"time"

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/modules/auth"
	"github.com/spaceuptech/space-cloud/modules/crud"
	"github.com/spaceuptech/space-cloud/utils"
)

// Subscribe performs the realtime subscribe operation.
func (m *Module) Subscribe(ctx context.Context, clientID string, auth *auth.Module, crud *crud.Module, data *model.RealtimeRequest, sendFeed SendFeed) ([]*model.FeedData, error) {

	// Check if the user is authenticated
	authObj, err := auth.IsAuthenticated(data.Token, data.DBType, data.Group, utils.Read)
	if err != nil {
		return nil, err
	}

	// Create an args object
	args := map[string]interface{}{
		"args":    map[string]interface{}{"find": data.Where, "op": utils.All, "auth": authObj},
		"project": data.Project, // Don't forget to do this for every request
	}

	// Check if user is authorized to make this request
	err = auth.IsAuthorized(data.Project, data.DBType, data.Group, utils.Read, args)
	if err != nil {
		return nil, err
	}

	return m.DoRealtimeSubscribe(ctx, clientID, crud, data, sendFeed)
}

// DoRealtimeSubscribe makes the realtime query
func (m *Module) DoRealtimeSubscribe(ctx context.Context, clientID string, crud *crud.Module, data *model.RealtimeRequest, sendFeed SendFeed) ([]*model.FeedData, error) {
	readReq := model.ReadRequest{Find: data.Where, Operation: utils.All}
	result, err := crud.Read(ctx, data.DBType, data.Project, data.Group, &readReq)
	if err != nil {
		return nil, err
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
	m.AddLiveQuery(data.ID, data.Group, clientID, data.Where, sendFeed)
	return feedData, nil
}

// Unsubscribe performs the realtime unsubscribe operation.
func (m *Module) Unsubscribe(clientID string, data *model.RealtimeRequest) {
	m.RemoveLiveQuery(data.Group, clientID, data.ID)
}
