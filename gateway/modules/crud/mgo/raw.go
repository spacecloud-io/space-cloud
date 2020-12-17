package mgo

import (
	"context"
	"errors"
	"fmt"

	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/model"
)

// RawBatch performs a batch operation for schema creation
// NOTE: not to be exposed externally
func (m *Mongo) RawBatch(ctx context.Context, queries []string) error {
	return errors.New("raw batch operation cannot be performed on mongo")
}

// RawQuery query document(s) from the database
func (m *Mongo) RawQuery(ctx context.Context, query string, args []interface{}) (int64, interface{}, *model.SQLMetaData, error) {
	return 0, "", nil, errors.New("error raw query operation cannot be performed on mongo")
}

// GetConnectionState : function to check connection state
func (m *Mongo) GetConnectionState(ctx context.Context) bool {
	if !m.enabled || m.client == nil {
		return false
	}

	// Ping to check if connection is established
	err := m.client.Ping(ctx, nil)
	if err != nil {
		_ = m.client.Disconnect(context.Background())
		m.client = nil
		_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to ping mongo database - %s", m.dbName), err, nil)
		return false
	}

	return true
}

// CreateDatabaseIfNotExist creates a database if not exist which has same name of project
func (m *Mongo) CreateDatabaseIfNotExist(ctx context.Context, project string) error {
	return errors.New("create project exists cannot be performed over mongo")
}
