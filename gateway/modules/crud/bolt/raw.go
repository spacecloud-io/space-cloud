package bolt

import "context"

func (b *Bolt) RawExec(ctx context.Context, project string) error {
	return nil
}

func (b *Bolt) CreateProjectIfNotExist(ctx context.Context, project string) error {
	return nil
}

func (b *Bolt) RawBatch(ctx context.Context, batchedQueries []string) error {
	return nil
}

// GetConnectionState : function to check connection state
func (b *Bolt) GetConnectionState(ctx context.Context) bool {
	if !b.enabled || b.client == nil {
		return false
	}

	// Ping to check if connection is established
	err := b.connect()
	return err == nil
}
