package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"runtime"
	"time"

	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/utils"
)

func currentTimeInMillis() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func (s *Server) generateMetricsRequest() (find, update map[string]interface{}) {
	// Get the cluster size
	clusterSize, err := s.syncMan.GetClusterSize(context.Background())
	if err != nil {
		clusterSize = 1
	}

	// Create the find and update clauses
	find = map[string]interface{}{"_id": s.nodeID}
	set := map[string]interface{}{
		"os":           runtime.GOOS,
		"isProd":       s.adminMan.LoadEnv(),
		"version":      utils.BuildVersion,
		"clusterSize":  clusterSize,
		"distribution": "ce",
		"lastUpdated":  currentTimeInMillis(),
	}
	min := map[string]interface{}{"startTime": currentTimeInMillis()}

	c := s.syncMan.GetGlobalConfig()
	if c != nil {
		set["sslEnabled"] = s.ssl != nil && s.ssl.Enabled
		if c.Admin != nil {
			set["mode"] = c.Admin.Operation.Mode
		}
		if c.Projects != nil && len(c.Projects) > 0 && c.Projects[0].Modules != nil {
			set["modules"] = getProjectInfo(c.Projects[0].Modules)
			set["projects"] = []string{c.Projects[0].ID}
		}
	}

	update = map[string]interface{}{"$set": set, "$min": min}
	return
}

func updateSCMetrics(find, update map[string]interface{}) error {

	req := map[string]interface{}{"find": find, "update": update, "op": "upsert"}
	jsonValue, err := json.Marshal(req)
	if err != nil {
		return err
	}
	resp, err := http.Post("https://api.spaceuptech.com/v1/api/space-cloud/crud/mongo/metrics/update", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New("Internal server error")
	}

	return nil
}

// RoutineMetrics routinely sends anonymous metrics
func (s *Server) RoutineMetrics() {
	ticker := time.NewTicker(time.Minute * 5)
	defer ticker.Stop()

	find, update := s.generateMetricsRequest()
	err := updateSCMetrics(find, update)
	if err != nil {
		// fmt.Println("Metrics Error -", err)
	}

	for range ticker.C {
		find, update := s.generateMetricsRequest()
		err := updateSCMetrics(find, update)
		if err != nil {
			// fmt.Println("Metrics Error -", err)
		}
	}
}

func getProjectInfo(config *config.Modules) map[string]interface{} {

	crudConfig := map[string]interface{}{"dbs": []string{}, "collections": 0}
	functionsConfig := map[string]interface{}{"enabled": false, "services": 0, "functions": 0}
	fileStoreConfig := map[string]interface{}{"enabled": false, "storeTypes": []string{}, "rules": 0}
	staticConfig := map[string]interface{}{"routes": 0, "internalRoutes": 0}
	auth := []string{}

	if config.Crud != nil {
		for k, v := range config.Crud {
			if v.Enabled {
				crudConfig["dbs"] = append(crudConfig["dbs"].([]string), k)
				if v.Collections != nil {
					crudConfig["collections"] = crudConfig["collections"].(int) + len(v.Collections)
				}
			}
		}
	}

	if config.Auth != nil {
		for k, v := range config.Auth {
			if v.Enabled {
				auth = append(auth, k)
			}
		}
	}

	if config.Services != nil {
		functionsConfig["enabled"] = true
		if config.Services.Services != nil {
			functionsConfig["services"] = functionsConfig["services"].(int) + len(config.Services.Services)
			for _, v := range config.Services.Services {
				if v != nil && v.Endpoints != nil {
					functionsConfig["functions"] = functionsConfig["functions"].(int) + len(v.Endpoints)
				}
			}
		}
	}

	if config.FileStore != nil && config.FileStore.Enabled {
		fileStoreConfig["enabled"] = true
		fileStoreConfig["storeTypes"] = []string{config.FileStore.StoreType}
		if config.FileStore.Rules != nil {
			fileStoreConfig["rules"] = len(config.FileStore.Rules)
		}
	}

	return map[string]interface{}{"crud": crudConfig, "functions": functionsConfig, "fileStore": fileStoreConfig, "auth": auth, "static": staticConfig}
}
