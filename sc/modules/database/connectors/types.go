package connectors

import (
	"context"

	"github.com/spacecloud-io/space-cloud/config"
	"github.com/spacecloud-io/space-cloud/model"
)

// Connector abstracts the implementation of crud operations of databases
type Connector interface {
	Create(ctx context.Context, col string, req *model.CreateRequest) (int64, error)
	Read(ctx context.Context, col string, req *model.ReadRequest) (int64, interface{}, map[string]map[string]string, *model.SQLMetaData, error)
	Update(ctx context.Context, col string, req *model.UpdateRequest) (int64, error)
	Delete(ctx context.Context, col string, req *model.DeleteRequest) (int64, error)
	Aggregate(ctx context.Context, col string, req *model.AggregateRequest) (interface{}, error)
	Batch(ctx context.Context, req *model.BatchRequest) ([]int64, error)
	DescribeTable(ctc context.Context, col string) ([]model.InspectorFieldType, []model.IndexType, error)
	RawQuery(ctx context.Context, query string, isDebug bool, args []interface{}) (int64, interface{}, *model.SQLMetaData, error)
	GetCollections(ctx context.Context) ([]model.DatabaseCollections, error)
	DeleteCollection(ctx context.Context, col string) error
	CreateDatabaseIfNotExist(ctx context.Context, name string) error
	RawBatch(ctx context.Context, batchedQueries []string) error
	GetDBType() model.DBType
	IsClientSafe(ctx context.Context) error
	IsSame(conn, dbName string, driverConf config.DriverConfig) bool
	Close() error
	GetConnectionState(ctx context.Context) bool
	SetQueryFetchLimit(limit int64)
	SetProjectAESKey(aesKey []byte)
}
