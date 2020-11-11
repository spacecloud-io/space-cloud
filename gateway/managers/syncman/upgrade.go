package syncman

import (
	"context"

	"github.com/spaceuptech/helpers"
)

func (s *Manager) SetOfflineLicense(ctx context.Context, license string) error {
	helpers.Logger.LogDebug(helpers.GetRequestID(ctx), `Upgrading gateway to enterprise...`, nil)

	oldConfig := s.adminMan.GetConfig()
	oldConfig.License = license
	return s.SetLicense(ctx, oldConfig)
}

func (s *Manager) RenewLicense(ctx context.Context, token string) error {
	return nil
}

func (s *Manager) ConvertToEnterprise(ctx context.Context, token, licenseKey, licenseValue, clusterName string) error {
	return nil
}
