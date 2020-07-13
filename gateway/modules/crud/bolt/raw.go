package bolt

import (
	"context"
	"errors"
	"fmt"
)

// RawQuery query document(s) from the database
func (b *Bolt) RawQuery(ctx context.Context, query string, args []interface{}) (int64, interface{}, error) {
	return 0, "", errors.New("error raw querry cannot be performed over embedded database")
}

// CreateDatabaseIfNotExist creates a project if none exist
func (b *Bolt) CreateDatabaseIfNotExist(ctx context.Context, project string) error {
	return fmt.Errorf("error create project operation cannot be performed over selected database")
}

// RawBatch performs a batch operation for schema creation
// NOTE: not to be exposed externally
func (b *Bolt) RawBatch(ctx context.Context, batchedQueries []string) error {
	return fmt.Errorf("error raw batchc cannot be performed over selected database")
}

// GetConnectionState : function to check connection state
func (b *Bolt) GetConnectionState(ctx context.Context) bool {
	if !b.enabled || b.client == nil {
		return false
	}

	return b.client.Info() != nil
}
