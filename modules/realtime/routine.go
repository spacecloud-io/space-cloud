package realtime

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"
)

func (m *Module) initWorkers(workerCount int) {
	for i := 0; i < workerCount; i++ {
		go m.worker()
	}
}

func (m *Module) worker() {
	if !m.enabled {
		return
	}

	for rawMsg := range m.feed {
		msg := new(Message)
		err := json.Unmarshal(rawMsg.Data, msg)
		if err != nil {
			log.Println("Realtime Worker Error:", err)
			continue
		}

		// Store the request if the msg type was intent
		if msg.Type == typeIntent {
			m.pendingRequests.Store(msg.ID, &pendingRequest{data: msg.Data, time: time.Now()})
			continue
		}

		// Delete the request if it failed (-ve ack)
		if !msg.Ack {
			m.pendingRequests.Delete(msg.ID)
			continue
		}

		tempReq, ok := m.pendingRequests.Load(msg.ID)
		if !ok {
			// Return since message was already flushed from queue
			continue
		}

		// Get the feedData
		req := tempReq.(*pendingRequest)
		m.helperSendFeed(req.data)
		m.pendingRequests.Delete(msg.ID)
	}
}

func (m *Module) removeStaleRequests() {
	ticker := time.NewTicker(2 * time.Minute)

	for range ticker.C {
		m.pendingRequests.Range(func(key interface{}, value interface{}) bool {
			req := value.(*pendingRequest)

			// Remove the request if its more than 30 seconds old
			diff := time.Now().Sub(req.time)
			if diff.Seconds() > 30 {
				idVar := "id"
				if req.data.DBType == string(utils.Mongo) {
					idVar = "_id"
				}
				switch req.data.Type {
				case utils.RealtimeInsert, utils.RealtimeUpdate:
					find := map[string]interface{}{idVar: req.data.DocID}
					data, err := m.crud.Read(context.TODO(), req.data.DBType, m.project, req.data.Group, &model.ReadRequest{Find: find, Operation: utils.One})
					if err == nil {
						// Send feed data if there is no error (doc found)
						req.data.Payload = data.(map[string]interface{})
						req.data.TimeStamp = time.Now().Unix()
						m.helperSendFeed(req.data)
					}

				case utils.RealtimeDelete:
					find := map[string]interface{}{idVar: req.data.DocID}
					_, err := m.crud.Read(context.TODO(), req.data.DBType, m.project, req.data.Group, &model.ReadRequest{Find: find, Operation: utils.One})
					if err != nil {
						// Send feed data if there is an error (no doc found)
						m.helperSendFeed(req.data)
					}
				}

				// Delete request
				m.pendingRequests.Delete(key)
			}

			return true
		})
	}
}

func (m *Module) helperSendFeed(data *model.FeedData) {
	clientsTemp, ok := m.groups.Load(data.Group)
	if !ok {
		return
	}

	clients := clientsTemp.(*clientsStub)
	clients.clients.Range(func(key interface{}, value interface{}) bool {
		queries := value.(*sync.Map)
		queries.Range(func(id interface{}, value interface{}) bool {
			query := value.(*queryStub)

			dataPoint := &model.FeedData{
				QueryID: id.(string), DocID: data.DocID, Group: data.Group, Payload: data.Payload,
				TimeStamp: data.TimeStamp, Type: data.Type, DBType: data.DBType,
			}

			switch data.Type {
			case utils.RealtimeDelete:
				query.sendFeed(dataPoint)

			case utils.RealtimeInsert:
				if utils.Validate(query.whereObj, data.Payload) {
					query.sendFeed(dataPoint)
				}

			case utils.RealtimeUpdate:
				idVar := "id"
				if data.DBType == string(utils.Mongo) {
					idVar = "_id"
				}

				// Fire a read request
				find := map[string]interface{}{idVar: data.DocID}
				d, err := m.crud.Read(context.TODO(), data.DBType, m.project, data.Group, &model.ReadRequest{Find: find, Operation: utils.One})
				if err == nil {
					// Send feed data if there is no error (doc found)
					dataPoint.Payload = d.(map[string]interface{})
					dataPoint.TimeStamp = time.Now().Unix()
					query.sendFeed(dataPoint)
				} else {
					// Send delete feed data on error
					dataPoint.Type = utils.RealtimeDelete
					dataPoint.TimeStamp = time.Now().Unix()
					query.sendFeed(dataPoint)
				}
			}
			return true
		})
		return true
	})
}
