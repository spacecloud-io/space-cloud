package mgo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/spaceuptech/space-cloud/utils"
)

// Mongo holds the mongo session
type Mongo struct {
	client  *mongo.Client
	timeOut time.Duration
}

// Init initialises a new mongo instance
func Init(connection string) (*Mongo, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(connection))
	if err != nil {
		return nil, err
	}

	timeOut := 5 * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &Mongo{client, timeOut}, nil
}

// Close gracefully the Mongo client
func (m *Mongo) Close() error {
	return m.client.Disconnect(context.TODO())
}

// GetDBType returns the dbType of the crud block
func (m *Mongo) GetDBType() utils.DBType {
	return utils.Mongo
}
