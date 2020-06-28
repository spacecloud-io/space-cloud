package metrics

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spaceuptech/space-api-go/types"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

func (m *Module) generateMetricsRequest(project *config.Project, ssl *config.SSL) (string, map[string]interface{}, map[string]interface{}) {
	m.lock.Lock()
	defer m.lock.Unlock()

	clusterID := m.adminMan.GetClusterID()
	projectID := project.ID
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

	modules := project.Modules
	// crud info
	set["crud"] = map[string]interface{}{"tables": map[string]interface{}{}}
	set["databases"] = map[string][]string{"databases": {}}
	if modules.Crud != nil {
		temps := map[string]interface{}{}
		dbs := []string{}
		for dbAlias, v := range modules.Crud {
			if v.Enabled {
				dbs = append(dbs, v.Type)
				temps[dbAlias] = map[string]interface{}{
					"tables": len(v.Collections) - 3, // NOTE : 2 is the number of tables used internally for eventing (invocation logs & event logs) + 1 which is the default table
				}
			}
		}
		set["crud"] = temps
		set["databases"] = map[string][]string{"databases": dbs}
	}

	set["file_store_store_type"] = "na"
	set["file_store_rules"] = 0
	if modules.FileStore != nil && modules.FileStore.Enabled {
		set["file_store_store_type"] = modules.FileStore.StoreType
		set["file_store_rules"] = len(modules.FileStore.Rules)
	}

	// auth info
	set["auth"] = map[string]interface{}{"providers": 0}
	if modules.Auth != nil {
		temps := []string{}
		for k, v := range modules.Auth {
			if v.Enabled {
				temps = append(temps, k)
			}
		}
		set["auth"] = map[string]interface{}{"providers": len(temps)}
	}

	// services info
	set["services"] = 0
	if modules.Services != nil {
		set["services"] = len(modules.Services.Services)
	}

	// let's encrypt info
	set["lets_encrypt"] = len(modules.LetsEncrypt.WhitelistedDomains)

	// routing info
	set["routes"] = 0
	if modules.Routes != nil {
		set["routes"] = len(modules.Routes)
	}

	// eventing info
	set["total_events"] = len(modules.Eventing.Rules)

	return fmt.Sprintf("%s--%s", clusterID, projectID), set, min
}

// NOTE: test not written for below function
func (m *Module) updateSCMetrics(id string, set, min map[string]interface{}) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, err := m.sink.Upsert("config_metrics").Where(types.Cond("id", "==", id)).Set(set).Min(min).Apply(ctx)
	if err != nil {
		logrus.Errorf("error querying database got error")
	}
	if result == nil {
		// when space api go is not able to connect to server, the result is empty
		return
	}
	if result.Status != http.StatusOK {
		logrus.Errorf("error querying database got status (%d) (%s)", result.Status, result.Error)
	}
}
