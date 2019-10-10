package eventing

import (
	"github.com/spaceuptech/space-cloud/model"
)

// ProcessTransmittedEvents processes the event received
func (m *Module) ProcessTransmittedEvents(eventDocs []*model.EventDocument) {

	// Get the assigned token range
	start, end := m.syncMan.GetAssignedTokens()

	for _, eventDoc := range eventDocs {
		if eventDoc.Token >= start && eventDoc.Token <= end {
			go m.processStagedEvent(eventDoc)
		}
	}
}
