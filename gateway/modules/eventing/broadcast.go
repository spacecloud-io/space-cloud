package eventing

import (
	"time"

	"github.com/spaceuptech/space-cloud/gateway/model"
)

// ProcessTransmittedEvents processes the event received
func (m *Module) ProcessTransmittedEvents(eventDocs []*model.EventDocument) {

	// Get the assigned token range
	start, end := m.syncMan.GetAssignedTokens()

	// Get current timestamp
	t := time.Now()
	currentTimestamp := t.UTC().UnixNano() / int64(time.Millisecond)

	for _, eventDoc := range eventDocs {
		if eventDoc.Token >= start && eventDoc.Token <= end {
			timestamp := eventDoc.Timestamp

			if currentTimestamp >= timestamp {
				go m.processStagedEvent(eventDoc)
			}
		}
	}
}
