package metrics

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/spaceuptech/space-api-go/types"
	"github.com/spaceuptech/space-cloud/gateway/utils"
	"net/http"
	"runtime"
	"time"
)

const metricsUpdaterInterval = 30 * time.Second

func currentTimeInMillis() int64 {
	// subtracting interval time make sures that multiple gateways in a cluster don't write to database frequently
	return time.Now().Add(-metricsUpdaterInterval).UnixNano()
}

func (m *Module) generateMetricsRequest() (string, map[string]interface{}, map[string]interface{}, bool) {
	m.lock.Lock()
	defer m.lock.Unlock()

	c := m.syncMan.GetGlobalConfig()
	if c == nil {
		return "", nil, nil, true
	}

	// Get the cluster size
	clusterSize, err := m.syncMan.GetClusterSize(context.Background())
	if err != nil {
		clusterSize = 1
	}
	find := m.syncMan.GetClusterID()
	min := map[string]interface{}{"start_time": currentTimeInMillis()}
	// Create the find and update clauses
	set := map[string]interface{}{
		"nodes":        m.syncMan.GetNodesInCluster(),
		"os":           runtime.GOOS,
		"is_prod":      m.isProd,
		"version":      utils.BuildVersion,
		"cluster_size": clusterSize,
		"distribution": "ce",
		"last_updated": time.Now().Format(time.RFC3339),
	}

	set["ssl_enabled"] = m.ssl != nil && m.ssl.Enabled
	set["mode"] = 0
	if c.Admin != nil {
		set["mode"] = c.Admin.Operation.Mode
	}
	if c.Projects != nil && len(c.Projects) > 0 && c.Projects[0].Modules != nil {
		modules := c.Projects[0].Modules
		set["project"] = c.Projects[0].ID

		// crud info
		set["crud"] = map[string]interface{}{"tables": map[string]interface{}{}}
		if modules.Crud != nil {
			temps := map[string]interface{}{}
			for _, v := range modules.Crud {
				if v.Enabled {
					temps[v.Type] = map[string]interface{}{
						"tables": len(v.Collections) - 2, // NOTE : 2 is the number of tables used internally for eventing (invocation logs & event logs)
					}
				}
			}
			set["crud"] = temps
		}

		// file store info
		set["file_store_store_type"] = ""
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
		temps := map[string]interface{}{}
		for k, v := range modules.Eventing.Rules {
			temps[k] = map[string]interface{}{
				"type": m.eventing[v.Type],
			}
		}
		m.eventing = map[string]int{}

		set["eventing"] = temps
	}

	return find, set, min, false
}

func (m *Module) updateSCMetrics(find string, set, min map[string]interface{}) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, err := m.sink.Upsert("config_metrics").Where(types.Cond("id", "==", find)).Set(set).Min(min).Apply(ctx)
	if err != nil {
		logrus.Errorf("error querying database got error")
	}
	if result.Status != http.StatusOK {
		logrus.Errorf("error querying database got status (%d) (%s)", result.Status, result.Error)
	}
}
