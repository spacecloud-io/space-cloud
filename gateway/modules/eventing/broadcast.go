package eventing

import (
	"context"
	"fmt"
	"time"

	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/model"
)

// ProcessTransmittedEvents processes the event received
func (m *Module) ProcessTransmittedEvents(eventDocs []*model.EventDocument) {

	// Get the assigned token range
	start, end := m.syncMan.GetAssignedTokens()

	// Get current timestamp
	currentTimestamp := time.Now()

	for _, eventDoc := range eventDocs {
		if eventDoc.Token >= start && eventDoc.Token <= end {
			timestamp, err := time.Parse(time.RFC3339, eventDoc.Timestamp)
			if err != nil {
				_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Could not parse (%s) in event doc (%s) as time", eventDoc.Timestamp, eventDoc.ID), err, nil)
				continue
			}

			if currentTimestamp.After(timestamp) || currentTimestamp.Equal(timestamp) {
				go m.processStagedEvent(eventDoc)
			}
		}
	}
}
