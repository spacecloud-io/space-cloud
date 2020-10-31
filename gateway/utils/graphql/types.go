package graphql

import (
	"context"

	"github.com/spaceuptech/space-cloud/gateway/model"
)

// CrudInterface is an interface consisting of functions of crud module used by graphql module
type CrudInterface interface {
	Create(ctx context.Context, dbAlias, collection string, request *model.CreateRequest, params model.RequestParams) error
	Read(ctx context.Context, dbAlias, collection string, request *model.ReadRequest, params model.RequestParams) (interface{}, error)
	Update(ctx context.Context, dbAlias, collection string, request *model.UpdateRequest, params model.RequestParams) error
	Delete(ctx context.Context, dbAlias, collection string, request *model.DeleteRequest, params model.RequestParams) error
	Batch(ctx context.Context, dbAlias string, req *model.BatchRequest, params model.RequestParams) error
	GetDBType(dbAlias string) (string, error)
	IsPreparedQueryPresent(directive, fieldName string) bool
	ExecPreparedQuery(ctx context.Context, dbAlias, id string, req *model.PreparedQueryRequest, params model.RequestParams) (interface{}, error)
}

// AuthInterface is an interface consisting of functions of auth module used by graphql module
type AuthInterface interface {
	ParseToken(ctx context.Context, token string) (map[string]interface{}, error)
	IsCreateOpAuthorised(ctx context.Context, project, dbAlias, col, token string, req *model.CreateRequest) (model.RequestParams, error)
	IsReadOpAuthorised(ctx context.Context, project, dbAlias, col, token string, req *model.ReadRequest, stub model.ReturnWhereStub) (*model.PostProcess, model.RequestParams, error)
	IsUpdateOpAuthorised(ctx context.Context, project, dbAlias, col, token string, req *model.UpdateRequest) (model.RequestParams, error)
	IsDeleteOpAuthorised(ctx context.Context, project, dbAlias, col, token string, req *model.DeleteRequest) (model.RequestParams, error)
	IsFuncCallAuthorised(ctx context.Context, project, service, function, token string, params interface{}) (*model.PostProcess, model.RequestParams, error)
	PostProcessMethod(ctx context.Context, postProcess *model.PostProcess, result interface{}) error
	IsPreparedQueryAuthorised(ctx context.Context, project, dbAlias, id, token string, req *model.PreparedQueryRequest) (*model.PostProcess, model.RequestParams, error)
}

// FunctionInterface is an interface consisting of functions of function module used by graphql module
type FunctionInterface interface {
	CallWithContext(ctx context.Context, service, function, token string, reqParams model.RequestParams, params interface{}) (int, interface{}, error)
}

// SchemaInterface is an interface consisting of functions of schema module used by graphql module
type SchemaInterface interface {
	GetSchema(dbAlias, col string) (model.Fields, bool)
}
