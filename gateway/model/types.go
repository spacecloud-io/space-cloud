package model

import (
	"context"
	"net/http"

	"github.com/spaceuptech/space-cloud/gateway/config"
)

// SchemaCrudInterface is an interface consisting of functions of schema module used by auth module
type SchemaCrudInterface interface {
	SetDatabaseSchema(dbSchemas config.DatabaseSchemas, project string) error
	ValidateCreateOperation(ctx context.Context, dbType, col string, req *CreateRequest) error
	ValidateUpdateOperation(ctx context.Context, dbType, col, op string, updateDoc, find map[string]interface{}) error
	CrudPostProcess(ctx context.Context, dbAlias, col string, result interface{}) error
	AdjustWhereClause(ctx context.Context, dbAlias string, dbType DBType, col string, find map[string]interface{}) error
}

// CrudAuthInterface is an interface consisting of functions of crud module used by auth module
type CrudAuthInterface interface {
	Read(ctx context.Context, dbAlias, col string, req *ReadRequest, params RequestParams) (interface{}, *SQLMetaData, error)
}

// SchemaEventingInterface is an interface consisting of functions of schema module used by eventing module
type SchemaEventingInterface interface {
	SchemaModifyAll(ctx context.Context, dbAlias, logicalDBName string, dbSchemas config.DatabaseSchemas) error
	SchemaInspection(ctx context.Context, dbAlias, project, col string, realSchema Collection) (string, error)
	GetSchema(dbAlias, col string) (Fields, bool)
	GetSchemaForDB(ctx context.Context, dbAlias, col, format string) ([]interface{}, error)
}

// CrudEventingInterface is an interface consisting of functions of crud module used by Eventing module
type CrudEventingInterface interface {
	InternalCreate(ctx context.Context, dbAlias, project, col string, req *CreateRequest, isIgnoreMetrics bool) error
	InternalUpdate(ctx context.Context, dbAlias, project, col string, req *UpdateRequest) error
	Read(ctx context.Context, dbAlias, col string, req *ReadRequest, params RequestParams) (interface{}, *SQLMetaData, error)
	GetDBType(dbAlias string) (string, error)
	GetSchema(dbAlias, col string) (Fields, bool)
}

// AuthEventingInterface is an interface consisting of functions of auth module used by Eventing module
type AuthEventingInterface interface {
	CreateToken(ctx context.Context, tokenClaims TokenClaims) (string, error)
	GetInternalAccessToken(ctx context.Context) (string, error)
	GetSCAccessToken(ctx context.Context) (string, error)
	IsEventingOpAuthorised(ctx context.Context, project, token string, event *QueueEventRequest) (RequestParams, error)
	MatchRule(ctx context.Context, project string, rule *config.Rule, args, auth map[string]interface{}, returnWhere ReturnWhereStub) (*PostProcess, error)
}

// AuthSyncManInterface is an interface consisting of functions of auth module used by sync man
type AuthSyncManInterface interface {
	GetIntegrationToken(ctx context.Context, id string) (string, error)
	GetMissionControlToken(ctx context.Context, claims map[string]interface{}) (string, error)
}

// FilestoreEventingInterface is an interface consisting of functions of Filestore module used by Eventing module
type FilestoreEventingInterface interface {
	DoesExists(ctx context.Context, project, token, path string) error
}

// AuthFilestoreInterface is an interface consisting of functions of auth module used by Filestore module
type AuthFilestoreInterface interface {
	IsFileOpAuthorised(ctx context.Context, project, token, path string, op FileOpType, args map[string]interface{}) (*PostProcess, error)
}

// AuthCrudInterface is an interface consisting of functions of auth module used by crud module
type AuthCrudInterface interface {
	PostProcessMethod(ctx context.Context, postProcess *PostProcess, result interface{}) error
}

// AuthFunctionInterface is an interface consisting of functions of auth module used by Function module
type AuthFunctionInterface interface {
	GetSCAccessToken(ctx context.Context) (string, error)
	Encrypt(value string) (string, error)
	CreateToken(ctx context.Context, tokenClaims TokenClaims) (string, error)
}

// EventingRealtimeInterface is an interface consisting of functions of Eventing module used by RealTime module
type EventingRealtimeInterface interface {
	SetRealtimeTriggers(eventingRules []*config.EventingTrigger)
}

// AuthRealtimeInterface is an interface consisting of functions of auth module used by RealTime module
type AuthRealtimeInterface interface {
	IsReadOpAuthorised(ctx context.Context, project, dbType, col, token string, req *ReadRequest, stub ReturnWhereStub) (*PostProcess, RequestParams, error)
	GetInternalAccessToken(ctx context.Context) (string, error)
	GetSCAccessToken(ctx context.Context) (string, error)
}

// CrudRealtimeInterface is an interface consisting of functions of crud module used by RealTime module
type CrudRealtimeInterface interface {
	Read(ctx context.Context, dbAlias, col string, req *ReadRequest, param RequestParams) (interface{}, *SQLMetaData, error)
}

// CrudSchemaInterface is an interface consisting of functions of crud module used by Schema module
type CrudSchemaInterface interface {
	GetDBType(dbAlias string) (string, error)
	// CreateProjectIfNotExists(ctx context.Context, project, dbAlias string) error
	RawBatch(ctx context.Context, dbAlias string, batchedQueries []string) error
	DescribeTable(ctx context.Context, dbAlias, col string) ([]InspectorFieldType, []IndexType, error)
}

// CrudUserInterface is an interface consisting of functions of crud module used by User module
type CrudUserInterface interface {
	GetDBType(dbAlias string) (string, error)
	Read(ctx context.Context, dbAlias, col string, req *ReadRequest, params RequestParams) (interface{}, *SQLMetaData, error)
	Create(ctx context.Context, dbAlias, col string, req *CreateRequest, params RequestParams) error
	Update(ctx context.Context, dbAlias, col string, req *UpdateRequest, params RequestParams) error
}

// AuthUserInterface is an interface consisting of functions of auth module used by User module
type AuthUserInterface interface {
	IsReadOpAuthorised(ctx context.Context, project, dbType, col, token string, req *ReadRequest, stub ReturnWhereStub) (*PostProcess, RequestParams, error)
	CreateToken(ctx context.Context, tokenClaims TokenClaims) (string, error)
	IsUpdateOpAuthorised(ctx context.Context, project, dbType, col, token string, req *UpdateRequest) (RequestParams, error)
}

// SyncmanEventingInterface is an interface consisting of functions of syncman module used by eventing module
type SyncmanEventingInterface interface {
	GetAssignedSpaceCloudID(ctx context.Context, project string, token int) (string, error)
	GetAssignedTokens() (start, end int)
	GetEventSource() string
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

// ReturnWhereStub describes return where stuff
type ReturnWhereStub struct {
	Where         map[string]interface{}
	ReturnWhere   bool
	Col           string
	PrefixColName bool
}
