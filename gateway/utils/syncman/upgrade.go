package syncman

import (
	"context"
	"fmt"
	"net/http"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

func (s *Manager) RenewLicense(ctx context.Context, token string) error {
	utils.LogDebug(`Force renewing the license key...`, "syncman", "RenewLicense", map[string]interface{}{})
	if !s.adminMan.IsRegistered() {
		return utils.LogError("Only registered clusters can force renew", "syncman", "RenewLicense", nil)
	}
	// A follower will forward this request to leader gateway
	if !s.checkIfLeaderGateway(s.nodeID) {
		service, err := s.getLeaderGateway()
		if err != nil {
			return err
		}

		url := fmt.Sprintf("http://%s/v1/config/renew-license", service.addr)
		params := map[string]string{}
		utils.LogDebug("Forwarding force renew request to leader", "syncman", "RenewLicense", map[string]interface{}{"leader": service.addr})
		return s.MakeHTTPRequest(ctx, http.MethodPost, url, token, "", params, &map[string]interface{}{})
	}

	if err := s.adminMan.RenewLicense(true); err != nil {
		return err
	}

	return s.SetAdminConfig(ctx, s.adminMan.GetConfig())
}

func (s *Manager) ConvertToEnterprise(ctx context.Context, token, clusterID, clusterKey string) error {
	utils.LogDebug(`Upgrading gateway to enterprise...`, "syncman", "ConvertToEnterprise", map[string]interface{}{"clusterId": clusterID, "clusterKey": clusterKey})
	if s.adminMan.IsRegistered() {
		return utils.LogError("Unable to upgrade, already running in enterprise mode", "syncman", "ConvertToEnterprise", nil)
	}

	// A follower will forward this request to leader gateway
	if !s.checkIfLeaderGateway(s.nodeID) {
		service, err := s.getLeaderGateway()
		if err != nil {
			return err
		}

		url := fmt.Sprintf("http://%s/v1/config/upgrade", service.addr)
		params := map[string]string{"clusterId": clusterID, "clusterKey": clusterKey}
		utils.LogDebug("Forwarding upgrade request to leader", "syncman", "ConvertToEnterprise", map[string]interface{}{"leader": service.addr})
		return s.MakeHTTPRequest(ctx, http.MethodPost, url, token, "", params, &map[string]interface{}{})
	}

	// send request to spaceuptech server for upgrade
	upgradeResponse := new(model.GraphqlFetchLicenseResponse)
	body := map[string]interface{}{
		"params":  &map[string]interface{}{"sessionId": s.adminMan.GetSessionID(), "clusterId": clusterID, "clusterKey": clusterKey},
		"timeout": 10,
	}
	if err := s.MakeHTTPRequest(ctx, http.MethodPost, "http://35.188.208.249/v1/api/spacecloud/services/backend/fetch_license", "", "", body, upgradeResponse); err != nil {
		return err
	}

	if upgradeResponse.Result.Status != http.StatusOK {
		return fmt.Errorf("%s--%s--%d", upgradeResponse.Result.Message, upgradeResponse.Result.Error, upgradeResponse.Result.Status)
	}

	// set updated admin config in config file
	clusterConfig := &config.Admin{ClusterID: upgradeResponse.Result.Result.ClusterID, ClusterKey: upgradeResponse.Result.Result.ClusterKey, License: upgradeResponse.Result.Result.License}
	if err := s.adminMan.SetConfig(clusterConfig, false); err != nil {
		return err
	}

	return s.SetAdminConfig(ctx, clusterConfig)
}
