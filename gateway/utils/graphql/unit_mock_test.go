package graphql_test

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/spaceuptech/space-cloud/gateway/model"
)

type mockGraphQLCrudInterface struct {
	mock.Mock
}

func (m *mockGraphQLCrudInterface) Create(ctx context.Context, dbAlias, collection string, request *model.CreateRequest, params model.RequestParams) error {
	args := m.Called(ctx, dbAlias, collection, request, params)
	return args.Error(0)
}
func (m *mockGraphQLCrudInterface) Read(ctx context.Context, dbAlias, collection string, request *model.ReadRequest, params model.RequestParams) (interface{}, error) {
	args := m.Called(ctx, dbAlias, collection, request, params)
	return args.Get(0), args.Error(1)
}
func (m *mockGraphQLCrudInterface) Update(ctx context.Context, dbAlias, collection string, request *model.UpdateRequest, params model.RequestParams) error {
	args := m.Called(ctx, dbAlias, collection, request, params)
	return args.Error(0)
}
func (m *mockGraphQLCrudInterface) Delete(ctx context.Context, dbAlias, collection string, request *model.DeleteRequest, params model.RequestParams) error {
	args := m.Called(ctx, dbAlias, collection, request, params)
	return args.Error(0)
}
func (m *mockGraphQLCrudInterface) Batch(ctx context.Context, dbAlias string, req *model.BatchRequest, params model.RequestParams) error {
	args := m.Called(ctx, dbAlias, req, params)
	return args.Error(0)
}
func (m *mockGraphQLCrudInterface) GetDBType(dbAlias string) (string, error) {
	args := m.Called(dbAlias)
	return args.String(0), args.Error(1)
}
func (m *mockGraphQLCrudInterface) IsPreparedQueryPresent(directive, fieldName string) bool {
	args := m.Called(directive, fieldName)
	return args.Bool(0)
}
func (m *mockGraphQLCrudInterface) ExecPreparedQuery(ctx context.Context, dbAlias, id string, req *model.PreparedQueryRequest, params model.RequestParams) (interface{}, error) {
	args := m.Called(ctx, dbAlias, id, req, params)
	return args.Get(0).(interface{}), args.Error(1)
}

type mockGraphQLAuthInterface struct {
	mock.Mock
}

func (m *mockGraphQLAuthInterface) ParseToken(ctx context.Context, token string) (map[string]interface{}, error) {
	args := m.Called(ctx, token)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *mockGraphQLAuthInterface) IsCreateOpAuthorised(ctx context.Context, project, dbAlias, col, token string, req *model.CreateRequest) (model.RequestParams, error) {
	args := m.Called(ctx, project, dbAlias, col, token, req)
	return args.Get(0).(model.RequestParams), args.Error(1)
}
func (m *mockGraphQLAuthInterface) IsReadOpAuthorised(ctx context.Context, project, dbAlias, col, token string, req *model.ReadRequest, stub model.ReturnWhereStub) (*model.PostProcess, model.RequestParams, error) {
	args := m.Called(ctx, project, dbAlias, col, token, req)
	return args.Get(0).(*model.PostProcess), args.Get(1).(model.RequestParams), args.Error(2)
}
func (m *mockGraphQLAuthInterface) IsUpdateOpAuthorised(ctx context.Context, project, dbAlias, col, token string, req *model.UpdateRequest) (model.RequestParams, error) {
	args := m.Called(ctx, project, dbAlias, col, token, req)
	return args.Get(0).(model.RequestParams), args.Error(1)
}
func (m *mockGraphQLAuthInterface) IsDeleteOpAuthorised(ctx context.Context, project, dbAlias, col, token string, req *model.DeleteRequest) (model.RequestParams, error) {
	args := m.Called(ctx, project, dbAlias, col, token, req)
	return args.Get(0).(model.RequestParams), args.Error(1)
}
func (m *mockGraphQLAuthInterface) IsFuncCallAuthorised(ctx context.Context, project, service, function, token string, params interface{}) (*model.PostProcess, model.RequestParams, error) {
	args := m.Called(ctx, project, service, function, token, params)
	return args.Get(0).(*model.PostProcess), args.Get(1).(model.RequestParams), args.Error(2)
}
func (m *mockGraphQLAuthInterface) PostProcessMethod(ctx context.Context, postProcess *model.PostProcess, result interface{}) error {
	args := m.Called(postProcess, result)
	return args.Error(0)
}
func (m *mockGraphQLAuthInterface) IsPreparedQueryAuthorised(ctx context.Context, project, dbAlias, id, token string, req *model.PreparedQueryRequest) (*model.PostProcess, model.RequestParams, error) {
	args := m.Called(ctx, project, dbAlias, id, token, req)
	return args.Get(0).(*model.PostProcess), args.Get(1).(model.RequestParams), args.Error(2)
}

type mockGraphQLFunctionInterface struct {
	mock.Mock
}

func (m *mockGraphQLFunctionInterface) CallWithContext(ctx context.Context, service, function, token string, reqParams model.RequestParams, params interface{}) (int, interface{}, error) {
	args := m.Called(ctx, service, function, token, reqParams, params)
	return 0, args.Get(0).(interface{}), args.Error(1)
}

type mockGraphQLSchemaInterface struct {
	mock.Mock
}

func (m *mockGraphQLSchemaInterface) GetSchema(dbAlias, col string) (model.Fields, bool) {
	args := m.Called(dbAlias, col)
	return args.Get(0).(model.Fields), args.Bool(1)
}
