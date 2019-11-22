package realtime

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/mitchellh/mapstructure"

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"
)

// Subscribe performs the realtime subscribe operation.
func (m *Module) Subscribe(ctx context.Context, clientID string, data *model.RealtimeRequest, sendFeed SendFeed) ([]*model.FeedData, error) {

	readReq := &model.ReadRequest{Find: data.Where, Operation: utils.All}

	// Check if the user is authorised to make the request
	_, err := m.auth.IsReadOpAuthorised(ctx, data.Project, data.DBType, data.Group, data.Token, readReq)
	if err != nil {
		return nil, err
	}

	if data.Options.SkipInitial {
		m.AddLiveQuery(data.ID, data.Project, data.DBType, data.Group, clientID, data.Where, sendFeed)
		return []*model.FeedData{}, nil
	}

	ctx2, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result, err := m.crud.Read(ctx2, data.DBType, data.Project, data.Group, readReq)
	if err != nil {
		return nil, err
	}

	feedData := make([]*model.FeedData, 0)
	array, ok := result.([]interface{})
	if ok {
		timeStamp := time.Now().Unix()
		for _, row := range array {
			payload := row.(map[string]interface{})
			idVar := "id"
			if data.DBType == string(utils.Mongo) {
				idVar = "_id"
			}
			if docID, ok := utils.AcceptableIDType(payload[idVar]); ok {
				feedData = append(feedData, &model.FeedData{
					Group:     data.Group,
					Type:      utils.RealtimeInitial,
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
	m.AddLiveQuery(data.ID, data.Project, data.DBType, data.Group, clientID, data.Where, sendFeed)
	return feedData, nil
}

// Unsubscribe performs the realtime unsubscribe operation.
func (m *Module) Unsubscribe(clientID string, data *model.RealtimeRequest) {
	m.RemoveLiveQuery(data.DBType, data.Group, clientID, data.ID)
}

// HandleRealtimeEvent handles an incoming realtime event from the eventing module
func (m *Module) HandleRealtimeEvent(ctxRoot context.Context, eventDoc *model.CloudEventPayload) error {

	urls := m.syncMan.GetSpaceCloudNodeURLs(m.project)

	// Create wait group
	var wg sync.WaitGroup
	wg.Add(len(urls))

	// Create success & error channels
	successCh := make(chan struct{}, 1)
	errCh := make(chan error, len(urls))

	ctx, cancel := context.WithTimeout(ctxRoot, 5*time.Second)
	defer cancel()

	for _, url := range urls {
		go func() {
			defer wg.Done()

			token, err := m.auth.GetInternalAccessToken()
			if err != nil {
				errCh <- err
				return
			}

			var res interface{}
			if err := m.syncMan.MakeHTTPRequest(ctx, "POST", url, token, eventDoc, &res); err != nil {
				errCh <- err
				return
			}
		}()
	}

	go func() {
		wg.Wait()
		successCh <- struct{}{}
	}()

	select {
	case err := <-errCh:
		cancel()
		log.Println("Realtime Module: Event handler error -", err)
		return err

	case <-successCh:
		return nil
	}
}

// ProcessRealtimeRequests handles an incoming realtime process event
func (m *Module) ProcessRealtimeRequests(eventDoc *model.CloudEventPayload) error {

	dbEvent := new(model.DatabaseEventMessage)
	if err := mapstructure.Decode(eventDoc.Data, dbEvent); err != nil {
		log.Println("Realtime Module Request Handler Error:", err)
		return err
	}

	t, _ := time.Parse(time.RFC3339, eventDoc.Time)

	feedData := &model.FeedData{
		DocID:     dbEvent.DocID,
		Type:      eventingToRealtimeEvent(eventDoc.Type),
		Payload:   dbEvent.Doc,
		TimeStamp: t.UnixNano() / int64(time.Millisecond),
		Group:     dbEvent.Col,
		DBType:    dbEvent.DBType,
	}

	m.helperSendFeed(feedData)

	return nil
}
