package bolt

import (
	"fmt"
	"log"

	"github.com/sirupsen/logrus"
	"go.etcd.io/bbolt"

	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// Bolt holds the bolt session
type Bolt struct {
	enabled    bool
	connection string
	bucketName string
	client     *bbolt.DB
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
func (b *Bolt) IsSame(conn, dbName string) bool {
	return b.connection == conn && dbName == b.bucketName
}

// IsClientSafe checks whether database is enabled and connected
func (b *Bolt) IsClientSafe() error {
	if !b.enabled {
		return utils.ErrDatabaseDisabled
	}

	if b.client == nil {
		if err := b.connect(); err != nil {
			logrus.Errorf("Error connecting to bboltdb - %v", err)
			return utils.ErrDatabaseConnection
		}
	}

	return nil
}

func (b *Bolt) connect() error {
	if err := b.Close(); err != nil {
		return fmt.Errorf("error closing previous database connection in bbolt db")
	}

	client, err := bbolt.Open(b.connection, 0600, bbolt.DefaultOptions)
	if err != nil {
		return err
	}
	log.Println("successfully connected to bbolt database")
	b.client = client
	return nil
}

// GetDBType returns the dbType of the crud block
func (b *Bolt) GetDBType() utils.DBType {
	return utils.EmbeddedDB
}
