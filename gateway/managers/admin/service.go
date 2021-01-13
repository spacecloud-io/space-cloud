package admin

import (
	"context"

	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
)

func (m *Manager) SetServices(eventType string, services model.ScServices) {
	m.lock.Lock()
	defer m.lock.Unlock()
	helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Setting services in admin", map[string]interface{}{"eventType": eventType, "services": services})

	m.services = services
	if eventType == config.ResourceDeleteEvent {
		if m.isEnterpriseMode() {
			if err := m.validateLicenseWithoutLock(services); err != nil {
				return
			}
		}
	}
}
