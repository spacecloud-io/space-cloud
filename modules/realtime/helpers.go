package realtime

import (
	"strings"

	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/utils"
)

var dbEvents = []string{utils.EventCreate, utils.EventUpdate, utils.EventDelete}

func eventingToRealtimeEvent(event string) string {
	switch event {
	case utils.EventCreate:
		return utils.RealtimeInsert

	case utils.EventUpdate:
		return utils.RealtimeUpdate

	case utils.EventDelete:
		return utils.RealtimeDelete

	default:
		return event
	}
}

func generateEventRules(crudConfig config.Crud) []config.EventingRule {

	var eventingRules []config.EventingRule

	// Iterate over all dbTypes
	for dbType, dbStub := range crudConfig {

		// Proceed only if db is enabled
		if dbStub.Enabled {

			// Iterate over all connections
			for col, colStub := range dbStub.Collections {

				// Check if realtime mode is enabled
				if colStub.IsRealTimeEnabled {

					// Add a new event for each db event type
					for _, eventType := range dbEvents {
						rule := config.EventingRule{
							Type:     eventType,
							Service:  serviceName,
							Function: funcName,
							Options:  map[string]string{"db": dbType, "col": col},
						}
						eventingRules = append(eventingRules, rule)
					}
				}
			}
		}
	}

	return eventingRules
}

func createGroupKey(dbType, col string) string {
	return dbType + "::" + col
}

func getSubjectName(project, dbType, col string) string {
	return "realtime:" + project + ":" + dbType + ":" + col
}

func getDBTypeAndColFromGroupKey(key string) (dbType string, col string) {
	array := strings.Split(key, "::")
	return array[0], array[1]
}
