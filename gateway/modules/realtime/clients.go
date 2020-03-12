package realtime

import (
	"sync"

	"github.com/spaceuptech/space-cloud/gateway/model"
)

type queryStub struct {
	sendFeed SendFeed
	whereObj map[string]interface{}
	actions  *model.PostProcess
}

type clientsStub struct {
	sync.Mutex
	clients sync.Map
}

// AddLiveQuery tracks a client for a live query
func (m *Module) AddLiveQuery(id, _, dbAlias, group, clientID string, whereObj map[string]interface{}, actions *model.PostProcess, sendFeed SendFeed) {
	// Load clients in a particular group
	clients := new(clientsStub)
	t, _ := m.groups.LoadOrStore(createGroupKey(dbAlias, group), clients)
	clients = t.(*clientsStub)

	// Load the queries of a particular client
	queries := new(sync.Map)
	t, _ = clients.clients.LoadOrStore(clientID, queries)
	queries = t.(*sync.Map)

	// Add the query
	queries.Store(id, &queryStub{sendFeed, whereObj, actions})
}

// RemoveLiveQuery removes a particular live query
func (m *Module) RemoveLiveQuery(dbAlias, group, clientID, queryID string) {
	// Load clients in a particular group
	clientsTemp, ok := m.groups.Load(createGroupKey(dbAlias, group))
	if !ok {
		return
	}
	clients := clientsTemp.(*clientsStub)

	// Load the queries of a particular client
	queriesTemp, ok := clients.clients.Load(clientID)
	if !ok {
		return
	}
	queries := queriesTemp.(*sync.Map)

	// Remove the query
	queries.Delete(queryID)

	// Delete client if it has no queries
	if mapLen(queries) == 0 {
		clients.clients.Delete(clientID)
	}

	// Delete group if no clients present
	if mapLen(&clients.clients) == 0 {
		m.groups.Delete(createGroupKey(dbAlias, group))
	}
}

// RemoveClient removes a client
func (m *Module) RemoveClient(clientID string) {
	// Delete the client from all groups
	m.groups.Range(func(key interface{}, value interface{}) bool {
		clients := value.(*clientsStub)
		clients.clients.Delete(clientID)
		if mapLen(&clients.clients) == 0 {
			m.groups.Delete(key)
		}
		return true
	})
}

func mapLen(m *sync.Map) int {
	counter := 0
	m.Range(func(k, v interface{}) bool {
		counter++
		return true
	})
	return counter
}
