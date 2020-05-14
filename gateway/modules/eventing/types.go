package eventing

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"

	"github.com/stretchr/testify/mock"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
)

type mockHTTPInterface struct {
	mock.Mock
}

func (m *mockHTTPInterface) Do(req *http.Request) (*http.Response, error) {
	c := m.Called(req)
	return &http.Response{StatusCode: 200, Body: ioutil.NopCloser(bytes.NewReader([]byte(`{"event": {"type": "someType"}, "response": "response"}`)))}, c.Error(1)
}

type mockCrudInterface struct {
	mock.Mock
}

func (m *mockCrudInterface) InternalCreate(ctx context.Context, dbAlias, project, col string, req *model.CreateRequest, isIgnoreMetrics bool) error {
	c := m.Called(ctx, dbAlias, project, col, req, isIgnoreMetrics)
	if err := c.Error(0); err != nil {
		return err
	}
	return nil
}

func (m *mockCrudInterface) Read(ctx context.Context, dbAlias, col string, req *model.ReadRequest) (interface{}, error) {
	c := m.Called(ctx, dbAlias, col, req)
	if len(c) > 1 {
		return c.Get(0).(interface{}), c.Error(1)
	}
	return c.Get(0).(interface{}), nil
}

func (m *mockCrudInterface) InternalUpdate(ctx context.Context, dbAlias, project, col string, req *model.UpdateRequest) error {
	c := m.Called(ctx, dbAlias, project, col, req)
	return c.Error(0)
}

type mockSyncmanEventingInterface struct {
	mock.Mock
}

func (m *mockSyncmanEventingInterface) GetAssignedSpaceCloudURL(ctx context.Context, project string, token int) (string, error) {
	c := m.Called(ctx, project, token)
	return c.String(0), c.Error(1)
}

func (m *mockSyncmanEventingInterface) GetAssignedTokens() (start, end int) {
	c := m.Called()
	return c.Int(0), c.Int(1)
}

func (m *mockSyncmanEventingInterface) GetEventSource() string {
	c := m.Called()
	return c.String(0)
}

func (m *mockSyncmanEventingInterface) GetNodeID() string {
	c := m.Called()
	return c.String(0)
}

func (m *mockSyncmanEventingInterface) GetSpaceCloudURLFromID(nodeID string) (string, error) {
	c := m.Called(nodeID)
	if len(c) > 1 {
		return c.String(0), c.Error(1)
	}
	return c.String(0), nil
}

func (m *mockSyncmanEventingInterface) MakeHTTPRequest(ctx context.Context, method, url, token, scToken string, params, vPtr interface{}) error {
	c := m.Called(ctx, method, url, token, scToken, params, vPtr)
	return c.Error(0)
}

type mockAdminEventingInterface struct {
	mock.Mock
}

func (m *mockAdminEventingInterface) GetInternalAccessToken() (string, error) {
	c := m.Called()
	return c.String(0), c.Error(1)
}

type mockAuthEventingInterface struct {
	mock.Mock
}

func (m *mockAuthEventingInterface) IsEventingOpAuthorised(ctx context.Context, project, token string, event *model.QueueEventRequest) error {
	c := m.Called(ctx, project, token, event)
	if err := c.Error(0); err != nil {
		return err
	}
	return nil
}

func (m *mockAuthEventingInterface) GetSCAccessToken() (string, error) {
	c := m.Called()
	return mock.Anything, c.Error(1)
}

func (m *mockAuthEventingInterface) GetInternalAccessToken() (string, error) {
	c := m.Called()
	return c.String(0), c.Error(1)
}

type mockSchemaEventingInterface struct {
	mock.Mock
}

func (m *mockSchemaEventingInterface) CheckIfEventingIsPossible(dbAlias, col string, obj map[string]interface{}, isFind bool) (findForUpdate map[string]interface{}, present bool) {
	c := m.Called(dbAlias, col, obj, isFind)
	return map[string]interface{}{}, c.Bool(1)
}

func (m *mockSchemaEventingInterface) Parser(crud config.Crud) (model.Type, error) {
	c := m.Called(crud)
	return nil, c.Error(1)
}

func (m *mockSchemaEventingInterface) SchemaValidator(col string, collectionFields model.Fields, doc map[string]interface{}) (map[string]interface{}, error) {
	c := m.Called(col, collectionFields, doc)
	return nil, c.Error(1)
}

func (m *mockSchemaEventingInterface) SchemaModifyAll(ctx context.Context, dbAlias, logicalDBName string, tables map[string]*config.TableRule) error {
	c := m.Called(ctx, dbAlias, logicalDBName, tables)
	return c.Error(0)
}

type mockFileStoreEventingInterface struct {
	mock.Mock
}

func (m *mockFileStoreEventingInterface) DoesExists(ctx context.Context, project, token, path string) error {
	c := m.Called(ctx, project, token, path)
	return c.Error(0)
}
