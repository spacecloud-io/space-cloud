package model

import (
	"context"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

type SchemaAuthInterface interface {
	SetConfig(conf config.Crud, project string) error
	ValidateCreateOperation(dbType, col string, req *CreateRequest) error
	ValidateUpdateOperation(dbType, col, op string, updateDoc, find map[string]interface{}) error
}

type CrudAuthInterface interface {
	Read(ctx context.Context, dbAlias, project, col string, req *ReadRequest) (interface{}, error)
}
type SchemaEventingInterface interface {
	CheckIfEventingIsPossible(dbAlias, col string, obj map[string]interface{}, isFind bool) (findForUpdate map[string]interface{}, present bool)
	Parser(crud config.Crud) (Type, error)
	SchemaValidator(col string, collectionFields Fields, doc map[string]interface{}) (map[string]interface{}, error)
}
type CrudEventingInterface interface {
	InternalCreate(ctx context.Context, dbAlias, project, col string, req *CreateRequest) error
	InternalUpdate(ctx context.Context, dbAlias, project, col string, req *UpdateRequest) error
	Read(ctx context.Context, dbAlias, project, col string, req *ReadRequest) (interface{}, error)
}

type AuthEventingInterface interface {
	GetInternalAccessToken() (string, error)
	GetSCAccessToken() (string, error)
	IsEventingOpAuthorised(ctx context.Context, project, token string, event *QueueEventRequest) error
}

type FilestoreEventingInterface interface {
	DoesExists(ctx context.Context, project, token, path string) error
}

type AuthFilestoreInterface interface {
	IsFileOpAuthorised(ctx context.Context, project, token, path string, op utils.FileOpType, args map[string]interface{}) (*PostProcess, error)
}

type AuthFunctionInterface interface {
	GetSCAccessToken() (string, error)
}

type EventingRealtimeInterface interface {
	AddInternalRules(eventingRules []config.EventingRule)
}

type AuthRealtimeInterface interface {
	IsReadOpAuthorised(ctx context.Context, project, dbType, col, token string, req *ReadRequest) (*PostProcess, int, error)
	PostProcessMethod(postProcess *PostProcess, result interface{}) error
	GetInternalAccessToken() (string, error)
	GetSCAccessToken() (string, error)
}

type CrudRealtimeInterface interface {
	Read(ctx context.Context, dbAlias, project, col string, req *ReadRequest) (interface{}, error)
}

type CrudSchemaInterface interface {
	GetDBType(dbAlias string) (string, error)
	//CreateProjectIfNotExists(ctx context.Context, project, dbAlias string) error
	CreateDatabaseIfNotExist(ctx context.Context, project, dbAlias string) error
	RawBatch(ctx context.Context, dbAlias string, batchedQueries []string) error
	DescribeTable(ctx context.Context, dbAlias, project, col string) ([]utils.FieldType, []utils.ForeignKeysType, []utils.IndexType, error)
}

type CrudUserInterface interface {
	GetDBType(dbAlias string) (string, error)
	Read(ctx context.Context, dbAlias, project, col string, req *ReadRequest) (interface{}, error)
	Create(ctx context.Context, dbAlias, project, col string, req *CreateRequest) error
	Update(ctx context.Context, dbAlias, project, col string, req *UpdateRequest) error
}

type AuthUserInterface interface {
	IsReadOpAuthorised(ctx context.Context, project, dbType, col, token string, req *ReadRequest) (*PostProcess, int, error)
	PostProcessMethod(postProcess *PostProcess, result interface{}) error
	CreateToken(tokenClaims TokenClaims) (string, error)
	IsUpdateOpAuthorised(ctx context.Context, project, dbType, col, token string, req *UpdateRequest) (int, error)
}

type PostProcess struct {
	PostProcessAction []PostProcessAction
}

type PostProcessAction struct {
	Action string
	Field  string
	Value  interface{}
}

type TokenClaims map[string]interface{}
