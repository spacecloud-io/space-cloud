package syncman

import (
	"context"
	"fmt"
	"net/http"

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

func (s *Manager) ConvertToEnterprise(ctx context.Context, token, licenseKey, licenseValue, clusterName string) error {
	utils.LogDebug(`Upgrading gateway to enterprise...`, "syncman", "convert-to-enterprise", map[string]interface{}{"LicenseKey": licenseKey, "LicenseValue": licenseValue, "clusterName": clusterName})
	if s.adminMan.IsRegistered() {
		return utils.LogError("Unable to upgrade, already running in enterprise mode", "syncman", "convert-to-enterprise", nil)
	}

	// A follower will forward this request to leader gateway
	if !s.checkIfLeaderGateway(s.nodeID) {
		service, err := s.getLeaderGateway()
		if err != nil {
			return err
		}

		url := fmt.Sprintf("http://%s/v1/config/upgrade", service.addr)
		params := map[string]string{"licenseKey": licenseKey, "licenseValue": licenseValue}
		utils.LogDebug("Forwarding upgrade request to leader", "syncman", "convert-to-enterprise", map[string]interface{}{"leader": service.addr})
		return s.MakeHTTPRequest(ctx, http.MethodPost, url, token, "", params, &map[string]interface{}{})
	}

	// send request to spaceuptech server for upgrade
	upgradeResponse := new(model.GraphqlFetchLicenseResponse)
	body := map[string]interface{}{
		"params":  &map[string]interface{}{"sessionId": s.adminMan.GetSessionID(), "licenseKey": licenseKey, "licenseValue": licenseValue, "clusterName": clusterName},
		"timeout": 10,
	}
	if err := s.MakeHTTPRequest(ctx, http.MethodPost, "https://api.spaceuptech.com/v1/api/spacecloud/services/billing/renewLicense", "", "", body, upgradeResponse); err != nil {
		return err
	}

	if upgradeResponse.Status != http.StatusOK {
		return fmt.Errorf("%s--%s--%d", upgradeResponse.Message, upgradeResponse.Result.Error, upgradeResponse.Status)
	}

	// set updated admin config in config file
	oldConfig := s.adminMan.GetConfig()
	oldConfig.LicenseKey = licenseKey
	oldConfig.LicenseValue = licenseValue
	oldConfig.License = upgradeResponse.Result.License

	return s.SetAdminConfig(ctx, oldConfig)
}
