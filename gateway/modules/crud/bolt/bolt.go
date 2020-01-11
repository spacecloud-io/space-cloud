package bolt

import (
	"github.com/boltdb/bolt"
	"github.com/sirupsen/logrus"

	"github.com/spaceuptech/space-cloud/gateway/utils"
)

type Bolt struct {
	enabled    bool
	connection string
	client     *bolt.DB
}

// Init initialises a new bolt instance
func Init(enabled bool, connection string) (mongoStub *Bolt, err error) {
	mongoStub = &Bolt{enabled: enabled, connection: connection, client: nil}

	if mongoStub.enabled {
		err = mongoStub.connect()
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
	client, err := bolt.Open(b.connection, 0600, bolt.DefaultOptions)
	if err != nil {
		return err
	}

	b.client = client
	return nil
}

// GetDBAlias returns the dbType of the crud block
func (b *Bolt) GetDBType() utils.DBType {
	return utils.BoltDb
}
