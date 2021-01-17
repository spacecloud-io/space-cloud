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

	count := 0
	fmt.Println("Length of process transmitted events", len(eventDocs))
	for _, eventDoc := range eventDocs {
		if eventDoc.Token >= start && eventDoc.Token <= end {
			timestamp, err := time.Parse(time.RFC3339Nano, eventDoc.Timestamp)
			if err != nil {
				_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Could not parse (%s) in event doc (%s) as time", eventDoc.Timestamp, eventDoc.ID), err, nil)
				continue
			}

			if currentTimestamp.After(timestamp) || currentTimestamp.Equal(timestamp) {
				count++
				fmt.Println("Staging event count", count)
				m.stageBufferedEvent(eventDoc)
			}
		}
	}
}
