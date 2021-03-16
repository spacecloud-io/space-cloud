package realtime

import (
	"github.com/segmentio/ksuid"

	"github.com/spaceuptech/space-api-go/transport/websocket"

	"github.com/spaceuptech/space-api-go/types"
)

// LiveQuery contains the methods for the liveQuery instance
type LiveQuery struct {
	appID   string
	db      string
	col     string
	client  *websocket.Socket
	store   types.DbStore
	options *types.LiveQueryOptions
	params  *types.RealtimeParams
}

// New returns a LiveQuery object
func New(appID, db, col string, client *websocket.Socket, store types.DbStore) *LiveQuery {
	return &LiveQuery{appID: appID, db: db, col: col, client: client, store: store, options: &types.LiveQueryOptions{}, params: &types.RealtimeParams{}}
}

// Where sets the where clause for the request
func (l *LiveQuery) Where(conds ...types.M) *LiveQuery {
	if len(conds) == 1 {
		l.params.Find = types.GenerateFind(conds[0])
	} else {
		l.params.Find = types.GenerateFind(types.And(conds...))
	}
	return l
}

// Options sets the live query options
func (l *LiveQuery) Options(options *types.LiveQueryOptions) *LiveQuery {
	l.options = &types.LiveQueryOptions{ChangesOnly: options.ChangesOnly, SkipInitial: options.ChangesOnly}
	return l
}

// Subscribe is used to subscribe to a new document
func (l *LiveQuery) Subscribe() *types.SubscriptionObject {

	id := ksuid.New().String()
	_, ok := l.store[l.db]
	if !ok {
		l.store[l.db] = types.ColStore{}
	}
	_, ok = l.store[l.db][l.col]
	if !ok {
		l.store[l.db][l.col] = types.IdStore{}
	}
	c := make(chan *types.SubscriptionEvent, 5)
	v := &types.Store{Snapshot: []*types.SnapshotData{}, C: c, Find: l.params.Find, Options: l.options}
	l.store[l.db][l.col][id] = v

	v.Unsubscribe = l.addSubscription(id, v.C)

	return l.subscribe(id, v)
}
