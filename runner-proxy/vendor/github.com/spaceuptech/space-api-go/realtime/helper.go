package realtime

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

	"github.com/spaceuptech/space-api-go/types"
)

// used internally
func (l *LiveQuery) addSubscription(id string, c chan *types.SubscriptionEvent) func() {
	return func() {
		_, _ = l.client.Request(types.TypeRealtimeUnsubscribe, types.RealtimeRequest{Group: l.col, ID: id, Options: l.options})
		delete(l.store[l.db][l.col], id)
		close(c)
	}
}

func (l *LiveQuery) subscribe(id string, store *types.Store) *types.SubscriptionObject {
	req := types.RealtimeRequest{DBType: l.db, Project: l.appID, Group: l.col, ID: id, Where: l.params.Find, Options: l.options}

	data, err := l.client.Request(types.TypeRealtimeSubscribe, req)
	if err != nil {
		store.C <- types.NewSubscriptionEvent("", nil, nil, fmt.Errorf("error unable to subscribe to realtime feature %v", err))
		store.Unsubscribe()
	}

	go func() {
		store.Snapshot = make([]*types.SnapshotData, 0)
		mapData, ok := data.(map[string]interface{})
		if ok {
			ack, ok := mapData["ack"]
			if ok && !ack.(bool) {
				err, ok := mapData["error"].(string)
				if ok {
					store.C <- types.NewSubscriptionEvent("", nil, nil, fmt.Errorf("error from server %v", err))
					store.Unsubscribe()
				}
			}
			docs := mapData["docs"].([]interface{})
			for _, doc := range docs {
				var feed feedData
				if err := mapstructure.Decode(doc, &feed); err != nil {
					continue
				}
				store.Snapshot = append(store.Snapshot, &types.SnapshotData{Find: feed.Find, Time: feed.TimeStamp, Payload: feed.Payload, IsDeleted: false})
				store.C <- types.NewSubscriptionEvent("initial", feed.Payload, l.params.Find, nil)
			}
		}
	}()

	return types.LiveQuerySubscriptionInit(store)
}

func snapshotCallback(store types.DbStore, rows []feedData) {
	if len(rows) == 0 {
		return
	}
	var obj = new(types.Store)
	var opts = types.LiveQueryOptions{}
	for _, data := range rows {
		obj = store[data.DBType][data.Group][data.QueryID]
		opts = obj.QueryOptions

		if opts.ChangesOnly {
			if !(opts.SkipInitial && data.Type == types.RealtimeInitial) {
				if data.Type != types.RealtimeDelete {
					obj.C <- types.NewSubscriptionEvent(data.Type, data.Payload, data.Find, nil)
				} else {
					obj.C <- types.NewSubscriptionEvent(data.Type, nil, data.Find, nil)
				}
			}
		} else {
			if data.Type == types.RealtimeInitial {
				obj.Snapshot = append(obj.Snapshot, &types.SnapshotData{Find: data.Find, Time: data.TimeStamp, Payload: data.Payload, IsDeleted: false})
				obj.C <- types.NewSubscriptionEvent(data.Type, data.Payload, data.Find, nil)
			} else if data.Type == types.RealtimeInsert || data.Type == types.RealtimeUpdate {
				isExisting := false
				for _, row := range obj.Snapshot {
					if validate(data.Find, row.Payload.(map[string]interface{})) {
						isExisting = true
						if row.Time <= data.TimeStamp {
							row.Time = data.TimeStamp
							row.Payload = data.Payload
							row.IsDeleted = false
							obj.C <- types.NewSubscriptionEvent(data.Type, data.Payload, data.Find, nil)
						}
					}
				}
				if !isExisting {
					obj.Snapshot = append(obj.Snapshot, &types.SnapshotData{Find: data.Find, Time: data.TimeStamp, Payload: data.Payload, IsDeleted: false})
					obj.C <- types.NewSubscriptionEvent(data.Type, data.Payload, data.Find, nil)
				}
			} else if data.Type == types.RealtimeDelete {
				for _, row := range obj.Snapshot {
					if validate(data.Find, row.Payload.(map[string]interface{})) {
						if row.Time <= data.TimeStamp {
							row.Time = data.TimeStamp
							row.Payload = map[string]interface{}{}
							row.IsDeleted = true
							obj.C <- types.NewSubscriptionEvent(data.Type, nil, data.Find, nil)
						}
					}
				}
			}
		}
	}
}

func validate(find map[string]interface{}, doc map[string]interface{}) bool {
	for k, v := range find {
		keyValue, p := doc[k]
		if !p {
			return false
		}

		if keyValue != v {
			return false
		}
	}
	return true
}
