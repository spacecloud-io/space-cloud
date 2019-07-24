package server

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"net/http"
	"runtime"
	"time"

	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/proto"
	"github.com/spaceuptech/space-cloud/utils"
)

type transport struct {
	conn *grpc.ClientConn
	stub proto.SpaceCloudClient
}

func (t *transport) insert(ctx context.Context, meta *proto.Meta, op string, obj interface{}) (int, error) {
	objJSON, err := json.Marshal(obj)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	req := proto.CreateRequest{Document: objJSON, Meta: meta, Operation: op}
	res, err := t.stub.Create(ctx, &req)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if res.Status >= 200 || res.Status < 300 {
		return int(res.Status), nil
	}

	return int(res.Status), nil
}

func getCurrentTimeinMillis() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func (t *transport) update(ctx context.Context, meta *proto.Meta, op string, find, update map[string]interface{}) (int, error) {
	updateJSON, err := json.Marshal(update)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	findJSON, err := json.Marshal(find)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	req := proto.UpdateRequest{Find: findJSON, Update: updateJSON, Meta: meta, Operation: op}
	res, err := t.stub.Update(ctx, &req)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if res.Status >= 200 || res.Status < 300 {
		return int(res.Status), nil
	}

	return int(res.Status), nil
}

// Init initialises a new transport
func newTransport(host, port string, sslEnabled bool) (*transport, error) {
	dialOptions := []grpc.DialOption{}

	if sslEnabled {
		dialOptions = append(dialOptions, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
	} else {
		dialOptions = append(dialOptions, grpc.WithInsecure())
	}

	conn, err := grpc.Dial(host+":"+port, dialOptions...)
	if err != nil {
		return nil, err
	}

	stub := proto.NewSpaceCloudClient(conn)
	return &transport{conn, stub}, nil
}

// RoutineMetrics routinely sends anonymous metrics
func (s *Server) RoutineMetrics() {
	ticker := time.NewTicker(time.Minute * 5)
	defer ticker.Stop()

	id := uuid.NewV1().String()
	col := "metrics"
	project := "space-cloud"
	m := &proto.Meta{Col: col, DbType: "mongo", Project: project}

	// Create the find and update clauses
	find := map[string]interface{}{"_id": id}
	update := map[string]interface{}{
		"$currentDate": map[string]interface{}{
			"lastUpdated": map[string]interface{}{"$type": "timestamp"},
			"startTime":   map[string]interface{}{"$type": "timestamp"},
		},
	}
	set := map[string]interface{}{
		"os":           runtime.GOOS,
		"isProd":       s.isProd,
		"version":      utils.BuildVersion,
		"clusterSize":  s.syncMan.GetClusterSize(),
		"distribution": "ce",
	}

	// Connect to metrics Server
	trans, err := newTransport("api.spaceuptech.com", "4128", true)
	if err != nil {
		//fmt.Println("Metrics Error -", err)
		return
	}

	c := s.syncMan.GetGlobalConfig()
	if c != nil {
		set["sslEnabled"] = s.ssl != nil && s.ssl.Enabled
		set["deployConfig"] = map[string]interface{}{"enabled": c.Deploy.Enabled, "orchestrator": c.Deploy.Orchestrator}
		if c.Admin != nil {
			set["mode"] = c.Admin.Operation.Mode
		}
		if c.Projects != nil && len(c.Projects) > 0 && c.Projects[0].Modules != nil {
			set["modules"] = getProjectInfo(c.Projects[0].Modules)
			set["projects"] = []string{c.Projects[0].ID}
		}
	}

	update["$set"] = set
	status, err := trans.update(context.TODO(), m, "upsert", find, update)
	if err != nil {
		//fmt.Println("Metrics Error -", err)
		return
	}

	if status != 200 {
		//fmt.Println("Metrics Error - Upsert failed: Invalid status code ", status)
		return
	}

	for range ticker.C {
		update := map[string]interface{}{
			"$currentDate": map[string]interface{}{"lastUpdated": map[string]interface{}{"$type": "timestamp"}},
			"clusterSize":  s.syncMan.GetClusterSize(),
		}

		c := s.syncMan.GetGlobalConfig()
		if c != nil {
			set["sslEnabled"] = s.ssl != nil && s.ssl.Enabled
			set["deployConfig"] = map[string]interface{}{"enabled": c.Deploy.Enabled, "orchestrator": c.Deploy.Orchestrator}
			if c.Admin != nil {
				set["mode"] = c.Admin.Operation.Mode
			}
			if c.Projects != nil && len(c.Projects) > 0 && c.Projects[0].Modules != nil {
				set["modules"] = getProjectInfo(c.Projects[0].Modules)
				set["projects"] = []string{c.Projects[0].ID}
			}
		}

		update["$set"] = set
		status, err := trans.update(context.TODO(), m, "one", find, update)
		if err != nil {
			//log.Println("Metrics Error -", err)
		}

		if status != 200 {
			//log.Println("Metrics Error - Invalid status code ", status)
		}
	}
}

func getProjectInfo(config *config.Modules) map[string]interface{} {

	crudConfig := map[string]interface{}{"dbs": []string{}, "collections": 0}
	functionsConfig := map[string]interface{}{"enabled": false, "services": 0, "functions": 0}
	realtimeConfig := map[string]interface{}{"enabled": false}
	fileStoreConfig := map[string]interface{}{"enabled": false, "storeTypes": []string{}, "rules": 0}
	staticConfig := map[string]interface{}{"enabled": false, "routes": 0}
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

	if config.Functions != nil && config.Functions.Enabled {
		functionsConfig["enabled"] = true
		if config.Functions.Rules != nil {
			functionsConfig["services"] = functionsConfig["services"].(int) + len(config.Functions.Rules)
			for _, v := range config.Functions.Rules {
				if v != nil {
					functionsConfig["functions"] = functionsConfig["functions"].(int) + len(v)
				}
			}
		}
	}

	if config.Realtime != nil && config.Realtime.Enabled {
		realtimeConfig["enabled"] = true
	}

	if config.FileStore != nil && config.FileStore.Enabled {
		fileStoreConfig["enabled"] = true
		fileStoreConfig["storeTypes"] = []string{config.FileStore.StoreType}
		if config.FileStore.Rules != nil {
			fileStoreConfig["rules"] = len(config.FileStore.Rules)
		}
	}

	if config.Static != nil && config.Static.Enabled {
		staticConfig["enabled"] = true
		if config.Static.Routes != nil {
			staticConfig["routes"] = len(config.Static.Routes)
		}
	}

	return map[string]interface{}{"crud": crudConfig, "functions": functionsConfig, "realtime": realtimeConfig, "fileStore": fileStoreConfig, "static": staticConfig, "auth": auth}
}
