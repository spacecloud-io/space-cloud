package cluster

import (
	"fmt"
	"net/http"

	"github.com/spaceuptech/space-cloud/space-cli/cmd/model"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils/transport"
)

type resp struct {
	Error  string      `json:"error,omitempty"`
	Result interface{} `json:"result,omitempty"`
}

// GetClusterConfig gets clusters config
func GetClusterConfig() ([]*model.SpecObject, error) {
	url := fmt.Sprintf("/v1/config/cluster")

	// Get the spec from the server
	payload := new(resp)
	if err := transport.Client.MakeHTTPRequest(http.MethodGet, url, map[string]string{}, payload); err != nil {
		return nil, err
	}

	var objs []*model.SpecObject

	spec := payload.Result.(map[string]interface{})

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

	for _, item := range payload.Result {
		spec := item.(map[string]interface{})
		meta := map[string]string{}
		// Printing the object on the screen
		s, err := utils.CreateSpecObject("/v1/config/integrations", "integrations", meta, spec)
		if err != nil {
			return nil, err
		}
		objs = append(objs, s)
	}

	return objs, nil
}
