package syncman

import (
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"

	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/utils"
)

func (s *Manager) watchService() error {
	watchParams := map[string]interface{}{
		"type":    "service",
		"service": utils.SpaceCloudServiceName,
		"tag":     s.clusterID,
	}
	p, err := watch.Parse(watchParams)
	if err != nil {
		return err
	}

	p.HybridHandler = func(val watch.BlockingParamVal, data interface{}) {
		spaceClouds := data.([]*api.ServiceEntry)
		s.setSpaceCloudInstances(spaceClouds)
	}

	go func() {
		if err := p.Run(""); err != nil {
			log.Println("Sync Manager: could not start watch -", err)
			os.Exit(-1)
		}
	}()

	return nil
}

func (s *Manager) watchProjects() error {
	watchParams := map[string]interface{}{
		"type":   "keyprefix",
		"prefix": "sc/projects/" + s.clusterID,
	}
	p, err := watch.Parse(watchParams)
	if err != nil {
		return err
	}

	p.HybridHandler = func(val watch.BlockingParamVal, data interface{}) {
		s.lock.Lock()
		defer s.lock.Unlock()

		kvPairs := data.(api.KVPairs)

		var projects []*config.Project

		for _, kv := range kvPairs {
			a := strings.Split(kv.Key, "/")
			if a[2] != s.clusterID {
				continue
			}

			project := new(config.Project)
			if err := json.Unmarshal(kv.Value, project); err != nil {
				log.Println("Sync manager: Could not parse project received -", err)
				continue
			}

			projects = append(projects, project)
		}

		s.projectConfig.Projects = projects
		config.StoreConfigToFile(s.projectConfig, s.configFile)

		if s.projectConfig.Projects != nil && len(s.projectConfig.Projects) > 0 {
			s.cb(s.projectConfig)
		}
	}

	go func() {
		if err := p.Run(""); err != nil {
			log.Println("Sync Manager: could not start watcher -", err)
			os.Exit(-1)
		}
	}()

	return nil
}
