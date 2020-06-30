package graphql_test

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/spaceuptech/space-cloud/gateway/model"
)

type mockGraphQLCrudInterface struct {
	mock.Mock
}

func (m *mockGraphQLCrudInterface) Create(ctx context.Context, dbAlias, collection string, request *model.CreateRequest) error {
	args := m.Called(ctx, dbAlias, collection, request)
	return args.Error(0)
}
func (m *mockGraphQLCrudInterface) Read(ctx context.Context, dbAlias, collection string, request *model.ReadRequest) (interface{}, error) {
	args := m.Called(ctx, dbAlias, collection, request)
	return args.Get(0), args.Error(1)
}
func (m *mockGraphQLCrudInterface) Update(ctx context.Context, dbAlias, collection string, request *model.UpdateRequest) error {
	args := m.Called(ctx, dbAlias, collection, request)
	return args.Error(0)
}
func (m *mockGraphQLCrudInterface) Delete(ctx context.Context, dbAlias, collection string, request *model.DeleteRequest) error {
	args := m.Called(ctx, dbAlias, collection, request)
	return args.Error(0)
}
func (m *mockGraphQLCrudInterface) Batch(ctx context.Context, dbAlias string, req *model.BatchRequest) error {
	args := m.Called(ctx, dbAlias, req)
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
func (m *mockGraphQLCrudInterface) ExecPreparedQuery(ctx context.Context, dbAlias, id string, req *model.PreparedQueryRequest, auth map[string]interface{}) (interface{}, error) {
	args := m.Called(ctx, dbAlias, id, req, auth)
	return args.Get(0).(interface{}), args.Error(1)
}

type mockGraphQLAuthInterface struct {
	mock.Mock
}

func (m *mockGraphQLAuthInterface) IsCreateOpAuthorised(ctx context.Context, project, dbAlias, col, token string, req *model.CreateRequest) (int, error) {
	args := m.Called(ctx, project, dbAlias, col, token, req)
	return args.Int(0), args.Error(1)
}
func (m *mockGraphQLAuthInterface) IsReadOpAuthorised(ctx context.Context, project, dbAlias, col, token string, req *model.ReadRequest) (*model.PostProcess, int, error) {
	args := m.Called(ctx, project, dbAlias, col, token, req)
	return args.Get(0).(*model.PostProcess), args.Int(1), args.Error(2)
}
func (m *mockGraphQLAuthInterface) IsUpdateOpAuthorised(ctx context.Context, project, dbAlias, col, token string, req *model.UpdateRequest) (int, error) {
	args := m.Called(ctx, project, dbAlias, col, token, req)
	return args.Int(0), args.Error(1)
}
func (m *mockGraphQLAuthInterface) IsDeleteOpAuthorised(ctx context.Context, project, dbAlias, col, token string, req *model.DeleteRequest) (int, error) {
	args := m.Called(ctx, project, dbAlias, col, token, req)
	return args.Int(0), args.Error(1)
}
func (m *mockGraphQLAuthInterface) IsFuncCallAuthorised(ctx context.Context, project, service, function, token string, params interface{}) (map[string]interface{}, error) {
	args := m.Called()
	return args.Get(0).(map[string]interface{}), args.Error(1)
}
func (m *mockGraphQLAuthInterface) PostProcessMethod(postProcess *model.PostProcess, result interface{}) error {
	args := m.Called()
	return args.Error(0)
}
func (m *mockGraphQLAuthInterface) IsPreparedQueryAuthorised(ctx context.Context, project, dbAlias, id, token string, req *model.PreparedQueryRequest) (*model.PostProcess, map[string]interface{}, int, error) {
	args := m.Called()
	return args.Get(0).(*model.PostProcess), args.Get(1).(map[string]interface{}), args.Int(2), args.Error(3)
}

type mockGraphQLFunctionInterface struct {
	mock.Mock
}

func (m *mockGraphQLFunctionInterface) CallWithContext(ctx context.Context, service, function, token string, auth, params interface{}) (interface{}, error) {
	args := m.Called(ctx, service, function, token, auth, params)
	return args.Get(0).(interface{}), args.Error(1)
}

type mockGraphQLSchemaInterface struct {
	mock.Mock
}

func (m *mockGraphQLSchemaInterface) GetSchema(dbAlias, col string) (model.Fields, bool) {
	args := m.Called(dbAlias, col)
	return args.Get(0).(model.Fields), args.Bool(1)
}
