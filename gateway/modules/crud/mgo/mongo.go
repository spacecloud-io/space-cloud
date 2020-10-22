package mgo

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/spaceuptech/helpers"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// Mongo holds the mongo session
type Mongo struct {
	enabled    bool
	connection string
	dbName     string
	client     *mongo.Client
	auth       model.AuthCrudInterface
}

// Init initialises a new mongo instance
func Init(enabled bool, connection, dbName string, auth model.AuthCrudInterface) (mongoStub *Mongo, err error) {
	mongoStub = &Mongo{dbName: dbName, enabled: enabled, connection: connection, client: nil, auth: auth}

	if mongoStub.enabled {
		err = mongoStub.connect()
	}

	return
}

// Close gracefully the Mongo client
func (m *Mongo) Close() error {
	if m.client != nil {
		return m.client.Disconnect(context.TODO())
	}

	return nil
}

// IsSame checks if we've got the same connection string
func (m *Mongo) IsSame(conn, dbName string) bool {
	return strings.HasPrefix(m.connection, conn) && dbName == m.dbName
}

// IsClientSafe checks whether database is enabled and connected
func (m *Mongo) IsClientSafe(context.Context) error {
	if !m.enabled {
		return utils.ErrDatabaseDisabled
	}

	if m.client == nil {
		if err := m.connect(); err != nil {
			helpers.Logger.LogInfo(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Error connecting to mongo %v", err.Error()), nil)
			return utils.ErrDatabaseConnection
		}
	}

	return nil
}

func (m *Mongo) connect() error {
	timeOut := 3 * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()

	client, err := mongo.NewClient(options.Client().ApplyURI(m.connection))
	if err != nil {
		return err
	}

	if err := client.Connect(ctx); err != nil {
		return err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return err
	}

	m.client = client
	return nil
}

// GetDBType returns the dbType of the crud block
func (m *Mongo) GetDBType() model.DBType {
	return model.Mongo
}
