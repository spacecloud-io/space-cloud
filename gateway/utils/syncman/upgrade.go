package syncman

import (
	"context"
	"fmt"
	"net/http"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
)

func (s *Manager) ConvertToEnterprise(ctx context.Context, userName, key, clusterID, clusterKey string, version int) (string, error) {
	// verify credentials
	_, _, err := s.adminMan.Login(userName, key)
	if err != nil {
		return "", err
	}

	// get cluster-type(docker or kubernetes) from runner
	token, err := s.adminMan.GetInternalAccessToken()
	if err != nil {
		return "", err
	}
	response := new(model.Response)
	if err := s.MakeHTTPRequest(ctx, http.MethodGet, fmt.Sprintf("http://%s/v1/runner/cluster-type", s.runnerAddr), token, "", map[string]interface{}{}, response); err != nil {
		return "", err
	}

	// set new cluster id & key in config file
	return response.Result.(string), s.setAdminConfig(ctx, &config.Admin{ClusterID: clusterID, ClusterKey: clusterKey, Version: version})
}
func (s *Manager) VersionUpgrade(ctx context.Context, version int) error {
	c := s.GetGlobalConfig().Admin
	c.Version = version
	return s.setAdminConfig(ctx, c)
}
