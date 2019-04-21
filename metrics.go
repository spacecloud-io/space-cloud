package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"time"

	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/proto"
)

type transport struct {
	stub proto.SpaceCloudClient
}

func (t *transport) insert(ctx context.Context, meta *proto.Meta, op string, obj interface{}) (int, error) {
	objJSON, err := json.Marshal(obj)
	if err != nil {
		return 500, err
	}

	req := proto.CreateRequest{Document: objJSON, Meta: meta, Operation: op}
	res, err := t.stub.Create(ctx, &req)
	if err != nil {
		return 500, err
	}

	if res.Status >= 200 || res.Status < 300 {
		return int(res.Status), nil
	}

	return int(res.Status), nil
}

func (t *transport) update(ctx context.Context, meta *proto.Meta, op string, find, update map[string]interface{}) (int, error) {
	updateJSON, err := json.Marshal(update)
	if err != nil {
		return 500, err
	}

	findJSON, err := json.Marshal(find)
	if err != nil {
		return 500, err
	}

	req := proto.UpdateRequest{Find: findJSON, Update: updateJSON, Meta: meta, Operation: op}
	res, err := t.stub.Update(ctx, &req)
	if err != nil {
		return 500, err
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
	return &transport{stub}, nil
}

func (s *server) routineMetrics() {
	ticker := time.NewTicker(time.Minute * 5)
	defer ticker.Stop()

	id := uuid.NewV1().String()
	col := "metrics"
	project := "crm"
	m := &proto.Meta{Col: col, DbType: "mongo", Project: project}

	trans, err := newTransport("spaceuptech.com", "11001", true)
	if err != nil {
		fmt.Println("Metrics Error -", err)
		return
	}

	startTime := time.Now().UTC()

	s.lock.Lock()
	obj := map[string]interface{}{"_id": id, "startTime": startTime, "lastUpdated": startTime}
	if s.config != nil && s.config.Modules != nil {
		obj["project"] = getProjectInfo(s.config.Modules)
	}
	s.lock.Unlock()

	status, err := trans.insert(context.TODO(), m, "one", obj)
	if err != nil {
		fmt.Println("Metrics Error2 -", err)
		return
	}

	if status != 200 {
		fmt.Println("Metrics Error3 -", status)
		return
	}

	for range ticker.C {
		update := map[string]interface{}{
			"$currentDate": map[string]interface{}{"lastUpdated": map[string]interface{}{"$type": "timestamp"}},
		}
		find := map[string]interface{}{"_id": id}

		s.lock.Lock()
		if s.config != nil && s.config.Modules != nil {
			update["$set"] = map[string]interface{}{"project": getProjectInfo(s.config.Modules)}
		}
		s.lock.Unlock()

		status, err := trans.update(context.TODO(), m, "one", find, update)
		if err != nil {
			return
		}

		if status != 200 {
			return
		}
	}
}

func getProjectInfo(config *config.Modules) map[string]interface{} {
	project := map[string]interface{}{
		"crud":      []string{},
		"faas":      map[string]interface{}{"enabled": false},
		"realtime":  map[string]interface{}{"enabled": false},
		"fileStore": map[string]interface{}{"enabled": false},
		"auth":      []string{},
	}

	if config.Crud != nil {
		crud := []string{}
		for k := range config.Crud {
			crud = append(crud, k)
		}
		project["crud"] = crud
	}

	if config.Auth != nil {
		auth := []string{}
		for k, v := range config.Auth {
			if v.Enabled {
				auth = append(auth, k)
			}
		}
		project["auth"] = auth
	}

	if config.FaaS != nil {
		project["faas"] = map[string]interface{}{"enabled": config.FaaS.Enabled}
	}

	if config.Realtime != nil {
		project["realtime"] = map[string]interface{}{"enabled": config.Realtime.Enabled}
	}

	if config.FileStore != nil {
		project["fileStore"] = map[string]interface{}{"enabled": config.FileStore.Enabled, "storeType": config.FileStore.StoreType}
	}

	return project
}
