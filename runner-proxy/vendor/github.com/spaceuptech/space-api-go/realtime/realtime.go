package realtime

import (
	"github.com/mitchellh/mapstructure"

	"github.com/spaceuptech/space-api-go/transport/websocket"
	"github.com/spaceuptech/space-api-go/types"
)

// Realtime handles all live query related tasks
type Realtime struct {
	appID  string
	store  types.DbStore
	client *websocket.Socket
}

// Init initialize the realtime module
func Init(appID string, client *websocket.Socket) *Realtime {
	r := &Realtime{appID: appID, client: client, store: make(types.DbStore)}

	// on reconnect register again according the value in store
	r.client.RegisterOnReconnectCallback(func() {
		for db, dbValue := range r.store {
			for col, colValue := range dbValue {
				for id := range colValue {
					obj := r.store[db][col][id]
					q := r.LiveQuery(db, col)
					q.options = obj.Options.(*types.LiveQueryOptions)
					q.params = &types.RealtimeParams{Find: obj.Find.(types.M)}
					q.subscribe(id, obj)
				}
			}
		}
	})

	// initialize the realtime on sc
	r.client.RegisterCallback(types.TypeRealtimeFeed, func(data interface{}) {
		var feed feedData
		if err := mapstructure.Decode(data, &feed); err != nil {
			return
		}
		snapshotCallback(r.store, []feedData{feed})
	})
	return r
}

// LiveQuery initialize the live query module
func (r *Realtime) LiveQuery(db, collection string) *LiveQuery {
	return New(r.appID, db, collection, r.client, r.store)
}
