package eventing

import (
	"encoding/json"

	"github.com/spaceuptech/space-cloud/model"
)

func (m *Module) broadcastEvents(eventDocs []*model.EventDocument) {
	data, err := json.Marshal(eventDocs)
	if err == nil {
		m.nc.Publish(internalEventingSubject, data)
	}
}

func (m *Module) processBroadcastEvents(eventDocs []*model.EventDocument) {

	// Get the assigned token range
	start, end := m.syncMan.GetAssignedTokens()

	for _, eventDoc := range eventDocs {
		if eventDoc.Token >= start && eventDoc.Token <= end {
			go m.processStagedEvent(eventDoc)
		}
	}
}
