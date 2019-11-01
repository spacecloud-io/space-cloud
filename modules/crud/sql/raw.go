package sql

import (
	"context"
	"database/sql"
)

// RawBatch performs a batch operation for schema creation
// NOTE: not to be exposed externally
func (s *SQL) RawBatch(ctx context.Context, queries []string) error {
	// Skip if length of queries == 0
	if len(queries) == 0 {
		return nil
	}

	tx, err := s.client.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	for _, query := range queries {
		_, err := tx.ExecContext(ctx, query)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

// RawExec performs an operation for schema creation
// NOTE: not to be exposed externally
func (s *SQL) RawExec(ctx context.Context, query string) error {
	_, err := s.client.ExecContext(ctx, query, []interface{}{}...)
	return err
}

// GetConnectionState : Function to get connection state
func (s *SQL) GetConnectionState(ctx context.Context, dbType string) bool {
	if !s.enabled || s.client == nil {
		return false
	}

	// Ping to check if connection is established
	err := s.client.PingContext(ctx)
	return err == nil
}
