package model

import (
	"context"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

//SchemaAuthInterface is an interface consisting of functions of schema module used by Auth module
type SchemaAuthInterface interface {
	SetConfig(conf config.Crud, project string) error
	ValidateCreateOperation(dbType, col string, req *CreateRequest) error
	ValidateUpdateOperation(dbType, col, op string, updateDoc, find map[string]interface{}) error
}

//CrudAuthInterface is an interface consisting of functions of crud module used by Auth module
type CrudAuthInterface interface {
	Read(ctx context.Context, dbAlias, project, col string, req *ReadRequest) (interface{}, error)
}

//SchemaEventingInterface is an interface consisting of functions of schema module used by eventing module
type SchemaEventingInterface interface {
	CheckIfEventingIsPossible(dbAlias, col string, obj map[string]interface{}, isFind bool) (findForUpdate map[string]interface{}, present bool)
	Parser(crud config.Crud) (Type, error)
	SchemaValidator(col string, collectionFields Fields, doc map[string]interface{}) (map[string]interface{}, error)
}

//CrudEventingInterface is an interface consisting of functions of crud module used by Eventing module
type CrudEventingInterface interface {
	InternalCreate(ctx context.Context, dbAlias, project, col string, req *CreateRequest, isIgnoreMetrics bool) error
	InternalUpdate(ctx context.Context, dbAlias, project, col string, req *UpdateRequest) error
	Read(ctx context.Context, dbAlias, project, col string, req *ReadRequest) (interface{}, error)
}

//AuthEventingInterface is an interface consisting of functions of Auth module used by Eventing module
type AuthEventingInterface interface {
	GetInternalAccessToken() (string, error)
	GetSCAccessToken() (string, error)
	IsEventingOpAuthorised(ctx context.Context, project, token string, event *QueueEventRequest) error
}

//FilestoreEventingInterface is an interface consisting of functions of Filestore module used by Eventing module
type FilestoreEventingInterface interface {
	DoesExists(ctx context.Context, project, token, path string) error
}

//AuthFilestoreInterface is an interface consisting of functions of Auth module used by Filestore module
type AuthFilestoreInterface interface {
	IsFileOpAuthorised(ctx context.Context, project, token, path string, op utils.FileOpType, args map[string]interface{}) (*PostProcess, error)
}

//AuthFunctionInterface is an interface consisting of functions of Auth module used by Function module
type AuthFunctionInterface interface {
	GetSCAccessToken() (string, error)
}

//EventingRealtimeInterface is an interface consisting of functions of Eventing module used by RealTime module
type EventingRealtimeInterface interface {
	AddInternalRules(eventingRules []config.EventingRule)
}

//AuthRealtimeInterface is an interface consisting of functions of Auth module used by RealTime module
type AuthRealtimeInterface interface {
	IsReadOpAuthorised(ctx context.Context, project, dbType, col, token string, req *ReadRequest) (*PostProcess, int, error)
	PostProcessMethod(postProcess *PostProcess, result interface{}) error
	GetInternalAccessToken() (string, error)
	GetSCAccessToken() (string, error)
}

//CrudRealtimeInterface is an interface consisting of functions of Crud module used by RealTime module
type CrudRealtimeInterface interface {
	Read(ctx context.Context, dbAlias, project, col string, req *ReadRequest) (interface{}, error)
}

//CrudSchemaInterface is an interface consisting of functions of Crud module used by Schema module
type CrudSchemaInterface interface {
	GetDBType(dbAlias string) (string, error)
	//CreateProjectIfNotExists(ctx context.Context, project, dbAlias string) error
	CreateDatabaseIfNotExist(ctx context.Context, project, dbAlias string) error
	RawBatch(ctx context.Context, dbAlias string, batchedQueries []string) error
	DescribeTable(ctx context.Context, dbAlias, project, col string) ([]utils.FieldType, []utils.ForeignKeysType, []utils.IndexType, error)
}

//CrudUserInterface is an interface consisting of functions of Crud module used by User module
type CrudUserInterface interface {
	GetDBType(dbAlias string) (string, error)
	Read(ctx context.Context, dbAlias, project, col string, req *ReadRequest) (interface{}, error)
	Create(ctx context.Context, dbAlias, project, col string, req *CreateRequest) error
	Update(ctx context.Context, dbAlias, project, col string, req *UpdateRequest) error
}

//AuthUserInterface is an interface consisting of functions of Auth module used by User module
type AuthUserInterface interface {
	IsReadOpAuthorised(ctx context.Context, project, dbType, col, token string, req *ReadRequest) (*PostProcess, int, error)
	PostProcessMethod(postProcess *PostProcess, result interface{}) error
	CreateToken(tokenClaims TokenClaims) (string, error)
	IsUpdateOpAuthorised(ctx context.Context, project, dbType, col, token string, req *UpdateRequest) (int, error)
}

//PostProcess filters the schema
type PostProcess struct {
	PostProcessAction []PostProcessAction
}

//PostProcessAction is struct of Action Field Value
type PostProcessAction struct {
	Action string
	Field  string
	Value  interface{}
}

//TokenClaims specifies the tokens and its claims
type TokenClaims map[string]interface{}
