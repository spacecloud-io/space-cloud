package admin

import (
	"context"
	"errors"

	"github.com/stretchr/testify/mock"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
)

type mockIntegrationManager struct {
	mock.Mock
}

func (m *mockIntegrationManager) HandleConfigAuth(_ context.Context, resource, op string, claims map[string]interface{}, attr map[string]string) config.IntegrationAuthResponse {
	return m.Called(resource, op, claims, attr).Get(0).(config.IntegrationAuthResponse)
}

func (m *mockIntegrationManager) InvokeHook(ctx context.Context, params model.RequestParams) config.IntegrationAuthResponse {
	return m.Called(params).Get(0).(config.IntegrationAuthResponse)
}

type mockIntegrationResponse struct {
	checkResponse bool
	err           string
	result        interface{}
	status        int
}

func (m mockIntegrationResponse) CheckResponse() bool {
	return m.checkResponse
}

func (m mockIntegrationResponse) Result() interface{} {
	return m.result
}

func (m mockIntegrationResponse) Status() int {
	return m.status
}

func (m mockIntegrationResponse) Error() error {
	if m.err == "" {
		return nil
	}
	return errors.New(m.err)
}
