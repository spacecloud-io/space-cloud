package syncman

import (
	"context"
	"fmt"
	"net/http"

	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
)

func (s *Manager) SetOfflineLicense(ctx context.Context, license string) error {
	helpers.Logger.LogDebug(helpers.GetRequestID(ctx), `Upgrading gateway to enterprise...`, nil)

	oldConfig := s.adminMan.GetConfig()
	oldConfig.License = license
	if err := s.adminMan.SetConfig(oldConfig); err != nil {
		return err
	}
	return s.SetLicense(ctx, oldConfig)
}

func (s *Manager) RenewLicense(ctx context.Context) error {
	helpers.Logger.LogDebug(helpers.GetRequestID(ctx), `Force renewing the license key...`, nil)
	if !s.adminMan.IsRegistered() {
		return helpers.Logger.LogError(helpers.GetRequestID(ctx), "Only registered clusters can force renew", nil, nil)
	}

	isLeader, err := s.leader.IsLeader(ctx, s.nodeID)
	if err != nil {
		return err
	}
	// A follower will forward this request to leader gateway
	if !isLeader {
		leaderNodeID, err := s.leader.GetLeaderNodeID(ctx)
		if err != nil {
			return helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to fetch leader node id from redis", err, nil)
		}
		if err := s.pubsubClient.Send(ctx, generatePubSubTopic(leaderNodeID, pubSubOperationRenew), ""); err != nil {
			return err
		}

		return nil
	}

	if err := s.adminMan.RenewLicense(true); err != nil {
		return err
	}

	return s.SetLicense(ctx, s.adminMan.GetConfig())
}

// GetLeaderGatewayID gets the current leader gateway node id
func (s *Manager) GetLeaderGatewayID(ctx context.Context) (string, error) {
	return s.leader.GetLeaderNodeID(ctx)
}

func (s *Manager) ConvertToEnterprise(ctx context.Context, req *model.LicenseUpgradeRequest) error {
	helpers.Logger.LogDebug(helpers.GetRequestID(ctx), `Upgrading gateway to enterprise...`, map[string]interface{}{"data": req})
	if s.adminMan.IsRegistered() {
		return helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to upgrade, already running in enterprise mode", nil, nil)
	}

	isLeader, err := s.leader.IsLeader(ctx, s.nodeID)
	if err != nil {
		return err
	}
	// A follower will forward this request to leader gateway
	if !isLeader {
		leaderNodeID, err := s.leader.GetLeaderNodeID(ctx)
		if err != nil {
			return helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to fetch leader node id from redis", err, nil)
		}
		if err := s.pubsubClient.Send(ctx, generatePubSubTopic(leaderNodeID, pubSubOperationUpgrade), req); err != nil {
			return err
		}

		return nil
	}

	// send request to spaceuptech server for upgrade
	sessionID, err := s.adminMan.GetSessionID()
	if err != nil {
		return err
	}
	upgradeResponse := new(model.GraphqlFetchLicenseResponse)
	body := map[string]interface{}{
		"params":  &map[string]interface{}{"sessionId": sessionID, "licenseKey": req.LicenseKey, "licenseValue": req.LicenseValue, "clusterName": req.ClusterName},
		"timeout": 10,
	}
	if err := s.MakeHTTPRequest(ctx, http.MethodPost, "https://api.spaceuptech.com/v1/api/spacecloud/services/billing/renewLicense", "", "", body, upgradeResponse); err != nil {
		return err
	}

	if upgradeResponse.Status != http.StatusOK {
		return fmt.Errorf("%s--%s--%d", upgradeResponse.Message, upgradeResponse.Error, upgradeResponse.Status)
	}

	// set updated admin config in config file
	if err := s.adminMan.SetConfig(&config.License{LicenseKey: req.LicenseKey, LicenseValue: req.LicenseValue, License: upgradeResponse.Result.License}); err != nil {
		return err
	}

	return s.SetLicense(ctx, s.adminMan.GetConfig())
}
