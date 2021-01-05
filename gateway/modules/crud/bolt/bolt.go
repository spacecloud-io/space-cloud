package bolt

import (
	"context"

	"github.com/spaceuptech/helpers"
	"go.etcd.io/bbolt"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// Bolt holds the bolt session
type Bolt struct {
	queryFetchLimit *int64
	enabled         bool
	connection      string
	bucketName      string
	client          *bbolt.DB
}

// Init initialises a new bolt instance
func Init(enabled bool, connection, bucketName string) (b *Bolt, err error) {
	b = &Bolt{enabled: enabled, connection: connection, bucketName: bucketName}

	if b.enabled {
		err = b.connect()
	}

	return
}

// Close gracefully the Bolt client
func (b *Bolt) Close() error {
	if b.client != nil {
		return b.client.Close()
	}

	return nil
}

// IsSame checks if we've got the same connection string
func (b *Bolt) IsSame(conn, dbName string, driverConf config.DriverConfig) bool {
	return b.connection == conn && dbName == b.bucketName // DriverConfig is not used for now.
}

// IsClientSafe checks whether database is enabled and connected
func (b *Bolt) IsClientSafe(ctx context.Context) error {
	if !b.enabled {
		return utils.ErrDatabaseDisabled
	}

	if b.client == nil {
		if err := b.connect(); err != nil {
			return helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to connect to bbbolt database", err, nil)
		}
	}

	return nil
}

func (b *Bolt) connect() error {
	if err := b.Close(); err != nil {
		return helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to close previous database connection in bbolt db", nil, nil)
	}

	client, err := bbolt.Open(b.connection, 0600, bbolt.DefaultOptions)
	if err != nil {
		return err
	}
	helpers.Logger.LogInfo(helpers.GetRequestID(context.TODO()), "Successfully connected to bbolt database", nil)
	b.client = client
	return nil
}

// GetDBType returns the dbType of the crud block
func (b *Bolt) GetDBType() model.DBType {
	return model.EmbeddedDB
}

// SetQueryFetchLimit sets data fetch limit
func (b *Bolt) SetQueryFetchLimit(limit int64) {
	b.queryFetchLimit = &limit
}

// SetProjectAESKey sets aes key
func (b *Bolt) SetProjectAESKey(aesKey []byte) {
}
