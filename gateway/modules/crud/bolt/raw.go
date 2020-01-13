package bolt

import (
	"context"
	"fmt"
)

func (b *Bolt) RawExec(ctx context.Context, project string) error {
	return fmt.Errorf("error raw exec cannot be performed over selected database")
}

func (b *Bolt) CreateProjectIfNotExist(ctx context.Context, project string) error {
	return fmt.Errorf("errocr create project operation cannot be performed over selected database")
}

func (b *Bolt) RawBatch(ctx context.Context, batchedQueries []string) error {
	return fmt.Errorf("error raw batchc cannot be performed over selected database")
}

// GetConnectionState : function to check connection state
func (b *Bolt) GetConnectionState(ctx context.Context) bool {
	if !b.enabled || b.client == nil {
		return false
	}

	// Ping to check if connection is established
	return b.connect() == nil
}
