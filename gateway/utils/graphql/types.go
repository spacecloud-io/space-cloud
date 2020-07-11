package graphql

import (
	"context"
	"github.com/spaceuptech/space-cloud/gateway/model"
)

// CrudInterface is an interface consisting of functions of crud module used by graphql module
type CrudInterface interface {
	Create(ctx context.Context, dbAlias, collection string, request *model.CreateRequest) error
	Read(ctx context.Context, dbAlias, collection string, request *model.ReadRequest) (interface{}, error)
	Update(ctx context.Context, dbAlias, collection string, request *model.UpdateRequest) error
	Delete(ctx context.Context, dbAlias, collection string, request *model.DeleteRequest) error
	Batch(ctx context.Context, dbAlias string, req *model.BatchRequest) error
	GetDBType(dbAlias string) (string, error)
	IsPreparedQueryPresent(directive, fieldName string) bool
	ExecPreparedQuery(ctx context.Context, dbAlias, id string, req *model.PreparedQueryRequest, auth map[string]interface{}) (interface{}, error)
}

// AuthInterface is an interface consisting of functions of auth module used by graphql module
type AuthInterface interface {
	IsCreateOpAuthorised(ctx context.Context, project, dbAlias, col, token string, req *model.CreateRequest) (int, error)
	IsReadOpAuthorised(ctx context.Context, project, dbAlias, col, token string, req *model.ReadRequest) (*model.PostProcess, int, error)
	IsUpdateOpAuthorised(ctx context.Context, project, dbAlias, col, token string, req *model.UpdateRequest) (int, error)
	IsDeleteOpAuthorised(ctx context.Context, project, dbAlias, col, token string, req *model.DeleteRequest) (int, error)
	IsFuncCallAuthorised(ctx context.Context, project, service, function, token string, params interface{}) (map[string]interface{}, error)
	PostProcessMethod(postProcess *model.PostProcess, result interface{}) error
	IsPreparedQueryAuthorised(ctx context.Context, project, dbAlias, id, token string, req *model.PreparedQueryRequest) (*model.PostProcess, map[string]interface{}, int, error)
}

// FunctionInterface is an interface consisting of functions of function module used by graphql module
type FunctionInterface interface {
	CallWithContext(ctx context.Context, service, function, token string, auth, params interface{}) (interface{}, error)
}

// SchemaInterface is an interface consisting of functions of schema module used by graphql module
type SchemaInterface interface {
	GetSchema(dbAlias, col string) (model.Fields, bool)
}
