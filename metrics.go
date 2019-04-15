package main

import (
	"fmt"
	"time"

	uuid "github.com/satori/go.uuid"

	api "github.com/spaceuptech/space-api-go"
	"github.com/spaceuptech/space-api-go/api/utils"

	"github.com/spaceuptech/space-cloud/config"
)

func (s *server) routineMetrics() {
	ticker := time.NewTicker(time.Minute * 5)
	defer ticker.Stop()

	id := uuid.NewV1().String()
	col := "metrics"

	a, err := api.New("crm", "spaceuptech.com", "11001", true)
	if err != nil {
		fmt.Println("Error -", err)
	}

	db := a.Mongo()
	startTime := time.Now().UTC()

	s.lock.Lock()
	obj := utils.M{"_id": id, "startTime": startTime, "lastUpdated": startTime}
	if s.config != nil && s.config.Modules != nil {
		obj["project"] = getProjectInfo(s.config.Modules)
	}
	s.lock.Unlock()

	res, err := db.Insert(col).Doc(obj).Apply()
	if err != nil {
		fmt.Println("Error -", err)
		return
	}

	if res.Status != 200 {
		fmt.Println("Error -", res.Error)
		return
	}

	for range ticker.C {
		query := db.Update(col).Where(utils.Cond("_id", "==", id)).CurrentTimestamp("lastUpdated")

		s.lock.Lock()
		if s.config != nil && s.config.Modules != nil {
			query = query.Set(utils.M{"project": getProjectInfo(s.config.Modules)})
		}
		s.lock.Unlock()

		res, err := query.Apply()
		if err != nil {
			fmt.Println("Error -", err)
			return
		}

		if res.Status != 200 {
			fmt.Println("Error -", res.Error)
			return
		}
	}
}

func getProjectInfo(config *config.Modules) utils.M {
	project := utils.M{
		"crud":      []string{},
		"faas":      utils.M{"enabled": false},
		"realtime":  utils.M{"enabled": false},
		"fileStore": utils.M{"enabled": false},
		"auth":      []string{},
	}

	if config.Crud != nil {
		crud := make([]string, len(config.Crud))
		for k := range config.Crud {
			crud = append(crud, k)
		}
		project["crud"] = crud
	}

	if config.Auth != nil {
		auth := make([]string, len(config.Auth))
		for k, v := range config.Auth {
			if v.Enabled {
				auth = append(auth, k)
			}
		}
		project["auth"] = auth
	}

	if config.FaaS != nil {
		project["faas"] = utils.M{"enabled": config.FaaS.Enabled}
	}

	if config.Realtime != nil {
		project["realtime"] = utils.M{"enabled": config.Realtime.Enabled}
	}

	if config.FileStore != nil {
		project["fileStore"] = utils.M{"enabled": config.FileStore.Enabled, "storeType": config.FileStore.StoreType}
	}

	return project
}
