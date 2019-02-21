package realtime

import (
	"sync"

	"github.com/satori/go.uuid"

	"github.com/spaceuptech/space-cloud/utils"
)

type queryStub struct {
	client   *utils.Client
	whereObj map[string]interface{}
}

// AddLiveQuery tracks a client for a live query
func (m *Module) AddLiveQuery(group string, client *utils.Client, whereObj map[string]interface{}) string {
	// Load clients in a particular group
	clients := new(sync.Map)
	t, _ := m.groups.LoadOrStore(group, clients)
	clients = t.(*sync.Map)

	// Load the queries of a particular client
	queries := new(sync.Map)
	t, _ = clients.LoadOrStore(client.ClientID(), queries)
	queries = t.(*sync.Map)

	// Add the query
	id := uuid.NewV1().String()
	queries.Store(id, &queryStub{client, whereObj})
	return id
}

// RemoveLiveQuery removes a particular live query
func (m *Module) RemoveLiveQuery(group, clientID, queryID string) {
	// Load clients in a particular group
	clientsTemp, ok := m.groups.Load(group)
	if !ok {
		return
	}
	clients := clientsTemp.(*sync.Map)

	// Load the queries of a particular client
	queriesTemp, ok := clients.Load(clientID)
	if !ok {
		return
	}
	queries := queriesTemp.(*sync.Map)

	// Remove the query
	queries.Delete(queryID)
}

// RemoveClient removes a client
func (m *Module) RemoveClient(group, clientID string) {
	// Load clients in a particular group
	clientsTemp, ok := m.groups.Load(group)
	if !ok {
		return
	}
	clients := clientsTemp.(*sync.Map)

	// Remove the client
	clients.Delete(group)
}
