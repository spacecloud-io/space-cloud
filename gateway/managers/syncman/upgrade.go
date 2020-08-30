package syncman

import (
	"context"
	"fmt"
	"net/http"

	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/model"
)

func (s *Manager) SetOfflineLicense(ctx context.Context, license string) error {
	helpers.Logger.LogDebug(helpers.GetRequestID(ctx), `Upgrading gateway to enterprise...`, nil)

	oldConfig := s.adminMan.GetConfig()
	oldConfig.License = license
	return s.SetAdminConfig(ctx, oldConfig)
}

func (s *Manager) RenewLicense(ctx context.Context, token string) error {
	helpers.Logger.LogDebug(helpers.GetRequestID(ctx), `Force renewing the license key...`, nil)
	if !s.adminMan.IsRegistered() {
		return helpers.Logger.LogError(helpers.GetRequestID(ctx), "Only registered clusters can force renew", nil, nil)
	}
	// A follower will forward this request to leader gateway
	if !s.checkIfLeaderGateway(s.nodeID) {
		service, err := s.getLeaderGateway()
		if err != nil {
			return err
		}

		url := fmt.Sprintf("http://%s/v1/config/renew-license", service.addr)
		params := map[string]string{}
		helpers.Logger.LogDebug(helpers.GetRequestID(ctx), "Forwarding force renew request to leader", map[string]interface{}{"leader": service.addr})
		return s.MakeHTTPRequest(ctx, http.MethodPost, url, token, "", params, &map[string]interface{}{})
	}

	if err := s.adminMan.RenewLicense(true); err != nil {
		return err
	}

	return s.SetAdminConfig(ctx, s.adminMan.GetConfig())
}

func (s *Manager) ConvertToEnterprise(ctx context.Context, token, licenseKey, licenseValue, clusterName string) error {
	helpers.Logger.LogDebug(helpers.GetRequestID(ctx), `Upgrading gateway to enterprise...`, map[string]interface{}{"LicenseKey": licenseKey, "LicenseValue": licenseValue, "clusterName": clusterName})
	if s.adminMan.IsRegistered() {
		return helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to upgrade, already running in enterprise mode", nil, nil)
	}

	// A follower will forward this request to leader gateway
	if !s.checkIfLeaderGateway(s.nodeID) {
		service, err := s.getLeaderGateway()
		if err != nil {
			return err
		}

		url := fmt.Sprintf("http://%s/v1/config/upgrade", service.addr)
		params := map[string]string{"licenseKey": licenseKey, "licenseValue": licenseValue}
		helpers.Logger.LogDebug(helpers.GetRequestID(ctx), "Forwarding upgrade request to leader", map[string]interface{}{"leader": service.addr})
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
		return fmt.Errorf("%s--%s--%d", upgradeResponse.Message, upgradeResponse.Error, upgradeResponse.Status)
	}

	// set updated admin config in config file
	oldConfig := s.adminMan.GetConfig()
	oldConfig.LicenseKey = licenseKey
	oldConfig.LicenseValue = licenseValue
	oldConfig.License = upgradeResponse.Result.License

	return s.SetAdminConfig(ctx, oldConfig)
}
