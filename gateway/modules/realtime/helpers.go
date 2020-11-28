package realtime

import (
	"fmt"

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

func isDBEnabled(dbAlias string, dbConfigs config.DatabaseConfigs) bool {
	for _, dbConfig := range dbConfigs {
		if dbConfig.DbAlias == dbAlias {
			return dbConfig.Enabled
		}
	}
	return false
}

func isRealTimeEnabled(dbAlias, table string, dbRules config.DatabaseRules) bool {
	for _, dbRule := range dbRules {
		if dbRule.DbAlias == dbAlias && dbRule.Table == table {
			return dbRule.IsRealTimeEnabled
		}
	}
	return false
}

func (m *Module) prepareFindObject(db, col string, row map[string]interface{}) map[string]interface{} {
	// Find the primary keys for the table
	primaryKeys := make([]string, 0)
	fields, p := m.schema.GetSchema(db, col)
	if p {
		for fieldName, value := range fields {
			if value.IsPrimary {
				primaryKeys = append(primaryKeys, fieldName)
			}
		}
	}

	// Extract primary keys from source and put it in find
	find := map[string]interface{}{}
	for _, key := range primaryKeys {
		if v, p := row[key]; p {
			find[key] = v
		}
	}

	return find
}

func generateEventRules(dbConfigs config.DatabaseConfigs, dbRules config.DatabaseRules, dbSchemas config.DatabaseSchemas, project, url string) []*config.EventingTrigger {

	var eventingRules []*config.EventingTrigger

	// Iterate over all dbTypes
	for _, dbSchema := range dbSchemas {

		// Proceed only if db is enabled
		if isDBEnabled(dbSchema.DbAlias, dbConfigs) {

			// Check if realtime mode is enabled
			if isRealTimeEnabled(dbSchema.DbAlias, dbSchema.Table, dbRules) {

				// Add a new event for each db event type
				for _, eventType := range dbEvents {
					rule := &config.EventingTrigger{
						Type:    eventType,
						URL:     url,
						Options: map[string]string{"db": dbSchema.DbAlias, "col": dbSchema.Table},
						Retries: 3,
						Timeout: 5000, // Timeout is in milliseconds
					}
					eventingRules = append(eventingRules, rule)
				}
			}
		}
	}

	return eventingRules
}

func createGroupKey(dbAlias, col string) string {
	return dbAlias + "::" + col
}

func getSendTopic(nodeID string) string {
	return fmt.Sprintf("realtime-%s", nodeID)
}
