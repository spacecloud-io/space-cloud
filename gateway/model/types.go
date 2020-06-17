package model

import (
	"context"
	"net/http"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// SchemaCrudInterface is an interface consisting of functions of schema module used by auth module
type SchemaCrudInterface interface {
	SetConfig(conf config.Crud, project string) error
	ValidateCreateOperation(dbType, col string, req *CreateRequest) error
	ValidateUpdateOperation(dbType, col, op string, updateDoc, find map[string]interface{}) error
	CrudPostProcess(ctx context.Context, dbAlias, col string, result interface{}) error
	AdjustWhereClause(dbAlias string, dbType utils.DBType, col string, find map[string]interface{}) error
}

// CrudAuthInterface is an interface consisting of functions of crud module used by auth module
type CrudAuthInterface interface {
	Read(ctx context.Context, dbAlias, col string, req *ReadRequest) (interface{}, error)
}

// SchemaEventingInterface is an interface consisting of functions of schema module used by eventing module
type SchemaEventingInterface interface {
	CheckIfEventingIsPossible(dbAlias, col string, obj map[string]interface{}, isFind bool) (findForUpdate map[string]interface{}, present bool)
	Parser(crud config.Crud) (Type, error)
	SchemaValidator(col string, collectionFields Fields, doc map[string]interface{}) (map[string]interface{}, error)
	SchemaModifyAll(ctx context.Context, dbAlias, logicalDBName string, tables map[string]*config.TableRule) error
}

// CrudEventingInterface is an interface consisting of functions of crud module used by Eventing module
type CrudEventingInterface interface {
	InternalCreate(ctx context.Context, dbAlias, project, col string, req *CreateRequest, isIgnoreMetrics bool) error
	InternalUpdate(ctx context.Context, dbAlias, project, col string, req *UpdateRequest) error
	Read(ctx context.Context, dbAlias, col string, req *ReadRequest) (interface{}, error)
}

// AuthEventingInterface is an interface consisting of functions of auth module used by Eventing module
type AuthEventingInterface interface {
	GetInternalAccessToken() (string, error)
	GetSCAccessToken() (string, error)
	IsEventingOpAuthorised(ctx context.Context, project, token string, event *QueueEventRequest) error
}

// FilestoreEventingInterface is an interface consisting of functions of Filestore module used by Eventing module
type FilestoreEventingInterface interface {
	DoesExists(ctx context.Context, project, token, path string) error
}

// AuthFilestoreInterface is an interface consisting of functions of auth module used by Filestore module
type AuthFilestoreInterface interface {
	IsFileOpAuthorised(ctx context.Context, project, token, path string, op utils.FileOpType, args map[string]interface{}) (*PostProcess, error)
}

// AuthFunctionInterface is an interface consisting of functions of auth module used by Function module
type AuthFunctionInterface interface {
	GetSCAccessToken() (string, error)
	Encrypt(value string) (string, error)
}

// EventingRealtimeInterface is an interface consisting of functions of Eventing module used by RealTime module
type EventingRealtimeInterface interface {
	SetRealtimeTriggers(eventingRules []config.EventingRule)
}

// AuthRealtimeInterface is an interface consisting of functions of auth module used by RealTime module
type AuthRealtimeInterface interface {
	IsReadOpAuthorised(ctx context.Context, project, dbType, col, token string, req *ReadRequest) (*PostProcess, int, error)
	PostProcessMethod(postProcess *PostProcess, result interface{}) error
	GetInternalAccessToken() (string, error)
	GetSCAccessToken() (string, error)
}

// CrudRealtimeInterface is an interface consisting of functions of crud module used by RealTime module
type CrudRealtimeInterface interface {
	Read(ctx context.Context, dbAlias, col string, req *ReadRequest) (interface{}, error)
}

// CrudSchemaInterface is an interface consisting of functions of crud module used by Schema module
type CrudSchemaInterface interface {
	GetDBType(dbAlias string) (string, error)
	// CreateProjectIfNotExists(ctx context.Context, project, dbAlias string) error
	RawBatch(ctx context.Context, dbAlias string, batchedQueries []string) error
	DescribeTable(ctx context.Context, dbAlias, col string) ([]utils.FieldType, []utils.ForeignKeysType, []utils.IndexType, error)
}

// CrudUserInterface is an interface consisting of functions of crud module used by User module
type CrudUserInterface interface {
	GetDBType(dbAlias string) (string, error)
	Read(ctx context.Context, dbAlias, col string, req *ReadRequest) (interface{}, error)
	Create(ctx context.Context, dbAlias, col string, req *CreateRequest) error
	Update(ctx context.Context, dbAlias, col string, req *UpdateRequest) error
}

// AuthUserInterface is an interface consisting of functions of auth module used by User module
type AuthUserInterface interface {
	IsReadOpAuthorised(ctx context.Context, project, dbType, col, token string, req *ReadRequest) (*PostProcess, int, error)
	PostProcessMethod(postProcess *PostProcess, result interface{}) error
	CreateToken(tokenClaims TokenClaims) (string, error)
	IsUpdateOpAuthorised(ctx context.Context, project, dbType, col, token string, req *UpdateRequest) (int, error)
}

// SyncmanEventingInterface is an interface consisting of functions of syncman module used by eventing module
type SyncmanEventingInterface interface {
	GetAssignedSpaceCloudURL(ctx context.Context, project string, token int) (string, error)
	GetAssignedTokens() (start, end int)
	GetEventSource() string
	GetSpaceCloudURLFromID(nodeID string) (string, error)
	GetNodeID() string
	MakeHTTPRequest(ctx context.Context, method, url, token, scToken string, params, vPtr interface{}) error
}

// AdminEventingInterface is an interface consisting of functions of admin module used by eventing module
type AdminEventingInterface interface {
	GetInternalAccessToken() (string, error)
}

// HTTPEventingInterface is an interface consisting of functions of a http client used by eventing module
type HTTPEventingInterface interface {
	Do(req *http.Request) (*http.Response, error)
}

// PostProcess filters the schema
type PostProcess struct {
	PostProcessAction []PostProcessAction
}

// PostProcessAction is struct of Action Field Value
type PostProcessAction struct {
	Action string
	Field  string
	Value  interface{}
}

// TokenClaims specifies the tokens and its claims
type TokenClaims map[string]interface{}

// Response is the object returned by every handler to client
type Response struct {
	Error  string      `json:"error,omitempty"`
	Result interface{} `json:"result,omitempty"`
}
