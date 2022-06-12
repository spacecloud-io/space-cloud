package common

import (
	"encoding/json"
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/spacecloud-io/space-cloud/config"
	"github.com/spacecloud-io/space-cloud/model"
)

func prepareRemoteServiceApp(fileConfig *model.SCConfig) json.RawMessage {
	config := map[string]interface{}{
		"services": prepareRemoteServices(fileConfig),
	}

	data, _ := json.Marshal(config)
	return data
}

func prepareRemoteServices(fileConfig *model.SCConfig) map[string]*config.Service {
	remoteServices := make(map[string]*config.Service)
	module, ok := fileConfig.Config["remote-service"]
	if !ok {
		return remoteServices
	}
	resourceObjects, ok := module["config"]
	if !ok {
		return remoteServices
	}

	for _, resourceObject := range resourceObjects {
		serviceConfig := new(config.Service)
		if err := mapstructure.Decode(resourceObject.Spec, serviceConfig); err != nil {
			return map[string]*config.Service{}
		}

		projectID := resourceObject.Meta.Parents["project"]
		serviceName := resourceObject.Meta.Name
		name := fmt.Sprintf("%s---%s", projectID, serviceName)
		remoteServices[name] = serviceConfig
	}

	return remoteServices
}
