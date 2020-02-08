package realtime

import (
	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

var dbEvents = []string{utils.EventDBCreate, utils.EventDBUpdate, utils.EventDBDelete}

func eventingToRealtimeEvent(event string) string {
	switch event {
	case utils.EventDBCreate:
		return utils.RealtimeInsert

	case utils.EventDBUpdate:
		return utils.RealtimeUpdate

	case utils.EventDBDelete:
		return utils.RealtimeDelete

	default:
		return event
	}
}

func generateEventRules(crudConfig config.Crud, project, url string) []config.EventingRule {

	var eventingRules []config.EventingRule

	// Iterate over all dbTypes
	for dbAlias, dbStub := range crudConfig {

		// Proceed only if db is enabled
		if dbStub.Enabled {

			// Iterate over all connections
			for col, colStub := range dbStub.Collections {

				// Check if realtime mode is enabled
				if colStub.IsRealTimeEnabled {

					// Add a new event for each db event type
					for _, eventType := range dbEvents {
						rule := config.EventingRule{
							Type:    eventType,
							URL:     url,
							Options: map[string]string{"db": dbAlias, "col": col},
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

// func getSubjectName(project, dbType, col string) string {
// 	return "realtime:" + project + ":" + dbType + ":" + col
// }

// func getDBTypeAndColFromGroupKey(key string) (dbType string, col string) {
// 	array := strings.Split(key, "::")
// 	return array[0], array[1]
// }
