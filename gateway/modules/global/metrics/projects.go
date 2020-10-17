package metrics

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/spaceuptech/helpers"
	"github.com/spaceuptech/space-api-go/types"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

func (m *Module) generateMetricsRequest(project *config.Project, ssl *config.SSL) (string, map[string]interface{}, map[string]interface{}) {
	m.lock.Lock()
	defer m.lock.Unlock()

	clusterID := m.adminMan.GetClusterID()
	projectID := project.ProjectConfig.ID
	min := map[string]interface{}{"start_time": time.Now().UnixNano() / int64(time.Millisecond)}
	// Create the find and update clauses
	set := map[string]interface{}{
		"nodes":        m.syncMan.GetNodesInCluster(),
		"os":           runtime.GOOS,
		"is_prod":      m.isProd,
		"version":      utils.BuildVersion,
		"distribution": "ce",
		"last_updated": time.Now().UnixNano() / int64(time.Millisecond),
		"project":      projectID,
		"cluster_id":   clusterID,
	}

	set["ssl_enabled"] = ssl != nil && ssl.Enabled

	// modules := project.Modules
	// crud info
	set["crud"] = map[string]interface{}{"tables": map[string]interface{}{}}
	set["databases"] = map[string][]string{"databases": {}}
	if project.DatabaseConfigs != nil {
		temps := map[string]interface{}{}
		dbs := []string{}
		for _, dbConfig := range project.DatabaseConfigs {
			dbs = append(dbs, dbConfig.Type)
			temps[dbConfig.DbAlias] = map[string]interface{}{"tables": getTablesCount(dbConfig.DbAlias, project.DatabaseSchemas, project.EventingConfig)}
		}
		set["crud"] = temps
		set["databases"] = map[string][]string{"databases": dbs}
	}

	set["file_store_store_type"] = "na"
	set["file_store_rules"] = 0
	if project.FileStoreConfig != nil && project.FileStoreConfig.Enabled {
		set["file_store_store_type"] = project.FileStoreConfig.StoreType
		set["file_store_rules"] = len(project.FileStoreRules)
	}

	// auth info
	set["auth"] = map[string]interface{}{"providers": 0}
	if project.Auths != nil {
		temps := []string{}
		for _, v := range project.Auths {
			if v.Enabled {
				temps = append(temps, v.ID)
			}
		}
		set["auth"] = map[string]interface{}{"providers": len(temps)}
	}

	// services info
	set["services"] = 0
	if project.RemoteService != nil {
		set["services"] = len(project.RemoteService)
	}

	// let's encrypt info
	set["lets_encrypt"] = len(project.LetsEncrypt.WhitelistedDomains)

	// routing info
	set["routes"] = 0
	if project.IngressRoutes != nil {
		set["routes"] = len(project.IngressRoutes)
	}

	// eventing info
	set["total_events"] = len(project.EventingTriggers)

	return fmt.Sprintf("%s--%s", clusterID, projectID), set, min
}

func getTablesCount(dbAlias string, dbSchemas config.DatabaseSchemas, eventConf *config.EventingConfig) int {
	count := 0
	for _, schema := range dbSchemas {
		if schema.DbAlias == dbAlias && schema.Table != "default" {
			count++
		}
	}
	if eventConf.Enabled && eventConf.DBAlias == dbAlias {
		// NOTE : 2 is the number of tables used internally for eventing (invocation logs & event logs)
		count -= 2
	}
	return count
}

// NOTE: test not written for below function
func (m *Module) updateSCMetrics(id string, set, min map[string]interface{}) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, err := m.sink.Upsert("config_metrics").Where(types.Cond("id", "==", id)).Set(set).Min(min).Apply(ctx)
	if err != nil {
		_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to push metrics", err, nil)
	}
	if result == nil {
		// when space api go is not able to connect to server, the result is empty
		return
	}
	if result.Status != http.StatusOK {
		_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to push metrics", fmt.Errorf(result.Error), map[string]interface{}{"statusCode": result.Status})
	}
}
