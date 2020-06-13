package syncman

import (
	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/stretchr/testify/mock"
)

type mockAdminSyncmanInterface struct {
	mock.Mock
}

func (m *mockAdminSyncmanInterface) IsTokenValid(token string) error {
	c := m.Called(token)
	return c.Error(0)
}

func (m *mockAdminSyncmanInterface) GetInternalAccessToken() (string, error) {
	c := m.Called()
	return c.String(0), c.Error(1)
}

func (m *mockAdminSyncmanInterface) ValidateSyncOperation(c *config.Config, project *config.Project) bool {
	a := m.Called(c, project)
	return a.Bool(0)
}
