package common

import (
	"encoding/json"

	"github.com/spacecloud-io/space-cloud/config"
	"github.com/spacecloud-io/space-cloud/modules/database"
)

func prepareDatabaseApp(scConfig *config.Config) []byte {
	dbConfigs := map[string]*database.Config{}
	for projectID, project := range scConfig.Projects {
		for _, dbConfig := range project.DatabaseConfigs {
			// Pick the db schemas and prepared queries specified for this database
			dbSchemas := make(config.DatabaseSchemas)
			dbPreparedQueries := make(config.DatabasePreparedQueries)

			// We are only interested in the schemas that belong to this database
			for _, schema := range project.DatabaseSchemas {
				if schema.DbAlias == dbConfig.DbAlias {
					dbSchemas[schema.Table] = schema
				}
			}

			// We are only interested in the prepared queries that belong to this database
			for _, query := range project.DatabasePreparedQueries {
				if query.DbAlias == dbConfig.DbAlias {
					dbPreparedQueries[query.ID] = query
				}
			}

			// We prefix the alias with project id to make sure we face no conflicts when too projects have the same dbAlias
			dbConfigs[database.CombineDBConfigKey(projectID, dbConfig.DbAlias)] = &database.Config{
				Connector:       dbConfig,
				Schemas:         dbSchemas,
				PreparedQueries: dbPreparedQueries,
			}
		}
	}

	data, _ := json.Marshal(map[string]interface{}{"dbConfigs": dbConfigs})
	return data
}
