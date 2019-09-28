package deploy

import (
	"context"
	"errors"
	"net/http"
)

// UploadAndDeploy uploads a service to the registry then deploys it
func (m *Module) UploadAndDeploy(r *http.Request) error {
	m.lock.RLock()
	defer m.lock.RUnlock()

	if m.driver == nil {
		return errors.New("Deploy: No driver initialised")
	}

	c, err := m.upload(*m.config.Registry.Token, r)
	if err != nil {
		return err
	}

	return m.driver.Deploy(context.TODO(), c)
}
