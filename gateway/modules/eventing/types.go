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

var limit int64 = 1000

type queueUpdateEvent struct {
	project, db, col string
	req              *model.UpdateRequest
	err              string
}

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

func (m *mockCrudInterface) GetDBType(dbAlias string) (string, error) {
	panic("implement me")
}

func (m *mockCrudInterface) InternalCreate(ctx context.Context, dbAlias, project, col string, req *model.CreateRequest, isIgnoreMetrics bool) error {
	c := m.Called(ctx, dbAlias, project, col, req, isIgnoreMetrics)
	if err := c.Error(0); err != nil {
		return err
	}
	return nil
}

func (m *mockCrudInterface) Read(ctx context.Context, dbAlias, col string, req *model.ReadRequest, params model.RequestParams) (interface{}, *model.SQLMetaData, error) {
	c := m.Called(ctx, dbAlias, col, req)
	if len(c) > 1 {
		return c.Get(0).(interface{}), c.Get(1).(*model.SQLMetaData), c.Error(2)
	}
	return c.Get(0).(interface{}), c.Get(1).(*model.SQLMetaData), nil
}

func (m *mockCrudInterface) InternalUpdate(ctx context.Context, dbAlias, project, col string, req *model.UpdateRequest) error {
	c := m.Called(ctx, dbAlias, project, col, req)
	return c.Error(0)
}

type mockSyncmanEventingInterface struct {
	mock.Mock
}

func (m *mockSyncmanEventingInterface) GetAssignedSpaceCloudID(ctx context.Context, project string, token int) (string, error) {
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

func (m *mockSyncmanEventingInterface) GetSpaceCloudURLFromID(ctx context.Context, nodeID string) (string, error) {
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

type mockAuthEventingInterface struct {
	mock.Mock
}

func (m *mockAuthEventingInterface) MatchRule(ctx context.Context, project string, rule *config.Rule, args, auth map[string]interface{}, returnWhere model.ReturnWhereStub) (*model.PostProcess, error) {
	return nil, nil
}

func (m *mockAuthEventingInterface) CreateToken(ctx context.Context, tokenClaims model.TokenClaims) (string, error) {
	c := m.Called(ctx, tokenClaims)
	return c.String(0), c.Error(1)
}

func (m *mockAuthEventingInterface) IsEventingOpAuthorised(ctx context.Context, project, token string, event *model.QueueEventRequest) (model.RequestParams, error) {
	c := m.Called(ctx, project, token, event)
	return c.Get(0).(model.RequestParams), c.Error(1)
}

func (m *mockAuthEventingInterface) GetSCAccessToken(context.Context) (string, error) {
	c := m.Called()
	return mock.Anything, c.Error(1)
}

func (m *mockAuthEventingInterface) GetInternalAccessToken(context.Context) (string, error) {
	c := m.Called()
	return c.String(0), c.Error(1)
}

type mockSchemaEventingInterface struct {
	mock.Mock
}

func (m *mockSchemaEventingInterface) GetSchemaForDB(ctx context.Context, dbAlias, col, format string) ([]interface{}, error) {
	c := m.Called(ctx, dbAlias, col, format)
	return c.Get(0).([]interface{}), c.Error(1)
}

func (m *mockSchemaEventingInterface) SchemaInspection(ctx context.Context, dbAlias, project, col string) (string, error) {
	c := m.Called(ctx, dbAlias, project, col)
	return c.String(0), c.Error(1)
}

func (m *mockSchemaEventingInterface) CheckIfEventingIsPossible(dbAlias, col string, obj map[string]interface{}, isFind bool) (findForUpdate map[string]interface{}, present bool) {
	c := m.Called(dbAlias, col, obj, isFind)
	return map[string]interface{}{}, c.Bool(1)
}

func (m *mockSchemaEventingInterface) Parser(dbSchemas config.DatabaseSchemas) (model.Type, error) {
	c := m.Called(dbSchemas)
	return nil, c.Error(1)
}

func (m *mockSchemaEventingInterface) SchemaValidator(ctx context.Context, col string, collectionFields model.Fields, doc map[string]interface{}) (map[string]interface{}, error) {
	c := m.Called(col, collectionFields, doc)
	return nil, c.Error(1)
}

func (m *mockSchemaEventingInterface) SchemaModifyAll(ctx context.Context, dbAlias, logicalDBName string, tables config.DatabaseSchemas) error {
	c := m.Called(ctx, dbAlias, logicalDBName, tables)
	return c.Error(0)
}
func (m *mockSchemaEventingInterface) GetSchema(dbAlias, col string) (model.Fields, bool) {
	c := m.Called(dbAlias, col)
	return c.Get(0).(model.Fields), c.Bool(1)
}

type mockFileStoreEventingInterface struct {
	mock.Mock
}

func (m *mockFileStoreEventingInterface) DoesExists(ctx context.Context, project, token, path string) error {
	c := m.Called(ctx, project, token, path)
	return c.Error(0)
}
