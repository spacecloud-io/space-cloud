package cluster

import (
	"fmt"
	"net/http"

	"github.com/spaceuptech/space-cloud/space-cli/cmd/model"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils/transport"
)

// GetClusterConfig gets clusters config
func GetClusterConfig() ([]*model.SpecObject, error) {
	url := fmt.Sprintf("/v1/config/cluster")

	// Get the spec from the server
	payload := new(model.Resp)
	if err := transport.Client.MakeHTTPRequest(http.MethodGet, url, map[string]string{}, payload); err != nil {
		return nil, err
	}

	var objs []*model.SpecObject

	spec := map[string]interface{}{"clusterConfig": payload.Result.(map[string]interface{})}

	// Printing the object on the screen
	s, err := utils.CreateSpecObject("/v1/config/cluster", "cluster-config", map[string]string{}, spec)
	if err != nil {
		return nil, err
	}
	objs = append(objs, s)

	return objs, nil
}

// GetIntegration gets integration
func GetIntegration() ([]*model.SpecObject, error) {
	url := fmt.Sprintf("/v1/config/integrations")

	// Get the spec from the server
	payload := new(model.Response)
	if err := transport.Client.MakeHTTPRequest(http.MethodGet, url, map[string]string{}, payload); err != nil {
		return nil, err
	}

	var objs []*model.SpecObject
	return objs, nil
}
