package bolt

import (
	"context"
	"errors"

	"github.com/spaceuptech/helpers"
)

// RawQuery query document(s) from the database
func (b *Bolt) RawQuery(ctx context.Context, query string, args []interface{}) (int64, interface{}, error) {
	return 0, "", errors.New("error raw querry cannot be performed over embedded database")
}

// CreateDatabaseIfNotExist creates a project if none exist
func (b *Bolt) CreateDatabaseIfNotExist(ctx context.Context, project string) error {
	return helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to create database operation cannot be performed over selected database", nil, nil)
}

// RawBatch performs a batch operation for schema creation
// NOTE: not to be exposed externally
func (b *Bolt) RawBatch(ctx context.Context, batchedQueries []string) error {
	return helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to create raw batch operation cannot be performed over selected database", nil, nil)
}

// GetConnectionState : function to check connection state
func (b *Bolt) GetConnectionState(ctx context.Context) bool {
	if !b.enabled || b.client == nil {
		return false
	}

	return b.client.Info() != nil
}
