package mgo

import (
	"context"
	"errors"
	"fmt"

	"github.com/spaceuptech/helpers"
)

// RawBatch performs a batch operation for schema creation
// NOTE: not to be exposed externally
func (m *Mongo) RawBatch(ctx context.Context, queries []string) error {
	return errors.New("raw batch operation cannot be performed on mongo")
}

// RawQuery query document(s) from the database
func (m *Mongo) RawQuery(ctx context.Context, query string, args []interface{}) (int64, interface{}, error) {
	return 0, "", errors.New("error raw querry operation cannot be performed on mongo")
}

// GetConnectionState : function to check connection state
func (m *Mongo) GetConnectionState(ctx context.Context) bool {
	if !m.enabled || m.getClient() == nil {
		return false
	}

	// Ping to check if connection is established
	err := m.getClient().Ping(ctx, nil)
	if err != nil {
		_ = m.getClient().Disconnect(context.Background())
		_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to ping mongo database - %s", m.dbName), err, nil)
		return false
	}

	return true
}

// CreateDatabaseIfNotExist creates a database if not exist which has same name of project
func (m *Mongo) CreateDatabaseIfNotExist(ctx context.Context, project string) error {
	return errors.New("create project exists cannot be performed over mongo")
}
